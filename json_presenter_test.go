package govalidator_test

import (
	"context"
	"encoding/json"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONPresenter(t *testing.T) {
	tests := []struct {
		name     string
		pathGlue string
		path     []string
		err      error
		want     map[string]string
	}{
		{
			name:     "simple error",
			pathGlue: ".",
			path:     []string{"$", "user", "age"},
			err:      govalidator.NotAnIntegerError{},
			want: map[string]string{
				"path":    "$.user.age",
				"message": "not an integer",
			},
		},
		{
			name:     "root error",
			pathGlue: ".",
			path:     []string{"$"},
			err:      govalidator.RequiredError{},
			want: map[string]string{
				"path":    "$",
				"message": "required",
			},
		},
		{
			name:     "array path",
			pathGlue: ".",
			path:     []string{"$", "users", "[0]", "name"},
			err:      govalidator.NotAStringError{},
			want: map[string]string{
				"path":    "$.users[0].name",
				"message": "not a string",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presenter := govalidator.JSONPresenter(tt.pathGlue)
			result := presenter(context.Background(), tt.path, tt.err)

			var actual map[string]string
			err := json.Unmarshal([]byte(result), &actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestJSONDetailedPresenter(t *testing.T) {
	tests := []struct {
		name     string
		pathGlue string
		path     []string
		err      error
		validate func(t *testing.T, result string)
	}{
		{
			name:     "MinSizeError with details",
			pathGlue: ".",
			path:     []string{"$", "items"},
			err: govalidator.MinSizeError{
				MinSize:    5,
				ActualSize: 2,
			},
			validate: func(t *testing.T, result string) {
				var actual map[string]any
				err := json.Unmarshal([]byte(result), &actual)
				assert.NoError(t, err)
				assert.Equal(t, "$.items", actual["path"])
				assert.Equal(t, "MinSizeError", actual["type"])
				assert.Equal(t, float64(5), actual["minSize"])
				assert.Equal(t, float64(2), actual["actualSize"])
			},
		},
		{
			name:     "FloatPrecisionError with details",
			pathGlue: ".",
			path:     []string{"$", "price"},
			err: govalidator.FloatPrecisionError{
				ExpectedPrecision: 2,
				ActualPrecision:   4,
			},
			validate: func(t *testing.T, result string) {
				var actual map[string]any
				err := json.Unmarshal([]byte(result), &actual)
				assert.NoError(t, err)
				assert.Equal(t, "$.price", actual["path"])
				assert.Equal(t, "FloatPrecisionError", actual["type"])
				assert.Equal(t, float64(2), actual["expectedPrecision"])
				assert.Equal(t, float64(4), actual["actualPrecision"])
			},
		},
		{
			name:     "StringTooShortError with details",
			pathGlue: ".",
			path:     []string{"$", "username"},
			err: govalidator.StringTooShortError{
				MinLength: 5,
			},
			validate: func(t *testing.T, result string) {
				var actual map[string]any
				err := json.Unmarshal([]byte(result), &actual)
				assert.NoError(t, err)
				assert.Equal(t, "$.username", actual["path"])
				assert.Equal(t, "StringTooShortError", actual["type"])
				assert.Equal(t, float64(5), actual["minLength"])
			},
		},
		{
			name:     "FieldNotDefinedError with field name",
			pathGlue: ".",
			path:     []string{"$", "user"},
			err: govalidator.FieldNotDefinedError{
				Field: "email",
			},
			validate: func(t *testing.T, result string) {
				var actual map[string]any
				err := json.Unmarshal([]byte(result), &actual)
				assert.NoError(t, err)
				assert.Equal(t, "$.user", actual["path"])
				assert.Equal(t, "FieldNotDefinedError", actual["type"])
				assert.Equal(t, "email", actual["field"])
			},
		},
		{
			name:     "simple error without extra fields",
			pathGlue: ".",
			path:     []string{"$", "value"},
			err:      govalidator.NotAStringError{},
			validate: func(t *testing.T, result string) {
				var actual map[string]any
				err := json.Unmarshal([]byte(result), &actual)
				assert.NoError(t, err)
				assert.Equal(t, "$.value", actual["path"])
				assert.Equal(t, "NotAStringError", actual["type"])
				assert.Equal(t, "not a string", actual["message"])
				// Should not have extra fields
				assert.NotContains(t, actual, "minSize")
				assert.NotContains(t, actual, "field")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presenter := govalidator.JSONDetailedPresenter(tt.pathGlue)
			result := presenter(context.Background(), tt.path, tt.err)
			tt.validate(t, result)
		})
	}
}
