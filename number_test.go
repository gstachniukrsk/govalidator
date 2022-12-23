package main_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"validator"
)

func TestNumberValidator(t *testing.T) {
	type args struct {
		ctx   context.Context
		value any
	}
	tests := []struct {
		name          string
		args          args
		input         any
		wantTwigBlock bool
		wantErrs      []error
	}{
		{
			name: "int",
			args: args{
				ctx: context.Background(),
			},
			input: 1,
		},
		{
			name: "float",
			args: args{
				ctx: context.Background(),
			},
			input: 1.0,
		},
		{
			name: "string",
			args: args{
				ctx: context.Background(),
			},
			input:         "1",
			wantTwigBlock: true,
			wantErrs: []error{
				main.NotANumberError{},
			},
		},
		{
			name: "bool",
			args: args{
				ctx: context.Background(),
			},
			input:         true,
			wantTwigBlock: true,
			wantErrs: []error{
				main.NotANumberError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := main.NumberValidator
			twigBlock, errs := v(tt.args.ctx, tt.input)
			assert.Equal(t, tt.wantTwigBlock, twigBlock)
			assert.Equal(t, tt.wantErrs, errs)
		})
	}
}
