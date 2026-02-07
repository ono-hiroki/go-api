package user

import (
	"time"

	"go-api/app/domain/user/valueobject"
)

// User はユーザーエンティティ。
type User struct {
	id        valueobject.UserID
	name      valueobject.UserName
	email     valueobject.Email
	createdAt time.Time
	updatedAt time.Time
}

// NewUser は新しいUserエンティティを生成する。IDとタイムスタンプは自動付与される。
func NewUser(name valueobject.UserName, email valueobject.Email) *User {
	now := time.Now()
	return &User{
		id:        valueobject.NewUserID(),
		name:      name,
		email:     email,
		createdAt: now,
		updatedAt: now,
	}
}

// Reconstruct は永続化層から読み出したデータでUserを復元する。
func Reconstruct(id valueobject.UserID, name valueobject.UserName, email valueobject.Email, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		name:      name,
		email:     email,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() valueobject.UserID     { return u.id }
func (u *User) Name() valueobject.UserName { return u.name }
func (u *User) Email() valueobject.Email   { return u.email }
func (u *User) CreatedAt() time.Time       { return u.createdAt }
func (u *User) UpdatedAt() time.Time       { return u.updatedAt }
