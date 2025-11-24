package govalidator_test

import (
	"context"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchemaValidator_SimpleValidation(t *testing.T) {
	t.Run("validates simple required string", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator).Required()
		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), "hello", schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("rejects nil for required field", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator).Required()
		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), nil, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$")
		assert.Contains(t, errs["$"][0], "required")
	})

	t.Run("accepts nil for optional field", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator).Optional()
		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), nil, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})
}

func TestSchemaValidator_ObjectValidation(t *testing.T) {
	t.Run("validates simple object with required fields", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("age").Required().WithValidators(govalidator.IsIntegerValidator),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"name": "John",
			"age":  30,
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("fails when required field is missing", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("age").Required().WithValidators(govalidator.IsIntegerValidator),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"name": "John",
			// age missing
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$.age")
	})

	t.Run("fails when field has wrong type", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("age").Required().WithValidators(govalidator.IsIntegerValidator),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"name": "John",
			"age":  "thirty", // wrong type
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$.age")
	})
}

func TestSchemaValidator_NestedObjects(t *testing.T) {
	t.Run("validates 2-level nested object", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("address").Required().WithSchema(
				govalidator.NewSchema().WithFields(
					govalidator.NewField("street").Required().WithValidators(govalidator.IsStringValidator),
					govalidator.NewField("city").Required().WithValidators(govalidator.IsStringValidator),
				),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"name": "John",
			"address": map[string]any{
				"street": "123 Main St",
				"city":   "Springfield",
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("validates 3-level nested object", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("user").Required().WithSchema(
				govalidator.NewSchema().WithFields(
					govalidator.NewField("profile").Required().WithSchema(
						govalidator.NewSchema().WithFields(
							govalidator.NewField("bio").Required().WithValidators(
								govalidator.IsStringValidator,
								govalidator.MinLengthValidator(10),
							),
						),
					),
				),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"user": map[string]any{
				"profile": map[string]any{
					"bio": "This is a long biography about the user",
				},
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("reports errors at correct nested path", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("user").Required().WithSchema(
				govalidator.NewSchema().WithFields(
					govalidator.NewField("profile").Required().WithSchema(
						govalidator.NewSchema().WithFields(
							govalidator.NewField("bio").Required().WithValidators(
								govalidator.IsStringValidator,
								govalidator.MinLengthValidator(10),
							),
						),
					),
				),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"user": map[string]any{
				"profile": map[string]any{
					"bio": "short", // too short
				},
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$.user.profile.bio")
	})

	t.Run("reports missing field at correct nested path", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("user").Required().WithSchema(
				govalidator.NewSchema().WithFields(
					govalidator.NewField("profile").Required().WithSchema(
						govalidator.NewSchema().WithFields(
							govalidator.NewField("bio").Required().WithValidators(govalidator.IsStringValidator),
						),
					),
				),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"user": map[string]any{
				"profile": map[string]any{
					// bio missing
				},
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$.user.profile.bio")
	})
}

func TestSchemaValidator_ArrayValidation(t *testing.T) {
	t.Run("validates simple array", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema(govalidator.IsStringValidator).Required(),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := []any{"one", "two", "three"}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("reports error at correct array index", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema(govalidator.IsStringValidator).Required(),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := []any{"one", 2, "three"} // index 1 is wrong type

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$[1]")
	})

	t.Run("validates array of objects", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema().WithFields(
				govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
				govalidator.NewField("age").Required().WithValidators(govalidator.IsIntegerValidator),
			).Required(),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := []any{
			map[string]any{"name": "John", "age": 30},
			map[string]any{"name": "Jane", "age": 25},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("reports error in nested object within array", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema().WithFields(
				govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
				govalidator.NewField("age").Required().WithValidators(govalidator.IsIntegerValidator),
			).Required(),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := []any{
			map[string]any{"name": "John", "age": 30},
			map[string]any{"name": "Jane", "age": "twenty-five"}, // wrong type
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$[1].age")
	})
}

func TestSchemaValidator_ComplexNesting(t *testing.T) {
	t.Run("validates array of objects with nested objects", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("users").Required().WithSchema(
				govalidator.Array(
					govalidator.NewSchema().WithFields(
						govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
						govalidator.NewField("address").Optional().WithSchema(
							govalidator.NewSchema().WithFields(
								govalidator.NewField("city").Required().WithValidators(govalidator.IsStringValidator),
							),
						),
					).Required(),
				),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"users": []any{
				map[string]any{
					"name": "John",
					"address": map[string]any{
						"city": "NYC",
					},
				},
				map[string]any{
					"name": "Jane",
					"address": map[string]any{
						"city": "LA",
					},
				},
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("reports error in deeply nested structure", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("users").Required().WithSchema(
				govalidator.Array(
					govalidator.NewSchema().WithFields(
						govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
						govalidator.NewField("address").Optional().WithSchema(
							govalidator.NewSchema().WithFields(
								govalidator.NewField("city").Required().WithValidators(govalidator.IsStringValidator),
							),
						),
					).Required(),
				),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"users": []any{
				map[string]any{
					"name": "John",
					"address": map[string]any{
						"city": "NYC",
					},
				},
				map[string]any{
					"name": "Jane",
					"address": map[string]any{
						"city": 12345, // wrong type
					},
				},
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$.users[1].address.city")
	})

	t.Run("validates nested arrays", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("matrix").Required().WithSchema(
				govalidator.Array(
					govalidator.Array(
						govalidator.NewSchema(govalidator.IsIntegerValidator).Required(),
					).Required(),
				),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"matrix": []any{
				[]any{1, 2, 3},
				[]any{4, 5, 6},
				[]any{7, 8, 9},
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("reports error in nested array", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("matrix").Required().WithSchema(
				govalidator.Array(
					govalidator.Array(
						govalidator.NewSchema(govalidator.IsIntegerValidator).Required(),
					).Required(),
				),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"matrix": []any{
				[]any{1, 2, 3},
				[]any{4, "five", 6}, // wrong type
				[]any{7, 8, 9},
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$.matrix[1][1]")
	})
}

func TestSchemaValidator_CustomValidators(t *testing.T) {
	t.Run("validates with multiple custom validators", func(t *testing.T) {
		schema := govalidator.NewSchema(
			govalidator.IsStringValidator,
			govalidator.MinLengthValidator(5),
			govalidator.MaxLengthValidator(20),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), "hello world", schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("fails when validator condition not met", func(t *testing.T) {
		schema := govalidator.NewSchema(
			govalidator.IsStringValidator,
			govalidator.MinLengthValidator(10),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), "short", schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$")
	})

	t.Run("validates with OneOfValidator", func(t *testing.T) {
		schema := govalidator.NewSchema(
			govalidator.IsStringValidator,
			govalidator.OneOfValidator("admin", "user", "guest"),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), "admin", schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("fails when OneOfValidator not satisfied", func(t *testing.T) {
		schema := govalidator.NewSchema(
			govalidator.IsStringValidator,
			govalidator.OneOfValidator("admin", "user", "guest"),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), "superuser", schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$")
	})

	t.Run("validates complex schema with multiple validators at different levels", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("username").Required().WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(3),
				govalidator.MaxLengthValidator(20),
			),
			govalidator.NewField("email").Required().WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(5),
			),
			govalidator.NewField("role").Required().WithValidators(
				govalidator.IsStringValidator,
				govalidator.OneOfValidator("admin", "user"),
			),
			govalidator.NewField("age").Optional().WithValidators(
				govalidator.IsIntegerValidator,
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"username": "john_doe",
			"email":    "john@example.com",
			"role":     "user",
			"age":      30,
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("reports multiple validation errors", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("username").Required().WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(3),
			),
			govalidator.NewField("age").Required().WithValidators(
				govalidator.IsIntegerValidator,
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"username": "jo",     // too short
			"age":      "thirty", // wrong type
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$.username")
		assert.Contains(t, errs, "$.age")
	})
}

func TestSchemaValidator_ExtraFields(t *testing.T) {
	t.Run("ExtraIgnore allows extra fields", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
		).WithExtra(govalidator.ExtraIgnore)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"name":  "John",
			"extra": "field",
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("ExtraForbid rejects extra fields", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
		).WithExtra(govalidator.ExtraForbid)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"name":  "John",
			"extra": "field",
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$")
	})

	t.Run("ExtraForbid works at nested levels", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("user").Required().WithSchema(
				govalidator.NewSchema().WithFields(
					govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
				).WithExtra(govalidator.ExtraForbid),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"user": map[string]any{
				"name":  "John",
				"extra": "field", // should be rejected
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$.user")
	})
}

func TestSchemaValidator_EdgeCases(t *testing.T) {
	t.Run("handles empty object", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("handles empty array", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema(govalidator.IsStringValidator).Required(),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := []any{}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("rejects wrong type for object", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), "not an object", schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$")
	})

	t.Run("rejects wrong type for array", func(t *testing.T) {
		schema := govalidator.Array(
			govalidator.NewSchema(govalidator.IsStringValidator).Required(),
		).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), "not an array", schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$")
	})

	t.Run("handles nil in nested structure", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("user").Optional().WithSchema(
				govalidator.NewSchema().WithFields(
					govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
				),
			),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"user": nil, // optional, so this is ok
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.True(t, valid)
		assert.Empty(t, errs)
	})
}

func TestSchemaValidator_CustomPresenters(t *testing.T) {
	t.Run("uses custom path presenter", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("user").Required().WithSchema(
				govalidator.NewSchema().WithFields(
					govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
				),
			),
		)

		// Use bracket notation for paths
		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("/"),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"user": map[string]any{
				// name missing
			},
		}

		valid, errs := validator.Validate(context.Background(), data, schema)

		assert.False(t, valid)
		assert.Contains(t, errs, "$/user/name")
	})

	t.Run("uses custom error presenter", func(t *testing.T) {
		schema := govalidator.NewSchema(govalidator.IsStringValidator).Required()

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.DetailedErrorPresenter(),
		)

		valid, errs := validator.Validate(context.Background(), 123, schema)

		assert.False(t, valid)
		// DetailedErrorPresenter gives more descriptive messages
		assert.NotEmpty(t, errs["$"])
	})
}

func TestSchemaValidator_ValidateFlat(t *testing.T) {
	t.Run("returns flat list of errors", func(t *testing.T) {
		schema := govalidator.NewSchema().WithFields(
			govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
			govalidator.NewField("age").Required().WithValidators(govalidator.IsIntegerValidator),
		)

		validator := govalidator.NewSchemaValidator(
			govalidator.PathPresenter("."),
			govalidator.SimpleErrorPresenter(),
		)

		data := map[string]any{
			"name": 123, // wrong type
			// age missing
		}

		valid, errs := validator.ValidateFlat(
			context.Background(),
			data,
			schema,
			govalidator.CombinedPresenter(".", ": "),
		)

		assert.False(t, valid)
		assert.Len(t, errs, 2)
		// Errors should be in format "path: message"
		assert.Contains(t, errs[0]+errs[1], "$.name")
		assert.Contains(t, errs[0]+errs[1], "$.age")
	})
}
