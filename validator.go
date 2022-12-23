package go_validator

import (
"context"
"fmt"
)
type validator struct {
	currentTree     []string
	errorCollectors []Collector
}

type Validator interface {
	Validate(ctx context.Context, key string, value any, validator Definition)
	WithCollector(collector Collector) Validator
	Copy() Validator
}

type BasicValidator interface {
	Validate(ctx context.Context, value any, def Definition) (bool, map[string][]string)
}

type basicValidator struct {
	pathPresenter PresenterFunc
	errPresenter  PresenterFunc
}

func NewBasicValidator(pathPresenter PresenterFunc, errPresenter PresenterFunc) BasicValidator {
	return &basicValidator{
		pathPresenter: pathPresenter,
		errPresenter:  errPresenter,
	}
}

func (bv *basicValidator) Validate(ctx context.Context, value any, def Definition) (bool, map[string][]string) {
	c := NewPathToErrCollector(
		PathPresenter("."),
		SimpleErrorPresenter(),
	)
	NewValidator().WithCollector(c.Collect).Validate(ctx, "$", value, def)

	errs := c.GetErrors()
	return len(errs) == 0, errs
}

func NewValidator() Validator {
	return &validator{
		errorCollectors: []Collector{},
	}
}

func (v *validator) Copy() Validator {
	currentTree := make([]string, len(v.currentTree))
	copy(currentTree, v.currentTree)

	errorCollectors := make([]Collector, len(v.errorCollectors))
	copy(errorCollectors, v.errorCollectors)

	return &validator{
		currentTree:     currentTree,
		errorCollectors: errorCollectors,
	}
}

func (v *validator) WithCollector(collector Collector) Validator {
	v2 := v.Copy().(*validator)
	v2.errorCollectors = append(v.errorCollectors, collector)
	return v2
}

func (v *validator) pushTree(value string) {
	v.currentTree = append(v.currentTree, value)
}

func (v *validator) popTree() {
	v.currentTree = v.currentTree[:len(v.currentTree)-1]
}

func (v *validator) addError(ctx context.Context, err error) {
	for _, collector := range v.errorCollectors {
		collector(ctx, v.currentTree, err)
	}
}

func (v *validator) Validate(ctx context.Context, key string, value any, def Definition) {
	v.pushTree(key)
	defer v.popTree()
	for _, validator := range def.Validator {
		twigBlock, errs := validator(ctx, value)
		for _, err := range errs {
			v.addError(ctx, err)
		}
		if twigBlock {
			return
		}
	}
	if def.Fields != nil {
		v.handleObject(ctx, value, def)
		return
	}

	if def.ListOf != nil {
		v.handleList(ctx, value, def)
		return
	}
}

func (v *validator) handleObject(ctx context.Context, value any, def Definition) {
	currentMap, ok := value.(map[string]interface{})

	if !ok || currentMap == nil {
		v.addError(ctx, NotAMapError{})
		return
	}

	for field, fDef := range *def.Fields {
		val, ok := currentMap[field]

		if !ok {
			if !def.AcceptNotDefinedProperty {
				v.addError(ctx, FieldNotDefinedError{
					Field: field,
				})
				continue
			}

			// check if field definition contains non-null validator,
			//	we need to fail that case
			for _, fCtxV := range fDef.Validator {
				if !fCtxV.AcceptsNull() {
					v.addError(ctx, FieldNotDefinedError{
						Field: field,
					})
					continue
				}
			}

			continue
		}

		v.Validate(ctx, field, val, fDef)
	}

	if def.AcceptExtraProperty {
		return
	}
	// get all keys from map
	keys := make([]string, 0, len(currentMap))
	for key := range currentMap {
		_, ok := (*def.Fields)[key]

		if ok {
			continue
		}

		keys = append(keys, key)
	}

	for _, k := range keys {
		v.addError(ctx, UnexpectedFieldError{
			Field: k,
		})
	}
}

func (v *validator) handleList(ctx context.Context, value any, def Definition) {
	l, ok := value.([]interface{})

	if !ok || l == nil {
		v.addError(ctx, NotAListError{})
		return
	}

	for i, lv := range l {
		v.Validate(ctx, fmt.Sprintf("[%d]", i), lv, *def.ListOf)
	}
}
