package user

import (
	"context"

	"go-api/internal/domain/user/valueobject"
)

//go:generate mockery

// UserRepository はユーザーの永続化インターフェース。
type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id valueobject.UserID) (*User, error)
	FindAll(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id valueobject.UserID) error
}
