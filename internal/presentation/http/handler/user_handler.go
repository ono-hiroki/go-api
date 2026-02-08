// Package handler はHTTP APIのハンドラーを提供する。
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-api/internal/application/user"
	"go-api/internal/domain"
	"go-api/internal/domain/user/valueobject"
	httperrors "go-api/internal/presentation/http/errors"
)

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

	if ve := validateCreateUserInput(input); ve != nil {
		httperrors.WriteError(w, r, ve, h.logger)
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

// validateCreateUserInput は入力値をバリデーションする。
// エラーがある場合は ValidationError を返し、問題なければ nil を返す。
func validateCreateUserInput(input user.CreateUserInput) *domain.ValidationError {
	ve := domain.NewValidationError()

	if _, err := valueobject.NewUserName(input.Name); err != nil {
		if fe := valueobject.ToFieldError("name", err); fe != nil {
			ve.Add(fe.Field, fe.Code, fe.Message)
		}
	}

	if _, err := valueobject.NewEmail(input.Email); err != nil {
		if fe := valueobject.ToFieldError("email", err); fe != nil {
			ve.Add(fe.Field, fe.Code, fe.Message)
		}
	}

	if ve.HasErrors() {
		return ve
	}

	return nil
}
