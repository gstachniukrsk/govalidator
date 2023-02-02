package govalidator

import (
	"context"
	"fmt"
	"reflect"
)

type InvalidOptionError struct {
	Options []any
	Actual  any
}

func (e InvalidOptionError) Error() string {
	return fmt.Sprintf("invalid option: %v, expected one of %v", e.Actual, e.Options)
}

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
