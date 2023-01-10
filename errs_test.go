package go_validator_test

import (
	"github.com/gstachniukrsk/go_validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	t.Run("not a float error", func(t *testing.T) {
		err := go_validator.NotAFloatError{}
		assert.EqualError(t, err, "not a float")
	})

	t.Run("not an int error", func(t *testing.T) {
		err := go_validator.NotAnIntegerError{}
		assert.EqualError(t, err, "not an integer")
	})

	t.Run("not a string error", func(t *testing.T) {
		err := go_validator.NotAStringError{}
		assert.EqualError(t, err, "not a string")
	})

	t.Run("not a boolean error", func(t *testing.T) {
		err := go_validator.NotABooleanError{}
		assert.EqualError(t, err, "not a boolean")
	})

	t.Run("not a map error", func(t *testing.T) {
		err := go_validator.NotAMapError{}
		assert.EqualError(t, err, "not a map")
	})

	t.Run("not an object error", func(t *testing.T) {
		err := go_validator.NotAnObjectError{}
		assert.EqualError(t, err, "not an object")
	})

	t.Run("not a value error", func(t *testing.T) {
		err := go_validator.NotAValueError{}
		assert.EqualError(t, err, "not a value")
	})

	t.Run("field not defined error", func(t *testing.T) {
		err := go_validator.FieldNotDefinedError{Field: "foo"}
		assert.EqualError(t, err, "field foo not defined")
	})

	t.Run("required error", func(t *testing.T) {
		err := go_validator.RequiredError{}
		assert.EqualError(t, err, "required")
	})

	t.Run("not a number error", func(t *testing.T) {
		err := go_validator.NotANumberError{}
		assert.EqualError(t, err, "not a number")
	})

	t.Run("float precision error", func(t *testing.T) {
		err := go_validator.FloatPrecisionError{ExpectedPrecision: 2, ActualPrecision: 3}
		assert.EqualError(t, err, "expected precision 2, actual precision 3")
	})

	t.Run("unexpected field error", func(t *testing.T) {
		err := go_validator.UnexpectedFieldError{Field: "foo"}
		assert.EqualError(t, err, "unexpected field foo")
	})

	t.Run("not a list error", func(t *testing.T) {
		err := go_validator.NotAListError{}
		assert.EqualError(t, err, "not a list")
	})

	t.Run("min size error", func(t *testing.T) {
		err := go_validator.MinSizeError{MinSize: 2, ActualSize: 1}
		assert.EqualError(t, err, "min size 2, actual size 1")
	})

	t.Run("max size error", func(t *testing.T) {
		err := go_validator.MaxSizeError{MaxSize: 2, ActualSize: 3}
		assert.EqualError(t, err, "max size 2, actual size 3")
	})
}
