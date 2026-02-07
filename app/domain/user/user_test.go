package user

import (
	"testing"
	"time"

	"go-api/app/domain/user/valueobject"
)

func TestNewUser(t *testing.T) {
	t.Run("正常系/IDとタイムスタンプが自動付与される", func(t *testing.T) {
		name, _ := valueobject.NewUserName("田中太郎")
		email, _ := valueobject.NewEmail("tanaka@example.com")

		u := NewUser(name, email)

		if u.ID().String() == "" {
			t.Error("IDが生成されるべき")
		}
		if u.Name().String() != "田中太郎" {
			t.Errorf("Name got %q, want %q", u.Name().String(), "田中太郎")
		}
		if u.Email().String() != "tanaka@example.com" {
			t.Errorf("Email got %q, want %q", u.Email().String(), "tanaka@example.com")
		}
		if u.CreatedAt().IsZero() {
			t.Error("CreatedAtが設定されるべき")
		}
		if u.UpdatedAt().IsZero() {
			t.Error("UpdatedAtが設定されるべき")
		}
	})
}

func TestReconstruct(t *testing.T) {
	t.Run("正常系/全フィールドが正しく復元される", func(t *testing.T) {
		id, _ := valueobject.ParseUserID("550e8400-e29b-41d4-a716-446655440000")
		name, _ := valueobject.NewUserName("佐藤花子")
		email, _ := valueobject.NewEmail("sato@example.com")
		createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

		u := Reconstruct(id, name, email, createdAt, updatedAt)

		if u.ID().String() != "550e8400-e29b-41d4-a716-446655440000" {
			t.Errorf("ID got %q, want %q", u.ID().String(), "550e8400-e29b-41d4-a716-446655440000")
		}
		if u.Name().String() != "佐藤花子" {
			t.Errorf("Name got %q, want %q", u.Name().String(), "佐藤花子")
		}
		if u.Email().String() != "sato@example.com" {
			t.Errorf("Email got %q, want %q", u.Email().String(), "sato@example.com")
		}
		if !u.CreatedAt().Equal(createdAt) {
			t.Errorf("CreatedAt got %v, want %v", u.CreatedAt(), createdAt)
		}
		if !u.UpdatedAt().Equal(updatedAt) {
			t.Errorf("UpdatedAt got %v, want %v", u.UpdatedAt(), updatedAt)
		}
	})
}
