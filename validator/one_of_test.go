package validator_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"validator/validator"
)

func TestOneOfValidator(t *testing.T) {
	type args struct {
		options []any
	}
	tests := []struct {
		name         string
		args         args
		input        any
		expectedErrs []error
	}{
		{
			name: "happy path",
			args: args{
				options: []any{
					"john",
					"doe",
				},
			},
			input: "john",
		},
		{
			name: "mixed types - string input",
			args: args{
				options: []any{
					"john",
					1,
				},
			},
			input: "john",
		},
		{
			name: "mixed types - int input",
			args: args{
				options: []any{
					"john",
					1,
				},
			},
			input: 1,
		},
		{
			name: "not one of",
			args: args{
				options: []any{
					"john",
					"doe",
					"jane",
					"fonda",
				},
			},
			input: "jim",
			expectedErrs: []error{
				validator.InvalidOptionError{
					Options: []any{
						"john",
						"doe",
						"jane",
						"fonda",
					},
					Actual: "jim",
				},
			},
		},
		{
			name: "one of - object",
			args: args{
				options: []any{
					map[string]interface{}{
						"foo": "bar",
					},
					map[string]interface{}{
						"foo": "baz",
					},
				},
			},
			input: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			name: "not one of - object",
			args: args{
				options: []any{
					map[string]interface{}{
						"foo": "bar",
					},
					map[string]interface{}{
						"foo": "baz",
					},
				},
			},
			input: map[string]interface{}{
				"foo": "qux",
			},
			expectedErrs: []error{
				validator.InvalidOptionError{
					Options: []any{
						map[string]interface{}{
							"foo": "bar",
						},
						map[string]interface{}{
							"foo": "baz",
						},
					},
					Actual: map[string]interface{}{
						"foo": "qux",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.OneOfValidator(tt.args.options...)

			twigBlock, errs := v(context.Background(), tt.input)

			assert.False(t, twigBlock)
			assert.Equal(t, tt.expectedErrs, errs)
		})
	}
}

func TestInvalidOptionError_Error(t *testing.T) {
	type fields struct {
		Options []any
		Actual  any
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "strings",
			fields: fields{
				Options: []any{
					"john",
					"doe",
				},
				Actual: "jim",
			},
			want: "invalid option: jim, expected one of [john doe]",
		},
		{
			name: "mixed types",
			fields: fields{
				Options: []any{
					"john",
					1,
				},
				Actual: "jim",
			},
			want: "invalid option: jim, expected one of [john 1]",
		},
		{
			name: "objects",
			fields: fields{
				Options: []any{
					map[string]interface{}{
						"foo": "bar",
					},
					map[string]interface{}{
						"foo": "baz",
					},
				},
				Actual: map[string]interface{}{
					"foo": "qux",
				},
			},
			want: "invalid option: map[foo:qux], expected one of [map[foo:bar] map[foo:baz]]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validator.InvalidOptionError{
				Options: tt.fields.Options,
				Actual:  tt.fields.Actual,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
		})
	}
}
