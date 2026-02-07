// Package middleware はHTTP APIのミドルウェアを提供する。
package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recover はパニックリカバリーを行うミドルウェア。
// HTTPハンドラー内でパニックが発生した場合、500エラーを返す。
// 注意: goroutine内でのパニックは別途recoverが必要。
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						"error", rec,
						"stack", string(debug.Stack()),
						"path", r.URL.Path,
						"method", r.Method,
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":{"code":"INTERNAL_ERROR","message":"internal server error"}}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter はレスポンスステータスを記録するためのラッパー。
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestLogger はリクエストログを記録するミドルウェア。
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rw, r)

			// エラーレスポンスの場合はログを記録
			if rw.status >= 400 {
				logger.Warn("error response",
					"status", rw.status,
					"path", r.URL.Path,
					"method", r.Method,
				)
			}
		})
	}
}
