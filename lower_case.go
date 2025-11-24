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
	if _, ok := value.(string); !ok {
		return true, []error{NotAStringError{}}
	}

	switch value.(type) {
	case string:
		if value.(string) != strings.ToLower(value.(string)) {
			errs = append(errs, NotLowerCasedError{
				Input: value.(string),
			})
		}
	}

	return
}
