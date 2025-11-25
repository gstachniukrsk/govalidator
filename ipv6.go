package govalidator

import (
	"context"
	"net"
)

// InvalidIPv6Error is returned when a string is not a valid IPv6 address.
type InvalidIPv6Error struct {
	Value string
}

// Error returns the error message.
func (e InvalidIPv6Error) Error() string {
	return "invalid IPv6 address"
}

// IPv6Validator validates that a string value is a valid IPv6 address.
// It accepts both full IPv6 notation (e.g., "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
// and short form (e.g., "::1", "fe80::1"), as well as CIDR notation (e.g., "2001:db8::/32").
func IPv6Validator(_ context.Context, value any) (twigBlock bool, errs []error) {
	str, ok := value.(string)
	if !ok {
		return true, []error{NotAStringError{}}
	}

	// Try parsing as CIDR first
	_, _, err := net.ParseCIDR(str)
	if err == nil {
		// Valid CIDR notation, verify it's IPv6
		ip, _, _ := net.ParseCIDR(str)
		if ip.To4() != nil {
			// It's an IPv4 CIDR, not IPv6
			return false, []error{InvalidIPv6Error{Value: str}}
		}
		return false, nil
	}

	// Parse as regular IP address
	ip := net.ParseIP(str)
	if ip == nil {
		return false, []error{InvalidIPv6Error{Value: str}}
	}

	// Check if it's IPv6 (To4() returns nil for IPv6)
	if ip.To4() != nil {
		return false, []error{InvalidIPv6Error{Value: str}}
	}

	return false, nil
}
