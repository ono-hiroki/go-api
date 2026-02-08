package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"go-api/internal/application/user"
	"go-api/internal/domain"
	httperrors "go-api/internal/presentation/http/errors"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

// createUserRequest はユーザー作成のJSONリクエスト。
type createUserRequest struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"required,max=255,email"`
}

func (r createUserRequest) toInput() user.CreateUserInput {
	return user.CreateUserInput{
		Name:  r.Name,
		Email: r.Email,
	}
}

// createUserResponse はユーザー作成のJSONレスポンス。
type createUserResponse struct {
	User createUserResponseUser `json:"user"`
}

type createUserResponseUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func newCreateUserResponse(output *user.CreateUserOutput) createUserResponse {
	return createUserResponse{
		User: createUserResponseUser{
			ID:    output.User.ID,
			Name:  output.User.Name,
			Email: output.User.Email,
		},
	}
}

// CreateHandler はユーザー作成のHTTPハンドラー。
type CreateHandler struct {
	uc     *user.CreateUserUsecase
	logger *slog.Logger
}

// NewCreateHandler は CreateHandler を生成する。
func NewCreateHandler(uc *user.CreateUserUsecase, logger *slog.Logger) *CreateHandler {
	return &CreateHandler{
		uc:     uc,
		logger: logger,
	}
}

// ServeHTTP はユーザーを作成する。
// POST /users
func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperrors.WriteError(w, r, domain.ErrInvalidInput, h.logger)
		return
	}

	if err := validate.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			httperrors.WriteError(w, r, toValidationError(ve), h.logger)
			return
		}
		httperrors.WriteError(w, r, domain.ErrInvalidInput, h.logger)
		return
	}

	output, err := h.uc.Execute(r.Context(), req.toInput())
	if err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newCreateUserResponse(output))
}

func toValidationError(ve validator.ValidationErrors) *domain.ValidationError {
	domainErr := domain.NewValidationError()
	for _, fe := range ve {
		domainErr.Add(fe.Field(), tagToCode(fe.Tag()), fieldErrorMessage(fe))
	}
	return domainErr
}

func tagToCode(tag string) string {
	switch tag {
	case "required":
		return "required"
	case "max":
		return "too_long"
	case "email":
		return "invalid_format"
	default:
		return tag
	}
}

func fieldErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "max":
		return fmt.Sprintf("%s must be %s characters or less", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
