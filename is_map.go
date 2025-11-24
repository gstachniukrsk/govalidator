package govalidator

import (
	"context"
)

// IsMapValidator is a validator that checks if the value is a map or pointer of a map.
func IsMapValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	// switch
	switch v := value.(type) {
	case map[string]any:
		return
	case *map[string]any:
		if v != nil {
			return
		}
	}

	twigBlock = true
	errs = append(errs, NotAMapError{})
	return
}
