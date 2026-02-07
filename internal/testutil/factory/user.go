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

	name, err := valueobject.NewUserName(p.name)
	if err != nil {
		panic(fmt.Sprintf("factory.NewUser: invalid name %q: %v", p.name, err))
	}
	email, err := valueobject.NewEmail(p.email)
	if err != nil {
		panic(fmt.Sprintf("factory.NewUser: invalid email %q: %v", p.email, err))
	}
	return user.NewUser(name, email)
}
