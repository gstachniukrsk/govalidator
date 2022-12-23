package go_validator_test

import (
"github.com/stretchr/testify/assert"
"testing"
"validator"
)

func TestMaxLengthValidator(t *testing.T) {
	type args struct {
		maxLength int
		input     any
	}
	tests := []struct {
		name         string
		args         args
		blocksTwig   bool
		expectedErrs []error
	}{
		{
			name: "empty",
			args: args{
				maxLength: 5,
				input:     "",
			},
		},
		{
			name: "eq",
			args: args{
				maxLength: 5,
				input:     "12345",
			},
		},
		{
			name: "lt",
			args: args{
				maxLength: 5,
				input:     "1234",
			},
		},
		{
			name: "gt",
			args: args{
				maxLength: 5,
				input:     "123456",
			},
			blocksTwig: false,
			expectedErrs: []error{
				main.StringTooLongError{
					MaxLength:    5,
					ActualLength: 6,
				},
			},
		},
		{
			name: "ptr ok",
			args: args{
				maxLength: 5,
				input:     strPtr("12345"),
			},
		},
		{
			name: "ptr too long",
			args: args{
				maxLength: 5,
				input:     strPtr("123456"),
			},
			blocksTwig: false,
			expectedErrs: []error{
				main.StringTooLongError{
					MaxLength:    5,
					ActualLength: 6,
				},
			},
		},
		{
			name: "ptr nil",
			args: args{
				maxLength: 5,
				input:     (*string)(nil),
			},
			blocksTwig: true,
			expectedErrs: []error{
				main.NotAStringError{},
			},
		},
		{
			name: "emoji with color suffix, special chars, exactly 11 chars",
			args: args{
				maxLength: 11,
				input:     "👍🏻ęąśłżźćńó",
			},
		},
		{
			name: "not a string",
			args: args{
				maxLength: 5,
				input:     123,
			},
			blocksTwig: true,
			expectedErrs: []error{
				main.NotAStringError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := main.MaxLengthValidator(tt.args.maxLength)
			twigBreak, errs := v(nil, tt.args.input)
			assert.Equal(t, tt.blocksTwig, twigBreak)
			assert.Equal(t, tt.expectedErrs, errs)
		})
	}
}

func TestStringTooLongError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  main.StringTooLongError
		want string
	}{
		{
			name: "happy path",
			err: main.StringTooLongError{
				MaxLength:    5,
				ActualLength: 6,
			},
			want: "expected at most 5 characters, got 6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}
