package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-api/internal/domain"
)

func TestValidationError(t *testing.T) {
	t.Run("errors.IsでErrInvalidInputとして判定できる", func(t *testing.T) {
		ve := domain.NewValidationError()
		ve.Add("email", "required", "email is required")

		assert.True(t, errors.Is(ve, domain.ErrInvalidInput))
	})

	t.Run("Errorは最初のフィールドエラーのメッセージを返す", func(t *testing.T) {
		ve := domain.NewValidationError()
		ve.Add("email", "required", "email is required")
		ve.Add("name", "too_long", "name is too long")

		assert.Equal(t, "email is required", ve.Error())
	})

	t.Run("エラーがない場合はデフォルトメッセージを返す", func(t *testing.T) {
		ve := domain.NewValidationError()

		assert.Equal(t, "validation failed", ve.Error())
	})

	t.Run("HasErrorsはエラーの有無を返す", func(t *testing.T) {
		ve := domain.NewValidationError()
		assert.False(t, ve.HasErrors())

		ve.Add("email", "required", "email is required")
		assert.True(t, ve.HasErrors())
	})
}

func TestDomainError(t *testing.T) {
	t.Run("NotFoundはErrNotFoundとして判定できる", func(t *testing.T) {
		err := domain.NotFound("user", "FindByID")

		assert.True(t, errors.Is(err, domain.ErrNotFound))
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("ConflictはErrConflictとして判定できる", func(t *testing.T) {
		originalErr := errors.New("duplicate key")
		err := domain.Conflict("user", "Save", originalErr)

		assert.True(t, errors.Is(err, domain.ErrConflict))
		assert.Equal(t, "user already exists", err.Error())
	})

	t.Run("Unwrapでラップされたエラーを取得できる", func(t *testing.T) {
		originalErr := errors.New("original error")
		err := domain.Conflict("user", "Save", originalErr)

		assert.True(t, errors.Is(err, originalErr))
	})

	t.Run("fmt.Errorfでラップしてもerrors.Isで判定できる", func(t *testing.T) {
		err := domain.NotFound("user", "FindByID")
		wrapped := fmt.Errorf("usecase: %w", err)

		assert.True(t, errors.Is(wrapped, domain.ErrNotFound))
	})
}
