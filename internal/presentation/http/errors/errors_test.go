package httperrors_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-api/internal/domain"
	httperrors "go-api/internal/presentation/http/errors"
	"go-api/internal/presentation/http/validation"
)

func TestStatusFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"ErrNotFound", domain.ErrNotFound, http.StatusNotFound},
		{"ErrConflict", domain.ErrConflict, http.StatusConflict},
		{"ErrInvalidInput", domain.ErrInvalidInput, http.StatusBadRequest},
		{"ErrUnauthorized", domain.ErrUnauthorized, http.StatusUnauthorized},
		{"ErrForbidden", domain.ErrForbidden, http.StatusForbidden},
		{"unknown error", errors.New("unknown"), http.StatusInternalServerError},
		{"DomainError NotFound", domain.NotFound("user", "FindByID"), http.StatusNotFound},
		{"DomainError Conflict", domain.Conflict("user", "Save", nil), http.StatusConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := httperrors.StatusFromError(tt.err)
			assert.Equal(t, tt.expected, status)
		})
	}
}

func TestCodeFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrNotFound", domain.ErrNotFound, "NOT_FOUND"},
		{"ErrConflict", domain.ErrConflict, "CONFLICT"},
		{"ErrInvalidInput", domain.ErrInvalidInput, "VALIDATION_ERROR"},
		{"ErrUnauthorized", domain.ErrUnauthorized, "UNAUTHORIZED"},
		{"ErrForbidden", domain.ErrForbidden, "FORBIDDEN"},
		{"unknown error", errors.New("unknown"), "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := httperrors.CodeFromError(tt.err)
			assert.Equal(t, tt.expected, code)
		})
	}
}

func TestWriteError(t *testing.T) {
	t.Run("NotFoundエラーのレスポンス", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users/123", http.NoBody)

		httperrors.WriteError(w, r, domain.NotFound("user", "FindByID"), nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp httperrors.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "NOT_FOUND", resp.Error.Code)
		assert.Equal(t, "user not found", resp.Error.Message)
		assert.Empty(t, resp.Error.Details)
	})

	t.Run("ValidationErrorのレスポンス", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/users", http.NoBody)

		// validator.ValidationErrors を生成するためにバリデーションを実行
		type testReq struct {
			Email string `json:"email" validate:"required,email"`
			Name  string `json:"name" validate:"required,max=100"`
		}
		req := testReq{Email: "", Name: strings.Repeat("a", 101)}
		ve := validation.Struct(req)

		httperrors.WriteError(w, r, ve, nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp httperrors.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "VALIDATION_ERROR", resp.Error.Code)
		assert.Len(t, resp.Error.Details, 2)
		assert.Equal(t, "email", resp.Error.Details[0].Field)
		assert.Equal(t, "required", resp.Error.Details[0].Code)
	})

	t.Run("内部エラーはメッセージを隠蔽する", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)

		httperrors.WriteError(w, r, errors.New("database connection failed"), nil)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp httperrors.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "INTERNAL_ERROR", resp.Error.Code)
		assert.Equal(t, "internal server error", resp.Error.Message)
	})
}
