// Package handler はHTTP APIのハンドラーを提供する。
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-api/internal/application/user"
	"go-api/internal/presentation/http/errors"
)

// UserHandler はユーザー関連のHTTPハンドラー。
type UserHandler struct {
	listUsersUC *user.ListUsersUsecase
	logger      *slog.Logger
}

// NewUserHandler は UserHandler を生成する。
func NewUserHandler(listUsersUC *user.ListUsersUsecase, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		listUsersUC: listUsersUC,
		logger:      logger,
	}
}

// ListUsers はユーザー一覧を取得するハンドラー。
// GET /users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	output, err := h.listUsersUC.Execute(r.Context())
	if err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(output)
}
