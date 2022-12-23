package validator

import (
	"context"
	"fmt"
	"strings"
)

type NotUpperCasedError struct {
	Input string
}

func (e NotUpperCasedError) Error() string {
	return fmt.Sprintf("\"%v\" is not upper cased", e.Input)
}

func UpperCaseValidator(_ context.Context, value any) (twigBlock bool, errs []error) {
	if _, ok := value.(string); !ok {
		return true, []error{NotAStringError{}}
	}

	switch value.(type) {
	case string:
		if value.(string) != strings.ToUpper(value.(string)) {
			errs = append(errs, NotUpperCasedError{
				Input: value.(string),
			})
		}
	}

	return
}
