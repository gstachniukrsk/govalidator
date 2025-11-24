package govalidator

import (
	"context"
	"strings"
)

// PathPresenter returns a PresenterFunc that joins path segments with a glue string.
func PathPresenter(glue string) PresenterFunc {
	return func(ctx context.Context, path []string, err error) string {
		return strings.Replace(strings.Join(path, glue), glue+"[", "[", -1)
	}
}
