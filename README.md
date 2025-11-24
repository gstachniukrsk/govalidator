# GoValidator

A modern, extensible Go library for validating JSON data with a clean, fluent API.

## Features

✨ **Modern Schema API** - Clean, fluent interface with builder pattern
🔧 **Extensible** - Easy to create custom validators
📝 **Type-safe** - Strong typing with Go generics support
🎯 **Precise Error Reporting** - Exact paths to validation errors
🚀 **High Performance** - Optimized validation engine
🔄 **Backward Compatible** - Legacy Definition API still supported

## Quick Start

### Installation

```bash
go get github.com/gstachniukrsk/govalidator
```

### Basic Example

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/gstachniukrsk/govalidator"
)

func main() {
    // Define validation schema using the modern builder pattern
    schema := govalidator.NewSchema().WithFields(
        govalidator.NewField("name").
            Required().
            WithValidators(
                govalidator.IsStringValidator,
                govalidator.MinLengthValidator(3),
            ),
        govalidator.NewField("age").
            Required().
            WithValidators(govalidator.IsIntegerValidator),
        govalidator.NewField("email").
            Required().
            WithValidators(govalidator.IsStringValidator),
    ).WithExtra(govalidator.ExtraForbid)

    // Parse JSON data
    jsonData := `{"name": "John", "age": 30, "email": "john@example.com"}`
    var data any
    json.Unmarshal([]byte(jsonData), &data)

    // Validate
    valid, errs := schema.Validate(context.Background(), data)

    if !valid {
        for path, messages := range errs {
            fmt.Printf("%s: %v\n", path, messages)
        }
    } else {
        fmt.Println("Validation passed!")
    }
}
```

## Examples

Check out the [examples/](examples/) directory for real-world, runnable examples:

- **[http-api/](examples/http-api/)** - REST API request validation with user registration and product creation endpoints
- **[csv-validator/](examples/csv-validator/)** - CSV file validation with row-by-row error reporting and JSON export
- **[config-loader/](examples/config-loader/)** - Application configuration validation with nested objects and fail-fast loading
- **[webhook-handler/](examples/webhook-handler/)** - GitHub and Stripe webhook validation with signature verification

Each example is self-contained and includes detailed documentation. See [examples/README.md](examples/README.md) for more information.

## Schema API (Recommended)

The modern Schema API provides a clean, intuitive way to define validation rules.

### Simple Field Validation

```go
// Required string field
schema := govalidator.NewSchema(govalidator.IsStringValidator).Required()

// Optional integer field
schema := govalidator.NewSchema(govalidator.IsIntegerValidator).Optional()

// Multiple validators
schema := govalidator.NewSchema(
    govalidator.IsStringValidator,
    govalidator.MinLengthValidator(5),
    govalidator.MaxLengthValidator(100),
).Required()
```

### Object Validation with Builder Pattern

The builder pattern provides the cleanest syntax for complex schemas:

```go
schema := govalidator.NewSchema().WithFields(
    govalidator.NewField("username").
        Required().
        WithValidators(
            govalidator.IsStringValidator,
            govalidator.MinLengthValidator(3),
            govalidator.MaxLengthValidator(20),
        ),
    govalidator.NewField("email").
        Required().
        WithValidators(govalidator.IsStringValidator),
    govalidator.NewField("role").
        Required().
        WithValidators(
            govalidator.IsStringValidator,
            govalidator.OneOfValidator("admin", "user", "guest"),
        ),
    govalidator.NewField("age").
        Optional().
        WithValidators(govalidator.IsIntegerValidator),
).WithExtra(govalidator.ExtraForbid)
```

### Nested Objects

```go
schema := govalidator.NewSchema().WithFields(
    govalidator.NewField("user").Required().WithSchema(
        govalidator.NewSchema().WithFields(
            govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
            govalidator.NewField("address").Optional().WithSchema(
                govalidator.NewSchema().WithFields(
                    govalidator.NewField("street").Required().WithValidators(govalidator.IsStringValidator),
                    govalidator.NewField("city").Required().WithValidators(govalidator.IsStringValidator),
                    govalidator.NewField("zip").Optional().WithValidators(govalidator.IsStringValidator),
                ),
            ),
        ),
    ),
)
```

### Array Validation

```go
// Simple array of strings
schema := govalidator.Array(
    govalidator.NewSchema(govalidator.IsStringValidator).Required(),
).Required()

// Array of objects
schema := govalidator.Array(
    govalidator.NewSchema().WithFields(
        govalidator.NewField("id").Required().WithValidators(govalidator.IsIntegerValidator),
        govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
    ).Required(),
).Required()
```

### Complex Example

```go
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
    govalidator.NewField("phones").
        Required().
        WithSchema(
            govalidator.Array(
                govalidator.NewSchema().WithFields(
                    govalidator.NewField("type").
                        Required().
                        WithValidators(
                            govalidator.IsStringValidator,
                            govalidator.OneOfValidator("home", "work", "mobile"),
                        ),
                    govalidator.NewField("number").
                        Required().
                        WithValidators(
                            govalidator.IsStringValidator,
                            govalidator.RegexpValidator(*regexp.MustCompile(`^\d{3}-\d{3}-\d{4}$`)),
                        ),
                ).Required(),
            ),
        ),
).WithExtra(govalidator.ExtraForbid)
```

### Alternative Syntax Options

The Schema API supports multiple syntax styles:

```go
// 1. Builder Pattern (Recommended)
schema := govalidator.NewSchema().WithFields(
    govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
    govalidator.NewField("age").Optional().WithValidators(govalidator.IsIntegerValidator),
)

// 2. Object Helper with Map
schema := govalidator.Object(map[string]*govalidator.Schema{
    "name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
    "age":  govalidator.NewSchema(govalidator.IsIntegerValidator).Optional(),
})

// 3. Struct Literal
schema := &govalidator.Schema{
    Fields: map[string]*govalidator.Schema{
        "name": govalidator.NewSchema(govalidator.IsStringValidator).Required(),
        "age":  govalidator.NewSchema(govalidator.IsIntegerValidator).Optional(),
    },
    Extra: govalidator.ExtraForbid,
}
```

### Extra Fields Control

```go
// Allow extra fields (default)
schema := govalidator.NewSchema().WithFields(
    govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
).WithExtra(govalidator.ExtraIgnore)

// Forbid extra fields
schema := govalidator.NewSchema().WithFields(
    govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
).WithExtra(govalidator.ExtraForbid)
```

### Validation Methods

```go
// Standard validation (returns map of errors)
valid, errs := schema.Validate(ctx, data)
if !valid {
    for path, messages := range errs {
        fmt.Printf("%s: %v\n", path, messages)
    }
}

// With custom presenters
valid, errs := schema.ValidateWithPresenter(
    ctx,
    data,
    govalidator.PathPresenter("."),
    govalidator.DetailedErrorPresenter(),
)

// Flat list of errors
valid, errs := schema.ValidateFlat(
    ctx,
    data,
    govalidator.CombinedPresenter(".", ": "),
)
// errs is []string: ["$.name: required", "$.age: not an integer"]
```

## Predefined Validators

| Validator | Description |
|-----------|-------------|
| `IsStringValidator` | Validates that value is a string |
| `IsIntegerValidator` | Validates that value is an integer or float with zero decimals |
| `IsBooleanValidator` | Validates that value is a boolean |
| `IsListValidator` | Validates that value is an array ([]any) |
| `IsMapValidator` | Validates that value is a map (map[string]any) |
| `FloatValidator(precision)` | Validates float with maximum precision |
| `MinLengthValidator(min)` | Validates string has minimum length (counts runes, not bytes) |
| `MaxLengthValidator(max)` | Validates string has maximum length (counts runes, not bytes) |
| `MinFloatValidator(min)` | Validates number is >= minimum |
| `MaxFloatValidator(max)` | Validates number is <= maximum |
| `MinSizeValidator(min, blocking)` | Validates array has minimum size |
| `MaxSizeValidator(max, blocking)` | Validates array has maximum size |
| `OneOfValidator(values...)` | Validates value is one of the allowed values |
| `RegexpValidator(pattern)` | Validates string matches regular expression |
| `UpperCaseValidator` | Validates string is uppercase |
| `LowerCaseValidator` | Validates string is lowercase |
| `NumberValidator` | Validates value is a number (int or float) |

## Custom Validators

Create custom validators by implementing the `ContextValidator` function signature:

```go
func MyCustomValidator(ctx context.Context, value any) (twigBreak bool, errs []error) {
    // Check type
    str, ok := value.(string)
    if !ok {
        return true, []error{errors.New("not a string")}
    }

    // Your validation logic
    if !strings.HasPrefix(str, "CUSTOM-") {
        return false, []error{errors.New("must start with CUSTOM-")}
    }

    return false, nil
}

// Use in schema
schema := govalidator.NewSchema(
    govalidator.IsStringValidator,
    MyCustomValidator,
).Required()
```

### Advanced Custom Validator Example

```go
func EmailDomainValidator(allowedDomains ...string) govalidator.ContextValidator {
    return func(ctx context.Context, value any) (bool, []error) {
        str, ok := value.(string)
        if !ok {
            return true, []error{errors.New("not a string")}
        }

        parts := strings.Split(str, "@")
        if len(parts) != 2 {
            return false, []error{errors.New("invalid email format")}
        }

        domain := parts[1]
        for _, allowed := range allowedDomains {
            if domain == allowed {
                return false, nil
            }
        }

        return false, []error{fmt.Errorf("domain must be one of: %v", allowedDomains)}
    }
}

// Use it
schema := govalidator.NewField("email").
    Required().
    WithValidators(
        govalidator.IsStringValidator,
        EmailDomainValidator("company.com", "example.com"),
    )
```

## Error Presenters

Customize how errors are formatted:

```go
// Simple error messages (default)
govalidator.SimpleErrorPresenter()

// Detailed error messages with context
govalidator.DetailedErrorPresenter()

// JSON format
govalidator.JSONPresenter(".")
govalidator.JSONDetailedPresenter(".")

// Combined path and error
govalidator.CombinedPresenter(".", ": ")
govalidator.CombinedBracketPresenter(".", ": ")

// Custom presenter
func MyPresenter() govalidator.PresenterFunc {
    return func(ctx context.Context, path []string, err error) string {
        return fmt.Sprintf("ERROR at %s: %v", strings.Join(path, " > "), err)
    }
}
```

## Migration from Definition API

If you're using the legacy Definition API, see [MIGRATION_SCHEMA_API.md](MIGRATION_SCHEMA_API.md) for a complete migration guide.

Quick comparison:

```go
// Old Definition API
def := govalidator.Definition{
    Validator: []govalidator.ContextValidator{
        govalidator.NonNullableValidator,
        govalidator.IsStringValidator,
    },
    Fields: &map[string]govalidator.Definition{
        "name": {
            Validator: []govalidator.ContextValidator{
                govalidator.NonNullableValidator,
                govalidator.IsStringValidator,
            },
        },
    },
}

// New Schema API
schema := govalidator.NewSchema().WithFields(
    govalidator.NewField("name").Required().WithValidators(govalidator.IsStringValidator),
)
```

---

## Legacy Definition API

> ⚠️ **Deprecated**: The Definition API is maintained for backward compatibility but is no longer recommended for new projects. Please use the Schema API instead.

<details>
<summary>Click to view legacy Definition API documentation</summary>

### Legacy Usage Example

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/gstachniukrsk/govalidator"
)

func getModel() govalidator.Definition {
    return govalidator.Definition{
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
            "gender": {
                Validator: []govalidator.ContextValidator{
                    govalidator.NonNullableValidator,
                    govalidator.IsStringValidator,
                    govalidator.OneOfValidator("male", "female"),
                },
            },
        },
    }
}

func main() {
    v := govalidator.NewBasicValidator(
        govalidator.PathPresenter("."),
        govalidator.SimpleErrorPresenter(),
    )

    jsonData := `{"name": "John", "age": 30, "gender": "male"}`
    var data any
    json.Unmarshal([]byte(jsonData), &data)

    valid, errs := v.Validate(context.Background(), data, getModel())

    if !valid {
        fmt.Printf("Errors: %v\n", errs)
    }
}
```

### Legacy Object Definition

```go
objModel := govalidator.Definition{
    Validator: []govalidator.ContextValidator{
        govalidator.NonNullableValidator,
        govalidator.IsMapValidator,
    },
    AcceptExtraProperty: true,
    AcceptNotDefinedProperty: true,
    Fields: &map[string]govalidator.Definition{
        "field1": {
            Validator: []govalidator.ContextValidator{
                govalidator.NonNullableValidator,
                govalidator.IsStringValidator,
            },
        },
    },
}
```

### Legacy Array Definition

```go
listModel := govalidator.Definition{
    Validator: []govalidator.ContextValidator{
        govalidator.NonNullableValidator,
        govalidator.IsListValidator,
    },
    ListOf: &govalidator.Definition{
        Validator: []govalidator.ContextValidator{
            govalidator.NonNullableValidator,
            govalidator.IsStringValidator,
        },
    },
}
```

</details>

---

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Quick Start for Contributors

```bash
# Clone the repository
git clone https://github.com/gstachniukrsk/govalidator.git
cd govalidator

# Set up commit message template (optional but recommended)
git config commit.template .gitmessage

# Run tests
go test -v ./...

# Run linters
golangci-lint run
```

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning:

- `feat:` - New feature (minor version bump)
- `fix:` - Bug fix (patch version bump)
- `feat!:` - Breaking change (major version bump)

See [RELEASE_QUICK_START.md](RELEASE_QUICK_START.md) or [RELEASING.md](RELEASING.md) for the complete release process.

## Versioning

This project follows [Semantic Versioning](https://semver.org/). Releases are automated based on commit messages when merging to `master`.

## Documentation

- [Migration Guide](MIGRATION_SCHEMA_API.md) - Migrate from Definition to Schema API
- [Release Guide](RELEASE_QUICK_START.md) - Quick release instructions
- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Releasing Guide](RELEASING.md) - Detailed release process

## License

[MIT License](LICENSE)
