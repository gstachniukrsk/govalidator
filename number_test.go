package govalidator_test

import (
	"context"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumberValidator(t *testing.T) {
	type args struct {
		ctx context.Context
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
				govalidator.NotANumberError{},
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
				govalidator.NotANumberError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := govalidator.NumberValidator
			twigBlock, errs := v(tt.args.ctx, tt.input)
			assert.Equal(t, tt.wantTwigBlock, twigBlock)
			assert.Equal(t, tt.wantErrs, errs)
		})
	}
}
