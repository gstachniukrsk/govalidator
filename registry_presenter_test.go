package go_validator_test

import (
"context"
"github.com/stretchr/testify/assert"
"testing"
"validator"
)

func TestNewRegistryPresenter(t *testing.T) {
	anyPresenter := main.SimpleErrorPresenter()
	t.Run("constructor", func(t *testing.T) {
		p := main.NewRegistryPresenter(anyPresenter, map[error]main.PresenterFunc{})
		assert.NotNil(t, p)
	})

	t.Run("nil error panics", func(t *testing.T) {
		p := main.NewRegistryPresenter(anyPresenter, map[error]main.PresenterFunc{})

		assert.Panics(t, func() {
			p.Present(context.Background(), []string{}, nil)
		})
	})

	t.Run("registered", func(t *testing.T) {
		p := main.NewRegistryPresenter(anyPresenter, map[error]main.PresenterFunc{
			main.RequiredError{}: func(ctx context.Context, path []string, err error) string {
				return "[required]"
			},
			main.MinSizeError{}: func(ctx context.Context, path []string, err error) string {
				return "[min_size]"
			},
		})

		out := p.Present(nil, []string{}, main.RequiredError{})
		assert.Equal(t, "[required]", out)

		out = p.Present(nil, []string{}, main.MinSizeError{
			MinSize:    5,
			ActualSize: 2,
		})
		assert.Equal(t, "[min_size]", out)
	})

	t.Run("fallback", func(t *testing.T) {
		p := main.NewRegistryPresenter(anyPresenter, map[error]main.PresenterFunc{
			main.RequiredError{}: func(ctx context.Context, path []string, err error) string {
				return "[required]"
			},
		})

		out := p.Present(nil, []string{}, main.NotAnIntegerError{})
		assert.Equal(t, "not an integer", out)
	})
}
