package go_validator_test

import (
"github.com/stretchr/testify/assert"
"regexp"
"testing"
"validator"
)

func TestRegexpValidator(t *testing.T) {
	type args struct {
		pattern regexp.Regexp
	}
	tests := []struct {
		name          string
		args          args
		input         any
		wantErrs      []error
		wantTwigBrake bool
	}{
		{
			name: "match",
			args: args{
				pattern: *regexp.MustCompile("^[a-z]+$"),
			},
			input: "john",
		},
		{
			name: "not a string",
			args: args{
				pattern: *regexp.MustCompile("^[a-z]+$"),
			},
			input:         1,
			wantErrs:      []error{main.NotAStringError{}},
			wantTwigBrake: true,
		},
		{
			name: "not match",
			args: args{
				pattern: *regexp.MustCompile("^[a-z]+$"),
			},
			input: "John",
			wantErrs: []error{
				main.ValueNotMatchingPatternError{
					Pattern: "^[a-z]+$",
					Actual:  "John",
				},
			},
			wantTwigBrake: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := main.RegexpValidator(tt.args.pattern)

			gotTwigBrake, gotErrs := v(nil, tt.input)

			assert.Equal(t, tt.wantTwigBrake, gotTwigBrake)
			assert.Equal(t, tt.wantErrs, gotErrs)
		})
	}
}

func TestValueNotMatchingPatternError_Error(t *testing.T) {
	t.Run("error message", func(t *testing.T) {
		err := main.ValueNotMatchingPatternError{
			Pattern: "^[a-z]+$",
			Actual:  "John",
		}

		assert.Equal(t, "\"John\" does not match \"^[a-z]+$\"", err.Error())
	})
}
