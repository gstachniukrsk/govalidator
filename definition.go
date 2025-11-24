package govalidator

// Definition describes validation rules for a value, including validators, fields, and list items.
type Definition struct {
	Validator []ContextValidator
	// for objects
	Fields *map[string]Definition
	// for lists
	ListOf *Definition
	// for object
	AcceptNotDefinedProperty bool
	// for object
	AcceptExtraProperty bool
}
