package govalidator

import (
	"context"
	"regexp"
)

// InvalidUUIDError is returned when a string is not a valid UUID.
type InvalidUUIDError struct {
	Value string
}

// Error returns the error message.
func (e InvalidUUIDError) Error() string {
	return "invalid UUID"
}

// UUIDValidator validates that a string value is a valid UUID (version 1-5).
// It accepts UUIDs in the standard format: 8-4-4-4-12 hex digits.
// Example: "550e8400-e29b-41d4-a716-446655440000"
func UUIDValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	str, ok := value.(string)
	if !ok {
		return true, []error{NotAStringError{}}
	}

	// UUID regex pattern matching standard format: 8-4-4-4-12
	// Matches UUIDv1-v5 with case-insensitive hex characters
	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

	if !uuidPattern.MatchString(str) {
		return false, []error{InvalidUUIDError{Value: str}}
	}

	return false, nil
}
