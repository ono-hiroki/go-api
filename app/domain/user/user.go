package user

import (
	"go-api/app/domain/user/valueobject"
)

// User はユーザーエンティティ。
type User struct {
	id    valueobject.UserID
	name  valueobject.UserName
	email valueobject.Email
}

// NewUser は新しいUserエンティティを生成する。IDは自動付与される。
func NewUser(name valueobject.UserName, email valueobject.Email) *User {
	return &User{
		id:    valueobject.NewUserID(),
		name:  name,
		email: email,
	}
}

// Reconstruct は永続化層から読み出したデータでUserを復元する。
func Reconstruct(id valueobject.UserID, name valueobject.UserName, email valueobject.Email) *User {
	return &User{
		id:    id,
		name:  name,
		email: email,
	}
}

func (u *User) ID() valueobject.UserID     { return u.id }
func (u *User) Name() valueobject.UserName { return u.name }
func (u *User) Email() valueobject.Email   { return u.email }
