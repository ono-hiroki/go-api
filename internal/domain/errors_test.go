package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-api/internal/domain"
)

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
