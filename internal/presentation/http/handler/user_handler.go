// Package handler はHTTP APIのハンドラーを提供する。
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go-api/internal/application/user"
	"go-api/internal/domain"
	httperrors "go-api/internal/presentation/http/errors"
)

var validate = validator.New()

// UserHandler はユーザー関連のHTTPハンドラー。
type UserHandler struct {
	listUsersUC  *user.ListUsersUsecase
	createUserUC *user.CreateUserUsecase
	logger       *slog.Logger
}

// NewUserHandler は UserHandler を生成する。
func NewUserHandler(
	listUsersUC *user.ListUsersUsecase,
	createUserUC *user.CreateUserUsecase,
	logger *slog.Logger,
) *UserHandler {
	return &UserHandler{
		listUsersUC:  listUsersUC,
		createUserUC: createUserUC,
		logger:       logger,
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

// CreateUser はユーザーを作成するハンドラー。
// POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input user.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	if err := validate.Struct(input); err != nil {
		httperrors.WriteError(w, r, domain.ErrInvalidInput, h.logger)
		return
	}

	output, err := h.createUserUC.Execute(r.Context(), input)
	if err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(output)
}
