package govalidator

import "context"

// PresenterFunc is a function that formats an error with its path for presentation.
type PresenterFunc func(ctx context.Context, path []string, err error) string
