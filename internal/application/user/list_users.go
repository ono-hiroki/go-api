package user

import (
	"context"

	"go-api/internal/domain/user"
)

// UserDTO はユーザー情報のDTO。
type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ListUsersOutput はユーザー一覧取得の出力。
type ListUsersOutput struct {
	Users []UserDTO `json:"users"`
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

	dtos := make([]UserDTO, len(users))
	for i, u := range users {
		dtos[i] = UserDTO{
			ID:    u.ID().String(),
			Name:  u.Name().String(),
			Email: u.Email().String(),
		}
	}

	return &ListUsersOutput{Users: dtos}, nil
}
