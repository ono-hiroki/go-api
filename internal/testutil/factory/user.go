package factory

import (
	"fmt"

	"go-api/internal/domain/user"
	"go-api/internal/domain/user/valueobject"
)

type UserOption func(*userParams)

type userParams struct {
	name  string
	email string
}

func WithName(name string) UserOption {
	return func(p *userParams) { p.name = name }
}

func WithEmail(email string) UserOption {
	return func(p *userParams) { p.email = email }
}

// NewUser はテスト用ユーザーを生成する（Functional Optionsパターン）
func NewUser(opts ...UserOption) *user.User {
	p := &userParams{
		name:  "テストユーザー",
		email: fmt.Sprintf("user-%s@example.com", valueobject.NewUserID().String()[:8]),
	}
	for _, opt := range opts {
		opt(p)
	}

	name, _ := valueobject.NewUserName(p.name)
	email, _ := valueobject.NewEmail(p.email)
	return user.NewUser(name, email)
}
