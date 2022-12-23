package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNullableValidator(t *testing.T) {
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
		},
		{
			name: "nil ptr map",
			args: args{
				ctx:   context.Background(),
				value: (*map[string]interface{})(nil),
			},
			wantTwigBlock: true,
		},
		{
			name: "empty map",
			args: args{
				ctx:   context.Background(),
				value: map[string]interface{}{},
			},
		},
		{
			name: "non-empty map",
			args: args{
				ctx:   context.Background(),
				value: map[string]interface{}{"foo": "bar"},
			},
		},
		{
			name: "nil slice",
			args: args{
				ctx:   context.Background(),
				value: []interface{}(nil),
			},
			wantTwigBlock: true,
		},
		{
			name: "nil ptr slice",
			args: args{
				ctx:   context.Background(),
				value: (*[]interface{})(nil),
			},
			wantTwigBlock: true,
		},
		{
			name: "empty slice",
			args: args{
				ctx:   context.Background(),
				value: []interface{}{},
			},
		},
		{
			name: "non-empty slice",
			args: args{
				ctx:   context.Background(),
				value: []interface{}{"foo", "bar"},
			},
		},
		{
			name: "non-empty nil ptr slice",
			args: args{
				ctx:   context.Background(),
				value: &[]interface{}{1, 2, 3},
			},
		},
		{
			name: "nil interface",
			args: args{
				ctx:   context.Background(),
				value: (interface{})(nil),
			},
			wantTwigBlock: true,
		},
		{
			name: "nil ptr interface",
			args: args{
				ctx:   context.Background(),
				value: (*interface{})(nil),
			},
			wantTwigBlock: true,
		},
		{
			name: "non-empty interface",
			args: args{
				ctx:   context.Background(),
				value: interface{}(map[string]string{"foo": "bar"}),
			},
		},
		{
			name: "non-empty ptr interface",
			args: args{
				ctx:   context.Background(),
				value: interface{}(&map[string]string{"foo": "bar"}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTwigBlock, gotErrs := NullableValidator(tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantTwigBlock, gotTwigBlock, "NullableValidator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equalf(t, tt.wantErrs, gotErrs, "NullableValidator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
