package govalidator_test

import (
	"context"
	"testing"

	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestJSONValidator(t *testing.T) {
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
			name: "valid JSON - object",
			args: args{
				ctx:   context.Background(),
				value: `{"name":"John","age":30}`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid JSON - array",
			args: args{
				ctx:   context.Background(),
				value: `[1,2,3,4,5]`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid JSON - string",
			args: args{
				ctx:   context.Background(),
				value: `"hello world"`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid JSON - number",
			args: args{
				ctx:   context.Background(),
				value: `42`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid JSON - boolean true",
			args: args{
				ctx:   context.Background(),
				value: `true`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid JSON - boolean false",
			args: args{
				ctx:   context.Background(),
				value: `false`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid JSON - null",
			args: args{
				ctx:   context.Background(),
				value: `null`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid JSON - nested object",
			args: args{
				ctx:   context.Background(),
				value: `{"user":{"name":"John","address":{"city":"NYC"}}}`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid JSON - array of objects",
			args: args{
				ctx:   context.Background(),
				value: `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "valid JSON - with whitespace",
			args: args{
				ctx: context.Background(),
				value: `{
					"name": "John",
					"age": 30
				}`,
			},
			wantTwigBlock: false,
			wantErrs:      nil,
		},
		{
			name: "invalid JSON - missing quote",
			args: args{
				ctx:   context.Background(),
				value: `{"name":"John,"age":30}`,
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidJSONError{Value: `{"name":"John,"age":30}`},
			},
		},
		{
			name: "invalid JSON - trailing comma",
			args: args{
				ctx:   context.Background(),
				value: `{"name":"John","age":30,}`,
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidJSONError{Value: `{"name":"John","age":30,}`},
			},
		},
		{
			name: "invalid JSON - missing colon",
			args: args{
				ctx:   context.Background(),
				value: `{"name" "John"}`,
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidJSONError{Value: `{"name" "John"}`},
			},
		},
		{
			name: "invalid JSON - single quotes",
			args: args{
				ctx:   context.Background(),
				value: `{'name':'John'}`,
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidJSONError{Value: `{'name':'John'}`},
			},
		},
		{
			name: "invalid JSON - unquoted key",
			args: args{
				ctx:   context.Background(),
				value: `{name:"John"}`,
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidJSONError{Value: `{name:"John"}`},
			},
		},
		{
			name: "invalid JSON - empty string",
			args: args{
				ctx:   context.Background(),
				value: ``,
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidJSONError{Value: ``},
			},
		},
		{
			name: "invalid JSON - plain text",
			args: args{
				ctx:   context.Background(),
				value: `hello world`,
			},
			wantTwigBlock: false,
			wantErrs: []error{
				govalidator.InvalidJSONError{Value: `hello world`},
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
			gotTwigBlock, gotErrs := govalidator.JSONValidator(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantTwigBlock, gotTwigBlock, "JSONValidator(%v, %v)", tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantErrs, gotErrs, "JSONValidator(%v, %v)", tt.args.ctx, tt.args.value)
		})
	}
}
