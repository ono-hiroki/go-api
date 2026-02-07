package valueobject

import (
	"errors"
	"github.com/google/uuid"
)

var ErrInvalidID = errors.New("invalid user id")

// UserID はユーザーIDを表す値オブジェクト。
type UserID struct {
	value string
}

// NewUserID はUUIDを新規生成してUserIDを返す。
func NewUserID() UserID {
	return UserID{value: uuid.New().String()}
}

// ParseUserID は文字列からUserIDを復元する。UUID形式でなければエラーを返す。
func ParseUserID(v string) (UserID, error) {
	if _, err := uuid.Parse(v); err != nil {
		return UserID{}, ErrInvalidID
	}
	return UserID{value: v}, nil
}

func (id UserID) String() string {
	return id.value
}
