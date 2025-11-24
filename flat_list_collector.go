package govalidator

import "context"

// flatListCollector collects errors as a flat list of strings rather than a map.
type flatListCollector struct {
	errors   []string
	combiner PresenterFunc
}

// FlatListCollector is an interface for collecting errors as a flat list.
type FlatListCollector interface {
	Collect(ctx context.Context, path []string, err error)
	GetErrors() []string
}

// NewFlatListCollector creates a new collector that returns errors as a flat string slice.
// The combiner function determines how each error is formatted.
//
// Example usage:
//
//	collector := NewFlatListCollector(CombinedPresenter(".", ": "))
//	// Results in: []string{"$.user.age: not an integer", "$.user.name: required"}
func NewFlatListCollector(combiner PresenterFunc) FlatListCollector {
	return &flatListCollector{
		errors:   make([]string, 0),
		combiner: combiner,
	}
}

// Collect adds an error to the flat list.
func (c *flatListCollector) Collect(ctx context.Context, path []string, err error) {
	errorStr := c.combiner(ctx, path, err)
	c.errors = append(c.errors, errorStr)
}

// GetErrors returns the collected errors as a flat slice of strings.
func (c *flatListCollector) GetErrors() []string {
	return c.errors
}

// FlatListValidator wraps the validation to return a flat list of error strings.
type FlatListValidator interface {
	Validate(ctx context.Context, value any, def Definition) (bool, []string)
}

type flatListValidator struct {
	combiner PresenterFunc
}

// NewFlatListValidator creates a validator that returns errors as a flat list of strings.
//
// Example usage:
//
//	v := NewFlatListValidator(CombinedPresenter(".", ": "))
//	valid, errs := v.Validate(ctx, data, definition)
//	// errs is []string{"$.user.age: not an integer", "$.user.name: required"}
func NewFlatListValidator(combiner PresenterFunc) FlatListValidator {
	return &flatListValidator{
		combiner: combiner,
	}
}

// Validate performs validation and returns a flat list of error strings.
func (v *flatListValidator) Validate(ctx context.Context, value any, def Definition) (bool, []string) {
	c := NewFlatListCollector(v.combiner)
	NewValidator().WithCollector(c.Collect).Validate(ctx, "$", value, def)

	errs := c.GetErrors()
	return len(errs) == 0, errs
}
