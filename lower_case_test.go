package govalidator_test

import (
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLowerCaseValidator(t *testing.T) {
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
				value: "john",
			},
		},
		{
			name: "lower case",
			args: args{
				value: "JOHN",
			},
			wantErrs: []error{govalidator.NotLowerCasedError{
				Input: "JOHN",
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
			wantErrs: []error{govalidator.NotLowerCasedError{
				Input: "JoHn",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTwigBlock, gotErrs := govalidator.LowerCaseValidator(nil, tt.args.value)
			assert.Equalf(t, tt.wantTwigBlock, gotTwigBlock, "LowerCaseValidator(%v, %v)", nil, tt.args.value)
			assert.Equalf(t, tt.wantErrs, gotErrs, "LowerCaseValidator(%v, %v)", nil, tt.args.value)
		})
	}
}

func TestNotLowerCasedError_Error(t *testing.T) {
	t.Run("error message", func(t *testing.T) {
		err := govalidator.NotLowerCasedError{
			Input: "john",
		}
		assert.Equal(t, "\"john\" is not lower cased", err.Error())
	})
}
