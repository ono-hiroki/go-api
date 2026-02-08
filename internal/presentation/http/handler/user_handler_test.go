package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	usecase "go-api/internal/application/user"
	"go-api/internal/domain/user"
	"go-api/internal/domain/user/mocks"
	"go-api/internal/presentation/http/handler"
	"go-api/internal/testutil/factory"
)

func TestUserHandler_ListUsers(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("ユーザー一覧を取得できる", func(t *testing.T) {
		testUser := factory.NewUser(factory.WithName("test"), factory.WithEmail("test@example.com"))

		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindAll(mock.Anything).Return([]*user.User{testUser}, nil)

		uc := usecase.NewListUsersUsecase(repo)
		h := handler.NewUserHandler(uc, logger)

		req := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)
		rec := httptest.NewRecorder()

		h.ListUsers(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp usecase.ListUsersOutput
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Len(t, resp.Users, 1)
		assert.Equal(t, testUser.ID().String(), resp.Users[0].ID)
		assert.Equal(t, "test", resp.Users[0].Name)
		assert.Equal(t, "test@example.com", resp.Users[0].Email)
	})

	t.Run("ユーザーが0件の場合は空配列を返す", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindAll(mock.Anything).Return([]*user.User{}, nil)

		uc := usecase.NewListUsersUsecase(repo)
		h := handler.NewUserHandler(uc, logger)

		req := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)
		rec := httptest.NewRecorder()

		h.ListUsers(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp usecase.ListUsersOutput
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Empty(t, resp.Users)
	})

	t.Run("ユースケースがエラーを返した場合は500エラーを返す", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindAll(mock.Anything).Return(nil, errors.New("db error"))

		uc := usecase.NewListUsersUsecase(repo)
		h := handler.NewUserHandler(uc, logger)

		req := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)
		rec := httptest.NewRecorder()

		h.ListUsers(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("コンテキストがキャンセルされた場合はエラーを返す", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindAll(mock.Anything).Return(nil, context.Canceled)

		uc := usecase.NewListUsersUsecase(repo)
		h := handler.NewUserHandler(uc, logger)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		req := httptest.NewRequest(http.MethodGet, "/users", http.NoBody).WithContext(ctx)
		rec := httptest.NewRecorder()

		h.ListUsers(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
