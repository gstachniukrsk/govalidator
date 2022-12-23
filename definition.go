package go_validator
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
