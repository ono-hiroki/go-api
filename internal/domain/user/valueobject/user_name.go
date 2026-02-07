package valueobject

import "errors"

var (
	ErrNameRequired = errors.New("name is required")
	ErrNameTooLong  = errors.New("name must be 100 characters or less")
)

// UserName はユーザー名を表す値オブジェクト。
type UserName struct {
	value string
}

// NewUserName はバリデーション付きでUserNameを生成する。
func NewUserName(v string) (UserName, error) {
	if v == "" {
		return UserName{}, ErrNameRequired
	}
	if len([]rune(v)) > 100 {
		return UserName{}, ErrNameTooLong
	}
	return UserName{value: v}, nil
}

func (n UserName) String() string {
	return n.value
}

func (n UserName) Equal(other UserName) bool {
	return n.value == other.value
}
