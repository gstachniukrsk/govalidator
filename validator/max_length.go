package validator

import (
	"context"
	"fmt"
)

// StringTooLongError is an error that is returned when a string is too long.
type StringTooLongError struct {
	MaxLength    int
	ActualLength int
}

// Error returns the error message.
func (e StringTooLongError) Error() string {
	return fmt.Sprintf("expected at most %v characters, got %v", e.MaxLength, e.ActualLength)
}

// MaxLengthValidator is a validator that checks if the value is a string and is less than or equal to the max.
func MaxLengthValidator(maxLength int) ContextValidator {
	return func(_ context.Context, value any) (twigBlock bool, errs []error) {
		switch value.(type) {
		case string:
			str := value.(string)
			// count utf runes
			return handleUtf8(str, maxLength)
		case *string:
			if value.(*string) == nil {
				return true, []error{NotAStringError{}}
			}
			return handleUtf8(*value.(*string), maxLength)
		}

		return true, []error{NotAStringError{}}
	}
}

func handleUtf8(str string, maxLength int) (twig bool, errs []error) {
	if len([]rune(str)) > maxLength {
		errs = append(errs, StringTooLongError{
			MaxLength:    maxLength,
			ActualLength: len([]rune(str)),
		})
	}

	return
}
