package validator

import "fmt"

type NotAFloatError struct {
}

func (e NotAFloatError) Error() string {
	return "not a float"
}

type NotANumberError struct {
}

func (e NotANumberError) Error() string {
	return "not a number"
}

type FloatPrecisionError struct {
	ExpectedPrecision int
	ActualPrecision   int
}

func (e FloatPrecisionError) Error() string {
	return fmt.Sprintf("expected precision %d, actual precision %d", e.ExpectedPrecision, e.ActualPrecision)
}

type NotAnIntegerError struct {
}

func (e NotAnIntegerError) Error() string {
	return "not an integer"
}

type NotAStringError struct {
}

func (e NotAStringError) Error() string {
	return "not a string"
}

type NotABooleanError struct {
}

func (e NotABooleanError) Error() string {
	return "not a boolean"
}

type NotAMapError struct {
}

func (e NotAMapError) Error() string {
	return "not a map"
}

type NotAnObjectError struct {
}

func (e NotAnObjectError) Error() string {
	return "not an object"
}

type NotAValueError struct {
}

func (e NotAValueError) Error() string {
	return "not a value"
}

type FieldNotDefinedError struct {
	Field string
}

func (e FieldNotDefinedError) Error() string {
	return fmt.Sprintf("field %s not defined", e.Field)
}

type UnexpectedFieldError struct {
	Field string
}

func (e UnexpectedFieldError) Error() string {
	return fmt.Sprintf("unexpected field %s", e.Field)
}

type NotAListError struct {
}

func (e NotAListError) Error() string {
	return "not a list"
}

type MinSizeError struct {
	MinSize    int
	ActualSize int
}

func (e MinSizeError) Error() string {
	return fmt.Sprintf("min size %d, actual size %d", e.MinSize, e.ActualSize)
}

type MaxSizeError struct {
	MaxSize    int
	ActualSize int
}

func (e MaxSizeError) Error() string {
	return fmt.Sprintf("max size %d, actual size %d", e.MaxSize, e.ActualSize)
}

type RequiredError struct {
}

func (e RequiredError) Error() string {
	return "required"
}
