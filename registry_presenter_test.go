package go_validator_test

import (
	"context"
	"github.com/gstachniukrsk/go_validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRegistryPresenter(t *testing.T) {
	anyPresenter := go_validator.SimpleErrorPresenter()
	t.Run("constructor", func(t *testing.T) {
		p := go_validator.NewRegistryPresenter(anyPresenter, map[error]go_validator.PresenterFunc{})
		assert.NotNil(t, p)
	})

	t.Run("nil error panics", func(t *testing.T) {
		p := go_validator.NewRegistryPresenter(anyPresenter, map[error]go_validator.PresenterFunc{})

		assert.Panics(t, func() {
			p.Present(context.Background(), []string{}, nil)
		})
	})

	t.Run("registered", func(t *testing.T) {
		p := go_validator.NewRegistryPresenter(anyPresenter, map[error]go_validator.PresenterFunc{
			go_validator.RequiredError{}: func(ctx context.Context, path []string, err error) string {
				return "[required]"
			},
			go_validator.MinSizeError{}: func(ctx context.Context, path []string, err error) string {
				return "[min_size]"
			},
		})

		out := p.Present(nil, []string{}, go_validator.RequiredError{})
		assert.Equal(t, "[required]", out)

		out = p.Present(nil, []string{}, go_validator.MinSizeError{
			MinSize:    5,
			ActualSize: 2,
		})
		assert.Equal(t, "[min_size]", out)
	})

	t.Run("fallback", func(t *testing.T) {
		p := go_validator.NewRegistryPresenter(anyPresenter, map[error]go_validator.PresenterFunc{
			go_validator.RequiredError{}: func(ctx context.Context, path []string, err error) string {
				return "[required]"
			},
		})

		out := p.Present(nil, []string{}, go_validator.NotAnIntegerError{})
		assert.Equal(t, "not an integer", out)
	})
}
