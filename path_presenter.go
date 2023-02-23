package govalidator

import (
	"context"
	"strings"
)

func PathPresenter(glue string) PresenterFunc {
	return func(ctx context.Context, path []string, err error) string {
		return strings.Replace(strings.Join(path, glue), glue+"[", "[", -1)
	}
}
