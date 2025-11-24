package govalidator

import (
	"context"
)

// FloatTooLargeError is an error that is returned when a float is greater than the max.
type FloatTooLargeError struct {
	MaxFloat float64
}

// Error is the error message.
func (e FloatTooLargeError) Error() string {
	return "value is greater than max"
}

// MaxFloatValidator is a validator that checks if the value is a float and is less than or equal to the max.
func MaxFloatValidator(maxFloat float64) ContextValidator {
	err := FloatTooLargeError{MaxFloat: maxFloat}
	return func(ctx context.Context, value any) (twigBlock bool, errs []error) {
		// get number
		floatValue, floatOk := value.(float64)
		intValue, intOk := value.(int)

		if !floatOk && !intOk {
			return true, []error{NotAFloatError{}}
		}

		if !floatOk && intOk {
			floatValue = float64(intValue)
		}

		if floatValue > maxFloat {
			return true, []error{err}
		}

		return
	}
}
