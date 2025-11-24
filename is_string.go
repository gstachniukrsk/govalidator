package govalidator

import (
	"context"
)

// IsStringValidator is a validator that checks if the value is a string or pointer of a string.
func IsStringValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	switch v := value.(type) {
	case string:
		return
	case *string:
		if v != nil {
			return
		}
	}

	twigBlock = true
	errs = append(errs, NotAStringError{})

	return
}
