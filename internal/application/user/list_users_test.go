package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	usecase "go-api/internal/application/user"
	"go-api/internal/domain/user"
	"go-api/internal/domain/user/mocks"
	"go-api/internal/testutil/factory"
)

func TestListUsersUsecase_Execute(t *testing.T) {
	t.Run("ユーザー一覧を取得できる", func(t *testing.T) {
		testUser := factory.NewUser(factory.WithName("test"))

		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindAll(mock.Anything).Return([]*user.User{testUser}, nil)

		uc := usecase.NewListUsersUsecase(repo)
		output, err := uc.Execute(context.Background())

		require.NoError(t, err)
		assert.Len(t, output.Users, 1)
		assert.Equal(t, "test", output.Users[0].Name().String())
	})

	t.Run("ユーザーが0件の場合は空のスライスを返す", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindAll(mock.Anything).Return([]*user.User{}, nil)

		uc := usecase.NewListUsersUsecase(repo)
		output, err := uc.Execute(context.Background())

		require.NoError(t, err)
		assert.Empty(t, output.Users)
	})

	t.Run("リポジトリがエラーを返した場合はエラーを返す", func(t *testing.T) {
		repo := mocks.NewMockUserRepository(t)
		repo.EXPECT().FindAll(mock.Anything).Return(nil, errors.New("db error"))

		uc := usecase.NewListUsersUsecase(repo)
		output, err := uc.Execute(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
