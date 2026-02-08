// Package httperrors はHTTP APIのエラーレスポンス処理を提供する。
package httperrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go-api/internal/domain"
	"go-api/internal/domain/user/valueobject"
)

// ErrorResponse はAPIエラーレスポンスのJSON構造。
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail はエラーの詳細情報。
type ErrorDetail struct {
	Code    string       `json:"code"`              // エラーコード（機械可読）
	Message string       `json:"message"`           // エラーメッセージ（人間可読）
	Details []FieldError `json:"details,omitempty"` // フィールドエラー詳細
}

// FieldError はフィールド単位のエラー（ValidationError用）。
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// StatusFromError はドメインエラーからHTTPステータスコードを導出する。
func StatusFromError(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInvalidInput),
		errors.Is(err, valueobject.ErrInvalidID),
		errors.Is(err, valueobject.ErrNameRequired),
		errors.Is(err, valueobject.ErrNameTooLong),
		errors.Is(err, valueobject.ErrEmailRequired),
		errors.Is(err, valueobject.ErrEmailTooLong),
		errors.Is(err, valueobject.ErrEmailInvalid):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// CodeFromError はエラーからエラーコード文字列を導出する。
func CodeFromError(err error) string {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return "NOT_FOUND"
	case errors.Is(err, domain.ErrConflict):
		return "CONFLICT"
	case errors.Is(err, domain.ErrInvalidInput),
		errors.Is(err, valueobject.ErrInvalidID),
		errors.Is(err, valueobject.ErrNameRequired),
		errors.Is(err, valueobject.ErrNameTooLong),
		errors.Is(err, valueobject.ErrEmailRequired),
		errors.Is(err, valueobject.ErrEmailTooLong),
		errors.Is(err, valueobject.ErrEmailInvalid):
		return "VALIDATION_ERROR"
	case errors.Is(err, domain.ErrUnauthorized):
		return "UNAUTHORIZED"
	case errors.Is(err, domain.ErrForbidden):
		return "FORBIDDEN"
	default:
		return "INTERNAL_ERROR"
	}
}

// WriteError はエラーレスポンスをJSONで書き込む。
// 500系エラーの場合は詳細をログに記録し、ユーザーには隠蔽する。
func WriteError(w http.ResponseWriter, r *http.Request, err error, logger *slog.Logger) {
	// validator.ValidationErrors の場合は特別処理
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		writeValidationError(w, ve)
		return
	}

	status := StatusFromError(err)
	code := CodeFromError(err)

	// 500系はログに詳細を記録
	if status >= 500 && logger != nil {
		logger.Error("internal error",
			"error", err.Error(),
			"path", r.URL.Path,
			"method", r.Method,
		)
	}

	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: userFacingMessage(err, status),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func writeValidationError(w http.ResponseWriter, ve validator.ValidationErrors) {
	details := make([]FieldError, len(ve))
	for i, fe := range ve {
		details[i] = FieldError{
			Field:   fe.Field(),
			Code:    tagToCode(fe.Tag()),
			Message: fieldErrorMessage(fe),
		}
	}

	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: "validation error",
			Details: details,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(resp)
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

// userFacingMessage は本番環境向けのエラーメッセージを返す。
// 内部エラーの詳細は隠蔽する。
func userFacingMessage(err error, status int) string {
	// 4xx系はエラーメッセージをそのまま返す
	if status >= 400 && status < 500 {
		return err.Error()
	}
	// 5xx系は詳細を隠蔽
	return "internal server error"
}
