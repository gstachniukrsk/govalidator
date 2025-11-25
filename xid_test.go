package govalidator_test

import (
	"context"
	"testing"

	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestXIDValidator(t *testing.T) {
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
			name: "valid XID - example 1",
			args: args{
				ctx:   context.Background(),
				value: "9m4e2mr0ui3e8a215n4g",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid XID - example 2",
			args: args{
				ctx:   context.Background(),
				value: "c0o9s9t62d9s72q7irug",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid XID - all numbers",
			args: args{
				ctx:   context.Background(),
				value: "01234567890123456789",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid XID - all letters (a-v)",
			args: args{
				ctx:   context.Background(),
				value: "abcdefghijklmnopqrst",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid XID - edge case with 'v'",
			args: args{
				ctx:   context.Background(),
				value: "vvvvvvvvvvvvvvvvvvvv",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "invalid XID - too short",
			args: args{
				ctx:   context.Background(),
				value: "9m4e2mr0ui3e8a215n4",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: "9m4e2mr0ui3e8a215n4"},
			},
		},
		{
			name: "invalid XID - too long",
			args: args{
				ctx:   context.Background(),
				value: "9m4e2mr0ui3e8a215n4gg",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: "9m4e2mr0ui3e8a215n4gg"},
			},
		},
		{
			name: "invalid XID - contains uppercase",
			args: args{
				ctx:   context.Background(),
				value: "9M4E2MR0UI3E8A215N4G",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: "9M4E2MR0UI3E8A215N4G"},
			},
		},
		{
			name: "invalid XID - contains invalid letter 'w'",
			args: args{
				ctx:   context.Background(),
				value: "9m4e2mr0ui3e8a215w4g",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: "9m4e2mr0ui3e8a215w4g"},
			},
		},
		{
			name: "invalid XID - contains invalid letter 'x'",
			args: args{
				ctx:   context.Background(),
				value: "9m4e2mr0ui3e8a215x4g",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: "9m4e2mr0ui3e8a215x4g"},
			},
		},
		{
			name: "invalid XID - contains invalid letter 'z'",
			args: args{
				ctx:   context.Background(),
				value: "9m4e2mr0ui3e8a215z4g",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: "9m4e2mr0ui3e8a215z4g"},
			},
		},
		{
			name: "invalid XID - contains special character",
			args: args{
				ctx:   context.Background(),
				value: "9m4e2mr0ui3e8a215-4g",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: "9m4e2mr0ui3e8a215-4g"},
			},
		},
		{
			name: "invalid XID - contains space",
			args: args{
				ctx:   context.Background(),
				value: "9m4e2mr0ui3e8a21 n4g",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: "9m4e2mr0ui3e8a21 n4g"},
			},
		},
		{
			name: "invalid XID - empty string",
			args: args{
				ctx:   context.Background(),
				value: "",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: ""},
			},
		},
		{
			name: "invalid XID - UUID format",
			args: args{
				ctx:   context.Background(),
				value: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidXIDError{Value: "550e8400-e29b-41d4-a716-446655440000"},
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
			gotTwigBlock, gotErrs := govalidator.XIDValidator(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock, "XIDValidator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantErrs, gotErrs, "XIDValidator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
