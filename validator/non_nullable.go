package validator

import (
	"context"
)

func NonNullableValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	if value == nil {
		errs = append(errs, RequiredError{})
		twigBlock = true
		return
	}

	switch value.(type) {
	case *interface{}:
		if value.(*interface{}) == nil {
			return fail()
		}
	case *[]interface{}:
		if value.(*[]interface{}) == nil {
			return fail()
		}
	case *map[string]interface{}:
		if value.(*map[string]interface{}) == nil {
			return fail()
		}
	case []interface{}:
		if value.([]interface{}) == nil {
			return fail()
		}
	case map[string]interface{}:
		if value.(map[string]interface{}) == nil {
			return fail()
		}
	}

	return
}

func fail() (bool, []error) {
	return true, []error{RequiredError{}}
}
