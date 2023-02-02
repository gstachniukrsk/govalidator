package govalidator

import "context"

type PresenterFunc func(ctx context.Context, path []string, err error) string
