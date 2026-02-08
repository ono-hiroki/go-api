// Package handler はHTTP APIのハンドラーを提供する。
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"

	"go-api/internal/application/user"
	"go-api/internal/domain"
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

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// validateCreateUserInput は入力値をバリデーションする。
// エラーがある場合は ValidationError を返し、問題なければ nil を返す。
func validateCreateUserInput(input user.CreateUserInput) *domain.ValidationError {
	ve := domain.NewValidationError()

	switch {
	case input.Name == "":
		ve.Add("name", "required", "name is required")
	case len([]rune(input.Name)) > 100:
		ve.Add("name", "too_long", "name must be 100 characters or less")
	}

	switch {
	case input.Email == "":
		ve.Add("email", "required", "email is required")
	case len(input.Email) > 255:
		ve.Add("email", "too_long", "email must be 255 characters or less")
	case !emailRegex.MatchString(input.Email):
		ve.Add("email", "invalid_format", "email format is invalid")
	}

	if ve.HasErrors() {
		return ve
	}

	return nil
}
