package main

import "context"

// IsListValidator is a validator that checks if the value is a list.
func IsListValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	switch value.(type) {
	case []any:
		return
	case *[]any:
		if value.(*[]any) != nil {
			return
		}
	}

	twigBlock = true
	errs = append(errs, NotAListError{})
	return
}
