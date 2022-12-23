package validator

import (
	"context"
)

func SimpleErrorPresenter() PresenterFunc {
	return func(_ context.Context, _ []string, err error) string {
		return err.Error()
	}
}
