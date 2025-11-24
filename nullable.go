package govalidator

import "context"

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
