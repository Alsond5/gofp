package must

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"text/template"
)

func Do[T any](value T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("must.Do: %v", err))
	}

	return value
}

func Require[T any](value T, err error, msg string) T {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err))
	}

	return value
}

func Be(condition bool, msg string) {
	if !condition {
		panic("must.Be: " + msg)
	}
}

func Bef(condition bool, format string, args ...any) {
	if !condition {
		panic("must.Be: " + fmt.Sprintf(format, args...))
	}
}

func Never(msg string) {
	panic("must.Never: unreachable code reached: " + msg)
}

func NotNil[T any](ptr *T, msg string) *T {
	if ptr == nil {
		panic("must.NotNil: " + msg)
	}

	return ptr
}

func Value[T any](ptr *T, msg string) T {
	if ptr == nil {
		panic("must.Value: " + msg)
	}

	return *ptr
}

func Env(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("must.Env: environment variable %q is not set", key))
	}

	return v
}

func EnvOr(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	return v
}

func EnvInt(key string) int {
	raw := Env(key)
	n, err := strconv.Atoi(raw)
	if err != nil {
		panic(fmt.Sprintf("must.EnvInt: %q=%q is not a valid integer: %v", key, raw, err))
	}

	return n
}

func EnvIntOr(key string, defaultValue int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}

	n, err := strconv.Atoi(raw)
	if err != nil {
		panic(fmt.Sprintf("must.EnvIntOr: %q=%q is not a valid integer: %v", key, raw, err))
	}

	return n
}

func EnvBool(key string) bool {
	raw := Env(key)
	b, err := strconv.ParseBool(raw)
	if err != nil {
		panic(fmt.Sprintf("must.EnvBool: %q=%q is not a valid boolean: %v", key, raw, err))
	}

	return b
}

func EnvBoolOr(key string, defaultValue bool) bool {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(raw)
	if err != nil {
		panic(fmt.Sprintf("must.EnvBoolOr: %q=%q is not a valid boolean: %v", key, raw, err))
	}

	return b
}

func Regexp(pattern string) *regexp.Regexp {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("must.Regexp: invalid pattern %q: %v", pattern, err))
	}

	return re
}

func RegexpPOSIX(pattern string) *regexp.Regexp {
	re, err := regexp.CompilePOSIX(pattern)
	if err != nil {
		panic(fmt.Sprintf("must.RegexpPOSIX: invalid pattern %q: %v", pattern, err))
	}

	return re
}

func Template(t *template.Template, err error) *template.Template {
	if err != nil {
		panic(fmt.Sprintf("must.Template: %v", err))
	}

	return t
}

func Index[T any](slice []T, i int, msg string) T {
	if i < 0 || i >= len(slice) {
		panic(fmt.Sprintf("must.Index: index %d out of bounds [0, %d): %s", i, len(slice), msg))
	}

	return slice[i]
}

func Key[K comparable, V any](m map[K]V, key K, msg string) V {
	v, ok := m[key]
	if !ok {
		panic(fmt.Sprintf("must.Key: key %v not found: %s", key, msg))
	}

	return v
}
