package valueobject

import (
	"strings"
	"testing"
)

func TestNewUserName(t *testing.T) {
	validCases := []struct {
		name  string
		input string
	}{
		{"正常系/1文字", "a"},
		{"正常系/日本語の名前", "田中太郎"},
		{"正常系/100文字ちょうど", strings.Repeat("a", 100)},
	}

	for _, tt := range validCases {
		t.Run(tt.name, func(t *testing.T) {
			n, err := NewUserName(tt.input)
			if err != nil {
				t.Fatalf("エラーが発生: %v", err)
			}
			if n.String() != tt.input {
				t.Errorf("got %q, want %q", n.String(), tt.input)
			}
		})
	}

	t.Run("異常系/空文字", func(t *testing.T) {
		_, err := NewUserName("")
		if err != ErrNameRequired {
			t.Errorf("got %v, want ErrNameRequired", err)
		}
	})

	t.Run("異常系/101文字以上", func(t *testing.T) {
		_, err := NewUserName(strings.Repeat("a", 101))
		if err != ErrNameTooLong {
			t.Errorf("got %v, want ErrNameTooLong", err)
		}
	})

	t.Run("異常系/マルチバイト101文字以上", func(t *testing.T) {
		_, err := NewUserName(strings.Repeat("あ", 101))
		if err != ErrNameTooLong {
			t.Errorf("got %v, want ErrNameTooLong", err)
		}
	})
}
