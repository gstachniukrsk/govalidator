package govalidator

import (
	"context"
	"net/url"
)

// InvalidURLError is returned when a string is not a valid URL.
type InvalidURLError struct {
	Value string
}

// Error returns the error message.
func (e InvalidURLError) Error() string {
	return "invalid URL"
}

// URLValidator validates that a string value is a valid URL.
// It accepts both HTTP and HTTPS URLs with proper scheme, host, and optional path/query.
func URLValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	str, ok := value.(string)
	if !ok {
		return true, []error{NotAStringError{}}
	}

	parsedURL, err := url.ParseRequestURI(str)
	if err != nil {
		return false, []error{InvalidURLError{Value: str}}
	}

	// Ensure it has a scheme and host
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false, []error{InvalidURLError{Value: str}}
	}

	return false, nil
}
