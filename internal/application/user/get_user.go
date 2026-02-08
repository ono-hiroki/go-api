package user

import (
	"context"

	"go-api/internal/domain/user"
	"go-api/internal/domain/user/valueobject"
)

// GetUserOutput はユーザー取得の出力。
type GetUserOutput struct {
	User UserDTO
}

// GetUserUsecase はユーザー取得のユースケース。
type GetUserUsecase struct {
	repo user.UserRepository
}

// NewGetUserUsecase は GetUserUsecase を生成する。
func NewGetUserUsecase(repo user.UserRepository) *GetUserUsecase {
	return &GetUserUsecase{repo: repo}
}

// Execute はユーザーを取得する。
func (uc *GetUserUsecase) Execute(ctx context.Context, id string) (*GetUserOutput, error) {
	userID, err := valueobject.ParseUserID(id)
	if err != nil {
		return nil, err
	}

	u, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{
		User: UserDTO{
			ID:    u.ID().String(),
			Name:  u.Name().String(),
			Email: u.Email().String(),
		},
	}, nil
}
