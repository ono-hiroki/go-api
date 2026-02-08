package user_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	usecase "go-api/internal/application/user"
	"go-api/internal/domain"
	"go-api/internal/domain/user/mocks"
	handler "go-api/internal/presentation/http/handler/user"
	"go-api/internal/testutil/factory"
)

func TestGetHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("ユーザーを取得できる", func(t *testing.T) {
		testUser := factory.NewUser(factory.WithName("test"), factory.WithEmail("test@example.com"))

		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testUser.ID()).Return(testUser, nil)

		uc := usecase.NewGetUserUsecase(repo)
		h := handler.NewGetHandler(uc, logger)

		req := httptest.NewRequest(http.MethodGet, "/users/"+testUser.ID().String(), http.NoBody)
		req.SetPathValue("id", testUser.ID().String())
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp usecase.GetUserOutput
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Equal(t, testUser.ID().String(), resp.User.ID)
		assert.Equal(t, "test", resp.User.Name)
		assert.Equal(t, "test@example.com", resp.User.Email)
	})

	t.Run("存在しないユーザーの場合は404エラーを返す", func(t *testing.T) {
		testUser := factory.NewUser()

		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testUser.ID()).Return(nil, domain.ErrNotFound)

		uc := usecase.NewGetUserUsecase(repo)
		h := handler.NewGetHandler(uc, logger)

		req := httptest.NewRequest(http.MethodGet, "/users/"+testUser.ID().String(), http.NoBody)
		req.SetPathValue("id", testUser.ID().String())
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("不正なIDの場合は400エラーを返す", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)

		uc := usecase.NewGetUserUsecase(repo)
		h := handler.NewGetHandler(uc, logger)

		req := httptest.NewRequest(http.MethodGet, "/users/invalid-id", http.NoBody)
		req.SetPathValue("id", "invalid-id")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
