package govalidator_test

import (
	"context"
	"errors"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCombinedPresenter(t *testing.T) {
	tests := []struct {
		name      string
		pathGlue  string
		separator string
		path      []string
		err       error
		want      string
	}{
		{
			name:      "simple path with error",
			pathGlue:  ".",
			separator: ": ",
			path:      []string{"$", "user", "age"},
			err:       errors.New("not an integer"),
			want:      "$.user.age: not an integer",
		},
		{
			name:      "root level error",
			pathGlue:  ".",
			separator: ": ",
			path:      []string{"$"},
			err:       errors.New("invalid data"),
			want:      "$: invalid data",
		},
		{
			name:      "empty path",
			pathGlue:  ".",
			separator: ": ",
			path:      []string{},
			err:       errors.New("error message"),
			want:      "error message",
		},
		{
			name:      "array path",
			pathGlue:  ".",
			separator: ": ",
			path:      []string{"$", "users", "[0]", "name"},
			err:       errors.New("required"),
			want:      "$.users[0].name: required",
		},
		{
			name:      "custom separator",
			pathGlue:  ".",
			separator: " -> ",
			path:      []string{"$", "field"},
			err:       errors.New("invalid"),
			want:      "$.field -> invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presenter := govalidator.CombinedPresenter(tt.pathGlue, tt.separator)
			result := presenter(context.Background(), tt.path, tt.err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestCombinedBracketPresenter(t *testing.T) {
	tests := []struct {
		name      string
		pathGlue  string
		separator string
		path      []string
		err       error
		want      string
	}{
		{
			name:      "simple path",
			pathGlue:  ".",
			separator: ": ",
			path:      []string{"$", "user", "age"},
			err:       errors.New("not an integer"),
			want:      "$.user.age: not an integer",
		},
		{
			name:      "array index",
			pathGlue:  ".",
			separator: ": ",
			path:      []string{"$", "users", "[0]", "name"},
			err:       errors.New("required"),
			want:      "$.users[0].name: required",
		},
		{
			name:      "nested arrays",
			pathGlue:  ".",
			separator: ": ",
			path:      []string{"$", "matrix", "[0]", "[1]"},
			err:       errors.New("invalid value"),
			want:      "$.matrix[0][1]: invalid value",
		},
		{
			name:      "empty path",
			pathGlue:  ".",
			separator: ": ",
			path:      []string{},
			err:       errors.New("error"),
			want:      "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presenter := govalidator.CombinedBracketPresenter(tt.pathGlue, tt.separator)
			result := presenter(context.Background(), tt.path, tt.err)
			assert.Equal(t, tt.want, result)
		})
	}
}
