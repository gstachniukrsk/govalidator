package validator_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"validator/validator"
)

func TestIsMapValidator(t *testing.T) {
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
				value: map[string]any{},
			},
		},
		{
			name: "ptr",
			args: args{
				value: &map[string]any{},
			},
		},
		{
			name:          "nil",
			args:          args{},
			wantTwigBlock: true,
			wantErrs: []error{
				validator.NotAMapError{},
			},
		},
		{
			name: "not a map",
			args: args{
				value: "not a map",
			},
			wantTwigBlock: true,
			wantErrs: []error{
				validator.NotAMapError{},
			},
		},
		{
			name: "nil pointer",
			args: args{
				value: (*map[string]any)(nil),
			},
			wantTwigBlock: true,
			wantErrs: []error{
				validator.NotAMapError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTwigBlock, gotErrs := validator.IsMapValidator(tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantTwigBlock, gotTwigBlock, "IsMapValidator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantErrs, gotErrs, "IsMapValidator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
