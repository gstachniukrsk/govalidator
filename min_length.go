package govalidator

import (
	"context"
	"fmt"
)

// StringTooShortError is an error that is returned when a string is too short.
type StringTooShortError struct {
	MinLength int
}

// Error returns the error message.
func (e StringTooShortError) Error() string {
	return fmt.Sprintf("expected at least %v characters", e.MinLength)
}

// MinLengthValidator is a validator that checks if the value is a string and is greater than or equal to the min.
func MinLengthValidator(minLength int) ContextValidator {
	return func(_ context.Context, value any) (twigBlock bool, errs []error) {
		str, ok := value.(string)
		if !ok {
			errs = append(errs, NotAStringError{})
			twigBlock = true
			return
		}
		if len([]rune(str)) < minLength {
			errs = append(errs, StringTooShortError{
				MinLength: minLength,
			})
			return
		}
		return
	}
}
