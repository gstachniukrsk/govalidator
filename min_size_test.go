package go_validator_test

import (
	"github.com/gstachniukrsk/go_validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinSizeValidator(t *testing.T) {
	type args struct {
		minSize int
		blocks  bool
	}
	tests := []struct {
		name          string
		args          args
		input         any
		wantTwigBlock bool
		wantErrs      []error
	}{
		{
			name: "empty non blocking",
			args: args{
				minSize: 5,
			},
			input:         []interface{}{},
			wantTwigBlock: false,
			wantErrs: []error{
				go_validator.MinSizeError{
					MinSize:    5,
					ActualSize: 0,
				},
			},
		},
		{
			name: "empty blocking",
			args: args{
				minSize: 5,
				blocks:  true,
			},
			input:         []interface{}{},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.MinSizeError{
					MinSize:    5,
					ActualSize: 0,
				},
			},
		},
		{
			name: "eq non blocking",
			args: args{
				minSize: 5,
				blocks:  false,
			},
			input: []interface{}{
				1, 2, 3, 4, 5,
			},
		},
		{
			name: "eq blocking",
			args: args{
				minSize: 5,
				blocks:  true,
			},
			input: []interface{}{
				1, 2, 3, 4, 5,
			},
		},
		{
			name: "lt non blocking",
			args: args{
				minSize: 5,
				blocks:  false,
			},
			input: []interface{}{
				1, 2, 3, 4,
			},
			wantTwigBlock: false,
			wantErrs: []error{
				go_validator.MinSizeError{
					MinSize:    5,
					ActualSize: 4,
				},
			},
		},
		{
			name: "lt blocking",
			args: args{
				minSize: 5,
				blocks:  true,
			},
			input: []interface{}{
				1, 2, 3, 4,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.MinSizeError{
					MinSize:    5,
					ActualSize: 4,
				},
			},
		},
		{
			name: "wrong type blocking",
			args: args{
				minSize: 5,
				blocks:  true,
			},
			input:         "",
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.NotAListError{},
			},
		},
		{
			name: "wrong type non blocking",
			args: args{
				minSize: 5,
				blocks:  false,
			},
			input:         "",
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.NotAListError{},
			},
		},
		{
			name: "ptr nil blocking",
			args: args{
				minSize: 5,
				blocks:  true,
			},
			input:         (*[]interface{})(nil),
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.NotAListError{},
			},
		},
		{
			name: "ptr nil non blocking",
			args: args{
				minSize: 5,
				blocks:  false,
			},
			input:         (*[]interface{})(nil),
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.NotAListError{},
			},
		},
		{
			name: "eq ptr non nil blocking",
			args: args{
				minSize: 5,
				blocks:  true,
			},
			input: &[]interface{}{
				1, 2, 3, 4, 5,
			},
		},
		{
			name: "eq ptr non nil non blocking",
			args: args{
				minSize: 5,
				blocks:  false,
			},
			input: &[]interface{}{
				1, 2, 3, 4, 5,
			},
		},
		{
			name: "lt ptr non nil blocking",
			args: args{
				minSize: 5,
				blocks:  true,
			},
			input: &[]interface{}{
				1, 2, 3, 4,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.MinSizeError{
					MinSize:    5,
					ActualSize: 4,
				},
			},
		},
		{
			name: "lt ptr non nil non blocking",
			args: args{
				minSize: 5,
				blocks:  false,
			},
			input: &[]interface{}{
				1, 2, 3, 4,
			},
			wantTwigBlock: false,
			wantErrs: []error{
				go_validator.MinSizeError{
					MinSize:    5,
					ActualSize: 4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := go_validator.MinSizeValidator(tt.args.minSize, tt.args.blocks)

			gotTwigBlock, gotErrs := v(nil, tt.input)
			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock)
			assert.Equal(t, tt.wantErrs, gotErrs)
		})
	}
}
