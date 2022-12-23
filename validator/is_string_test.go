package validator_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"validator/validator"
)

func TestStringValidator(t *testing.T) {
	type args struct {
		ctx   context.Context
		value any
	}
	tests := []struct {
		name          string
		args          args
		wantTwigBlock bool
		wantErrs      []error
	}{
		{
			name: "ok",
			args: args{
				ctx:   nil,
				value: "ok",
			},
		},
		{
			name: "ptr",
			args: args{
				ctx:   nil,
				value: strPtr("ok"),
			},
		},
		{
			name:          "nil",
			args:          args{},
			wantTwigBlock: true,
			wantErrs: []error{
				validator.NotAStringError{},
			},
		},
		{
			name: "not a string",
			args: args{
				ctx:   nil,
				value: 1,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				validator.NotAStringError{},
			},
		},
		{
			name: "nil pointer",
			args: args{
				ctx:   nil,
				value: (*string)(nil),
			},
			wantTwigBlock: true,
			wantErrs: []error{
				validator.NotAStringError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTwigBlock, gotErrs := validator.StringValidator(tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantTwigBlock, gotTwigBlock, "StringValidator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantErrs, gotErrs, "StringValidator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
