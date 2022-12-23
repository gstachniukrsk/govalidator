package validator

import (
	"context"
)

// FloatIsLesserThanError is an error that is returned when the value is less than the min.
type FloatIsLesserThanError struct {
	MinFloat float64
}

// Error is the error message.
func (e FloatIsLesserThanError) Error() string {
	return "value is less than min"
}

// MinFloatValidator is a validator that checks if the value is a float and is greater than or equal to the min.
func MinFloatValidator(minFloat float64) ContextValidator {
	return func(ctx context.Context, value any) (twigBlock bool, errs []error) {
		floatValue, floatOk := value.(float64)
		intValue, intOk := value.(int)

		if !floatOk && !intOk {
			return true, []error{NotAFloatError{}}
		}

		if !floatOk && intOk {
			floatValue = float64(intValue)
		}

		if floatValue < minFloat {
			return false, []error{FloatIsLesserThanError{MinFloat: minFloat}}
		}

		return
	}
}
