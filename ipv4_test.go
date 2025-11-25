package govalidator_test

import (
	"context"
	"testing"

	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestIPv4Validator(t *testing.T) {
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
			name: "valid IPv4 - standard",
			args: args{
				ctx:   context.Background(),
				value: "192.168.1.1",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv4 - localhost",
			args: args{
				ctx:   context.Background(),
				value: "127.0.0.1",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv4 - zeros",
			args: args{
				ctx:   context.Background(),
				value: "0.0.0.0",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv4 - max values",
			args: args{
				ctx:   context.Background(),
				value: "255.255.255.255",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv4 - with CIDR",
			args: args{
				ctx:   context.Background(),
				value: "192.168.1.0/24",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid IPv4 - CIDR /32",
			args: args{
				ctx:   context.Background(),
				value: "10.0.0.1/32",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "invalid IPv4 - out of range octet",
			args: args{
				ctx:   context.Background(),
				value: "256.168.1.1",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv4Error{Value: "256.168.1.1"},
			},
		},
		{
			name: "invalid IPv4 - too few octets",
			args: args{
				ctx:   context.Background(),
				value: "192.168.1",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv4Error{Value: "192.168.1"},
			},
		},
		{
			name: "invalid IPv4 - too many octets",
			args: args{
				ctx:   context.Background(),
				value: "192.168.1.1.1",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv4Error{Value: "192.168.1.1.1"},
			},
		},
		{
			name: "invalid IPv4 - letters",
			args: args{
				ctx:   context.Background(),
				value: "192.168.a.1",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv4Error{Value: "192.168.a.1"},
			},
		},
		{
			name: "invalid IPv4 - negative number",
			args: args{
				ctx:   context.Background(),
				value: "192.168.-1.1",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv4Error{Value: "192.168.-1.1"},
			},
		},
		{
			name: "invalid IPv4 - empty string",
			args: args{
				ctx:   context.Background(),
				value: "",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv4Error{Value: ""},
			},
		},
		{
			name: "invalid - IPv6 address",
			args: args{
				ctx:   context.Background(),
				value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv4Error{Value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			},
		},
		{
			name: "invalid - IPv6 short form",
			args: args{
				ctx:   context.Background(),
				value: "::1",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidIPv4Error{Value: "::1"},
			},
		},
		{
			name: "not a string - integer",
			args: args{
				ctx:   context.Background(),
				value: 192168001001,
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
			gotTwigBlock, gotErrs := govalidator.IPv4Validator(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock, "IPv4Validator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantErrs, gotErrs, "IPv4Validator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
