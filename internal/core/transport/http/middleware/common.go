package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/rallaverdi/golang-todoapp/internal/core/logger"
	core_http_response "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r) // call original handler
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)
			ctx := core_logger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle http request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)
			before := time.Now()

			log.Debug(
				">>> incoming http request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done http request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:5050": {},
			}

			origin := r.Header.Get("Origin")
			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

/*
Если мы хотим всегда знать айди запроса, то пришлось бы копировать эту логику в каждый хендлер.
Но хендлеров у нас может быть много, поэтому, чтобы не копировать одно и тоже пишем мидлвар.
Мидлвар - это отдельная функция принимющая любой HTTP хендлер который мы в нее передаем, и в мидлваре описана логика,
которую нужно выполнить до и после выполнения основного  HTTP хендлера. Таким образом нам не надо делать рутину
в хендлерах. Флоу такой - HTTP обработчик (хендлер) оборачивается в миддлвар, когда приходит запрос в хендлер
он попададает в миддлвар, пишется заголовок, потом вызывается сам хендлер, и на выходе снова попадаем в миддлвар.
Запрос может быть обёрнут в несколько миддлваров.
Схема примерно такая |M|-> M|-> M|->|H|->|M|->|M|->|M| (M- middleware, H- HTTP handler)



Итого приходит HTTP запрос:
1) попадаем в мидлвар RequestID() - в ней мы пытаемся либо получить айди запрсоа либо дообогощаем его сами

2) попадаем в мидлвар Logger() - после того как айди запроса получен или дообогащен, мы получаем
пре-конфигурированный логгер который пишет айди запроса и урл, и прокидывается дальше через контекст

3) попадаем в мидлвар Panic() - тут мы пытаемся отловить панику если она возникла во время вызова HTTP хендлера

4) попадаем в мидлвар Trace() - тут мы уже занем что паник нет, и логируем что запрос был, пишем метрики запроса
которые получили до выполнения HTTP хендлера

5) попадаем в сам хендлер который имеет уже кучу данных из мидлваров

6) хендлер отдает ответ и мы снова попадаем в мидлвар Trace() и пишем метрики по окончанию запроса

7) идем назад вплоть до 1ого мидлвара Trace() -> Panic() -> Logger() -> RequestID()
*/
