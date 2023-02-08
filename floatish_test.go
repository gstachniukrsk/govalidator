package govalidator_test

import (
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloatishValidator(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		twig, errs := govalidator.FloatishValidator(1)(nil, 1.0)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("int", func(t *testing.T) {
		twig, errs := govalidator.FloatishValidator(1)(nil, 1)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("string", func(t *testing.T) {
		twig, errs := govalidator.FloatishValidator(1)(nil, "1")
		assert.True(t, twig)
		assert.Equal(t, []error{govalidator.NotAFloatError{}}, errs)
	})

	t.Run("lower precision", func(t *testing.T) {
		twig, errs := govalidator.FloatishValidator(2)(nil, 1.1)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("higher precision", func(t *testing.T) {
		twig, errs := govalidator.FloatishValidator(2)(nil, 1.1230)
		assert.False(t, twig)
		assert.Equal(t, []error{govalidator.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 3}}, errs)
	})

	t.Run("equal precision", func(t *testing.T) {
		twig, errs := govalidator.FloatishValidator(2)(nil, 1.12)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("very high precision", func(t *testing.T) {
		twig, errs := govalidator.FloatishValidator(4)(nil, 1.123456789)
		assert.False(t, twig)
		assert.Equal(t, []error{govalidator.FloatPrecisionError{ExpectedPrecision: 4, ActualPrecision: 9}}, errs)
	})
}
