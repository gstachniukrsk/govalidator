package validator_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"validator/validator"
)

func TestNewRegistryPresenter(t *testing.T) {
	anyPresenter := validator.SimpleErrorPresenter()
	t.Run("constructor", func(t *testing.T) {
		p := validator.NewRegistryPresenter(anyPresenter, map[error]validator.PresenterFunc{})
		assert.NotNil(t, p)
	})

	t.Run("nil error panics", func(t *testing.T) {
		p := validator.NewRegistryPresenter(anyPresenter, map[error]validator.PresenterFunc{})

		assert.Panics(t, func() {
			p.Present(context.Background(), []string{}, nil)
		})
	})

	t.Run("registered", func(t *testing.T) {
		p := validator.NewRegistryPresenter(anyPresenter, map[error]validator.PresenterFunc{
			validator.RequiredError{}: func(ctx context.Context, path []string, err error) string {
				return "[required]"
			},
			validator.MinSizeError{}: func(ctx context.Context, path []string, err error) string {
				return "[min_size]"
			},
		})

		out := p.Present(nil, []string{}, validator.RequiredError{})
		assert.Equal(t, "[required]", out)

		out = p.Present(nil, []string{}, validator.MinSizeError{
			MinSize:    5,
			ActualSize: 2,
		})
		assert.Equal(t, "[min_size]", out)
	})

	t.Run("fallback", func(t *testing.T) {
		p := validator.NewRegistryPresenter(anyPresenter, map[error]validator.PresenterFunc{
			validator.RequiredError{}: func(ctx context.Context, path []string, err error) string {
				return "[required]"
			},
		})

		out := p.Present(nil, []string{}, validator.NotAnIntegerError{})
		assert.Equal(t, "not an integer", out)
	})
}
