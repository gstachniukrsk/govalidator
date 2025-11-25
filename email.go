package govalidator

import (
	"context"
	"net/mail"
)

// InvalidEmailError is returned when a string is not a valid email address.
type InvalidEmailError struct {
	Value string
	Err   error
}

// Error returns the error message.
func (e InvalidEmailError) Error() string {
	if e.Err != nil {
		return "invalid email address: " + e.Err.Error()
	}
	return "invalid email address"
}

// Unwrap returns the underlying error.
func (e InvalidEmailError) Unwrap() error {
	return e.Err
}

// EmailValidator validates that a string value is a valid email address.
// It uses Go's net/mail.ParseAddress which implements RFC 5322.
func EmailValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	str, ok := value.(string)
	if !ok {
		return true, []error{NotAStringError{}}
	}

	// Use Go's built-in RFC 5322 compliant email parser
	_, err := mail.ParseAddress(str)
	if err != nil {
		return false, []error{InvalidEmailError{Value: str, Err: err}}
	}

	return false, nil
}
