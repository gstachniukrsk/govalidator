package govalidator_test

import (
	"context"
	"testing"

	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestEmailValidator(t *testing.T) {
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
			name: "valid simple email",
			args: args{
				ctx:   context.Background(),
				value: "user@example.com",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid email with subdomain",
			args: args{
				ctx:   context.Background(),
				value: "user@mail.example.com",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid email with dots",
			args: args{
				ctx:   context.Background(),
				value: "first.last@example.com",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid email with plus",
			args: args{
				ctx:   context.Background(),
				value: "user+tag@example.com",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid email with numbers",
			args: args{
				ctx:   context.Background(),
				value: "user123@example123.com",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid email with hyphen in domain",
			args: args{
				ctx:   context.Background(),
				value: "user@my-company.com",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "invalid email - no @",
			args: args{
				ctx:   context.Background(),
				value: "userexample.com",
			},
			wantTwigBlock: false,
		},
		{
			name: "invalid email - multiple @",
			args: args{
				ctx:   context.Background(),
				value: "user@@example.com",
			},
			wantTwigBlock: false,
		},
		{
			name: "invalid email - no domain",
			args: args{
				ctx:   context.Background(),
				value: "user@",
			},
			wantTwigBlock: false,
		},
		{
			name: "invalid email - no local part",
			args: args{
				ctx:   context.Background(),
				value: "@example.com",
			},
			wantTwigBlock: false,
		},
		{
			name: "invalid email - spaces",
			args: args{
				ctx:   context.Background(),
				value: "user @example.com",
			},
			wantTwigBlock: false,
		},
		{
			name: "invalid email - empty string",
			args: args{
				ctx:   context.Background(),
				value: "",
			},
			wantTwigBlock: false,
		},
		{
			name: "not a string - integer",
			args: args{
				ctx:   context.Background(),
				value: 123,
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
			gotTwigBlock, gotErrs := govalidator.EmailValidator(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock, "EmailValidator(%v, %v)", tt.args.ctx, tt.args.value)

			if tt.wantErrs != nil {
				assert.Equal(t, tt.wantErrs, gotErrs, "EmailValidator(%v, %v)", tt.args.ctx, tt.args.value)
			}
		})
	}
}
