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

The project uses **fully automated releases** based on conventional commits.

```bash
# Create a release (automatic on push to master)
git add .
git commit -m "feat: add new feature"  # Minor version bump (0.x.0)
git commit -m "fix: bug fix"           # Patch version bump (0.0.x)
git commit -m "feat!: breaking change" # Major version bump (x.0.0)
git push origin master
# Release created automatically in ~2-3 minutes

# Manual release (emergency only)
git tag v1.2.3
git push origin v1.2.3
gh release create v1.2.3 --generate-notes
```

**Commit Message Format** (see [RELEASE_QUICK_START.md](RELEASE_QUICK_START.md)):
- `feat:` - New feature (minor bump: 0.x.0)
- `fix:` - Bug fix (patch bump: 0.0.x)
- `perf:` - Performance improvement (patch bump: 0.0.x)
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
- Constraint errors: `MinSizeError`, `MaxSizeError`, `FloatPrecisionError`, `FloatTooSmallError`, `FloatTooLargeError`
- String errors: `StringTooShortError`, `StringTooLongError`
- Schema errors: `FieldNotDefinedError`, `UnexpectedFieldError`
- `RequiredError`: Special error indicating non-nullable constraint failure

**Twig Breaking**: Validators can return `twigBlock = true` to stop validation of the current branch. This is used by:
- `NonNullableValidator`: Stops validation if value is null (returns `RequiredError`)
- `NullableValidator`: Stops validation but accepts null values
- Type validators: Stop validation when type mismatch occurs

### Built-in Validators

Current validator naming convention uses `Is` prefix for type checks:

**Type Validators** (using `Is` prefix):
- `IsStringValidator`: Checks if value is string or *string (is_string.go)
- `IsIntegerValidator`: Checks if value is integer or float with no decimals (is_int.go)
- `IsBooleanValidator`: Checks if value is boolean (is_boolean.go)
- `IsMapValidator`: Checks if value is map[string]any (is_map.go)
- `IsListValidator`: Checks if value is []any (is_list.go)

**Numeric Validators**:
- `FloatValidator(maxPrecision)`: Validates float with precision limit (floatish.go)
- `MinFloatValidator(min)`: Minimum float value (min_float.go)
- `MaxFloatValidator(max)`: Maximum float value (max_float.go)
- `NumberValidator`: Basic number check (number.go)

**String Validators**:
- `MinLengthValidator(len)`: Minimum string length in runes (min_length.go)
- `MaxLengthValidator(len)`: Maximum string length in runes (max_length.go)
- `UpperCaseValidator`: String must be uppercase (upper_case.go)
- `LowerCaseValidator`: String must be lowercase (lower_case.go)
- `RegexpValidator(pattern)`: String must match regex (regex.go)

**Collection Validators**:
- `MinSizeValidator(size, blocks)`: Minimum collection size (min_size.go)
- `MaxSizeValidator(size, blocks)`: Maximum collection size (max_size.go)

**Value Validators**:
- `NonNullableValidator`: Rejects null values (non_nullable.go)
- `NullableValidator`: Accepts null values (nullable.go)
- `OneOfValidator(values...)`: Enum validation using reflect.DeepEqual (one_of.go)

### Creating Custom Validators

Implement the `ContextValidator` function signature:
```go
func CustomValidator(params) ContextValidator {
    return func(ctx context.Context, value any) (twigBreak bool, errs []error) {
        // Type assertion
        v, ok := value.(expectedType)
        if !ok {
            return true, []error{NotATypeError{}}
        }

        // Validation logic
        if !isValid(v) {
            return false, []error{YourCustomError{}}
        }

        // Return twigBreak=true to stop branch validation
        return false, nil
    }
}
```

Key patterns:
- Return `twigBreak=true` when type assertion fails or further validation is meaningless
- Use typed error structs from errs.go for consistency
- Accept configuration parameters via closure for reusable validators
- Use rune length (`len([]rune(str))`) for Unicode string validation

### Path Representation

The validation tree is represented as a string slice where:
- Root is `$`
- Object fields are appended as keys: `["$", "address", "city"]`
- List items use bracket notation: `["$", "phone", "[0]", "number"]`

Presenters convert this to user-friendly formats (e.g., `$.phone[0].number`).

### Error Presenters

The library provides several built-in error presenters for different use cases:

**Basic Presenters**:
- `SimpleErrorPresenter()`: Returns only the error message (simple_error_presenter.go)
- `PathPresenter(glue)`: Returns only the path with configurable separator (path_presenter.go)
- `RegistryPresenter()`: Routes different error types to different presenters (registry_presenter.go)

**Combined Presenters** (combined_presenter.go):
- `CombinedPresenter(pathGlue, separator)`: Combines path and error in one string
  - Example: `"$.user.age: not an integer"`
- `CombinedBracketPresenter(pathGlue, separator)`: Smart bracket handling for array paths
  - Example: `"$.items[0].name: required"`

**JSON Presenters** (json_presenter.go):
- `JSONPresenter(pathGlue)`: Formats errors as JSON with path and message fields
  - Example: `{"path":"$.age","message":"not an integer"}`
- `JSONDetailedPresenter(pathGlue)`: Includes structured error details (minSize, maxLength, etc.)
  - Example: `{"path":"$.name","message":"too short","type":"StringTooShortError","minLength":5}`

**User-Friendly Presenters** (detailed_error_presenter.go):
- `DetailedErrorPresenter()`: Human-readable messages for end users
  - Example: "text must be at least 5 characters long"
- `VerboseErrorPresenter()`: Technical messages with full error type information
  - Example: "StringTooShortError: minimum length is 5 characters"

**Flat List Collector** (flat_list_collector.go):
- `NewFlatListCollector(combiner)`: Collects errors as `[]string` instead of `map[string][]string`
- `NewFlatListValidator(combiner)`: Convenience validator that returns flat error lists

Example usage:
```go
// Combined presenter with detailed errors for user-facing validation
validator := NewBasicValidator(
    CombinedPresenter(".", ": "),
    DetailedErrorPresenter(),
)

// JSON output for REST APIs
jsonValidator := NewFlatListValidator(
    JSONDetailedPresenter("."),
)
valid, errs := jsonValidator.Validate(ctx, data, definition)
// errs is []string with JSON-formatted errors

// Map-based errors with custom path format
collector := NewPathToErrCollector(
    PathPresenter("/"),  // Use / instead of .
    VerboseErrorPresenter(),
)
```

## Code Quality

### Linting Configuration

- **golangci-lint**: Configured in `.golangci.yaml` with 30+ linters enabled
- **revive**: Configured in `revive.toml` for additional code style checks
- Both run automatically in CI via `.github/workflows/lint.yml`

### Testing Best Practices

- Use table-driven tests for validators
- Test both success and failure cases
- Test edge cases: nil, empty, boundary values, Unicode strings
- Aim for >80% code coverage
- Place tests in `*_test.go` files with same name as implementation

## Release Automation

Releases are fully automated based on conventional commits:
- Every push to `master` triggers `.github/workflows/release.yml`
- Tests and linters must pass before release
- Version is calculated from commit messages
- CHANGELOG.md is automatically updated
- Git tag and GitHub release are created automatically

See [RELEASE_QUICK_START.md](RELEASE_QUICK_START.md) for a quick reference guide.

## Documentation Files

- **README.md**: User-facing documentation with examples
- **RELEASING.md**: Complete release process documentation
- **RELEASE_QUICK_START.md**: Quick reference for releases
- **CONTRIBUTING.md**: Guidelines for contributors
- **CHANGELOG.md**: Auto-generated from commits
- **.gitmessage**: Commit message template (activate with `git config commit.template .gitmessage`)
