// Package domain はドメイン層の共通エラー型を提供する。
package domain

import "errors"

// センチネルエラー（errors.Is で判定可能）
var (
	// ErrNotFound はリソースが見つからない場合のエラー。
	ErrNotFound = errors.New("resource not found")

	// ErrConflict は一意制約違反等の競合エラー。
	ErrConflict = errors.New("resource conflict")

	// ErrInvalidInput は入力値が不正な場合のエラー。
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized は認証が必要な場合のエラー。
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden は権限がない場合のエラー。
	ErrForbidden = errors.New("forbidden")
)

// DomainError はドメインエラーにコンテキストを付与するラッパー。
type DomainError struct {
	Kind    error  // 元のセンチネルエラー（ErrNotFound等）
	Entity  string // エンティティ名（例: "user"）
	Op      string // 操作名（例: "FindByID"）
	Message string // 詳細メッセージ
	Err     error  // ラップされた元エラー
}

// Error はerrorインターフェースを実装する。
func (e *DomainError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Kind.Error()
}

// Unwrap は errors.Is/As のチェーンを維持するために必須。
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is は errors.Is でセンチネルエラーとの比較を可能にする。
func (e *DomainError) Is(target error) bool {
	return errors.Is(e.Kind, target)
}

// NotFound は ErrNotFound をラップしたエラーを生成する。
func NotFound(entity, op string) *DomainError {
	return &DomainError{
		Kind:    ErrNotFound,
		Entity:  entity,
		Op:      op,
		Message: entity + " not found",
	}
}

// Conflict は ErrConflict をラップしたエラーを生成する。
func Conflict(entity, op string, err error) *DomainError {
	return &DomainError{
		Kind:    ErrConflict,
		Entity:  entity,
		Op:      op,
		Message: entity + " already exists",
		Err:     err,
	}
}
