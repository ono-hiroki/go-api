package user

import (
	"context"

	"go-api/internal/domain/user"
	"go-api/internal/domain/user/valueobject"
)

// DeleteUserUsecase はユーザー削除のユースケース。
type DeleteUserUsecase struct {
	repo user.UserRepository
}

// NewDeleteUserUsecase は DeleteUserUsecase を生成する。
func NewDeleteUserUsecase(repo user.UserRepository) *DeleteUserUsecase {
	return &DeleteUserUsecase{repo: repo}
}

// Execute はユーザーを削除する。
func (uc *DeleteUserUsecase) Execute(ctx context.Context, id string) error {
	userID, err := valueobject.ParseUserID(id)
	if err != nil {
		return err
	}

	// 存在確認
	if _, err := uc.repo.FindByID(ctx, userID); err != nil {
		return err
	}

	return uc.repo.Delete(ctx, userID)
}
