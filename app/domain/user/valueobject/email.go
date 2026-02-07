package valueobject

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var (
	ErrEmailRequired = errors.New("email is required")
	ErrEmailTooLong  = errors.New("email must be 255 characters or less")
	ErrEmailInvalid  = errors.New("email format is invalid")
)

// Email はメールアドレスを表す値オブジェクト。
type Email struct {
	value string
}

// NewEmail はバリデーション付きでEmailを生成する。
func NewEmail(v string) (Email, error) {
	if v == "" {
		return Email{}, ErrEmailRequired
	}
	if len(v) > 255 {
		return Email{}, ErrEmailTooLong
	}
	if !emailRegex.MatchString(v) {
		return Email{}, ErrEmailInvalid
	}
	return Email{value: v}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equal(other Email) bool {
	return e.value == other.value
}
