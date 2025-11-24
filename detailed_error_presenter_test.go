package govalidator_test

import (
	"context"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerboseErrorPresenter(t *testing.T) {
	t.Run("VerboseErrorPresenter formats errors", func(t *testing.T) {
		presenter := govalidator.VerboseErrorPresenter()

		result := presenter(context.Background(), []string{"$", "name"}, govalidator.RequiredError{})
		assert.NotEmpty(t, result)
	})

	t.Run("VerboseErrorPresenter with different errors", func(t *testing.T) {
		presenter := govalidator.VerboseErrorPresenter()

		errors := []error{
			govalidator.RequiredError{},
			govalidator.NotAStringError{},
			govalidator.NotAnIntegerError{},
			govalidator.StringTooShortError{MinLength: 5},
			govalidator.StringTooLongError{MaxLength: 10},
		}

		for _, err := range errors {
			result := presenter(context.Background(), []string{"$"}, err)
			assert.NotEmpty(t, result)
		}
	})
}

func TestVerboseErrorPresenter_AllCases(t *testing.T) {
	presenter := govalidator.VerboseErrorPresenter()

	testCases := []struct {
		name  string
		error error
	}{
		{"RequiredError", govalidator.RequiredError{}},
		{"NotAStringError", govalidator.NotAStringError{}},
		{"NotAnIntegerError", govalidator.NotAnIntegerError{}},
		{"NotABooleanError", govalidator.NotABooleanError{}},
		{"NotAFloatError", govalidator.NotAFloatError{}},
		{"NotAMapError", govalidator.NotAMapError{}},
		{"NotAnObjectError", govalidator.NotAnObjectError{}},
		{"NotAListError", govalidator.NotAListError{}},
		{"NotANumberError", govalidator.NotANumberError{}},
		{"NotAValueError", govalidator.NotAValueError{}},
		{"StringTooShortError", govalidator.StringTooShortError{MinLength: 5}},
		{"StringTooLongError", govalidator.StringTooLongError{MaxLength: 10}},
		{"FloatTooSmallError", govalidator.FloatTooSmallError{MinFloat: 0.0}},
		{"FloatTooLargeError", govalidator.FloatTooLargeError{MaxFloat: 100.0}},
		{"FloatPrecisionError", govalidator.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 4}},
		{"MinSizeError", govalidator.MinSizeError{MinSize: 1, ActualSize: 0}},
		{"MaxSizeError", govalidator.MaxSizeError{MaxSize: 10, ActualSize: 15}},
		{"FieldNotDefinedError", govalidator.FieldNotDefinedError{Field: "test"}},
		{"UnexpectedFieldError", govalidator.UnexpectedFieldError{Field: "extra"}},
		{"InvalidOptionError", govalidator.InvalidOptionError{Options: []any{"a", "b", "c"}}},
		{"ValueNotMatchingPatternError", govalidator.ValueNotMatchingPatternError{Pattern: "^test$"}},
		{"NotLowerCasedError", govalidator.NotLowerCasedError{}},
		{"NotUpperCasedError", govalidator.NotUpperCasedError{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := presenter(context.Background(), []string{"$"}, tc.error)
			assert.NotEmpty(t, result)
		})
	}
}

func TestDetailedErrorPresenter_Coverage(t *testing.T) {
	t.Run("DetailedErrorPresenter with all error types", func(t *testing.T) {
		presenter := govalidator.DetailedErrorPresenter()

		errorTypes := []error{
			govalidator.RequiredError{},
			govalidator.NotAStringError{},
			govalidator.NotAnIntegerError{},
			govalidator.NotABooleanError{},
			govalidator.NotAFloatError{},
			govalidator.NotAMapError{},
			govalidator.NotAnObjectError{},
			govalidator.NotAListError{},
			govalidator.NotANumberError{},
			govalidator.NotAValueError{},
			govalidator.StringTooShortError{MinLength: 5},
			govalidator.StringTooLongError{MaxLength: 100},
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
		}

		for _, err := range errorTypes {
			result := presenter(context.Background(), []string{"$"}, err)
			assert.NotEmpty(t, result)
		}
	})
}
