package valueobject

import "go-api/internal/domain"

// エラーコード定数
const (
	CodeRequired      = "required"
	CodeTooLong       = "too_long"
	CodeInvalidFormat = "invalid_format"
)

// ToFieldError は既存のセンチネルエラーをFieldErrorに変換する。
// 未知のエラーの場合はnilを返す。
func ToFieldError(field string, err error) *domain.FieldError {
	switch err {
	case ErrEmailRequired:
		return &domain.FieldError{Field: field, Code: CodeRequired, Message: "email is required"}
	case ErrEmailTooLong:
		return &domain.FieldError{Field: field, Code: CodeTooLong, Message: "email must be 255 characters or less"}
	case ErrEmailInvalid:
		return &domain.FieldError{Field: field, Code: CodeInvalidFormat, Message: "email format is invalid"}
	case ErrNameRequired:
		return &domain.FieldError{Field: field, Code: CodeRequired, Message: "name is required"}
	case ErrNameTooLong:
		return &domain.FieldError{Field: field, Code: CodeTooLong, Message: "name must be 100 characters or less"}
	case ErrInvalidID:
		return &domain.FieldError{Field: field, Code: CodeInvalidFormat, Message: "invalid user id"}
	default:
		return nil
	}
}
