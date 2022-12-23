package go_validator

import (
"context"
)

// StringValidator is a validator that checks if the value is a string or pointer of a string.
func StringValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	switch value.(type) {
	case string:
		return
	case *string:
		if value.(*string) != nil {
			return
		}
	}

	twigBlock = true
	errs = append(errs, NotAStringError{})

	return
}
