// Package api はHTTP APIのエラーレスポンス処理を提供する。
package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"go-api/internal/domain"
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
	case errors.Is(err, domain.ErrInvalidInput):
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
	case errors.Is(err, domain.ErrInvalidInput):
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

	// ValidationErrorの場合はdetailsを追加
	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		resp.Error.Details = make([]FieldError, len(validationErr.Errors))
		for i, fe := range validationErr.Errors {
			resp.Error.Details[i] = FieldError{
				Field:   fe.Field,
				Code:    fe.Code,
				Message: fe.Message,
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
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
