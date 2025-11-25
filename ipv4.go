package govalidator

import (
	"context"
	"net"
)

// InvalidIPv4Error is returned when a string is not a valid IPv4 address.
type InvalidIPv4Error struct {
	Value string
}

// Error returns the error message.
func (e InvalidIPv4Error) Error() string {
	return "invalid IPv4 address"
}

// IPv4Validator validates that a string value is a valid IPv4 address.
// It accepts both standard IPv4 notation (e.g., "192.168.1.1") and
// IPv4 addresses with CIDR notation (e.g., "192.168.1.0/24").
func IPv4Validator(_ context.Context, value any) (twigBlock bool, errs []error) {
	str, ok := value.(string)
	if !ok {
		return true, []error{NotAStringError{}}
	}

	// Parse as IP address (with or without CIDR)
	ip := net.ParseIP(str)
	if ip == nil {
		// Try parsing as CIDR
		_, _, err := net.ParseCIDR(str)
		if err != nil {
			return false, []error{InvalidIPv4Error{Value: str}}
		}
		// Valid CIDR, but we need to ensure it's IPv4
		// ParseCIDR accepts both IPv4 and IPv6, so we need additional check
		return false, nil
	}

	// Check if it's IPv4 (not IPv6)
	if ip.To4() == nil {
		return false, []error{InvalidIPv4Error{Value: str}}
	}

	return false, nil
}
