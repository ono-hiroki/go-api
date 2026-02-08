package user

import (
	"context"
	"fmt"

	"go-api/internal/domain/user"
	"go-api/internal/domain/user/valueobject"
)

// UpdateUserInput はユーザー更新の入力。
type UpdateUserInput struct {
	Name  string
	Email string
}

// UpdateUserOutput はユーザー更新の出力。
type UpdateUserOutput struct {
	User UserDTO
}

// UpdateUserUsecase はユーザー更新のユースケース。
type UpdateUserUsecase struct {
	repo user.UserRepository
}

// NewUpdateUserUsecase は UpdateUserUsecase を生成する。
func NewUpdateUserUsecase(repo user.UserRepository) *UpdateUserUsecase {
	return &UpdateUserUsecase{repo: repo}
}

// Execute はユーザーを更新する。
func (uc *UpdateUserUsecase) Execute(ctx context.Context, id string, input UpdateUserInput) (*UpdateUserOutput, error) {
	userID, err := valueobject.ParseUserID(id)
	if err != nil {
		return nil, err
	}

	u, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	name, err := valueobject.NewUserName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("unexpected name validation error: %w", err)
	}

	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("unexpected email validation error: %w", err)
	}

	u.ChangeName(name)
	u.ChangeEmail(email)

	if err := uc.repo.Save(ctx, u); err != nil {
		return nil, err
	}

	return &UpdateUserOutput{
		User: UserDTO{
			ID:    u.ID().String(),
			Name:  u.Name().String(),
			Email: u.Email().String(),
		},
	}, nil
}
