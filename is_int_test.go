package main_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"validator"
)

func TestIntValidator(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ds := []any{
			1,
			float64(1),
		}

		for i, v := range ds {
			twig, errs := main.IntValidator(nil, v)
			require.False(t, twig, "case #%d", i)
			require.Emptyf(t, errs, "expected no errors, got %v", errs)
		}
	})

	t.Run("not a int - string", func(t *testing.T) {
		twig, errs := main.IntValidator(nil, "1")
		assert.True(t, twig)
		assert.Equal(t, []error{main.NotAnIntegerError{}}, errs)
	})

	t.Run("not a int - float with precision", func(t *testing.T) {
		var f interface{}
		f = 1.1
		twig, errs := main.IntValidator(nil, f)
		assert.True(t, twig)
		assert.Equal(t, []error{main.NotAnIntegerError{}}, errs)
	})

	t.Run("int - float without precision", func(t *testing.T) {
		var f interface{}
		f = 1.00000000
		twig, errs := main.IntValidator(nil, f)
		assert.False(t, twig)
		assert.Empty(t, errs)
	})
}
