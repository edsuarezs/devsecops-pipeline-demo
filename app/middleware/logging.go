package middleware

import (
	"log/slog"
	"net/http"
	"time"
	chimw "github.com/go-chi/chi/v5/middleware"

)

// Logger is a structured logging middleware using slog.
// Logs method, path, status, duration, and request ID for every request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"bytes_written", ww.BytesWritten(),
			"request_id", chimw.GetReqID(r.Context()),
			"remote_addr", r.RemoteAddr,
		)
	})
}
