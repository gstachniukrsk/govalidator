package govalidator

import (
	"context"
	"fmt"
)

// DetailedErrorPresenter creates a presenter that provides detailed, human-readable error messages.
// It extracts structured information from typed errors and formats them in a user-friendly way.
//
// Example outputs:
//   - "value must be at least 5 characters (got 2)"
//   - "value must be between 0.00 and 100.00"
//   - "field 'email' is required"
//
//nolint:cyclop // High complexity is acceptable for comprehensive error formatting
func DetailedErrorPresenter() PresenterFunc {
	return func(_ context.Context, _ []string, err error) string {
		switch e := err.(type) { //nolint:errorlint // Errors are not wrapped in this codebase
		case RequiredError:
			return "this field is required"

		case NotAStringError:
			return "value must be a string"

		case NotAnIntegerError:
			return "value must be a whole number"

		case NotAFloatError:
			return "value must be a number"

		case NotABooleanError:
			return "value must be true or false"

		case NotAMapError, NotAnObjectError:
			return "value must be an object"

		case NotAListError:
			return "value must be a list"

		case NotANumberError:
			return "value must be numeric"

		case MinSizeError:
			return fmt.Sprintf("list must contain at least %d item(s) (got %d)", e.MinSize, e.ActualSize)

		case MaxSizeError:
			return fmt.Sprintf("list must contain at most %d item(s) (got %d)", e.MaxSize, e.ActualSize)

		case FloatPrecisionError:
			return fmt.Sprintf("number must have at most %d decimal place(s) (got %d)", e.ExpectedPrecision, e.ActualPrecision)

		case FloatTooSmallError:
			return fmt.Sprintf("value must be at least %.2f", e.MinFloat)

		case FloatTooLargeError:
			return fmt.Sprintf("value must be at most %.2f", e.MaxFloat)

		case StringTooShortError:
			return fmt.Sprintf("text must be at least %d character(s) long", e.MinLength)

		case StringTooLongError:
			return fmt.Sprintf("text must be at most %d character(s) long", e.MaxLength)

		case FieldNotDefinedError:
			return fmt.Sprintf("field '%s' is required", e.Field)

		case UnexpectedFieldError:
			return fmt.Sprintf("unexpected field '%s'", e.Field)

		case NotAValueError:
			return "value is required"

		default:
			// fallback to default error message
			return err.Error()
		}
	}
}

// VerboseErrorPresenter creates a presenter that provides very detailed technical error messages.
// It includes all available information from structured errors.
//
// Example outputs:
//   - "MinSizeError: expected minimum size 5, actual size 2"
//   - "FloatPrecisionError: expected precision 2, actual precision 4"
func VerboseErrorPresenter() PresenterFunc {
	return func(_ context.Context, _ []string, err error) string {
		switch e := err.(type) { //nolint:errorlint // Errors are not wrapped in this codebase
		case MinSizeError:
			return fmt.Sprintf("MinSizeError: expected minimum size %d, actual size %d", e.MinSize, e.ActualSize)

		case MaxSizeError:
			return fmt.Sprintf("MaxSizeError: expected maximum size %d, actual size %d", e.MaxSize, e.ActualSize)

		case FloatPrecisionError:
			return fmt.Sprintf("FloatPrecisionError: expected precision %d, actual precision %d", e.ExpectedPrecision, e.ActualPrecision)

		case FloatTooSmallError:
			return fmt.Sprintf("FloatTooSmallError: minimum allowed value is %f", e.MinFloat)

		case FloatTooLargeError:
			return fmt.Sprintf("FloatTooLargeError: maximum allowed value is %f", e.MaxFloat)

		case StringTooShortError:
			return fmt.Sprintf("StringTooShortError: minimum length is %d characters", e.MinLength)

		case StringTooLongError:
			return fmt.Sprintf("StringTooLongError: maximum length is %d characters", e.MaxLength)

		case FieldNotDefinedError:
			return fmt.Sprintf("FieldNotDefinedError: field '%s' is not defined", e.Field)

		case UnexpectedFieldError:
			return fmt.Sprintf("UnexpectedFieldError: field '%s' is not allowed", e.Field)

		default:
			// For other errors, return their type and message
			return fmt.Sprintf("%T: %s", err, err.Error())
		}
	}
}
