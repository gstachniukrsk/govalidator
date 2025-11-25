package govalidator_test

import (
	"context"
	"testing"

	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestUUIDValidator(t *testing.T) {
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
			name: "valid UUID v1",
			args: args{
				ctx:   context.Background(),
				value: "550e8400-e29b-11d4-a716-446655440000",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid UUID v4",
			args: args{
				ctx:   context.Background(),
				value: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid UUID v5",
			args: args{
				ctx:   context.Background(),
				value: "886313e1-3b8a-5372-9b90-0c9aee199e5d",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid UUID - uppercase",
			args: args{
				ctx:   context.Background(),
				value: "550E8400-E29B-41D4-A716-446655440000",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid UUID - mixed case",
			args: args{
				ctx:   context.Background(),
				value: "550e8400-E29b-41D4-a716-446655440000",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "invalid UUID - wrong format",
			args: args{
				ctx:   context.Background(),
				value: "550e8400e29b41d4a716446655440000",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidUUIDError{Value: "550e8400e29b41d4a716446655440000"},
			},
		},
		{
			name: "invalid UUID - too short",
			args: args{
				ctx:   context.Background(),
				value: "550e8400-e29b-41d4-a716",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidUUIDError{Value: "550e8400-e29b-41d4-a716"},
			},
		},
		{
			name: "invalid UUID - wrong separator",
			args: args{
				ctx:   context.Background(),
				value: "550e8400_e29b_41d4_a716_446655440000",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidUUIDError{Value: "550e8400_e29b_41d4_a716_446655440000"},
			},
		},
		{
			name: "invalid UUID - invalid characters",
			args: args{
				ctx:   context.Background(),
				value: "550e8400-e29b-41d4-a716-44665544000g",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidUUIDError{Value: "550e8400-e29b-41d4-a716-44665544000g"},
			},
		},
		{
			name: "invalid UUID - wrong version",
			args: args{
				ctx:   context.Background(),
				value: "550e8400-e29b-61d4-a716-446655440000",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidUUIDError{Value: "550e8400-e29b-61d4-a716-446655440000"},
			},
		},
		{
			name: "invalid UUID - wrong variant",
			args: args{
				ctx:   context.Background(),
				value: "550e8400-e29b-41d4-c716-446655440000",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidUUIDError{Value: "550e8400-e29b-41d4-c716-446655440000"},
			},
		},
		{
			name: "invalid UUID - empty string",
			args: args{
				ctx:   context.Background(),
				value: "",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidUUIDError{Value: ""},
			},
		},
		{
			name: "not a string - integer",
			args: args{
				ctx:   context.Background(),
				value: 12345,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				govalidator.NotAStringError{},
			},
		},
		{
			name: "not a string - nil",
			args: args{
				ctx:   context.Background(),
				value: nil,
			},
			wantTwigBlock: true,
			wantErrs: []error{
				govalidator.NotAStringError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTwigBlock, gotErrs := govalidator.UUIDValidator(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock, "UUIDValidator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantErrs, gotErrs, "UUIDValidator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
