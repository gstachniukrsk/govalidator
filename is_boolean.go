package govalidator

import (
	"context"
)

// IsBooleanValidator is a validator that checks if the value is a boolean.
func IsBooleanValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	switch value.(type) {
	case bool, *bool:
		return
	}
	twigBlock = true
	errs = append(errs, NotABooleanError{})
	return
}
