package go_validator_test

import (
	"github.com/gstachniukrsk/go_validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinLengthValidator(t *testing.T) {
	type args struct {
		minLength int
	}
	tests := []struct {
		name          string
		args          args
		input         any
		wantTwigBlock bool
		wantErrs      []error
	}{
		{
			name: "empty",
			args: args{
				minLength: 5,
			},
			input:         "",
			wantTwigBlock: false,
			wantErrs: []error{
				go_validator.StringTooShortError{
					MinLength: 5,
				},
			},
		},
		{
			name: "eq",
			args: args{
				minLength: 5,
			},
			input: "12345",
		},
		{
			name: "lt",
			args: args{
				minLength: 5,
			},
			input:         "1234",
			wantTwigBlock: false,
			wantErrs: []error{
				go_validator.StringTooShortError{
					MinLength: 5,
				},
			},
		},
		{
			name: "gt",
			args: args{
				minLength: 5,
			},
			input: "123456",
		},
		{
			name: "empty string on min 0",
			args: args{
				minLength: 0,
			},
			input: "",
		},
		{
			name: "nil fail",
			args: args{
				minLength: 5,
			},
			input:         nil,
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.NotAStringError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := go_validator.MinLengthValidator(tt.args.minLength)
			gotTwigBlock, gotErrs := v(nil, tt.input)

			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock)
			assert.Equal(t, tt.wantErrs, gotErrs)
		})
	}
}

func TestStringTooShortError_Error(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		err := go_validator.StringTooShortError{
			MinLength: 5,
		}
		assert.Equal(t, "expected at least 5 characters", err.Error())
	})

	t.Run("non empty", func(t *testing.T) {
		err := go_validator.StringTooShortError{
			MinLength: 5,
		}
		assert.Equal(t, "expected at least 5 characters", err.Error())
	})
}
