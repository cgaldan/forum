package middleware

import (
	"net/http"
	"real-time-forum/packages/logger"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func LoggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			start := time.Now()

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			logger.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"ip", getIP(r),
				"duration", duration.String(),
			)
		})
	}
}
