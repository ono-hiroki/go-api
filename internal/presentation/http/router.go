package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	userhandler "go-api/internal/presentation/http/handler/user"
	"go-api/internal/presentation/http/middleware"
)

// Dependencies はルーター構築に必要な依存関係のインターフェース。
// di.Container がこのインターフェースを満たす。
type Dependencies interface {
	ListUserHandler() *userhandler.ListHandler
	CreateUserHandler() *userhandler.CreateHandler
	GetUserHandler() *userhandler.GetHandler
	UpdateUserHandler() *userhandler.UpdateHandler
	DeleteUserHandler() *userhandler.DeleteHandler
	Logger() *slog.Logger
}

// NewRouter はHTTPルーターを生成する。
func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	// 基本エンドポイント
	mux.HandleFunc("/health", handleHealth)

	// ユーザー
	mux.Handle("GET /users", deps.ListUserHandler())
	mux.Handle("POST /users", deps.CreateUserHandler())
	mux.Handle("GET /users/{id}", deps.GetUserHandler())
	mux.Handle("PUT /users/{id}", deps.UpdateUserHandler())
	mux.Handle("DELETE /users/{id}", deps.DeleteUserHandler())

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
