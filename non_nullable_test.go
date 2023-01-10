package go_validator_test

import (
	"context"
	"github.com/gstachniukrsk/go_validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNonNullableValidator(t *testing.T) {
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
			name: "nil",
			args: args{
				ctx:   context.Background(),
				value: nil,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.RequiredError{},
			},
		},
		{
			name: "empty string",
			args: args{
				ctx:   context.Background(),
				value: "",
			},
		},
		{
			name: "non-empty string",
			args: args{
				ctx:   context.Background(),
				value: "foo",
			},
		},
		{
			name: "nil map",
			args: args{
				ctx:   context.Background(),
				value: map[string]interface{}(nil),
			},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.RequiredError{},
			},
		},
		{
			name: "nil ptr map",
			args: args{
				ctx:   context.Background(),
				value: (*map[string]interface{})(nil),
			},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.RequiredError{},
			},
		},
		{
			name: "nil interface",
			args: args{
				ctx:   context.Background(),
				value: (*interface{})(nil),
			},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.RequiredError{},
			},
		},
		{
			name: "nil ptr interface",
			args: args{
				ctx:   context.Background(),
				value: (*interface{})(nil),
			},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.RequiredError{},
			},
		},
		{
			name: "nil slice",
			args: args{
				ctx:   context.Background(),
				value: []interface{}(nil),
			},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.RequiredError{},
			},
		},
		{
			name: "nil ptr slice",
			args: args{
				ctx:   context.Background(),
				value: (*[]interface{})(nil),
			},
			wantTwigBlock: true,
			wantErrs: []error{
				go_validator.RequiredError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTwigBlock, gotErrs := go_validator.NonNullableValidator(tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantTwigBlock, gotTwigBlock, "NonNullableValidator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantErrs, gotErrs, "NonNullableValidator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
