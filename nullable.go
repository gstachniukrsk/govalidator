package govalidator

import "context"

// NullableValidator allows null values and stops validation chain when encountering null.
func NullableValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	if value == nil {
		twigBlock = true
		return
	}

	switch v := value.(type) {
	case *interface{}:
		if v == nil {
			twigBlock = true
		}
	case *[]interface{}:
		if v == nil {
			twigBlock = true
		}
	case *map[string]interface{}:
		if v == nil {
			twigBlock = true
		}
	case []interface{}:
		if v == nil {
			twigBlock = true
		}
	case map[string]interface{}:
		if v == nil {
			twigBlock = true
		}
	}

	return
}
