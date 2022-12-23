# Validator

It's a go library that aims to validate json data unmarshalled to `any`.

## Usage

```go
package go-validator

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"validator/validator"
)

const validInput = `{
		"name": "John", 
		"age": 42, 
		"gender": "male",
		"address": {
			"street": "123 Main St",
			"city": "Anytown",
			"state": "CA",
			"zip": "12345"
		},
		"phone": [
			{
				"type": "home",
				"number": "123-456-7890"
			},
			{
				"type": "work",
				"number": "123-456-7890"
			}
		]
	}`

const invalidInput = `{
		"name": "",
		"age": "42",
		"gender": "yes",
		"address": {
			"street": "1/2",
			"city": "",
			"state": "California",
			"zip": "12345"
		},
		"phone": [
			null,
			{
				"type": "office",
				"number": "123-456-"
			},
			{
				"type": "home",
				"number": "123-456-7890"
			}
		]
	}`

func getModel() validator.Definition {
	return validator.Definition{
		Validator: []validator.ContextValidator{},
		Fields: &map[string]validator.Definition{
			"name": {
				Validator: []validator.ContextValidator{
					validator.NonNullableValidator,
					validator.StringValidator,
				},
			},
			"age": {
				Validator: []validator.ContextValidator{
					validator.NonNullableValidator,
					validator.IntValidator,
				},
			},
			"gender": {
				Validator: []validator.ContextValidator{
					validator.NonNullableValidator,
					validator.StringValidator,
					validator.OneOfValidator("male", "female"),
				},
			},
			"address": {
				Validator: []validator.ContextValidator{},
				Fields: &map[string]validator.Definition{
					"street": {
						Validator: []validator.ContextValidator{
							validator.NonNullableValidator,
							validator.StringValidator,
							validator.MinLengthValidator(3),
							validator.MaxLengthValidator(100),
						},
					},
					"city": {
						Validator: []validator.ContextValidator{
							validator.NonNullableValidator,
							validator.StringValidator,
							validator.MinLengthValidator(3),
							validator.MaxLengthValidator(100),
						},
					},
					"state": {
						Validator: []validator.ContextValidator{
							validator.NonNullableValidator,
							validator.StringValidator,
							validator.UpperCaseValidator,
							validator.MinLengthValidator(2),
							validator.MaxLengthValidator(2),
						},
					},
					"zip": {
						Validator: []validator.ContextValidator{
							validator.NonNullableValidator,
							validator.StringValidator,
							validator.MinLengthValidator(5),
							validator.MaxLengthValidator(5),
						},
					},
				},
			},
			"phone": {
				Validator: []validator.ContextValidator{
					validator.NonNullableValidator,
					validator.IsListValidator,
					validator.MinSizeValidator(1, false),
				},
				ListOf: &validator.Definition{
					Validator: []validator.ContextValidator{
						validator.NonNullableValidator,
					},
					Fields: &map[string]validator.Definition{
						"type": {
							Validator: []validator.ContextValidator{
								validator.NonNullableValidator,
								validator.StringValidator,
								validator.OneOfValidator("home", "work"),
							},
						},
						"number": {
							Validator: []validator.ContextValidator{
								validator.NonNullableValidator,
								validator.RegexpValidator(
									*regexp.MustCompile("^[0-9]{3}-[0-9]{3}-[0-9]{4}$"),
								),
							},
						},
					},
				},
			},
		},
	}
}

func main() {
	v := validator.NewBasicValidator(validator.PathPresenter("."), validator.SimpleErrorPresenter())

	var target any
	err := json.Unmarshal([]byte(validInput), &target)

	if err != nil {
		panic(err)
	}

	valid, errs := v.Validate(context.Background(), target, getModel())

	fmt.Printf("----Valid Example----\n")
	fmt.Printf("valid: %v\n", valid)
	fmt.Printf("errors: %v\n", mustFormatErrs(errs))

	err = json.Unmarshal([]byte(invalidInput), &target)

	if err != nil {
		panic(err)
	}

	valid, errs = v.Validate(context.Background(), target, getModel())

	fmt.Printf("----Invalid Example----\n")
	fmt.Printf("valid: %v\n", valid)
	fmt.Printf("errors: %v\n", mustFormatErrs(errs))
}

func mustFormatErrs(errs map[string][]string) string {
	b, err := json.MarshalIndent(errs, "", "  ")

	if err != nil {
		panic(err)
	}

	return string(b)
}

```
## Predefined Validators
| Name                 | Description                                                                                               |
|----------------------|-----------------------------------------------------------------------------------------------------------|
| FloatishValidator    | checks for float or int, if float also verifies if float number is provided with predefined max precision |
| IsBooleanValidator   | if value is a boolean                                                                                     |
| IsIntegerValidator   | if value is integer or float with only 0 decimals                                                         |
| IsListValidator      | if value is of type []interface{}                                                                         |
| IsMapValidator       | if value is of type map[string]interface{}                                                                |
| IsStringValidator    | if value is of type string                                                                                |
| LowerCasValidator    | if value if string and lowercase                                                                          |
| MaxFloatValidator    | if value is float and lower than provided expectation                                                     |
| MaxLengthValidator   | if value is a string and shorter or eq than provided expectation                                          |
| MinFloatValidator    | if value is a float and higher or eq to than provided expectation                                         |
| MinLengthValidator   | if value is a string and longer or eq than provided expectation                                           |
| MinSizeValidator     | if value is a list and its count is higher or eq than provided expectation                                |
| NonNullableValidator | if value is null, breaks validation on current value                                                      |
| NullableValidator    | accepts null, breaks current value validation                                                             |
| NumberValidator      | is number, breaks validation if not, fails                                                                |
| OneOfValidator       | is one of provided values, can be of mixed types, checks with reflect.DeepEqual                           |
| RegexpValidator      | matches regexp, if not a string, fails and breaks twig                                                    |
| UpperCaseValidator   | if value not a string, breaks twig and fails, if a string but not upper cased just fails                  |

## Custom Validators
You can create your own validators by implementing the ContextValidator interface. The interface is defined as follows:

```go
package go-validator

import (
	"context"
	"errors"
	"fmt"
	"validator/validator"
)

func SumOfMapPropertiesValidator(propertyName string, expectedSum int) validator.ContextValidator {
	return func(ctx context.Context, value any) (twigBreak bool, errs []error) {
		// check if type is correct
		if _, ok := value.([]map[string]interface{}); !ok {
			errs = append(errs, errors.New("value is not a list of maps"))
			return true, errs
		}

		list := value.([]map[string]interface{})

		// do your validation here
		sum := 0

		for _, item := range list {
			if _, ok := item[propertyName]; !ok {
				errs = append(errs, validator.FieldNotDefinedError{Field: propertyName})
				return true, errs
			}

			v, ok := item[propertyName].(int)

			if !ok {
				errs = append(errs, fmt.Errorf("property %s is not an int", propertyName))
				return true, errs
			}

			sum += v
		}

		if sum != expectedSum {
			errs = append(errs, fmt.Errorf("sum of %s is not %d, got %d instead", propertyName, expectedSum, sum))
		}

		return false, errs
	}
}

```

## Model definition

### Map/Object

```go
package go-validator

import (
    "validator/validator"
)

func main() {
    objModel := validator.Definition{
		Validator: []validator.ContextValidator{
			// dont accept null value
			validator.NonNullableValidator,
			// its a map
			validator.IsMapValidator,
		},
		// can have any number of properties, not defined in Fields property
		AcceptExtraProperty: true,
		// if property not provided for example {"name": null} is not same as {} in terms of `name` property, 
		//   but if `name` property is defined in `Fields` and has a validator that fails with
		//   `validator.RequiredError` then it will fail on it anyway
		AcceptNotDefinedProperty: true,
		// it's map only if it has defined fields, can be empty, but can't be null
		Fields: &map[string]validator.Definition{
			"field1": {
                Validator: []validator.ContextValidator{
                    validator.NonNullableValidator,
                    validator.StringValidator,
                },
            },
		},
    }
	
	listModel := validator.Definition{
        Validator: []validator.ContextValidator{
            // dont accept null value
            validator.NonNullableValidator,
            // it's a list
            validator.IsListValidator,
        },
		// it's list only if ListOf is defined
		ListOf: &objModel,
    }
}

```