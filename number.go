package govalidator

import "context"

// NumberValidator validates that a value is a number (float64 or int).
func NumberValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	_, ok := value.(float64)
	_, ok2 := value.(int)

	if !ok && !ok2 {
		return true, []error{NotANumberError{}}
	}

	return
}
