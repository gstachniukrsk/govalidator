package govalidator

import (
	"context"
	"fmt"
	"reflect"
)

// InvalidOptionError represents an error when a value does not match any of the allowed options.
type InvalidOptionError struct {
	Options []any
	Actual  any
}

func (e InvalidOptionError) Error() string {
	return fmt.Sprintf("invalid option: %v, expected one of %v", e.Actual, e.Options)
}

// OneOfValidator validates that a value matches one of the provided options.
func OneOfValidator(options ...any) ContextValidator {
	return func(_ context.Context, value any) (twigBlock bool, errs []error) {
		for _, option := range options {
			if reflect.DeepEqual(value, option) {
				return
			}
		}

		errs = append(errs, InvalidOptionError{Options: options, Actual: value})
		return
	}
}
