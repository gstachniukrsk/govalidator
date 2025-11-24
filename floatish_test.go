package govalidator_test

import (
	"context"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloatValidator(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		twig, errs := govalidator.FloatValidator(1)(nil, 1.0)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("int", func(t *testing.T) {
		twig, errs := govalidator.FloatValidator(1)(nil, 1)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("string", func(t *testing.T) {
		twig, errs := govalidator.FloatValidator(1)(nil, "1")
		assert.True(t, twig)
		assert.Equal(t, []error{govalidator.NotAFloatError{}}, errs)
	})

	t.Run("lower precision", func(t *testing.T) {
		twig, errs := govalidator.FloatValidator(2)(nil, 1.1)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("higher precision", func(t *testing.T) {
		twig, errs := govalidator.FloatValidator(2)(nil, 1.1230)
		assert.False(t, twig)
		assert.Equal(t, []error{govalidator.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 3}}, errs)
	})

	t.Run("equal precision", func(t *testing.T) {
		twig, errs := govalidator.FloatValidator(2)(nil, 1.12)
		assert.False(t, twig)
		assert.Emptyf(t, errs, "expected no errors, got %v", errs)
	})

	t.Run("very high precision", func(t *testing.T) {
		twig, errs := govalidator.FloatValidator(4)(nil, 1.123456789)
		assert.False(t, twig)
		assert.Equal(t, []error{govalidator.FloatPrecisionError{ExpectedPrecision: 4, ActualPrecision: 9}}, errs)
	})
}

func TestFloatValidator_Coverage(t *testing.T) {
	t.Run("FloatValidator with high precision", func(t *testing.T) {
		validator := govalidator.FloatValidator(10)

		// Test within precision
		shouldBlock, errs := validator(context.Background(), 1.123456789)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test exceeding precision
		shouldBlock, errs = validator(context.Background(), 1.12345678901)
		assert.False(t, shouldBlock)
		assert.NotEmpty(t, errs)
	})
}

func TestFloatValidator_AllPrecisionCases(t *testing.T) {
	t.Run("handles various numeric types", func(t *testing.T) {
		validator := govalidator.FloatValidator(2)

		// Test with float64
		shouldBlock, errs := validator(context.Background(), 1.23)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with int
		shouldBlock, errs = validator(context.Background(), 123)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with wrong type (string)
		shouldBlock, errs = validator(context.Background(), "not a float")
		assert.True(t, shouldBlock)
		assert.NotEmpty(t, errs)
		assert.IsType(t, govalidator.NotAFloatError{}, errs[0])

		// Test with wrong type (float32 - not accepted)
		var f32 float32 = 1.23
		shouldBlock, errs = validator(context.Background(), f32)
		assert.True(t, shouldBlock)
		assert.NotEmpty(t, errs)
		assert.IsType(t, govalidator.NotAFloatError{}, errs[0])
	})

	t.Run("handles trailing zeros in decimal", func(t *testing.T) {
		validator := govalidator.FloatValidator(2)

		// Test with trailing zeros (should be stripped and pass)
		shouldBlock, errs := validator(context.Background(), 1.230000)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with float that has zeros in middle (1.2003 -> 4 decimals)
		shouldBlock, errs = validator(context.Background(), 1.2003)
		assert.False(t, shouldBlock)
		assert.NotEmpty(t, errs) // Should fail - 4 > 2
	})

	t.Run("handles floats without decimal part", func(t *testing.T) {
		validator := govalidator.FloatValidator(2)

		// Test with whole number as float (5.0)
		shouldBlock, errs := validator(context.Background(), 5.0)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)
	})

	t.Run("handles precision edge cases with strings", func(t *testing.T) {
		validator := govalidator.FloatValidator(3)

		// Test floats with different numbers of trailing zeros
		testCases := []struct {
			value float64
			valid bool
		}{
			{1.100, true},   // exactly 1 significant decimal (after removing trailing zeros)
			{1.120, true},   // exactly 2 significant decimals
			{1.123, true},   // exactly 3 significant decimals
			{1.1234, false}, // 4 significant decimals - too many
			{10.00, true},   // all trailing zeros
			{10.01, true},   // 1 significant decimal (after removing trailing zero)
			{10.010, true},  // 2 significant decimals (after removing trailing zero)
		}

		for _, tc := range testCases {
			shouldBlock, errs := validator(context.Background(), tc.value)
			if tc.valid {
				assert.False(t, shouldBlock)
				assert.Empty(t, errs)
			} else {
				assert.False(t, shouldBlock)
				assert.NotEmpty(t, errs)
			}
		}
	})
}
