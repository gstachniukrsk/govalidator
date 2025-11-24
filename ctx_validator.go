package govalidator

import "context"

// ContextValidator is a function that validates a value and returns validation errors.
type ContextValidator func(ctx context.Context, value any) (twigBlock bool, errs []error)

// Validate executes the validator function with the given context and value.
func (v ContextValidator) Validate(ctx context.Context, value any) (twigBlock bool, errs []error) {
	return v(ctx, value)
}

// AcceptsNull checks if the validator accepts null values without returning a RequiredError.
func (v ContextValidator) AcceptsNull() bool {
	_, errs := v.Validate(context.Background(), nil)

	for _, err := range errs {
		if _, ok := err.(RequiredError); ok { //nolint:errorlint // RequiredError is not wrapped in this codebase
			return false
		}
	}

	return true
}
