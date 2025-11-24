package govalidator

import "context"

// ExtraFieldsMode defines how additional properties are handled in objects.
type ExtraFieldsMode int

// Schema represents a validation schema with a simpler, more intuitive structure.
// Fields are public for easy construction, but fluent methods are also provided.
//
// Example (struct literal):
//
//	schema := &Schema{
//	    Validators: []ContextValidator{IsStringValidator, MinLengthValidator(5)},
//	}
//
// Example (fluent):
//
//	schema := NewSchema(IsStringValidator, MinLengthValidator(5)).Required()
type Schema struct {
	// Validators applied to the current value (applied in order)
	Validators []ContextValidator

	// Fields defines the schema for object/map properties
	// nil means this is not an object schema
	Fields map[string]*Schema

	// Items defines the schema for array elements
	// nil means this is not an array schema
	Items *Schema

	// Extra controls how unexpected fields are handled in objects
	// Default is ExtraIgnore (allow extra fields)
	Extra ExtraFieldsMode

	// required tracks if this value must be non-null
	// Use Required() and Optional() methods to set this
	required bool
}

// Field represents a field definition with its name and schema.
// This is used with the builder pattern for defining object fields.
type Field struct {
	name   string
	schema *Schema
}

const (
	// ExtraForbid causes validation to fail if extra fields are present.
	ExtraForbid ExtraFieldsMode = iota
	// ExtraIgnore allows extra fields without validation or errors.
	ExtraIgnore
)

// NewSchema creates a Schema with the given validators.
// By default, schemas are optional (nullable) and allow extra fields.
//
// Example:
//
//	schema := NewSchema(
//	    IsStringValidator,
//	    MinLengthValidator(5),
//	)
func NewSchema(validators ...ContextValidator) *Schema {
	return &Schema{
		Validators: validators,
		Extra:      ExtraIgnore, // Default: allow extra fields
		required:   false,       // Default: optional (nullable)
	}
}

// Object creates a schema for validating objects/maps with specific fields.
// This is a convenience function equivalent to &Schema{Fields: fields}.
//
// Example:
//
//	schema := Object(map[string]*Schema{
//	    "name": NewSchema(IsStringValidator).Required(),
//	    "age":  NewSchema(IsIntegerValidator),
//	})
func Object(fields map[string]*Schema) *Schema {
	return &Schema{
		Fields: fields,
		Extra:  ExtraIgnore,
	}
}

// Array creates a schema for validating arrays where each item matches the given schema.
// This is a convenience function equivalent to &Schema{Validators: []ContextValidator{IsListValidator}, Items: itemSchema}.
//
// Example:
//
//	schema := Array(NewSchema(IsStringValidator).Required())
func Array(itemSchema *Schema) *Schema {
	return &Schema{
		Validators: []ContextValidator{IsListValidator},
		Items:      itemSchema,
	}
}

// NewField creates a new field definition with the given name.
// Use the fluent methods to configure the field's validation.
//
// Example:
//
//	field := NewField("email").Required().WithValidators(IsStringValidator)
func NewField(name string) *Field {
	return &Field{
		name:   name,
		schema: NewSchema(),
	}
}

// Required marks this field as required (non-null).
// Returns the field for method chaining.
func (f *Field) Required() *Field {
	f.schema.Required()
	return f
}

// Optional marks this field as optional (nullable).
// Returns the field for method chaining.
func (f *Field) Optional() *Field {
	f.schema.Optional()
	return f
}

// WithValidators sets the validators for this field.
// Returns the field for method chaining.
//
// Example:
//
//	NewField("age").Required().WithValidators(IsIntegerValidator, MinValidator(0))
func (f *Field) WithValidators(validators ...ContextValidator) *Field {
	f.schema.Validators = validators
	return f
}

// WithSchema sets a complete schema for this field.
// This is useful for nested objects or arrays.
// Returns the field for method chaining.
//
// Example:
//
//	NewField("address").WithSchema(
//	    Object(map[string]*Schema{
//	        "street": NewSchema(IsStringValidator).Required(),
//	    }),
//	)
func (f *Field) WithSchema(schema *Schema) *Field {
	f.schema = schema
	return f
}

// WithFields sets the Fields map (fluent style).
// Returns the schema for method chaining.
//
// Example with map:
//
//	schema := NewSchema().WithFields(map[string]*Schema{
//	    "name": NewSchema(IsStringValidator).Required(),
//	})
//
// Example with builder pattern:
//
//	schema := NewSchema().WithFields(
//	    NewField("name").Required().WithValidators(IsStringValidator),
//	    NewField("age").Optional().WithValidators(IsIntegerValidator),
//	)
func (s *Schema) WithFields(fields ...any) *Schema {
	// Handle two cases:
	// 1. Single map[string]*Schema argument (backward compatible)
	// 2. Multiple Field arguments (new builder pattern)

	if len(fields) == 1 {
		if fieldMap, ok := fields[0].(map[string]*Schema); ok {
			s.Fields = fieldMap
			return s
		}
	}

	// Convert Field builders to map
	fieldMap := make(map[string]*Schema)
	for _, f := range fields {
		if field, ok := f.(*Field); ok {
			fieldMap[field.name] = field.schema
		}
	}

	s.Fields = fieldMap
	return s
}

// WithItems sets the Items schema (fluent style).
// Returns the schema for method chaining.
//
// Example:
//
//	schema := NewSchema(IsListValidator).WithItems(
//	    NewSchema(IsStringValidator),
//	)
func (s *Schema) WithItems(itemSchema *Schema) *Schema {
	s.Items = itemSchema
	return s
}

// WithExtra sets the extra fields mode (fluent style).
// Returns the schema for method chaining.
//
// Example:
//
//	schema := Object(fields).WithExtra(ExtraForbid)
func (s *Schema) WithExtra(mode ExtraFieldsMode) *Schema {
	s.Extra = mode
	return s
}

// Required marks this schema as required (non-null).
// Returns the schema for method chaining.
//
// Example:
//
//	schema := NewSchema(IsStringValidator).Required()
func (s *Schema) Required() *Schema {
	s.required = true
	return s
}

// Optional marks this schema as optional (nullable).
// This is the default, so it's mainly useful for clarity or overriding.
// Returns the schema for method chaining.
//
// Example:
//
//	schema := NewSchema(IsStringValidator).Optional()
func (s *Schema) Optional() *Schema {
	s.required = false
	return s
}

// IsRequired returns true if this schema requires a non-null value.
func (s *Schema) IsRequired() bool {
	return s.required
}

// ToDefinition converts a Schema to the legacy Definition format.
// This provides backward compatibility during migration to the new API.
func (s *Schema) ToDefinition() Definition {
	def := Definition{
		Validator: make([]ContextValidator, len(s.Validators)),
	}
	copy(def.Validator, s.Validators)

	// Handle required/optional by prepending appropriate validator
	if s.required {
		def.Validator = append([]ContextValidator{NonNullableValidator}, def.Validator...)
	} else {
		def.Validator = append([]ContextValidator{NullableValidator}, def.Validator...)
	}

	// Handle fields (object validation)
	if s.Fields != nil {
		fieldDefs := make(map[string]Definition)
		for name, fieldSchema := range s.Fields {
			fieldDefs[name] = fieldSchema.ToDefinition()
		}
		def.Fields = &fieldDefs
		def.AcceptExtraProperty = s.Extra == ExtraIgnore
		def.AcceptNotDefinedProperty = true // Always allow undefined fields for now
	}

	// Handle items (array validation)
	if s.Items != nil {
		itemDef := s.Items.ToDefinition()
		def.ListOf = &itemDef
	}

	return def
}

// Validate validates a value against this schema using default presenters.
// Uses the modern SchemaValidator for direct, non-legacy validation.
//
// Example:
//
//	valid, errs := schema.Validate(ctx, data)
//	if !valid {
//	    for path, messages := range errs {
//	        fmt.Printf("%s: %v\n", path, messages)
//	    }
//	}
func (s *Schema) Validate(ctx context.Context, value any) (bool, map[string][]string) {
	v := NewSchemaValidator(PathPresenter("."), SimpleErrorPresenter())
	return v.Validate(ctx, value, s)
}

// ValidateWithPresenter validates a value with custom error presentation.
// Uses the modern SchemaValidator for direct, non-legacy validation.
//
// Example:
//
//	valid, errs := schema.ValidateWithPresenter(
//	    ctx,
//	    data,
//	    PathPresenter("."),
//	    DetailedErrorPresenter(),
//	)
func (s *Schema) ValidateWithPresenter(
	ctx context.Context,
	value any,
	pathPresenter PresenterFunc,
	errorPresenter PresenterFunc,
) (bool, map[string][]string) {
	v := NewSchemaValidator(pathPresenter, errorPresenter)
	return v.Validate(ctx, value, s)
}

// ValidateFlat validates and returns errors as a flat list of strings.
// Uses the modern SchemaValidator for direct, non-legacy validation.
//
// Example:
//
//	valid, errs := schema.ValidateFlat(
//	    ctx,
//	    data,
//	    CombinedPresenter(".", ": "),
//	)
//	// errs is []string: ["$.name: required", "$.age: not an integer"]
func (s *Schema) ValidateFlat(
	ctx context.Context,
	value any,
	combiner PresenterFunc,
) (bool, []string) {
	v := NewSchemaValidator(PathPresenter("."), SimpleErrorPresenter())
	return v.ValidateFlat(ctx, value, s, combiner)
}
