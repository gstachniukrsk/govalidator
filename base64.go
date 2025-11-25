package govalidator

import (
	"context"
	"encoding/base64"
)

// InvalidBase64Error is returned when a string is not valid base64.
type InvalidBase64Error struct {
	Value string
}

// Error returns the error message.
func (e InvalidBase64Error) Error() string {
	return "invalid base64"
}

// Base64Validator validates that a string value is valid base64 encoded data.
// It accepts both standard and URL-safe base64 encoding with or without padding.
func Base64Validator(_ context.Context, value any) (twigBlock bool, errs []error) {
	str, ok := value.(string)
	if !ok {
		return true, []error{NotAStringError{}}
	}

	// Empty string is not valid base64
	if str == "" {
		return false, []error{InvalidBase64Error{Value: str}}
	}

	// Try standard base64 first
	_, err := base64.StdEncoding.DecodeString(str)
	if err == nil {
		return false, nil
	}

	// Try URL-safe base64
	_, err = base64.URLEncoding.DecodeString(str)
	if err == nil {
		return false, nil
	}

	// Try without padding
	_, err = base64.RawStdEncoding.DecodeString(str)
	if err == nil {
		return false, nil
	}

	// Try URL-safe without padding
	_, err = base64.RawURLEncoding.DecodeString(str)
	if err == nil {
		return false, nil
	}

	return false, []error{InvalidBase64Error{Value: str}}
}
