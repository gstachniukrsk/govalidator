package govalidator_test

import (
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsIntegerValidator(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ds := []any{
			1,
			float64(1),
		}

		for i, v := range ds {
			twig, errs := govalidator.IsIntegerValidator(nil, v)
			require.False(t, twig, "case #%d", i)
			require.Emptyf(t, errs, "expected no errors, got %v", errs)
		}
	})

	t.Run("not a int - string", func(t *testing.T) {
		twig, errs := govalidator.IsIntegerValidator(nil, "1")
		assert.True(t, twig)
		assert.Equal(t, []error{govalidator.NotAnIntegerError{}}, errs)
	})

	t.Run("not a int - float with precision", func(t *testing.T) {
		var f interface{}
		f = 1.1
		twig, errs := govalidator.IsIntegerValidator(nil, f)
		assert.True(t, twig)
		assert.Equal(t, []error{govalidator.NotAnIntegerError{}}, errs)
	})

	t.Run("int - float without precision", func(t *testing.T) {
		var f interface{}
		f = 1.00000000
		twig, errs := govalidator.IsIntegerValidator(nil, f)
		assert.False(t, twig)
		assert.Empty(t, errs)
	})
}
