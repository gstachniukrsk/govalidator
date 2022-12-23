package validator

import "context"

func NumberValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	_, ok := value.(float64)
	_, ok2 := value.(int)

	if !ok && !ok2 {
		return true, []error{NotANumberError{}}
	}

	return
}
