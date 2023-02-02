package govalidator

import (
	"context"
	"fmt"
	"strings"
)

// FloatishValidator is a validator that checks if the value is a float64 or an int,
//
//	if float checks against maximal precision.
func FloatishValidator(maxPrecision int) ContextValidator {
	return func(_ context.Context, value any) (twigBlock bool, errs []error) {
		f, ok := value.(float64)
		_, ok2 := value.(int)

		if !ok && !ok2 {
			return true, []error{NotAFloatError{}}
		}

		// get a number
		n := f
		if ok2 {
			// it is int, no precision at all
			return
		}
		// cast to string
		s := fmt.Sprintf("%f", n)
		// get the decimal part
		decimal := strings.Split(s, ".")[1]

		// remove zeroes from the end
		for i := len(decimal) - 1; i >= 0; i-- {
			if decimal[i] != '0' {
				break
			}
			decimal = decimal[0:i]
		}

		// validate length
		if len(decimal) > maxPrecision {
			return false, []error{
				FloatPrecisionError{
					ExpectedPrecision: maxPrecision,
					ActualPrecision:   len(decimal),
				},
			}
		}

		return
	}
}
