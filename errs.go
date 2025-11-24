package govalidator

import "fmt"

// NotAFloatError is returned when a value is not a float64 or int.
type NotAFloatError struct {
}

func (e NotAFloatError) Error() string {
	return "not a float"
}

// NotANumberError is returned when a value is not a number (int, float, or numeric string).
type NotANumberError struct {
}

func (e NotANumberError) Error() string {
	return "not a number"
}

// FloatPrecisionError is returned when a float has more decimal places than allowed.
type FloatPrecisionError struct {
	ExpectedPrecision int
	ActualPrecision   int
}

func (e FloatPrecisionError) Error() string {
	return fmt.Sprintf("expected precision %d, actual precision %d", e.ExpectedPrecision, e.ActualPrecision)
}

// NotAnIntegerError is returned when a value is not an integer.
type NotAnIntegerError struct {
}

func (e NotAnIntegerError) Error() string {
	return "not an integer"
}

// NotAStringError is returned when a value is not a string.
type NotAStringError struct {
}

func (e NotAStringError) Error() string {
	return "not a string"
}

// NotABooleanError is returned when a value is not a boolean.
type NotABooleanError struct {
}

func (e NotABooleanError) Error() string {
	return "not a boolean"
}

// NotAMapError is returned when a value is not a map.
type NotAMapError struct {
}

func (e NotAMapError) Error() string {
	return "not a map"
}

// NotAnObjectError is returned when a value is not an object (map[string]any).
type NotAnObjectError struct {
}

func (e NotAnObjectError) Error() string {
	return "not an object"
}

// NotAValueError is returned when a value is not a valid value type.
type NotAValueError struct {
}

func (e NotAValueError) Error() string {
	return "not a value"
}

// FieldNotDefinedError is returned when a field is not defined in the schema.
type FieldNotDefinedError struct {
	Field string
}

func (e FieldNotDefinedError) Error() string {
	return fmt.Sprintf("field %s not defined", e.Field)
}

// UnexpectedFieldError is returned when a field is present but not allowed in the schema.
type UnexpectedFieldError struct {
	Field string
}

func (e UnexpectedFieldError) Error() string {
	return fmt.Sprintf("unexpected field %s", e.Field)
}

// NotAListError is returned when a value is not a list/array.
type NotAListError struct {
}

func (e NotAListError) Error() string {
	return "not a list"
}

// MinSizeError is returned when a collection is smaller than the minimum size.
type MinSizeError struct {
	MinSize    int
	ActualSize int
}

func (e MinSizeError) Error() string {
	return fmt.Sprintf("min size %d, actual size %d", e.MinSize, e.ActualSize)
}

// MaxSizeError is returned when a collection is larger than the maximum size.
type MaxSizeError struct {
	MaxSize    int
	ActualSize int
}

func (e MaxSizeError) Error() string {
	return fmt.Sprintf("max size %d, actual size %d", e.MaxSize, e.ActualSize)
}

// RequiredError is returned when a required field is missing or nil.
type RequiredError struct {
}

func (e RequiredError) Error() string {
	return "required"
}
