# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library for validating JSON data unmarshalled to `any` type. The library provides a declarative validation framework using a Definition-based approach with composable validators.

## Common Commands

### Testing
```bash
# Run all tests
go test -v ./...

# Run a specific test
go test -v -run TestName

# Run tests with coverage
go test -v -cover ./...
```

### Building
```bash
# Build the library
go build -v ./...

# Verify the module
go mod verify
go mod tidy
```

### Linting
```bash
# Run golangci-lint
golangci-lint run

# Run revive
revive -config revive.toml ./...

# Run all linters (via GitHub Actions locally)
act -j code-styling
act -j revive
```

### Releasing

The project uses automated semantic versioning based on conventional commits.

```bash
# Create a release (via conventional commit to master)
git add .
git commit -m "feat: add new feature"  # Minor version bump
git commit -m "fix: bug fix"           # Patch version bump
git commit -m "feat!: breaking change" # Major version bump
git push origin master

# Manual release (emergency only)
git tag v1.2.3
git push origin v1.2.3
gh release create v1.2.3 --generate-notes

# View release workflows
cat RELEASING.md  # Complete release documentation
```

**Commit Message Format:**
- `feat:` - New feature (minor bump: 0.x.0)
- `fix:` - Bug fix (patch bump: 0.0.x)
- `feat!:` or `BREAKING CHANGE:` - Breaking change (major bump: x.0.0)
- `docs:`, `style:`, `refactor:`, `test:`, `chore:` - No version bump

See [RELEASING.md](RELEASING.md) for complete release process documentation.

## Architecture

### Core Concepts

**Definition-Based Validation**: Validation rules are declared as `Definition` structs that describe the expected shape and constraints of data. Definitions are composable and can be extended using the `ExtendedWith` method.

**ContextValidator**: The fundamental building block - a function type that validates a value and returns:
- `twigBlock bool`: Whether to stop validation of the current branch
- `errs []error`: List of validation errors

**Validator Pattern**: Two-level validator interface:
- `Validator`: Low-level interface that walks the tree structure (internal use)
- `BasicValidator`: High-level interface for simple validation calls (public API)

### Key Components

**validator.go**: Contains the core validation engine
- `validator` struct: Internal validator that walks the validation tree, maintains path context (`currentTree`), and manages error collectors
- `basicValidator` struct: Public-facing validator interface
- Tree traversal logic handles both object fields (via `Fields`) and list items (via `ListOf`)
- Fields are always validated in sorted order for consistent error reporting

**definition.go**: Defines validation schemas
- `Validator`: Slice of `ContextValidator` functions applied sequentially
- `Fields`: For object/map validation - defines expected properties
- `ListOf`: For array validation - defines the schema for list items
- `AcceptNotDefinedProperty`: Controls whether missing fields are allowed
- `AcceptExtraProperty`: Controls whether extra fields are allowed

**ctx_validator.go**: Core validator function type
- `AcceptsNull()` method checks if a validator accepts null by testing with nil and checking for `RequiredError`

**extend.go**: Provides `ExtendedWith()` method for merging definitions
- Combines validators from both definitions
- Merges field definitions recursively
- Handles ListOf merging

**Collector Pattern**: Error collection is abstracted via `Collector` function type
- `PathToErrCollector`: Aggregates errors by path using customizable presenters
- Presenters format paths (e.g., dot notation: `$.address.city`) and error messages
- Error collectors receive: context, path (as string slice), and error

### Error Handling

**Structured Errors** (errs.go): All validation errors are typed structs implementing `error`
- Type errors: `NotAStringError`, `NotAnIntegerError`, `NotAMapError`, `NotAListError`, `NotABooleanError`
- Constraint errors: `MinSizeError`, `MaxSizeError`, `FloatPrecisionError`
- Schema errors: `FieldNotDefinedError`, `UnexpectedFieldError`
- `RequiredError`: Special error indicating non-nullable constraint failure

**Twig Breaking**: Validators can return `twigBlock = true` to stop validation of the current branch. This is used by:
- `NonNullableValidator`: Stops validation if value is null (returns `RequiredError`)
- `NullableValidator`: Stops validation but accepts null values
- Type validators: Stop validation when type mismatch occurs

### Built-in Validators

Located in individual files, categorized by purpose:
- **Type checks**: `is_string.go`, `is_int.go`, `is_boolean.go`, `is_map.go`, `is_list.go`
- **Null handling**: `nullable.go`, `non_nullable.go`
- **String constraints**: `min_length.go`, `max_length.go`, `upper_case.go`, `lower_case.go`, `regex.go`
- **Numeric constraints**: `min_float.go`, `max_float.go`, `floatish.go` (validates float precision), `number.go`
- **Collection constraints**: `min_size.go`, `max_size.go`
- **Value constraints**: `one_of.go` (enum validation using `reflect.DeepEqual`)

### Creating Custom Validators

Implement the `ContextValidator` function signature:
```go
func CustomValidator(params) ContextValidator {
    return func(ctx context.Context, value any) (twigBreak bool, errs []error) {
        // Type assertion
        // Validation logic
        // Return twigBreak=true to stop branch validation
        // Return errors slice
    }
}
```

Key patterns:
- Return `twigBreak=true` when type assertion fails or further validation is meaningless
- Use typed error structs from errs.go for consistency
- Accept configuration parameters via closure for reusable validators

### Path Representation

The validation tree is represented as a string slice where:
- Root is `$`
- Object fields are appended as keys: `["$", "address", "city"]`
- List items use bracket notation: `["$", "phone", "[0]", "number"]`

Presenters convert this to user-friendly formats (e.g., `$.phone[0].number`).

### Error Presenters

The library provides several built-in error presenters for different use cases:

**Basic Presenters**:
- `SimpleErrorPresenter()`: Returns only the error message
- `PathPresenter(glue)`: Returns only the path with configurable separator
- `RegistryPresenter()`: Routes different error types to different presenters

**Combined Presenters** (combined_presenter.go):
- `CombinedPresenter(pathGlue, separator)`: Combines path and error in one string (e.g., "$.age: not an integer")
- `CombinedBracketPresenter(pathGlue, separator)`: Smart bracket handling for array paths

**JSON Presenters** (json_presenter.go):
- `JSONPresenter(pathGlue)`: Formats errors as JSON with path and message fields
- `JSONDetailedPresenter(pathGlue)`: Includes structured error details (minSize, maxLength, etc.)

**User-Friendly Presenters** (detailed_error_presenter.go):
- `DetailedErrorPresenter()`: Human-readable messages ("text must be at least 5 characters")
- `VerboseErrorPresenter()`: Technical messages with full error type information

**Flat List Collector** (flat_list_collector.go):
- `NewFlatListCollector(combiner)`: Collects errors as []string instead of map[string][]string
- `NewFlatListValidator(combiner)`: Convenience validator that returns flat error lists

Example usage:
```go
// Combined presenter with custom formatting
validator := NewBasicValidator(
    CombinedPresenter(".", ": "),
    DetailedErrorPresenter(),
)

// JSON output for APIs
jsonValidator := NewFlatListValidator(
    JSONDetailedPresenter("."),
)
valid, errs := jsonValidator.Validate(ctx, data, definition)
// errs is []string with JSON-formatted errors
```
