package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsListValidator(t *testing.T) {
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
				value: []interface{}{},
			},
			wantTwigBlock: false,
		},
		{
			name: "nil",
			args: args{
				ctx: nil,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				NotAListError{},
			},
		},
		{
			name: "not a list",
			args: args{
				ctx:   nil,
				value: "not a list",
			},
			wantTwigBlock: true,
			wantErrs: []error{
				NotAListError{},
			},
		},
		{
			name: "nil pointer",
			args: args{
				ctx:   nil,
				value: (*[]interface{})(nil),
			},
			wantTwigBlock: true,
			wantErrs: []error{
				NotAListError{},
			},
		},
		{
			name: "ptr",
			args: args{
				ctx:   nil,
				value: &[]interface{}{},
			},
			wantTwigBlock: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTwigBlock, gotErrs := IsListValidator(tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantTwigBlock, gotTwigBlock, "IsListValidator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantErrs, gotErrs, "IsListValidator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
