package validator

import "context"

// MaxSizeValidator is a validator that checks if the value is a string or pointer of a string and if it's length is less than or equal to the given max size.
func MaxSizeValidator(maxSize int, blocks bool) ContextValidator {
	return func(_ context.Context, value any) (twigBlock bool, errs []error) {
		if value == nil {
			return true, []error{NotAListError{}}
		}

		switch v := value.(type) {
		case []interface{}:
			if v == nil {
				return true, []error{NotAListError{}}
			}

			if len(v) > maxSize {
				twigBlock = blocks
				errs = append(errs, MaxSizeError{
					MaxSize:    maxSize,
					ActualSize: len(v),
				})
			}
		case *[]interface{}:
			if v == nil {
				return true, []error{NotAListError{}}
			}

			if len(*v) > maxSize {
				twigBlock = blocks
				errs = append(errs, MaxSizeError{
					MaxSize:    maxSize,
					ActualSize: len(*v),
				})
			}
		}
		return
	}
}
