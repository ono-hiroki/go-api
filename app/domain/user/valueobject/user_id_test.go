package valueobject

import (
	"testing"
)

func TestNewUserID(t *testing.T) {
	t.Run("正常系/空でないIDが生成される", func(t *testing.T) {
		id := NewUserID()
		if id.String() == "" {
			t.Error("空でないIDが生成されるべき")
		}
	})

	t.Run("正常系/毎回異なるIDが生成される", func(t *testing.T) {
		id1 := NewUserID()
		id2 := NewUserID()
		if id1.String() == id2.String() {
			t.Error("異なるIDが生成されるべき")
		}
	})
}

func TestParseUserID(t *testing.T) {
	t.Run("正常系/有効なUUIDを復元できる", func(t *testing.T) {
		input := "550e8400-e29b-41d4-a716-446655440000"

		parsed, err := ParseUserID(input)
		if err != nil {
			t.Fatalf("有効なUUIDでエラーが発生: %v", err)
		}
		if parsed.String() != input {
			t.Errorf("got %s, want %s", parsed.String(), input)
		}
	})

	invalidCases := []struct {
		name  string
		input string
	}{
		{"異常系/空文字", ""},
		{"異常系/通常の文字列", "not-a-uuid"},
		{"異常系/不完全なUUID", "550e8400-e29b-41d4-a716"},
	}

	for _, tt := range invalidCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseUserID(tt.input)
			if err != ErrInvalidID {
				t.Errorf("got %v, want ErrInvalidID", err)
			}
		})
	}
}
