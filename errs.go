package govalidator

import "fmt"

// NotAFloatError is returned when a value is not a float64 or int.
type NotAFloatError struct {
}

// NotANumberError is returned when a value is not a number (int, float, or numeric string).
type NotANumberError struct {
}

// FloatPrecisionError is returned when a float has more decimal places than allowed.
type FloatPrecisionError struct {
	ExpectedPrecision int
	ActualPrecision   int
}

// RequiredError is returned when a required field is missing or nil.
type RequiredError struct {
}

// MaxSizeError is returned when a collection is larger than the maximum size.
type MaxSizeError struct {
	MaxSize    int
	ActualSize int
}

// MinSizeError is returned when a collection is smaller than the minimum size.
type MinSizeError struct {
	MinSize    int
	ActualSize int
}

// NotAListError is returned when a value is not a list/array.
type NotAListError struct {
}

// UnexpectedFieldError is returned when a field is present but not allowed in the schema.
type UnexpectedFieldError struct {
	Field string
}

// FieldNotDefinedError is returned when a field is not defined in the schema.
type FieldNotDefinedError struct {
	Field string
}

// NotAValueError is returned when a value is not a valid value type.
type NotAValueError struct {
}

// NotAnObjectError is returned when a value is not an object (map[string]any).
type NotAnObjectError struct {
}

// NotAMapError is returned when a value is not a map.
type NotAMapError struct {
}

// NotABooleanError is returned when a value is not a boolean.
type NotABooleanError struct {
}

// NotAStringError is returned when a value is not a string.
type NotAStringError struct {
}

// NotAnIntegerError is returned when a value is not an integer.
type NotAnIntegerError struct {
}

func (e NotAFloatError) Error() string {
	return "not a float"
}

func (e NotANumberError) Error() string {
	return "not a number"
}

func (e FloatPrecisionError) Error() string {
	return fmt.Sprintf("expected precision %d, actual precision %d", e.ExpectedPrecision, e.ActualPrecision)
}

func (e NotAnIntegerError) Error() string {
	return "not an integer"
}

func (e NotAStringError) Error() string {
	return "not a string"
}

func (e NotABooleanError) Error() string {
	return "not a boolean"
}

func (e NotAMapError) Error() string {
	return "not a map"
}

func (e NotAnObjectError) Error() string {
	return "not an object"
}

func (e NotAValueError) Error() string {
	return "not a value"
}

func (e FieldNotDefinedError) Error() string {
	return fmt.Sprintf("field %s not defined", e.Field)
}

func (e UnexpectedFieldError) Error() string {
	return fmt.Sprintf("unexpected field %s", e.Field)
}

func (e NotAListError) Error() string {
	return "not a list"
}

func (e MinSizeError) Error() string {
	return fmt.Sprintf("min size %d, actual size %d", e.MinSize, e.ActualSize)
}

func (e MaxSizeError) Error() string {
	return fmt.Sprintf("max size %d, actual size %d", e.MaxSize, e.ActualSize)
}

func (e RequiredError) Error() string {
	return "required"
}
