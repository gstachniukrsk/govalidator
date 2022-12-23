package main

import "context"

type ContextValidator func(ctx context.Context, value any) (twigBlock bool, errs []error)

func (v ContextValidator) Validate(ctx context.Context, value any) (twigBlock bool, errs []error) {
	return v(ctx, value)
}

func (v ContextValidator) AcceptsNull() bool {
	_, errs := v.Validate(context.Background(), nil)

	for _, err := range errs {
		if _, ok := err.(RequiredError); ok {
			return false
		}
	}

	return true
}
