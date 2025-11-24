package govalidator

import "context"

// Collector is a function that collects validation errors with their paths.
type Collector = func(ctx context.Context, path []string, err error)
