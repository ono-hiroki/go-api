// Package handler はHTTP APIのハンドラーを提供する。
package handler

import (
	"encoding/json"
	"errors"
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
	err := validate.Struct(input)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	ve := domain.NewValidationError()
	for _, fe := range validationErrors {
		ve.Add(fe.Field(), tagToCode(fe.Tag()), tagToMessage(fe))
	}

	return ve
}

// tagToCode はバリデーションタグをエラーコードに変換する。
func tagToCode(tag string) string {
	switch tag {
	case "required":
		return "required"
	case "max":
		return "too_long"
	case "email":
		return "invalid_format"
	default:
		return "invalid"
	}
}

// tagToMessage はバリデーションエラーからメッセージを生成する。
func tagToMessage(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "max":
		return field + " must be " + fe.Param() + " characters or less"
	case "email":
		return field + " format is invalid"
	default:
		return field + " is invalid"
	}
}
