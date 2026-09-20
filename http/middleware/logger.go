package middleware

import (
	"net/http"
	"time"

	apphttp "go.chrastecky.dev/repolock/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Logger(log *zap.Logger) apphttp.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			t1 := time.Now()
			wrappedWriter := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)
			defer func() {
				log.Info("HTTP request",
					zap.Int("status", wrappedWriter.Status()),
					zap.Int("bytesWritten", wrappedWriter.BytesWritten()),
					zap.Duration("duration", time.Since(t1)),
					zap.String("request", request.RequestURI),
					zap.String("remote", request.RemoteAddr),
					zap.String("method", request.Method),
					zap.String("protocol", request.Proto),
				)
			}()

			next.ServeHTTP(wrappedWriter, request)
		})
	}
}
