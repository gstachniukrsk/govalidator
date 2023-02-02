package govalidator_test

import (
	"context"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRegistryPresenter(t *testing.T) {
	anyPresenter := govalidator.SimpleErrorPresenter()
	t.Run("constructor", func(t *testing.T) {
		p := govalidator.NewRegistryPresenter(anyPresenter, map[error]govalidator.PresenterFunc{})
		assert.NotNil(t, p)
	})

	t.Run("nil error panics", func(t *testing.T) {
		p := govalidator.NewRegistryPresenter(anyPresenter, map[error]govalidator.PresenterFunc{})

		assert.Panics(t, func() {
			p.Present(context.Background(), []string{}, nil)
		})
	})

	t.Run("registered", func(t *testing.T) {
		p := govalidator.NewRegistryPresenter(anyPresenter, map[error]govalidator.PresenterFunc{
			govalidator.RequiredError{}: func(ctx context.Context, path []string, err error) string {
				return "[required]"
			},
			govalidator.MinSizeError{}: func(ctx context.Context, path []string, err error) string {
				return "[min_size]"
			},
		})

		out := p.Present(nil, []string{}, govalidator.RequiredError{})
		assert.Equal(t, "[required]", out)

		out = p.Present(nil, []string{}, govalidator.MinSizeError{
			MinSize:    5,
			ActualSize: 2,
		})
		assert.Equal(t, "[min_size]", out)
	})

	t.Run("fallback", func(t *testing.T) {
		p := govalidator.NewRegistryPresenter(anyPresenter, map[error]govalidator.PresenterFunc{
			govalidator.RequiredError{}: func(ctx context.Context, path []string, err error) string {
				return "[required]"
			},
		})

		out := p.Present(nil, []string{}, govalidator.NotAnIntegerError{})
		assert.Equal(t, "not an integer", out)
	})
}
