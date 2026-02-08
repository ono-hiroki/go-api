package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-api/internal/presentation/http/handler"
	"go-api/internal/presentation/http/middleware"
)

// Dependencies はルーター構築に必要な依存関係のインターフェース。
// di.Container がこのインターフェースを満たす。
type Dependencies interface {
	UserHandler() *handler.UserHandler
	Logger() *slog.Logger
}

// NewRouter はHTTPルーターを生成する。
func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	// 基本エンドポイント
	mux.HandleFunc("/health", handleHealth)

	// ユーザー
	mux.HandleFunc("GET /users", deps.UserHandler().ListUsers)

	// ミドルウェア適用
	var h http.Handler = mux
	h = middleware.RequestLogger(deps.Logger())(h)
	h = middleware.Recover(deps.Logger())(h)

	return h
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
