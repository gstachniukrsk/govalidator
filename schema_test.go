package govalidator_test

import (
	"context"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSchema(t *testing.T) {
	t.Run("creates schema with validators", func(t *testing.T) {
		schema := govalidator.NewSchema(
			govalidator.IsStringValidator,
			govalidator.MinLengthValidator(5),
		)
		assert.NotNil(t, schema)
	})

	t.Run("creates schema without validators", func(t *testing.T) {
		schema := govalidator.NewSchema()
		assert.NotNil(t, schema)
	})
}

func TestSchema_Required(t *testing.T) {
	t.Run("required schema rejects null", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator).Required()

		valid, errs := schema.Validate(context.Background(), nil)

		assert.False(t, valid)
		assert.NotEmpty(t, errs)
	})

	t.Run("required schema accepts valid value", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator).Required()

		valid, errs := schema.Validate(context.Background(), "hello")

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("IsRequired returns true", func(t *testing.T) {
		schema := govalidator.NewSchema().Required()
		assert.True(t, schema.IsRequired())
	})
}

func TestSchema_Optional(t *testing.T) {
	t.Run("optional schema accepts null", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator).Optional()

		valid, errs := schema.Validate(context.Background(), nil)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("optional is default", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator)

		valid, errs := schema.Validate(context.Background(), nil)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("IsRequired returns false", func(t *testing.T) {
		schema := govalidator.NewSchema().Optional()
		assert.False(t, schema.IsRequired())
	})
}

func TestSchema_Object(t *testing.T) {
	t.Run("Object helper creates object schema", func(t *testing.T) {
		schema := govalidator.Object(map[string]*govalidator.Schema{
			"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
			"age":  govalidator.NewSchema(govalidator.IsIntegerValidator).Required(),
		})

		assert.NotNil(t, schema.Fields)
	})
}

func TestSchema_Array(t *testing.T) {
	t.Run("Array helper creates array schema", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema(govalidator.IsStringValidator).Required(),
		)

		assert.NotNil(t, schema.Items)
		assert.Len(t, schema.Validators, 1)
	})
}

func TestSchema_StructLiteralSyntax(t *testing.T) {
	t.Run("can create schema with struct literal", func(t *testing.T) {
		schema := &govalidator.Schema{
			Validators: []govalidator.ContextValidator{
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(5),
			},
		}

		// Mark as required
		schema.Required()

		valid, errs := schema.Validate(context.Background(), "hello")

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("can create nested schema with struct literals", func(t *testing.T) {
		nameSchema := &govalidator.Schema{
			Validators: []govalidator.ContextValidator{govalidator.IsStringValidator},
		}
		nameSchema.Required()

		ageSchema := &govalidator.Schema{
			Validators: []govalidator.ContextValidator{govalidator.IsIntegerValidator},
		}
		ageSchema.Required()

		schema := &govalidator.Schema{
			Fields: map[string]*govalidator.Schema{
				"name": nameSchema,
				"age":  ageSchema,
			},
		}

		data := map[string]any{
			"name": "John",
			"age":  30,
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})
}

func TestSchema_WithFields(t *testing.T) {
	t.Run("validates object with required fields", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(map[string]*govalidator.Schema{
			"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
			"age":  govalidator.NewSchema(govalidator.IsIntegerValidator).Required(),
		})

		data := map[string]any{
			"name": "John",
			"age":  30,
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("fails when required field is missing", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(map[string]*govalidator.Schema{
			"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
			"age":  govalidator.NewSchema(govalidator.IsIntegerValidator).Required(),
		})

		data := map[string]any{
			"name": "John",
			// age is missing
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.False(t, valid)
		assert.NotEmpty(t, errs)
	})

	t.Run("allows optional fields to be missing", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(map[string]*govalidator.Schema{
			"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
			"age":  govalidator.NewSchema(govalidator.IsIntegerValidator).Optional(),
		})

		data := map[string]any{
			"name": "John",
			// age is optional and missing
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})
}

func TestSchema_WithExtra(t *testing.T) {
	t.Run("ExtraIgnore permits extra fields", func(t *testing.T) {
		schema := govalidator.NewSchema().
			WithFields(map[string]*govalidator.Schema{
				"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
			}).
			WithExtra(govalidator.ExtraIgnore)

		data := map[string]any{
			"name":  "John",
			"extra": "field",
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("ExtraForbid rejects extra fields", func(t *testing.T) {
		schema := govalidator.NewSchema().
			WithFields(map[string]*govalidator.Schema{
				"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
			}).
			WithExtra(govalidator.ExtraForbid)

		data := map[string]any{
			"name":  "John",
			"extra": "field",
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.False(t, valid)
		assert.NotEmpty(t, errs)
	})

	t.Run("default allows extra fields", func(t *testing.T) {
		schema := govalidator.NewSchema().
			WithFields(map[string]*govalidator.Schema{
				"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
			})

		data := map[string]any{
			"name":  "John",
			"extra": "field",
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})
}

func TestSchema_WithItems(t *testing.T) {
	t.Run("validates array of strings", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsListValidator).
			WithItems(govalidator.NewSchema(govalidator.IsStringValidator).Required())

		data := []any{"one", "two", "three"}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("fails when array item is wrong type", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsListValidator).
			WithItems(govalidator.NewSchema(govalidator.IsStringValidator).Required())

		data := []any{"one", 2, "three"}

		valid, errs := schema.Validate(context.Background(), data)

		assert.False(t, valid)
		assert.NotEmpty(t, errs)
	})

	t.Run("validates array of objects", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsListValidator).
			WithItems(
				govalidator.NewSchema().WithFields(map[string]*govalidator.Schema{
					"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
					"age":  govalidator.NewSchema(govalidator.IsIntegerValidator).Required(),
				}),
			)

		data := []any{
			map[string]any{"name": "John", "age": 30},
			map[string]any{"name": "Jane", "age": 25},
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})
}

func TestSchema_FluentChaining(t *testing.T) {
	t.Run("complex nested schema with fluent API", func(t *testing.T) {
		schema := govalidator.NewSchema().
			WithFields(map[string]*govalidator.Schema{
				"name": govalidator.NewSchema(
					govalidator.IsStringValidator,
					govalidator.MinLengthValidator(3),
				).Required(),
				"email": govalidator.NewSchema(
					govalidator.IsStringValidator,
				).Required(),
				"age": govalidator.NewSchema(
					govalidator.IsIntegerValidator,
				).Optional(),
				"tags": govalidator.NewSchema(govalidator.IsListValidator).
					WithItems(govalidator.NewSchema(govalidator.IsStringValidator).Required()).
					Optional(),
				"address": govalidator.NewSchema().
					WithFields(map[string]*govalidator.Schema{
						"street": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
						"city":   govalidator.NewSchema(govalidator.IsStringValidator).Required(),
						"zip":    govalidator.NewSchema(govalidator.IsStringValidator).Optional(),
					}).
					Optional(),
			}).
			WithExtra(govalidator.ExtraForbid)

		data := map[string]any{
			"name":  "John Doe",
			"email": "john@example.com",
			"tags":  []any{"developer", "golang"},
			"address": map[string]any{
				"street": "123 Main St",
				"city":   "Springfield",
			},
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})
}

func TestSchema_ValidateFlat(t *testing.T) {
	t.Run("returns flat list of errors", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(map[string]*govalidator.Schema{
			"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
			"age":  govalidator.NewSchema(govalidator.IsIntegerValidator).Required(),
		})

		data := map[string]any{
			"name": 123, // wrong type
			// age missing
		}

		valid, errs := schema.ValidateFlat(
			context.Background(),
			data,
			govalidator.CombinedPresenter(".", ": "),
		)

		assert.False(t, valid)
		assert.NotEmpty(t, errs)
		// errs is []string instead of map[string][]string
		assert.IsType(t, []string{}, errs)
	})
}

func TestSchema_ValidateWithPresenter(t *testing.T) {
	t.Run("uses custom presenters", func(t *testing.T) {
		schema := govalidator.NewSchema(
			govalidator.IsStringValidator,
			govalidator.MinLengthValidator(5),
		).Required()

		valid, errs := schema.ValidateWithPresenter(
			context.Background(),
			"hi",
			govalidator.PathPresenter("."),
			govalidator.DetailedErrorPresenter(),
		)

		assert.False(t, valid)
		assert.NotEmpty(t, errs)
		// Should have detailed error message
		assert.Contains(t, errs["$"][0], "at least 5")
	})
}

func TestSchema_ToDefinition(t *testing.T) {
	t.Run("converts simple required schema", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator).Required()
		def := schema.ToDefinition()

		assert.NotNil(t, def.Validator)
		assert.Len(t, def.Validator, 2) // NonNullableValidator + IsStringValidator
	})

	t.Run("converts simple optional schema", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator).Optional()
		def := schema.ToDefinition()

		assert.NotNil(t, def.Validator)
		assert.Len(t, def.Validator, 2) // NullableValidator + IsStringValidator
	})

	t.Run("converts object schema with fields", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("age").Optional().WithValidators(govalidator.IsIntegerValidator),
		)
		def := schema.ToDefinition()

		assert.NotNil(t, def.Fields)
		assert.Len(t, *def.Fields, 2)
		assert.Contains(t, *def.Fields, "name")
		assert.Contains(t, *def.Fields, "age")
	})

	t.Run("converts array schema with items", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema(govalidator.IsStringValidator).Required(),
		).Required()
		def := schema.ToDefinition()

		assert.NotNil(t, def.ListOf)
	})

	t.Run("converts ExtraIgnore mode", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
		).WithExtra(govalidator.ExtraIgnore)
		def := schema.ToDefinition()

		assert.True(t, def.AcceptExtraProperty)
	})

	t.Run("converts ExtraForbid mode", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
		).WithExtra(govalidator.ExtraForbid)
		def := schema.ToDefinition()

		assert.False(t, def.AcceptExtraProperty)
	})

	t.Run("converts nested objects", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("user").Required().WithSchema(
				govalidator.NewSchema().WithFields(
					govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
				),
			),
		)
		def := schema.ToDefinition()

		assert.NotNil(t, def.Fields)
		userDef := (*def.Fields)["user"]
		assert.NotNil(t, userDef.Fields)
		assert.Contains(t, *userDef.Fields, "name")
	})
}

func TestNewField(t *testing.T) {
	t.Run("creates field with name", func(t *testing.T) {
		field := govalidator.NewField("email")
		assert.NotNil(t, field)
	})

	t.Run("field can be marked required", func(t *testing.T) {
		field := govalidator.NewField("email").Required()
		assert.NotNil(t, field)
	})

	t.Run("field can have validators", func(t *testing.T) {
		field := govalidator.NewField("email").
			Required().
			WithValidators(govalidator.IsStringValidator)
		assert.NotNil(t, field)
	})
}

func TestSchema_WithFieldsBuilder(t *testing.T) {
	t.Run("can use builder pattern with NewField", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("age").Optional().WithValidators(govalidator.IsIntegerValidator),
		)

		data := map[string]any{
			"name": "John",
			"age":  30,
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("builder pattern validates required fields", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("age").Required().WithValidators(govalidator.IsIntegerValidator),
		)

		data := map[string]any{
			"name": "John",
			// age is missing
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.False(t, valid)
		assert.NotEmpty(t, errs)
	})

	t.Run("builder pattern with nested schema", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("address").Optional().WithSchema(
				govalidator.Object(map[string]*govalidator.Schema{
					"street": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
					"city":   govalidator.NewSchema(govalidator.IsStringValidator).Required(),
				}),
			),
		)

		data := map[string]any{
			"name": "John",
			"address": map[string]any{
				"street": "123 Main St",
				"city":   "Springfield",
			},
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("builder pattern maintains backward compatibility with map", func(t *testing.T) {
		// Old map-based API should still work
		schema := govalidator.NewSchema().WithFields(map[string]*govalidator.Schema{
			"name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
			"age":  govalidator.NewSchema(govalidator.IsIntegerValidator).Optional(),
		})

		data := map[string]any{
			"name": "John",
			"age":  30,
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("complex nested schema with builder pattern", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").
				Required().
				WithValidators(
					govalidator.IsStringValidator,
					govalidator.MinLengthValidator(3),
				),
			govalidator.NewField("email").
				Required().
				WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("age").
				Optional().
				WithValidators(govalidator.IsIntegerValidator),
			govalidator.NewField("tags").
				Optional().
				WithSchema(
					govalidator.Array(
						govalidator.NewSchema(govalidator.IsStringValidator).Required(),
					),
				),
			govalidator.NewField("address").
				Optional().
				WithSchema(
					govalidator.Object(map[string]*govalidator.Schema{
						"street": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
						"city":   govalidator.NewSchema(govalidator.IsStringValidator).Required(),
						"zip":    govalidator.NewSchema(govalidator.IsStringValidator).Optional(),
					}),
				),
		).WithExtra(govalidator.ExtraForbid)

		data := map[string]any{
			"name":  "John Doe",
			"email": "john@example.com",
			"tags":  []any{"developer", "golang"},
			"address": map[string]any{
				"street": "123 Main St",
				"city":   "Springfield",
			},
		}

		valid, errs := schema.Validate(context.Background(), data)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})
}
