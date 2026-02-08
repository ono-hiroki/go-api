package user_test

import (
	"encoding/json"
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
	"go-api/internal/domain"
	"go-api/internal/domain/user/mocks"
	handler "go-api/internal/presentation/http/handler/user"
	"go-api/internal/testutil/factory"
)

func TestUpdateHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("ユーザーを更新できる", func(t *testing.T) {
		testUser := factory.NewUser(factory.WithName("old"), factory.WithEmail("old@example.com"))

		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testUser.ID()).Return(testUser, nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		uc := usecase.NewUpdateUserUsecase(repo)
		h := handler.NewUpdateHandler(uc, logger)

		body := `{"name": "new", "email": "new@example.com"}`
		req := httptest.NewRequest(http.MethodPut, "/users/"+testUser.ID().String(), strings.NewReader(body))
		req.SetPathValue("id", testUser.ID().String())
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp usecase.UpdateUserOutput
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Equal(t, testUser.ID().String(), resp.User.ID)
		assert.Equal(t, "new", resp.User.Name)
		assert.Equal(t, "new@example.com", resp.User.Email)
	})

	t.Run("存在しないユーザーの場合は404エラーを返す", func(t *testing.T) {
		testUser := factory.NewUser()

		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testUser.ID()).Return(nil, domain.ErrNotFound)

		uc := usecase.NewUpdateUserUsecase(repo)
		h := handler.NewUpdateHandler(uc, logger)

		body := `{"name": "new", "email": "new@example.com"}`
		req := httptest.NewRequest(http.MethodPut, "/users/"+testUser.ID().String(), strings.NewReader(body))
		req.SetPathValue("id", testUser.ID().String())
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("バリデーションエラーの場合は400エラーを返す", func(t *testing.T) {
		testUser := factory.NewUser()

		repo := mocks.NewMockUserRepository(t)

		uc := usecase.NewUpdateUserUsecase(repo)
		h := handler.NewUpdateHandler(uc, logger)

		body := `{"name": "", "email": "invalid"}`
		req := httptest.NewRequest(http.MethodPut, "/users/"+testUser.ID().String(), strings.NewReader(body))
		req.SetPathValue("id", testUser.ID().String())
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errResp map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		require.NoError(t, err)

		errorObj := errResp["error"].(map[string]interface{})
		assert.Equal(t, "VALIDATION_ERROR", errorObj["code"])
	})

	t.Run("不正なIDの場合は400エラーを返す", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)

		uc := usecase.NewUpdateUserUsecase(repo)
		h := handler.NewUpdateHandler(uc, logger)

		body := `{"name": "test", "email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPut, "/users/invalid-id", strings.NewReader(body))
		req.SetPathValue("id", "invalid-id")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("不正なJSONの場合は400エラーを返す", func(t *testing.T) {
		testUser := factory.NewUser()

		repo := mocks.NewMockUserRepository(t)

		uc := usecase.NewUpdateUserUsecase(repo)
		h := handler.NewUpdateHandler(uc, logger)

		body := `{invalid json}`
		req := httptest.NewRequest(http.MethodPut, "/users/"+testUser.ID().String(), strings.NewReader(body))
		req.SetPathValue("id", testUser.ID().String())
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
