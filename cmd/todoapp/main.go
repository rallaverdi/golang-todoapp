package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/rallaverdi/golang-todoapp/docs"
	todov1 "github.com/rallaverdi/golang-todoapp/gen/go/todo/v1"
	core_config "github.com/rallaverdi/golang-todoapp/internal/core/config"
	core_logger "github.com/rallaverdi/golang-todoapp/internal/core/logger"
	core_metrics "github.com/rallaverdi/golang-todoapp/internal/core/metrics"
	core_pgx_pool "github.com/rallaverdi/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_redis "github.com/rallaverdi/golang-todoapp/internal/core/repository/redis"
	core_grpc_server "github.com/rallaverdi/golang-todoapp/internal/core/transport/grpc/server"
	core_http_middleware "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/rallaverdi/golang-todoapp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/rallaverdi/golang-todoapp/internal/features/statistics/service"
	statistics_transport_http "github.com/rallaverdi/golang-todoapp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/rallaverdi/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/rallaverdi/golang-todoapp/internal/features/tasks/service"
	tasks_transport_grpc "github.com/rallaverdi/golang-todoapp/internal/features/tasks/transport/grpc"
	tasks_transport_http "github.com/rallaverdi/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/rallaverdi/golang-todoapp/internal/features/users/repository/postgres"
	users_redis_cache "github.com/rallaverdi/golang-todoapp/internal/features/users/repository/redis"
	users_service "github.com/rallaverdi/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/rallaverdi/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// @title Golang 	Todo API
// @version 		1.0
// @description 	Todo Application REST-API schema
// @host 			127.0.0.1:5050
// @BasePath 		/api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("application time zone", zap.Any("zone", time.Local))
	logger.Debug("initializing postgres connection pool")

	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing redis cache")
	redisClient, err := core_redis.NewRedisClient(ctx, core_redis.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init redis cache", zap.Error(err))
	}
	defer redisClient.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	filterCache := users_redis_cache.NewFilterCache(redisClient, time.Minute*2)
	usersService := users_service.NewUsersService(usersRepository, filterCache)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)
	tasksTransportGRPC := tasks_transport_grpc.NewTasksGRPCHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing application metrics")
	appMetrics := core_metrics.NewMetrics()

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1, appMetrics.HTTPMiddleware(string(core_http_server.ApiVersion1)))
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)
	/*
			Example of usage apoVersionRouterV2 with separate middlewares

		//apoVersionRouterV2 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion2, core_http_middleware.Dummy("api v2 middleware"))
		//apoVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)
		//httpServer.RegisterAPIRouters(apoVersionRouterV2)
	*/
	httpServer.RegisterMetrics(appMetrics.Handler())
	httpServer.RegisterSwagger()

	logger.Debug("initializing gRPC server")
	grpcServer := core_grpc_server.NewGRPCServer(
		core_grpc_server.NewConfigMust(),
		logger,
	)

	todov1.RegisterTaskServiceServer(
		grpcServer,
		tasksTransportGRPC,
	)

	serversGroup, serversContext := errgroup.WithContext(ctx)

	serversGroup.Go(func() error {
		return httpServer.Run(serversContext)
	})

	serversGroup.Go(func() error {
		return grpcServer.Run(serversContext)
	})

	if err := serversGroup.Wait(); err != nil {
		logger.Error("application server error", zap.Error(err))
	}

}
