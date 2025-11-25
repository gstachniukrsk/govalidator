package govalidator

import (
	"context"
	"regexp"
)

// InvalidXIDError is returned when a string is not a valid XID.
type InvalidXIDError struct {
	Value string
}

// Error returns the error message.
func (e InvalidXIDError) Error() string {
	return "invalid XID"
}

// XIDValidator validates that a string value is a valid XID (Globally Unique ID).
// XID is a 12-byte globally unique id that uses base32 hex encoding (20 chars).
// The format is: 20 lowercase characters consisting of 'a'-'v' and '0'-'9'.
// Example: "9m4e2mr0ui3e8a215n4g"
//
// XID uses the Mongo Object ID algorithm:
// - 4-byte timestamp (seconds since Unix epoch)
// - 3-byte machine identifier
// - 2-byte process id
// - 3-byte counter
func XIDValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	str, ok := value.(string)
	if !ok {
		return true, []error{NotAStringError{}}
	}

	// XID format: exactly 20 chars, all lowercase, only '0-9' and 'a-v'
	xidPattern := regexp.MustCompile(`^[0-9a-v]{20}$`)

	if !xidPattern.MatchString(str) {
		return false, []error{InvalidXIDError{Value: str}}
	}

	return false, nil
}
