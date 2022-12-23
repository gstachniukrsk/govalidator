package validator

import (
	"context"
	"fmt"
	"strings"
)

func PathPresenter(glue string) PresenterFunc {
	return func(ctx context.Context, path []string, err error) string {
		return fmt.Sprintf("%s", strings.Join(path, glue))
	}
}
