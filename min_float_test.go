package govalidator_test

import (
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinFloatValidator(t *testing.T) {
	type args struct {
		minFloat float64
	}
	tests := []struct {
		name          string
		args          args
		input         any
		wantTwigBlock bool
		wantErrs      []error
	}{
		{
			name: "eq positive",
			args: args{
				minFloat: 1.0,
			},
			input: 1.0,
		},
		{
			name: "eq zero",
			args: args{
				minFloat: 0.0,
			},
			input: 0.0,
		},
		{
			name: "eq negative",
			args: args{
				minFloat: -1.0,
			},
			input: -1.0,
		},
		{
			name: "gt positive",
			args: args{
				minFloat: 1.0,
			},
			input: 1.5,
		},
		{
			name: "gt zero",
			args: args{
				minFloat: 0.0,
			},
			input: 0.5,
		},
		{
			name: "gt negative",
			args: args{
				minFloat: -1.0,
			},
			input: -0.5,
		},
		{
			name: "lt positive",
			args: args{
				minFloat: 1.0,
			},
			input:         0.5,
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.FloatIsLesserThanError{
					MinFloat: 1.0,
				},
			},
		},
		{
			name: "lt zero",
			args: args{
				minFloat: 0.0,
			},
			input:         -0.5,
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.FloatIsLesserThanError{
					MinFloat: 0.0,
				},
			},
		},
		{
			name: "lt negative",
			args: args{
				minFloat: -1.0,
			},
			input:         -1.5,
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.FloatIsLesserThanError{
					MinFloat: -1.0,
				},
			},
		},
		{
			name: "not float",
			args: args{
				minFloat: 1.0,
			},
			input:         "not float",
			wantTwigBlock: true,
			wantErrs: []error{
				govalidator.NotAFloatError{},
			},
		},
		{
			name: "int",
			args: args{
				minFloat: 1.0,
			},
			input: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := govalidator.MinFloatValidator(tt.args.minFloat)
			gotTwigBlock, gotErrs := v(nil, tt.input)

			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock)
			assert.Equal(t, tt.wantErrs, gotErrs)
		})
	}
}

func Test_FloatIsLesserThanError(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "error message",
			want: "value is less than min",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := govalidator.FloatIsLesserThanError{
				MinFloat: 1.0,
			}
			assert.Equal(t, tt.want, e.Error())
		})
	}
}
