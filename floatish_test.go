package main_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"validator"
)

func TestFloatishValidator(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		twig, errs := main.FloatishValidator(1)(nil, 1.0)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("int", func(t *testing.T) {
		twig, errs := main.FloatishValidator(1)(nil, 1)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("string", func(t *testing.T) {
		twig, errs := main.FloatishValidator(1)(nil, "1")
		assert.True(t, twig)
		assert.Equal(t, []error{main.NotAFloatError{}}, errs)
	})

	t.Run("lower precision", func(t *testing.T) {
		twig, errs := main.FloatishValidator(2)(nil, 1.1)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("higher precision", func(t *testing.T) {
		twig, errs := main.FloatishValidator(2)(nil, 1.1230)
		assert.False(t, twig)
		assert.Equal(t, []error{main.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 3}}, errs)
	})

	t.Run("equal precision", func(t *testing.T) {
		twig, errs := main.FloatishValidator(2)(nil, 1.12)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})
}
