package go_validator

import "context"

func NullableValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	if value == nil {
		twigBlock = true
		return
	}

	switch value.(type) {
	case *interface{}:
		if value.(*interface{}) == nil {
			twigBlock = true
		}
	case *[]interface{}:
		if value.(*[]interface{}) == nil {
			twigBlock = true
		}
	case *map[string]interface{}:
		if value.(*map[string]interface{}) == nil {
			twigBlock = true
		}
	case []interface{}:
		if value.([]interface{}) == nil {
			twigBlock = true
		}
	case map[string]interface{}:
		if value.(map[string]interface{}) == nil {
			twigBlock = true
		}
	}

	return
}
