package user_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	usecase "go-api/internal/application/user"
	"go-api/internal/domain/user/mocks"
	handler "go-api/internal/presentation/http/handler/user"
)

func TestCreateHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("ユーザーを作成できる", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		uc := usecase.NewCreateUserUsecase(repo)
		h := handler.NewCreateHandler(uc, logger)

		body := `{"name": "test", "email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp usecase.CreateUserOutput
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		assert.NotEmpty(t, resp.User.ID)
		assert.Equal(t, "test", resp.User.Name)
		assert.Equal(t, "test@example.com", resp.User.Email)
	})

	t.Run("不正なJSONの場合は500エラーを返す", func(t *testing.T) {
		uc := usecase.NewCreateUserUsecase(nil)
		h := handler.NewCreateHandler(uc, logger)

		body := `{invalid json}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		// JSONデコードエラーはドメインエラーではないため500になる
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("バリデーションエラーの場合は400エラーを返す", func(t *testing.T) {
		uc := usecase.NewCreateUserUsecase(nil)
		h := handler.NewCreateHandler(uc, logger)

		body := `{"name": "", "email": "invalid"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ユースケースがエラーを返した場合は500エラーを返す", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db error"))

		uc := usecase.NewCreateUserUsecase(repo)
		h := handler.NewCreateHandler(uc, logger)

		body := `{"name": "test", "email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
