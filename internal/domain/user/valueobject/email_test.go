package valueobject

import (
	"strings"
	"testing"
)

func TestNewEmail(t *testing.T) {
	validCases := []struct {
		name  string
		input string
	}{
		{"正常系/基本形式", "test@example.com"},
		{"正常系/サブドメイン付き", "user@mail.example.com"},
		{"正常系/プラス記号付き", "user+tag@example.com"},
	}

	for _, tt := range validCases {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEmail(tt.input)
			if err != nil {
				t.Fatalf("エラーが発生: %v", err)
			}
			if e.String() != tt.input {
				t.Errorf("got %q, want %q", e.String(), tt.input)
			}
		})
	}

	t.Run("異常系/空文字", func(t *testing.T) {
		_, err := NewEmail("")
		if err != ErrEmailRequired {
			t.Errorf("got %v, want ErrEmailRequired", err)
		}
	})

	t.Run("異常系/256文字以上", func(t *testing.T) {
		long := strings.Repeat("a", 244) + "@example.com" // 256 chars
		_, err := NewEmail(long)
		if err != ErrEmailTooLong {
			t.Errorf("got %v, want ErrEmailTooLong", err)
		}
	})

	invalidFormatCases := []struct {
		name  string
		input string
	}{
		{"異常系/アットマークなし", "testexample.com"},
		{"異常系/ドメインなし", "test@"},
		{"異常系/ローカルパートなし", "@example.com"},
		{"異常系/スペース含む", "test @example.com"},
	}

	for _, tt := range invalidFormatCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmail(tt.input)
			if err != ErrEmailInvalid {
				t.Errorf("got %v, want ErrEmailInvalid", err)
			}
		})
	}
}
