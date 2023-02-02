package govalidator

import "context"

type Collector = func(ctx context.Context, path []string, err error)
