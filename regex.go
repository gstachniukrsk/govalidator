package govalidator

import (
	"context"
	"fmt"
	"regexp"
)

// ValueNotMatchingPatternError represents an error when a string does not match the expected pattern.
type ValueNotMatchingPatternError struct {
	Pattern string
	Actual  string
}

func (e ValueNotMatchingPatternError) Error() string {
	return fmt.Sprintf("\"%v\" does not match \"%v\"", e.Actual, e.Pattern)
}

// RegexpValidator validates that a string value matches the given regular expression pattern.
func RegexpValidator(pattern regexp.Regexp) ContextValidator {
	return func(_ context.Context, value any) (twigBlock bool, errs []error) {
		// we can't proceed if the value is not a string
		if _, ok := value.(string); !ok {
			return true, []error{NotAStringError{}}
		}

		if !pattern.MatchString(value.(string)) {
			errs = append(errs, ValueNotMatchingPatternError{Pattern: pattern.String(), Actual: value.(string)})
		}

		return
	}
}
