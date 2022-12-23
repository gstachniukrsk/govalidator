package go_validator

import (
"github.com/stretchr/testify/assert"
"testing"
)

func TestMaxSizeValidator(t *testing.T) {
	type args struct {
		maxSize int
		blocks  bool
	}
	tests := []struct {
		name         string
		args         args
		input        any
		blocksTwig   bool
		expectedErrs []error
	}{
		{
			name: "empty",
			args: args{
				maxSize: 5,
				blocks:  false,
			},
			input: []interface{}{},
		},
		{
			name: "eq",
			args: args{
				maxSize: 5,
				blocks:  false,
			},
			input: []interface{}{1, 2, 3, 4, 5},
		},
		{
			name: "lt",
			args: args{
				maxSize: 5,
				blocks:  false,
			},
			input: []interface{}{1, 2, 3, 4},
		},
		{
			name: "gt",
			args: args{
				maxSize: 5,
				blocks:  true,
			},
			input: []interface{}{1, 2, 3, 4, 5, 6},
			expectedErrs: []error{
				MaxSizeError{
					MaxSize:    5,
					ActualSize: 6,
				},
			},
			blocksTwig: true,
		},
		{
			name: "gt, non blocking",
			args: args{
				maxSize: 5,
				blocks:  false,
			},
			input: []interface{}{1, 2, 3, 4, 5, 6},
			expectedErrs: []error{
				MaxSizeError{
					MaxSize:    5,
					ActualSize: 6,
				},
			},
			blocksTwig: false,
		},
		{
			name: "nil",
			args: args{
				maxSize: 5,
				blocks:  false,
			},
			input: nil,
			expectedErrs: []error{
				NotAListError{},
			},
			blocksTwig: true,
		},
		{
			name: "nil list",
			args: args{
				maxSize: 5,
				blocks:  false,
			},
			input: []interface{}(nil),
			expectedErrs: []error{
				NotAListError{},
			},
			blocksTwig: true,
		},
		{
			name: "nil ptr list",
			args: args{
				maxSize: 5,
				blocks:  false,
			},
			input: (*[]interface{})(nil),
			expectedErrs: []error{
				NotAListError{},
			},
			blocksTwig: true,
		},
		{
			name: "ptr, eq",
			args: args{
				maxSize: 5,
				blocks:  false,
			},
			input: &[]interface{}{1, 2, 3, 4, 5},
		},
		{
			name: "ptr, gt, blocking",
			args: args{
				maxSize: 5,
				blocks:  true,
			},
			input: &[]interface{}{1, 2, 3, 4, 5, 6},
			expectedErrs: []error{
				MaxSizeError{
					MaxSize:    5,
					ActualSize: 6,
				},
			},
			blocksTwig: true,
		},
		{
			name: "ptr, gt, non blocking",
			args: args{
				maxSize: 5,
				blocks:  false,
			},
			input: &[]interface{}{1, 2, 3, 4, 5, 6},
			expectedErrs: []error{
				MaxSizeError{
					MaxSize:    5,
					ActualSize: 6,
				},
			},
			blocksTwig: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := MaxSizeValidator(tt.args.maxSize, tt.args.blocks)
			twigBlock, errs := v(nil, tt.input)

			assert.Equal(t, tt.blocksTwig, twigBlock)
			assert.Equal(t, tt.expectedErrs, errs)
		})
	}
}
