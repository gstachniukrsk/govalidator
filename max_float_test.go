package govalidator

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaxFloatValidator(t *testing.T) {
	type args struct {
		maxFloat float64
		input    float64
	}
	tests := []struct {
		name          string
		args          args
		wantTwigBlock bool
		wantErrs      []error
	}{
		{
			name: "eq positive",
			args: args{
				maxFloat: 1.0,
				input:    1.0,
			},
		},
		{
			name: "eq zero",
			args: args{
				maxFloat: 0.0,
				input:    0.0,
			},
		},
		{
			name: "eq negative",
			args: args{
				maxFloat: -1.0,
				input:    -1.0,
			},
		},
		{
			name: "lt positive",
			args: args{
				maxFloat: 1.0,
				input:    0.5,
			},
		},
		{
			name: "lt zero",
			args: args{
				maxFloat: 0.0,
				input:    -0.5,
			},
		},
		{
			name: "lt negative",
			args: args{
				maxFloat: -1.0,
				input:    -1.5,
			},
		},
		{
			name: "gt positive",
			args: args{
				maxFloat: 2.0,
				input:    2.5,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				FloatTooLargeError{MaxFloat: 2.0},
			},
		},
		{
			name: "gt zero",
			args: args{
				maxFloat: 0.0,
				input:    0.5,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				FloatTooLargeError{MaxFloat: 0.0},
			},
		},
		{
			name: "gt negative",
			args: args{
				maxFloat: -.5,
				input:    0.0,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				FloatTooLargeError{MaxFloat: -.5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := MaxFloatValidator(tt.args.maxFloat)
			twigBlock, errs := v(nil, tt.args.input)
			assert.Equal(t, tt.wantTwigBlock, twigBlock)
			assert.Equal(t, tt.wantErrs, errs)
		})
	}
}

func TestMaxFloatError(t *testing.T) {
	t.Run("MaxFloatError Error method", func(t *testing.T) {
		err := FloatTooLargeError{MaxFloat: 100.0}
		assert.Equal(t, "value is greater than max", err.Error())
		assert.Equal(t, 100.0, err.MaxFloat)
	})
}

func TestMaxFloatValidator_Coverage(t *testing.T) {
	t.Run("MaxFloatValidator with edge cases", func(t *testing.T) {
		validator := MaxFloatValidator(100.0)

		// Test with float below max
		shouldBlock, errs := validator(context.Background(), 50.0)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with float above max - should block
		shouldBlock, errs = validator(context.Background(), 150.0)
		assert.True(t, shouldBlock)
		assert.NotEmpty(t, errs)

		// Test with int below max
		shouldBlock, errs = validator(context.Background(), 50)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with wrong type - should block
		shouldBlock, errs = validator(context.Background(), "not a number")
		assert.True(t, shouldBlock)
		assert.NotEmpty(t, errs)
	})
}
