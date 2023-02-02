package govalidator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsBooleanValidator(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		twig, errs := IsBooleanValidator(nil, true)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("pointer", func(t *testing.T) {
		b := true
		twig, errs := IsBooleanValidator(nil, &b)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("not a boolean", func(t *testing.T) {
		twig, errs := IsBooleanValidator(nil, "true")
		assert.True(t, twig)
		assert.Equal(t, []error{NotABooleanError{}}, errs)
	})
}
