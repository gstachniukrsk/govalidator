package validator

import "context"

// MinSizeValidator is a validator that checks if the value is a list or map with a minimum size.
func MinSizeValidator(minSize int, blocks bool) ContextValidator {
	return func(_ context.Context, value any) (twigBlock bool, errs []error) {
		switch v := value.(type) {
		case []interface{}:
			if len(v) < minSize {
				twigBlock = blocks
				errs = append(errs, MinSizeError{
					MinSize:    minSize,
					ActualSize: len(v),
				})
			}
		case *[]interface{}:
			if v == nil {
				return true, []error{NotAListError{}}
			}

			if len(*v) < minSize {
				twigBlock = blocks
				errs = append(errs, MinSizeError{
					MinSize:    minSize,
					ActualSize: len(*v),
				})
			}
			return
		default:
			twigBlock = true
			errs = append(errs, NotAListError{})
		}
		return
	}
}
