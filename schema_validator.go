package govalidator

import (
	"context"
	"fmt"
	"sort"
)

// ErrorCollector collects validation errors during the validation process.
// This interface allows for different error collection strategies.
type ErrorCollector interface {
	// Collect adds an error at the given path
	Collect(path []string, err error)

	// GetErrors returns all collected errors
	GetErrors() map[string][]string

	// HasErrors returns true if any errors were collected
	HasErrors() bool
}

// SchemaValidator validates values against Schema definitions.
// This is the modern validator that works directly with Schema,
// providing better extensibility and cleaner architecture than the legacy validator.
type SchemaValidator struct {
	pathPresenter  PresenterFunc
	errorPresenter PresenterFunc
}

// ValidationContext holds state during validation traversal.
// This allows validators to access parent context and pass data down the validation tree.
type ValidationContext struct {
	ctx            context.Context
	path           []string
	errorCollector ErrorCollector
	pathPresenter  PresenterFunc
	errorPresenter PresenterFunc
}

// MapErrorCollector collects errors into a map[string][]string structure.
type MapErrorCollector struct {
	ctx            context.Context
	errors         map[string][]string
	pathPresenter  PresenterFunc
	errorPresenter PresenterFunc
}

// FlatErrorCollector collects errors into a flat string slice.
type FlatErrorCollector struct {
	ctx      context.Context
	errors   []string
	combiner PresenterFunc
}

// NewSchemaValidator creates a new Schema-based validator with the given presenters.
//
// Example:
//
//	validator := NewSchemaValidator(
//	    PathPresenter("."),
//	    SimpleErrorPresenter(),
//	)
func NewSchemaValidator(pathPresenter PresenterFunc, errorPresenter PresenterFunc) *SchemaValidator {
	return &SchemaValidator{
		pathPresenter:  pathPresenter,
		errorPresenter: errorPresenter,
	}
}

// Validate validates a value against a schema and returns whether validation passed
// and a map of errors grouped by path.
//
// Example:
//
//	validator := NewSchemaValidator(PathPresenter("."), SimpleErrorPresenter())
//	valid, errs := validator.Validate(ctx, data, schema)
//	if !valid {
//	    for path, messages := range errs {
//	        fmt.Printf("%s: %v\n", path, messages)
//	    }
//	}
func (sv *SchemaValidator) Validate(ctx context.Context, value any, schema *Schema) (bool, map[string][]string) {
	collector := NewMapErrorCollector(ctx, sv.pathPresenter, sv.errorPresenter)

	valCtx := &ValidationContext{
		ctx:            ctx,
		path:           []string{"$"},
		errorCollector: collector,
		pathPresenter:  sv.pathPresenter,
		errorPresenter: sv.errorPresenter,
	}

	sv.validateValue(valCtx, value, schema)

	return !collector.HasErrors(), collector.GetErrors()
}

// ValidateFlat validates a value and returns errors as a flat list of strings.
//
// Example:
//
//	validator := NewSchemaValidator(PathPresenter("."), SimpleErrorPresenter())
//	valid, errs := validator.ValidateFlat(ctx, data, schema, CombinedPresenter(".", ": "))
func (sv *SchemaValidator) ValidateFlat(
	ctx context.Context,
	value any,
	schema *Schema,
	combiner PresenterFunc,
) (bool, []string) {
	collector := NewFlatErrorCollector(ctx, combiner)

	valCtx := &ValidationContext{
		ctx:            ctx,
		path:           []string{"$"},
		errorCollector: collector,
		pathPresenter:  sv.pathPresenter,
		errorPresenter: sv.errorPresenter,
	}

	sv.validateValue(valCtx, value, schema)

	return !collector.HasErrors(), collector.GetFlatErrors()
}

// validateValue is the core validation logic that handles a single value.
func (sv *SchemaValidator) validateValue(valCtx *ValidationContext, value any, schema *Schema) {
	// Step 1: Check required/optional
	if !sv.validateRequired(valCtx, value, schema) {
		return // If required check fails, stop validation
	}

	// Step 2: If value is nil and optional, skip further validation
	if value == nil && !schema.required {
		return
	}

	// Step 3: Run all validators in order
	if !sv.runValidators(valCtx, value, schema) {
		return // If validator blocks, stop validation
	}

	// Step 4: Handle nested structures
	if schema.Fields != nil {
		sv.validateObject(valCtx, value, schema)
		return
	}

	if schema.Items != nil {
		sv.validateArray(valCtx, value, schema)
		return
	}
}

// validateRequired checks if a required field is present.
func (sv *SchemaValidator) validateRequired(valCtx *ValidationContext, value any, schema *Schema) bool {
	if schema.required && value == nil {
		valCtx.errorCollector.Collect(valCtx.path, RequiredError{})
		return false
	}
	return true
}

// runValidators executes all validators in the schema.
// Returns false if a validator blocks further validation.
func (sv *SchemaValidator) runValidators(valCtx *ValidationContext, value any, schema *Schema) bool {
	for _, validator := range schema.Validators {
		shouldBlock, errs := validator(valCtx.ctx, value)

		for _, err := range errs {
			valCtx.errorCollector.Collect(valCtx.path, err)
		}

		if shouldBlock {
			return false
		}
	}
	return true
}

// validateObject validates an object/map against field schemas.
func (sv *SchemaValidator) validateObject(valCtx *ValidationContext, value any, schema *Schema) {
	// Type check
	currentMap, ok := value.(map[string]any)
	if !ok || currentMap == nil {
		valCtx.errorCollector.Collect(valCtx.path, NotAMapError{})
		return
	}

	// Get sorted field names for consistent ordering
	fieldNames := sv.getSortedFieldNames(schema)

	// Validate each defined field
	for _, fieldName := range fieldNames {
		fieldSchema := schema.Fields[fieldName]
		fieldValue, exists := currentMap[fieldName]

		if !exists {
			// Field not present - validate as nil
			childCtx := sv.pushPath(valCtx, fieldName)
			sv.validateValue(childCtx, nil, fieldSchema)
			continue
		}

		// Field present - validate its value
		childCtx := sv.pushPath(valCtx, fieldName)
		sv.validateValue(childCtx, fieldValue, fieldSchema)
	}

	// Check for extra fields
	sv.validateExtraFields(valCtx, currentMap, schema)
}

// validateArray validates an array against item schema.
func (sv *SchemaValidator) validateArray(valCtx *ValidationContext, value any, schema *Schema) {
	// Type check
	list, ok := value.([]any)
	if !ok || list == nil {
		valCtx.errorCollector.Collect(valCtx.path, NotAListError{})
		return
	}

	// Validate each item
	for i, item := range list {
		childCtx := sv.pushPath(valCtx, fmt.Sprintf("[%d]", i))
		sv.validateValue(childCtx, item, schema.Items)
	}
}

// validateExtraFields checks for unexpected fields in objects.
func (sv *SchemaValidator) validateExtraFields(valCtx *ValidationContext, currentMap map[string]any, schema *Schema) {
	if schema.Extra == ExtraIgnore {
		return
	}

	// Find fields not in schema
	for fieldName := range currentMap {
		if _, defined := schema.Fields[fieldName]; !defined {
			valCtx.errorCollector.Collect(valCtx.path, UnexpectedFieldError{
				Field: fieldName,
			})
		}
	}
}

// getSortedFieldNames returns field names in sorted order for consistent validation.
func (sv *SchemaValidator) getSortedFieldNames(schema *Schema) []string {
	names := make([]string, 0, len(schema.Fields))
	for name := range schema.Fields {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// pushPath creates a child context with an additional path segment.
func (sv *SchemaValidator) pushPath(parent *ValidationContext, segment string) *ValidationContext {
	newPath := make([]string, len(parent.path)+1)
	copy(newPath, parent.path)
	newPath[len(parent.path)] = segment

	return &ValidationContext{
		ctx:            parent.ctx,
		path:           newPath,
		errorCollector: parent.errorCollector,
		pathPresenter:  parent.pathPresenter,
		errorPresenter: parent.errorPresenter,
	}
}

// NewMapErrorCollector creates a new map-based error collector.
func NewMapErrorCollector(ctx context.Context, pathPresenter PresenterFunc, errorPresenter PresenterFunc) *MapErrorCollector {
	return &MapErrorCollector{
		ctx:            ctx,
		errors:         make(map[string][]string),
		pathPresenter:  pathPresenter,
		errorPresenter: errorPresenter,
	}
}

// Collect adds an error at the given path.
func (c *MapErrorCollector) Collect(path []string, err error) {
	pathStr := c.pathPresenter(c.ctx, path, err)
	errStr := c.errorPresenter(c.ctx, path, err)
	c.errors[pathStr] = append(c.errors[pathStr], errStr)
}

// GetErrors returns all collected errors.
func (c *MapErrorCollector) GetErrors() map[string][]string {
	return c.errors
}

// HasErrors returns true if any errors were collected.
func (c *MapErrorCollector) HasErrors() bool {
	return len(c.errors) > 0
}

// NewFlatErrorCollector creates a new flat error collector.
func NewFlatErrorCollector(ctx context.Context, combiner PresenterFunc) *FlatErrorCollector {
	return &FlatErrorCollector{
		ctx:      ctx,
		errors:   make([]string, 0),
		combiner: combiner,
	}
}

// Collect adds an error at the given path.
func (c *FlatErrorCollector) Collect(path []string, err error) {
	combined := c.combiner(c.ctx, path, err)
	c.errors = append(c.errors, combined)
}

// GetErrors returns errors as a map (for interface compatibility).
func (c *FlatErrorCollector) GetErrors() map[string][]string {
	// Return a single entry with all errors
	if len(c.errors) == 0 {
		return map[string][]string{}
	}
	return map[string][]string{
		"errors": c.errors,
	}
}

// GetFlatErrors returns errors as a flat slice.
func (c *FlatErrorCollector) GetFlatErrors() []string {
	return c.errors
}

// HasErrors returns true if any errors were collected.
func (c *FlatErrorCollector) HasErrors() bool {
	return len(c.errors) > 0
}
