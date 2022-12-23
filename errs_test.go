package go_validator_test

import (
"github.com/stretchr/testify/assert"
"testing"
"validator"
)

func TestErrors(t *testing.T) {
	t.Run("not a float error", func(t *testing.T) {
		err := main.NotAFloatError{}
		assert.EqualError(t, err, "not a float")
	})

	t.Run("not an int error", func(t *testing.T) {
		err := main.NotAnIntegerError{}
		assert.EqualError(t, err, "not an integer")
	})

	t.Run("not a string error", func(t *testing.T) {
		err := main.NotAStringError{}
		assert.EqualError(t, err, "not a string")
	})

	t.Run("not a boolean error", func(t *testing.T) {
		err := main.NotABooleanError{}
		assert.EqualError(t, err, "not a boolean")
	})

	t.Run("not a map error", func(t *testing.T) {
		err := main.NotAMapError{}
		assert.EqualError(t, err, "not a map")
	})

	t.Run("not an object error", func(t *testing.T) {
		err := main.NotAnObjectError{}
		assert.EqualError(t, err, "not an object")
	})

	t.Run("not a value error", func(t *testing.T) {
		err := main.NotAValueError{}
		assert.EqualError(t, err, "not a value")
	})

	t.Run("field not defined error", func(t *testing.T) {
		err := main.FieldNotDefinedError{Field: "foo"}
		assert.EqualError(t, err, "field foo not defined")
	})

	t.Run("required error", func(t *testing.T) {
		err := main.RequiredError{}
		assert.EqualError(t, err, "required")
	})

	t.Run("not a number error", func(t *testing.T) {
		err := main.NotANumberError{}
		assert.EqualError(t, err, "not a number")
	})

	t.Run("float precision error", func(t *testing.T) {
		err := main.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 3}
		assert.EqualError(t, err, "expected precision 2, actual precision 3")
	})

	t.Run("unexpected field error", func(t *testing.T) {
		err := main.UnexpectedFieldError{Field: "foo"}
		assert.EqualError(t, err, "unexpected field foo")
	})

	t.Run("not a list error", func(t *testing.T) {
		err := main.NotAListError{}
		assert.EqualError(t, err, "not a list")
	})

	t.Run("min size error", func(t *testing.T) {
		err := main.MinSizeError{MinSize: 2, ActualSize: 1}
		assert.EqualError(t, err, "min size 2, actual size 1")
	})

	t.Run("max size error", func(t *testing.T) {
		err := main.MaxSizeError{MaxSize: 2, ActualSize: 3}
		assert.EqualError(t, err, "max size 2, actual size 3")
	})
}
