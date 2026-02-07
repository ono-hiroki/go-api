package user

import (
	"context"

	"go-api/internal/domain/user"
)

// ListUsersOutput はユーザー一覧取得の出力。
type ListUsersOutput struct {
	Users []*user.User
}

// ListUsersUsecase はユーザー一覧取得のユースケース。
type ListUsersUsecase struct {
	repo user.UserRepository
}

// NewListUsersUsecase は ListUsersUsecase を生成する。
func NewListUsersUsecase(repo user.UserRepository) *ListUsersUsecase {
	return &ListUsersUsecase{repo: repo}
}

// Execute はユーザー一覧を取得する。
func (uc *ListUsersUsecase) Execute(ctx context.Context) (*ListUsersOutput, error) {
	users, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return &ListUsersOutput{Users: users}, nil
}
