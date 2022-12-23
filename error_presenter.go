package go_validator

import "context"
type PresenterFunc func(ctx context.Context, path []string, err error) string
