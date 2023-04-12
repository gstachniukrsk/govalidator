package govalidator_test

import (
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDefinition_Extend(t *testing.T) {
	t.Run("should extend definition", func(t *testing.T) {
		a := govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.NonNullableValidator,
			},
		}
		b := govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.NullableValidator,
			},
		}

		c := a.ExtendedWith(b)

		assert.Equal(t, 2, len(c.Validator))

		assert.True(t, isReflectValuePointerEq(c.Validator[0], govalidator.NonNullableValidator))
		assert.True(t, isReflectValuePointerEq(c.Validator[1], govalidator.NullableValidator))
	})

	t.Run("should extend definition with fields", func(t *testing.T) {
		a := govalidator.Definition{
			Fields: &map[string]govalidator.Definition{
				"foo": {
					Fields: &map[string]govalidator.Definition{
						"bar": {
							Validator: []govalidator.ContextValidator{
								govalidator.NonNullableValidator,
							},
						},
					},
				},
			},
		}

		b := govalidator.Definition{
			Fields: &map[string]govalidator.Definition{
				"baz": {
					Validator: []govalidator.ContextValidator{
						govalidator.NonNullableValidator,
					},
				},
			},
		}

		c := a.ExtendedWith(b)

		assert.Equal(t, 2, len(*c.Fields))

		assert.True(t, isReflectValuePointerEq((*c.Fields)["foo"].Validator, (*a.Fields)["foo"].Validator))
		assert.True(t, isReflectValuePointerEq((*c.Fields)["baz"].Validator, (*b.Fields)["baz"].Validator))

		assert.True(t, isReflectValuePointerEq((*(*c.Fields)["foo"].Fields)["bar"].Validator, (*(*a.Fields)["foo"].Fields)["bar"].Validator))
	})

	t.Run("should override booleans", func(t *testing.T) {
		a := govalidator.Definition{
			AcceptNotDefinedProperty: true,
			AcceptExtraProperty:      true,
		}

		b := a.ExtendedWith(govalidator.Definition{})

		assert.True(t, b.AcceptNotDefinedProperty)
		assert.True(t, b.AcceptExtraProperty)
	})

	t.Run("should override booleans, reverse", func(t *testing.T) {
		a := govalidator.Definition{}

		b := a.ExtendedWith(govalidator.Definition{
			AcceptNotDefinedProperty: true,
			AcceptExtraProperty:      true,
		})

		assert.True(t, b.AcceptNotDefinedProperty)
		assert.True(t, b.AcceptExtraProperty)
	})

	t.Run("ListOf left", func(t *testing.T) {
		lo := &govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.NonNullableValidator,
			},
		}
		a := govalidator.Definition{
			ListOf: lo,
		}

		b := govalidator.Definition{}

		c := a.ExtendedWith(b)

		assert.True(t, isReflectValuePointerEq(c.ListOf, lo))
	})

	t.Run("ListOf right", func(t *testing.T) {
		lo := &govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.NonNullableValidator,
			},
		}
		a := govalidator.Definition{}

		b := govalidator.Definition{
			ListOf: lo,
		}

		c := a.ExtendedWith(b)

		assert.True(t, isReflectValuePointerEq(c.ListOf, lo))
	})

	t.Run("ListOf Conflict", func(t *testing.T) {
		loL := &govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.NonNullableValidator,
			},
		}

		loR := &govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.NullableValidator,
			},
		}
		a := govalidator.Definition{
			ListOf: loL,
		}

		b := govalidator.Definition{
			ListOf: loR,
		}

		c := a.ExtendedWith(b)

		assert.True(t, isReflectValuePointerEq(c.ListOf.Validator[0], loL.Validator[0]))
		assert.True(t, isReflectValuePointerEq(c.ListOf.Validator[1], loR.Validator[0]))
	})
}

func isReflectValuePointerEq(a, b interface{}) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}
