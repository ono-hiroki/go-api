package user

import (
	"context"
	"fmt"

	"go-api/internal/domain/user"
	"go-api/internal/domain/user/valueobject"
)

// CreateUserInput はユーザー作成の入力。
type CreateUserInput struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"required,max=255,email"`
}

// CreateUserOutput はユーザー作成の出力。
type CreateUserOutput struct {
	User UserDTO `json:"user"`
}

// CreateUserUsecase はユーザー作成のユースケース。
type CreateUserUsecase struct {
	repo user.UserRepository
}

// NewCreateUserUsecase は CreateUserUsecase を生成する。
func NewCreateUserUsecase(repo user.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{repo: repo}
}

// Execute はユーザーを作成する。
// 入力バリデーションはハンドラー層で実施済みの前提。
// VO生成エラーは防御的チェックとして扱い、発生時はシステムエラーとする。
func (uc *CreateUserUsecase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	name, err := valueobject.NewUserName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("unexpected name validation error: %w", err)
	}

	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("unexpected email validation error: %w", err)
	}

	u := user.NewUser(name, email)

	if err := uc.repo.Save(ctx, u); err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		User: UserDTO{
			ID:    u.ID().String(),
			Name:  u.Name().String(),
			Email: u.Email().String(),
		},
	}, nil
}
