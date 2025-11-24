package govalidator

import (
	"context"
	"fmt"
	"strings"
)

// NotUpperCasedError represents an error when a string value is not in uppercase.
type NotUpperCasedError struct {
	Input string
}

func (e NotUpperCasedError) Error() string {
	return fmt.Sprintf("\"%v\" is not upper cased", e.Input)
}

// UpperCaseValidator validates that a string value is in uppercase.
func UpperCaseValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	if _, ok := value.(string); !ok {
		return true, []error{NotAStringError{}}
	}

	switch value.(type) {
	case string:
		if value.(string) != strings.ToUpper(value.(string)) {
			errs = append(errs, NotUpperCasedError{
				Input: value.(string),
			})
		}
	}

	return
}
