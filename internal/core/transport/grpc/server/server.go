package core_grpc_server

import (
	"context"
	"fmt"
	"net"
	"time"

	core_logger "github.com/rallaverdi/golang-todoapp/internal/core/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	config Config
	log    *core_logger.Logger
	server *grpc.Server
}

func NewGRPCServer(
	config Config,
	log *core_logger.Logger,
	options ...grpc.ServerOption,
) *GRPCServer {
	return &GRPCServer{
		config: config,
		log:    log,
		server: grpc.NewServer(options...),
	}
}

func (s *GRPCServer) RegisterService(
	serviceDescription *grpc.ServiceDesc,
	implementation any,
) {
	s.server.RegisterService(serviceDescription, implementation)
}

func (s *GRPCServer) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return fmt.Errorf(
			"listen on gRPC address %q: %w",
			s.config.Addr,
			err,
		)
	}

	serveResult := make(chan error, 1)

	go func() {
		s.log.Warn(
			"starting gRPC server",
			zap.String("addr", s.config.Addr),
		)

		serveResult <- s.server.Serve(listener)
	}()

	select {
	case err := <-serveResult:
		if err != nil {
			return fmt.Errorf("serve gRPC: %w", err)
		}

		return nil

	case <-ctx.Done():
		s.log.Warn("shutting down gRPC server gracefully")
	}

	gracefulStopCompleted := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(gracefulStopCompleted)
	}()

	timer := time.NewTimer(s.config.ShutdownTimeout)
	defer timer.Stop()

	select {
	case <-gracefulStopCompleted:
		s.log.Warn("gRPC server stopped gracefully")

	case <-timer.C:
		s.log.Warn("gRPC graceful shutdown timeout exceeded")
		s.server.Stop()
		<-gracefulStopCompleted
	}

	return nil
}
