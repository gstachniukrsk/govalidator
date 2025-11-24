package govalidator

import (
	"context"
	"fmt"
	"strings"
)

// NotLowerCasedError represents an error when a string value is not in lowercase.
type NotLowerCasedError struct {
	Input string
}

func (e NotLowerCasedError) Error() string {
	return fmt.Sprintf("\"%v\" is not lower cased", e.Input)
}

// LowerCaseValidator validates that a string value is in lowercase.
func LowerCaseValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	str, ok := value.(string)
	if !ok {
		return true, []error{NotAStringError{}}
	}

	if str != strings.ToLower(str) {
		errs = append(errs, NotLowerCasedError{
			Input: str,
		})
	}

	return
}
