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

func TestJSONPresenter_Coverage(t *testing.T) {
	t.Run("JSONPresenter with various errors", func(t *testing.T) {
		presenter := govalidator.JSONPresenter(".")

		errors := []error{
			govalidator.RequiredError{},
			govalidator.StringTooShortError{MinLength: 5},
			govalidator.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 4},
		}

		for _, err := range errors {
			result := presenter(context.Background(), []string{"$", "field"}, err)
			assert.Contains(t, result, "{")
			assert.Contains(t, result, "}")
		}
	})

	t.Run("JSONDetailedPresenter with errors", func(t *testing.T) {
		presenter := govalidator.JSONDetailedPresenter(".")

		errors := []error{
			govalidator.MinSizeError{MinSize: 1, ActualSize: 0},
			govalidator.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 4},
			govalidator.StringTooShortError{MinLength: 5},
			govalidator.FieldNotDefinedError{Field: "test"},
			govalidator.RequiredError{},
		}

		for _, err := range errors {
			result := presenter(context.Background(), []string{"$"}, err)
			assert.Contains(t, result, "{")
			assert.Contains(t, result, "}")
		}
	})
}

func TestJSONPresenter_AllErrors(t *testing.T) {
	presenter := govalidator.JSONPresenter(".")

	errors := []error{
		govalidator.RequiredError{},
		govalidator.NotAStringError{},
		govalidator.NotAnIntegerError{},
		govalidator.NotAFloatError{},
		govalidator.NotABooleanError{},
		govalidator.NotAMapError{},
		govalidator.NotAListError{},
		govalidator.StringTooShortError{MinLength: 5},
		govalidator.StringTooLongError{MaxLength: 10},
		govalidator.FloatTooSmallError{MinFloat: 0.0},
		govalidator.FloatTooLargeError{MaxFloat: 100.0},
		govalidator.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 4},
		govalidator.MinSizeError{MinSize: 1, ActualSize: 0},
		govalidator.MaxSizeError{MaxSize: 10, ActualSize: 15},
		govalidator.FieldNotDefinedError{Field: "test"},
		govalidator.UnexpectedFieldError{Field: "extra"},
		govalidator.InvalidOptionError{Options: []any{"a", "b"}},
		govalidator.ValueNotMatchingPatternError{Pattern: "^test$"},
		govalidator.NotLowerCasedError{},
		govalidator.NotUpperCasedError{},
		govalidator.NotANumberError{},
		govalidator.NotAValueError{},
		govalidator.NotAnObjectError{},
	}

	for _, err := range errors {
		result := presenter(context.Background(), []string{"$"}, err)
		assert.Contains(t, result, "{")
		assert.Contains(t, result, "}")
	}
}

func TestJSONDetailedPresenter_AllErrors(t *testing.T) {
	presenter := govalidator.JSONDetailedPresenter(".")

	errors := []error{
		govalidator.RequiredError{},
		govalidator.NotAStringError{},
		govalidator.NotAnIntegerError{},
		govalidator.NotAFloatError{},
		govalidator.NotABooleanError{},
		govalidator.NotAMapError{},
		govalidator.NotAListError{},
		govalidator.StringTooShortError{MinLength: 5},
		govalidator.StringTooLongError{MaxLength: 10},
		govalidator.FloatTooSmallError{MinFloat: 0.0},
		govalidator.FloatTooLargeError{MaxFloat: 100.0},
		govalidator.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 4},
		govalidator.MinSizeError{MinSize: 1, ActualSize: 0},
		govalidator.MaxSizeError{MaxSize: 10, ActualSize: 15},
		govalidator.FieldNotDefinedError{Field: "test"},
		govalidator.UnexpectedFieldError{Field: "extra"},
		govalidator.InvalidOptionError{Options: []any{"a", "b"}},
		govalidator.ValueNotMatchingPatternError{Pattern: "^test$"},
		govalidator.NotLowerCasedError{},
		govalidator.NotUpperCasedError{},
		govalidator.NotANumberError{},
		govalidator.NotAValueError{},
		govalidator.NotAnObjectError{},
	}

	for _, err := range errors {
		result := presenter(context.Background(), []string{"$"}, err)
		assert.Contains(t, result, "{")
		assert.Contains(t, result, "}")
	}
}

func TestJSONPresenter_MissingErrorTypes(t *testing.T) {
	t.Run("covers all error types in getErrorType", func(t *testing.T) {
		presenter := govalidator.JSONPresenter(".")

		// Test errors that might be missing from coverage
		errors := []error{
			govalidator.NotAValueError{},
			govalidator.NotAnObjectError{},
			govalidator.InvalidOptionError{Options: []any{"a", "b"}},
			govalidator.ValueNotMatchingPatternError{Pattern: "test"},
			govalidator.NotLowerCasedError{},
			govalidator.NotUpperCasedError{},
		}

		for _, err := range errors {
			result := presenter(context.Background(), []string{"$"}, err)
			assert.Contains(t, result, "{")
			assert.Contains(t, result, "message")
		}
	})

	t.Run("handles edge cases in path rendering", func(t *testing.T) {
		presenter := govalidator.JSONPresenter(".")

		// Test with various path structures
		testCases := []struct {
			path []string
			err  error
		}{
			{[]string{}, govalidator.RequiredError{}},
			{[]string{"$"}, govalidator.RequiredError{}},
			{[]string{"$", "field", "[0]", "nested"}, govalidator.RequiredError{}},
		}

		for _, tc := range testCases {
			result := presenter(context.Background(), tc.path, tc.err)
			assert.Contains(t, result, "{")
			assert.Contains(t, result, "path")
		}
	})
}
