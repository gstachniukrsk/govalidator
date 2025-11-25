package govalidator_test

import (
	"context"
	"testing"

	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestBase64Validator(t *testing.T) {
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
			name: "valid base64 - standard with padding",
			args: args{
				ctx:   context.Background(),
				value: "SGVsbG8gV29ybGQ=",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid base64 - standard without padding",
			args: args{
				ctx:   context.Background(),
				value: "SGVsbG8gV29ybGQ",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid base64 - URL-safe with padding",
			args: args{
				ctx:   context.Background(),
				value: "aHR0cHM6Ly9leGFtcGxlLmNvbS9wYXRoP3F1ZXJ5PXZhbHVl",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid base64 - URL-safe without padding",
			args: args{
				ctx:   context.Background(),
				value: "aHR0cHM6Ly9leGFtcGxlLmNvbS9wYXRoP3F1ZXJ5PXZhbHVl",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid base64 - single character",
			args: args{
				ctx:   context.Background(),
				value: "YQ==",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid base64 - numbers",
			args: args{
				ctx:   context.Background(),
				value: "MTIzNDU2Nzg5MA==",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid base64 - binary data",
			args: args{
				ctx:   context.Background(),
				value: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid base64 - long string",
			args: args{
				ctx:   context.Background(),
				value: "VGhpcyBpcyBhIHZlcnkgbG9uZyBzdHJpbmcgdGhhdCB3aWxsIGJlIGVuY29kZWQgaW4gYmFzZTY0IGZvcm1hdA==",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "invalid base64 - invalid characters",
			args: args{
				ctx:   context.Background(),
				value: "Hello@World!",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidBase64Error{Value: "Hello@World!"},
			},
		},
		{
			name: "invalid base64 - spaces",
			args: args{
				ctx:   context.Background(),
				value: "SGVs bG8g V29y bGQ=",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidBase64Error{Value: "SGVs bG8g V29y bGQ="},
			},
		},
		{
			name: "invalid base64 - wrong padding",
			args: args{
				ctx:   context.Background(),
				value: "SGVsbG8gV29ybGQ===",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidBase64Error{Value: "SGVsbG8gV29ybGQ==="},
			},
		},
		{
			name: "invalid base64 - plain text",
			args: args{
				ctx:   context.Background(),
				value: "This is not base64!",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidBase64Error{Value: "This is not base64!"},
			},
		},
		{
			name: "invalid base64 - empty string",
			args: args{
				ctx:   context.Background(),
				value: "",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidBase64Error{Value: ""},
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
			gotTwigBlock, gotErrs := govalidator.Base64Validator(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock, "Base64Validator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantErrs, gotErrs, "Base64Validator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
