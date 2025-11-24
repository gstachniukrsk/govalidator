package govalidator

import (
	"context"
)

func NonNullableValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	if value == nil {
		errs = append(errs, RequiredError{})
		twigBlock = true
		return
	}

	switch v := value.(type) {
	case *interface{}:
		if v == nil {
			return fail()
		}
	case *[]interface{}:
		if v == nil {
			return fail()
		}
	case *map[string]interface{}:
		if v == nil {
			return fail()
		}
	case []interface{}:
		if v == nil {
			return fail()
		}
	case map[string]interface{}:
		if v == nil {
			return fail()
		}
	}

	return
}

func fail() (bool, []error) {
	return true, []error{RequiredError{}}
}
