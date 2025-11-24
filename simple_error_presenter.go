package govalidator

import (
	"context"
)

// SimpleErrorPresenter returns a PresenterFunc that presents errors using their Error() method.
func SimpleErrorPresenter() PresenterFunc {
	return func(_ context.Context, _ []string, err error) string {
		return err.Error()
	}
}
