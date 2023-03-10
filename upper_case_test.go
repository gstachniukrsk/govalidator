package govalidator_test

import (
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpperCaseValidator(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name          string
		args          args
		wantTwigBlock bool
		wantErrs      []error
	}{
		{
			name: "happy path",
			args: args{
				value: "JOHN",
			},
		},
		{
			name: "lower case",
			args: args{
				value: "john",
			},
			wantErrs: []error{govalidator.NotUpperCasedError{
				Input: "john",
			}},
		},
		{
			name: "not a string",
			args: args{
				value: 1,
			},
			wantErrs:      []error{govalidator.NotAStringError{}},
			wantTwigBlock: true,
		},
		{
			name: "mixed case",
			args: args{
				value: "JoHn",
			},
			wantErrs: []error{govalidator.NotUpperCasedError{
				Input: "JoHn",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTwigBlock, gotErrs := govalidator.UpperCaseValidator(nil, tt.args.value)
			assert.Equalf(t, tt.wantTwigBlock, gotTwigBlock, "UpperCaseValidator(%v, %v)", nil, tt.args.value)
			assert.Equalf(t, tt.wantErrs, gotErrs, "UpperCaseValidator(%v, %v)", nil, tt.args.value)
		})
	}
}

func TestNotUpperCasedError_Error(t *testing.T) {
	t.Run("error message", func(t *testing.T) {
		err := govalidator.NotUpperCasedError{
			Input: "john",
		}
		assert.Equal(t, "\"john\" is not upper cased", err.Error())
	})
}
