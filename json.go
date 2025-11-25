package govalidator

import (
	"context"
	"encoding/json"
)

// InvalidJSONError is returned when a string is not valid JSON.
type InvalidJSONError struct {
	Value string
}

// Error returns the error message.
func (e InvalidJSONError) Error() string {
	return "invalid JSON"
}

// JSONValidator validates that a string value contains valid JSON.
// It accepts any valid JSON: objects, arrays, strings, numbers, booleans, or null.
func JSONValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	str, ok := value.(string)
	if !ok {
		return true, []error{NotAStringError{}}
	}

	// Try to unmarshal to verify it's valid JSON
	var js any
	if err := json.Unmarshal([]byte(str), &js); err != nil {
		return false, []error{InvalidJSONError{Value: str}}
	}

	return false, nil
}
