package user_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	usecase "go-api/internal/application/user"
	"go-api/internal/domain"
	"go-api/internal/domain/user/mocks"
	handler "go-api/internal/presentation/http/handler/user"
	"go-api/internal/testutil/factory"
)

func TestDeleteHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("ユーザーを削除できる", func(t *testing.T) {
		testUser := factory.NewUser()

		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testUser.ID()).Return(testUser, nil)
		repo.EXPECT().Delete(mock.Anything, testUser.ID()).Return(nil)

		uc := usecase.NewDeleteUserUsecase(repo)
		h := handler.NewDeleteHandler(uc, logger)

		req := httptest.NewRequest(http.MethodDelete, "/users/"+testUser.ID().String(), http.NoBody)
		req.SetPathValue("id", testUser.ID().String())
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())
	})

	t.Run("存在しないユーザーの場合は404エラーを返す", func(t *testing.T) {
		testUser := factory.NewUser()

		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testUser.ID()).Return(nil, domain.ErrNotFound)

		uc := usecase.NewDeleteUserUsecase(repo)
		h := handler.NewDeleteHandler(uc, logger)

		req := httptest.NewRequest(http.MethodDelete, "/users/"+testUser.ID().String(), http.NoBody)
		req.SetPathValue("id", testUser.ID().String())
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("不正なIDの場合は400エラーを返す", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)

		uc := usecase.NewDeleteUserUsecase(repo)
		h := handler.NewDeleteHandler(uc, logger)

		req := httptest.NewRequest(http.MethodDelete, "/users/invalid-id", http.NoBody)
		req.SetPathValue("id", "invalid-id")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
