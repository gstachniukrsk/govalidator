package govalidator_test

import (
	"context"
	"testing"

	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestURLValidator(t *testing.T) {
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
			name: "valid HTTP URL",
			args: args{
				ctx:   context.Background(),
				value: "http://example.com",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid HTTPS URL",
			args: args{
				ctx:   context.Background(),
				value: "https://example.com",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid URL with path",
			args: args{
				ctx:   context.Background(),
				value: "https://example.com/path/to/resource",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid URL with query",
			args: args{
				ctx:   context.Background(),
				value: "https://example.com?param=value&foo=bar",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid URL with port",
			args: args{
				ctx:   context.Background(),
				value: "https://example.com:8080/path",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid URL with subdomain",
			args: args{
				ctx:   context.Background(),
				value: "https://api.example.com/v1/users",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid URL with fragment",
			args: args{
				ctx:   context.Background(),
				value: "https://example.com/page#section",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid URL with authentication",
			args: args{
				ctx:   context.Background(),
				value: "https://user:pass@example.com",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid FTP URL",
			args: args{
				ctx:   context.Background(),
				value: "ftp://ftp.example.com/file.txt",
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "invalid URL - no scheme",
			args: args{
				ctx:   context.Background(),
				value: "example.com",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidURLError{Value: "example.com"},
			},
		},
		{
			name: "invalid URL - missing host",
			args: args{
				ctx:   context.Background(),
				value: "https://",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidURLError{Value: "https://"},
			},
		},
		{
			name: "invalid URL - spaces",
			args: args{
				ctx:   context.Background(),
				value: "https://example .com",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidURLError{Value: "https://example .com"},
			},
		},
		{
			name: "invalid URL - empty string",
			args: args{
				ctx:   context.Background(),
				value: "",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidURLError{Value: ""},
			},
		},
		{
			name: "invalid URL - malformed",
			args: args{
				ctx:   context.Background(),
				value: "ht!tp://example.com",
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidURLError{Value: "ht!tp://example.com"},
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
			gotTwigBlock, gotErrs := govalidator.URLValidator(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock, "URLValidator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantErrs, gotErrs, "URLValidator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
