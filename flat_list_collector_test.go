package govalidator_test

import (
	"context"
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlatListCollector(t *testing.T) {
	t.Run("NewFlatListCollector creates collector", func(t *testing.T) {
		collector := govalidator.NewFlatListCollector(govalidator.CombinedPresenter(".", ": "))
		assert.NotNil(t, collector)
	})

	t.Run("Collect adds errors", func(t *testing.T) {
		collector := govalidator.NewFlatListCollector(govalidator.CombinedPresenter(".", ": "))
		collector.Collect(context.Background(), []string{"$", "name"}, govalidator.RequiredError{})

		errs := collector.GetErrors()
		assert.Len(t, errs, 1)
	})

	t.Run("GetErrors returns flat list", func(t *testing.T) {
		collector := govalidator.NewFlatListCollector(govalidator.CombinedPresenter(".", ": "))
		collector.Collect(context.Background(), []string{"$", "name"}, govalidator.RequiredError{})
		collector.Collect(context.Background(), []string{"$", "age"}, govalidator.NotAnIntegerError{})

		errs := collector.GetErrors()
		assert.Len(t, errs, 2)
	})
}

func TestFlatListValidator(t *testing.T) {
	t.Run("NewFlatListValidator creates validator", func(t *testing.T) {
		validator := govalidator.NewFlatListValidator(govalidator.CombinedPresenter(".", ": "))
		assert.NotNil(t, validator)
	})

	t.Run("Validate returns flat errors", func(t *testing.T) {
		def := govalidator.Definition{
			Validator: []govalidator.ContextValidator{
				govalidator.NonNullableValidator,
				govalidator.IsStringValidator,
			},
		}

		validator := govalidator.NewFlatListValidator(govalidator.CombinedPresenter(".", ": "))
		valid, errs := validator.Validate(context.Background(), 123, def)

		assert.False(t, valid)
		assert.NotEmpty(t, errs)
		assert.IsType(t, []string{}, errs)
	})

	t.Run("Validate with complex schema", func(t *testing.T) {
		def := govalidator.Definition{
			Validator: []govalidator.ContextValidator{},
			Fields: &map[string]govalidator.Definition{
				"name": {
					Validator: []govalidator.ContextValidator{
						govalidator.NonNullableValidator,
						govalidator.IsStringValidator,
					},
				},
				"age": {
					Validator: []govalidator.ContextValidator{
						govalidator.NonNullableValidator,
						govalidator.IsIntegerValidator,
					},
				},
			},
		}

		validator := govalidator.NewFlatListValidator(govalidator.CombinedPresenter(".", ": "))
		data := map[string]any{
			"name": 123,      // wrong type
			"age":  "thirty", // wrong type
		}

		valid, errs := validator.Validate(context.Background(), data, def)

		assert.False(t, valid)
		assert.Len(t, errs, 2)
	})
}

func TestFlatErrorCollector_GetErrors(t *testing.T) {
	t.Run("GetErrors returns map format", func(t *testing.T) {
		collector := govalidator.NewFlatErrorCollector(context.Background(), govalidator.CombinedPresenter(".", ": "))
		collector.Collect([]string{"$", "name"}, govalidator.RequiredError{})

		errs := collector.GetErrors()
		assert.NotNil(t, errs)
		assert.Contains(t, errs, "errors")
	})

	t.Run("GetErrors returns empty map when no errors", func(t *testing.T) {
		collector := govalidator.NewFlatErrorCollector(context.Background(), govalidator.CombinedPresenter(".", ": "))

		errs := collector.GetErrors()
		assert.Empty(t, errs)
	})
}
