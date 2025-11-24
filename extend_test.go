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
