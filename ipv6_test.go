package govalidator_test

import (
	"context"
	"testing"

	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestIPv6Validator(t *testing.T) {
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
			name: "valid IPv6 - full form",
			args: args{
				ctx:   context.Background(),
				value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv6 - short form",
			args: args{
				ctx:   context.Background(),
				value: "2001:db8:85a3::8a2e:370:7334",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv6 - localhost",
			args: args{
				ctx:   context.Background(),
				value: "::1",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv6 - all zeros",
			args: args{
				ctx:   context.Background(),
				value: "::",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv6 - link local",
			args: args{
				ctx:   context.Background(),
				value: "fe80::1",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv6 - with CIDR",
			args: args{
				ctx:   context.Background(),
				value: "2001:db8::/32",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv6 - CIDR /128",
			args: args{
				ctx:   context.Background(),
				value: "2001:db8::1/128",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv6 - uppercase",
			args: args{
				ctx:   context.Background(),
				value: "2001:0DB8:85A3:0000:0000:8A2E:0370:7334",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "invalid IPv6 - too many groups",
			args: args{
				ctx:   context.Background(),
				value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334:extra",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv6Error{Value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334:extra"},
			},
		},
		{
			name: "invalid IPv6 - invalid hex characters",
			args: args{
				ctx:   context.Background(),
				value: "2001:0db8:85g3::1",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv6Error{Value: "2001:0db8:85g3::1"},
			},
		},
		{
			name: "invalid IPv6 - too many colons",
			args: args{
				ctx:   context.Background(),
				value: "2001:::1",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv6Error{Value: "2001:::1"},
			},
		},
		{
			name: "invalid IPv6 - empty string",
			args: args{
				ctx:   context.Background(),
				value: "",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv6Error{Value: ""},
			},
		},
		{
			name: "invalid - IPv4 address",
			args: args{
				ctx:   context.Background(),
				value: "192.168.1.1",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv6Error{Value: "192.168.1.1"},
			},
		},
		{
			name: "invalid - IPv4 CIDR",
			args: args{
				ctx:   context.Background(),
				value: "192.168.1.0/24",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv6Error{Value: "192.168.1.0/24"},
			},
		},
		{
			name: "not a string - integer",
			args: args{
				ctx:   context.Background(),
				value: 123456,
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
			gotTwigBlock, gotErrs := govalidator.IPv6Validator(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock, "IPv6Validator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantErrs, gotErrs, "IPv6Validator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
