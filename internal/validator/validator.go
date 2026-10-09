package validator

import (
	"cmp"
	"context"
	"net/mail"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

var PhoneRX = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

type Validator interface {
	Valid(ctx context.Context) Evaluator
}

type Evaluator map[string]string

func (e *Evaluator) addFieldError(key, message string) {
	if *e == nil {
		*e = make(map[string]string)
	}

	if _, exists := (*e)[key]; !exists {
		(*e)[key] = message
	}
}

func (e *Evaluator) CheckField(valid bool, key, message string) {
	if !valid {
		e.addFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func IsEmail(value string) bool {
	addr, err := mail.ParseAddress(value)
	return err == nil && addr.Address == value
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Min[T cmp.Ordered](value, min T) bool {
	return value >= min
}

func Max[T cmp.Ordered](value, max T) bool {
	return value <= max
}

func Between[T cmp.Ordered](value, min, max T) bool {
	return value >= min && value <= max
}

func GreaterThan[T cmp.Ordered](value, limit T) bool {
	return value > limit
}

func IsValidUUID(value string) bool {
	err := uuid.Validate(value)
	return err == nil
}
