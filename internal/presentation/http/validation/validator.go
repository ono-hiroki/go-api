// Package validation はHTTPリクエストのバリデーション機能を提供する。
package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var instance = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

// Struct はバリデーションを実行する。
// 失敗時は validator.ValidationErrors を返す。
func Struct(s any) error {
	return instance.Struct(s)
}
