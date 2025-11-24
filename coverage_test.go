package govalidator_test

import (
	"context"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Tests to achieve 100% code coverage

func TestFlatListCollector(t *testing.T) {
	t.Run("NewFlatListCollector creates collector", func(t *testing.T) {
		collector := govalidator.NewFlatListCollector(govalidator.CombinedPresenter(".", ": "))
		assert.NotNil(t, collector)
	})

	t.Run("Collect adds errors", func(t *testing.T) {
		collector := govalidator.NewFlatListCollector(govalidator.CombinedPresenter(".", ": "))
		collector.Collect(context.Background(), []string{"$", "name"}, govalidator.RequiredError{})

		errs := collector.GetErrors()
		assert.Len(t, errs, 1)
	})

	t.Run("GetErrors returns flat list", func(t *testing.T) {
		collector := govalidator.NewFlatListCollector(govalidator.CombinedPresenter(".", ": "))
		collector.Collect(context.Background(), []string{"$", "name"}, govalidator.RequiredError{})
		collector.Collect(context.Background(), []string{"$", "age"}, govalidator.NotAnIntegerError{})

		errs := collector.GetErrors()
		assert.Len(t, errs, 2)
	})
}

func TestFlatListValidator(t *testing.T) {
	t.Run("NewFlatListValidator creates validator", func(t *testing.T) {
		validator := govalidator.NewFlatListValidator(govalidator.CombinedPresenter(".", ": "))
		assert.NotNil(t, validator)
	})

	t.Run("Validate returns flat errors", func(t *testing.T) {
		def := govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.NonNullableValidator,
				govalidator.IsStringValidator,
			},
		}

		validator := govalidator.NewFlatListValidator(govalidator.CombinedPresenter(".", ": "))
		valid, errs := validator.Validate(context.Background(), 123, def)

		assert.False(t, valid)
		assert.NotEmpty(t, errs)
		assert.IsType(t, []string{}, errs)
	})

	t.Run("Validate with complex schema", func(t *testing.T) {
		def := govalidator.Definition{
			Validator: []govalidator.ContextValidator{},
			Fields: &map[string]govalidator.Definition{
				"name": {
					Validator: []govalidator.ContextValidator{
						govalidator.NonNullableValidator,
						govalidator.IsStringValidator,
					},
				},
				"age": {
					Validator: []govalidator.ContextValidator{
						govalidator.NonNullableValidator,
						govalidator.IsIntegerValidator,
					},
				},
			},
		}

		validator := govalidator.NewFlatListValidator(govalidator.CombinedPresenter(".", ": "))
		data := map[string]any{
			"name": 123,      // wrong type
			"age":  "thirty", // wrong type
		}

		valid, errs := validator.Validate(context.Background(), data, def)

		assert.False(t, valid)
		assert.Len(t, errs, 2)
	})
}

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

func TestMaxFloatError(t *testing.T) {
	t.Run("MaxFloatError Error method", func(t *testing.T) {
		err := govalidator.FloatTooLargeError{MaxFloat: 100.0}
		assert.Equal(t, "value is greater than max", err.Error())
		assert.Equal(t, 100.0, err.MaxFloat)
	})
}

func TestSchemaValidator_ArrayEdgeCases(t *testing.T) {
	t.Run("validates nil array", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema(govalidator.IsStringValidator).Required(),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), nil, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$")
	})
}

func TestFlatErrorCollector_GetErrors(t *testing.T) {
	t.Run("GetErrors returns map format", func(t *testing.T) {
		collector := govalidator.NewFlatErrorCollector(context.Background(), govalidator.CombinedPresenter(".", ": "))
		collector.Collect([]string{"$", "name"}, govalidator.RequiredError{})

		errs := collector.GetErrors()
		assert.NotNil(t, errs)
		assert.Contains(t, errs, "errors")
	})

	t.Run("GetErrors returns empty map when no errors", func(t *testing.T) {
		collector := govalidator.NewFlatErrorCollector(context.Background(), govalidator.CombinedPresenter(".", ": "))

		errs := collector.GetErrors()
		assert.Empty(t, errs)
	})
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

func TestMaxFloatValidator_Coverage(t *testing.T) {
	t.Run("MaxFloatValidator with edge cases", func(t *testing.T) {
		validator := govalidator.MaxFloatValidator(100.0)

		// Test with float below max
		shouldBlock, errs := validator(context.Background(), 50.0)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with float above max - should block
		shouldBlock, errs = validator(context.Background(), 150.0)
		assert.True(t, shouldBlock)
		assert.NotEmpty(t, errs)

		// Test with int below max
		shouldBlock, errs = validator(context.Background(), 50)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with wrong type - should block
		shouldBlock, errs = validator(context.Background(), "not a number")
		assert.True(t, shouldBlock)
		assert.NotEmpty(t, errs)
	})
}

func TestFloatValidator_Coverage(t *testing.T) {
	t.Run("FloatValidator with high precision", func(t *testing.T) {
		validator := govalidator.FloatValidator(10)

		// Test within precision
		shouldBlock, errs := validator(context.Background(), 1.123456789)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test exceeding precision
		shouldBlock, errs = validator(context.Background(), 1.12345678901)
		assert.False(t, shouldBlock)
		assert.NotEmpty(t, errs)
	})
}

func TestExtendedWith_Coverage(t *testing.T) {
	t.Run("ExtendedWith merges fields", func(t *testing.T) {
		base := govalidator.Definition{
			Fields: &map[string]govalidator.Definition{
				"field1": {
					Validator: []govalidator.ContextValidator{
						govalidator.IsStringValidator,
					},
				},
			},
		}

		extension := govalidator.Definition{
			Fields: &map[string]govalidator.Definition{
				"field2": {
					Validator: []govalidator.ContextValidator{
						govalidator.IsIntegerValidator,
					},
				},
			},
		}

		result := base.ExtendedWith(extension)
		assert.Len(t, *result.Fields, 2)
		assert.Contains(t, *result.Fields, "field1")
		assert.Contains(t, *result.Fields, "field2")
	})

	t.Run("ExtendedWith merges conflicting ListOf recursively", func(t *testing.T) {
		base := govalidator.Definition{
			ListOf: &govalidator.Definition{
				Validator: []govalidator.ContextValidator{
					govalidator.IsStringValidator,
				},
			},
		}

		extension := govalidator.Definition{
			ListOf: &govalidator.Definition{
				Validator: []govalidator.ContextValidator{
					govalidator.MinLengthValidator(3),
				},
			},
		}

		// Should merge ListOf recursively
		result := base.ExtendedWith(extension)
		assert.NotNil(t, result.ListOf)
		// Should have both validators merged
		assert.Len(t, result.ListOf.Validator, 2)
	})
}

func TestSchemaValidator_ValidateArray_WrongType(t *testing.T) {
	t.Run("rejects non-array value", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema(govalidator.IsStringValidator).Required(),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		// Pass a string instead of an array
		valid, errs := validator.Validate(context.Background(), "not an array", schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$")
	})

	t.Run("rejects nil typed as []any", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema(govalidator.IsStringValidator).Required(),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		// Pass nil with type []any to trigger the nil check in validateArray
		var nilArray []any = nil
		valid, errs := validator.Validate(context.Background(), nilArray, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$")
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

func TestFloatValidator_AllPrecisionCases(t *testing.T) {
	t.Run("handles various numeric types", func(t *testing.T) {
		validator := govalidator.FloatValidator(2)

		// Test with float64
		shouldBlock, errs := validator(context.Background(), 1.23)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with int
		shouldBlock, errs = validator(context.Background(), 123)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with wrong type (string)
		shouldBlock, errs = validator(context.Background(), "not a float")
		assert.True(t, shouldBlock)
		assert.NotEmpty(t, errs)
		assert.IsType(t, govalidator.NotAFloatError{}, errs[0])

		// Test with wrong type (float32 - not accepted)
		var f32 float32 = 1.23
		shouldBlock, errs = validator(context.Background(), f32)
		assert.True(t, shouldBlock)
		assert.NotEmpty(t, errs)
		assert.IsType(t, govalidator.NotAFloatError{}, errs[0])
	})

	t.Run("handles trailing zeros in decimal", func(t *testing.T) {
		validator := govalidator.FloatValidator(2)

		// Test with trailing zeros (should be stripped and pass)
		shouldBlock, errs := validator(context.Background(), 1.230000)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)

		// Test with float that has zeros in middle (1.2003 -> 4 decimals)
		shouldBlock, errs = validator(context.Background(), 1.2003)
		assert.False(t, shouldBlock)
		assert.NotEmpty(t, errs) // Should fail - 4 > 2
	})

	t.Run("handles floats without decimal part", func(t *testing.T) {
		validator := govalidator.FloatValidator(2)

		// Test with whole number as float (5.0)
		shouldBlock, errs := validator(context.Background(), 5.0)
		assert.False(t, shouldBlock)
		assert.Empty(t, errs)
	})

	t.Run("handles precision edge cases with strings", func(t *testing.T) {
		validator := govalidator.FloatValidator(3)

		// Test floats with different numbers of trailing zeros
		testCases := []struct {
			value float64
			valid bool
		}{
			{1.100, true},   // exactly 1 significant decimal (after removing trailing zeros)
			{1.120, true},   // exactly 2 significant decimals
			{1.123, true},   // exactly 3 significant decimals
			{1.1234, false}, // 4 significant decimals - too many
			{10.00, true},   // all trailing zeros
			{10.01, true},   // 1 significant decimal (after removing trailing zero)
			{10.010, true},  // 2 significant decimals (after removing trailing zero)
		}

		for _, tc := range testCases {
			shouldBlock, errs := validator(context.Background(), tc.value)
			if tc.valid {
				assert.False(t, shouldBlock)
				assert.Empty(t, errs)
			} else {
				assert.False(t, shouldBlock)
				assert.NotEmpty(t, errs)
			}
		}
	})
}

func TestExtendedWith_AllCases(t *testing.T) {
	t.Run("extends validators", func(t *testing.T) {
		base := govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.IsStringValidator,
			},
		}

		extension := govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.MinLengthValidator(5),
			},
		}

		result := base.ExtendedWith(extension)
		assert.Len(t, result.Validator, 2)
	})

	t.Run("extends with boolean flags", func(t *testing.T) {
		base := govalidator.Definition{
			AcceptExtraProperty:      false,
			AcceptNotDefinedProperty: false,
		}

		extension := govalidator.Definition{
			AcceptExtraProperty:      true,
			AcceptNotDefinedProperty: true,
		}

		result := base.ExtendedWith(extension)
		assert.True(t, result.AcceptExtraProperty)
		assert.True(t, result.AcceptNotDefinedProperty)
	})

	t.Run("extends fields recursively", func(t *testing.T) {
		base := govalidator.Definition{
			Fields: &map[string]govalidator.Definition{
				"name": {
					Validator: []govalidator.ContextValidator{
						govalidator.IsStringValidator,
					},
				},
			},
		}

		extension := govalidator.Definition{
			Fields: &map[string]govalidator.Definition{
				"name": {
					Validator: []govalidator.ContextValidator{
						govalidator.MinLengthValidator(3),
					},
				},
			},
		}

		result := base.ExtendedWith(extension)
		assert.NotNil(t, result.Fields)
		assert.Contains(t, *result.Fields, "name")
		// Should have merged validators
		assert.Len(t, (*result.Fields)["name"].Validator, 2)
	})

	t.Run("handles nil base Fields", func(t *testing.T) {
		base := govalidator.Definition{
			Fields: nil,
		}

		extension := govalidator.Definition{
			Fields: &map[string]govalidator.Definition{
				"age": {
					Validator: []govalidator.ContextValidator{
						govalidator.IsIntegerValidator,
					},
				},
			},
		}

		result := base.ExtendedWith(extension)
		assert.NotNil(t, result.Fields)
		assert.Contains(t, *result.Fields, "age")
	})

	t.Run("handles nil extension Fields", func(t *testing.T) {
		base := govalidator.Definition{
			Fields: &map[string]govalidator.Definition{
				"name": {
					Validator: []govalidator.ContextValidator{
						govalidator.IsStringValidator,
					},
				},
			},
		}

		extension := govalidator.Definition{
			Fields: nil,
		}

		result := base.ExtendedWith(extension)
		assert.NotNil(t, result.Fields)
		assert.Contains(t, *result.Fields, "name")
	})

	t.Run("handles only extension ListOf", func(t *testing.T) {
		base := govalidator.Definition{
			ListOf: nil,
		}

		extension := govalidator.Definition{
			ListOf: &govalidator.Definition{
				Validator: []govalidator.ContextValidator{
					govalidator.IsStringValidator,
				},
			},
		}

		result := base.ExtendedWith(extension)
		assert.NotNil(t, result.ListOf)
		assert.Equal(t, extension.ListOf, result.ListOf)
	})

	t.Run("handles only base ListOf", func(t *testing.T) {
		base := govalidator.Definition{
			ListOf: &govalidator.Definition{
				Validator: []govalidator.ContextValidator{
					govalidator.IsStringValidator,
				},
			},
		}

		extension := govalidator.Definition{
			ListOf: nil,
		}

		result := base.ExtendedWith(extension)
		assert.NotNil(t, result.ListOf)
		assert.Equal(t, base.ListOf, result.ListOf)
	})
}
