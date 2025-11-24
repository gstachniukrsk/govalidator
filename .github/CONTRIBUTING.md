# Contributing to govalidator

Thank you for your interest in contributing to govalidator! This document provides guidelines and instructions for contributing.

## Development Setup

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/govalidator.git
   cd govalidator
   ```

2. **Ensure you have Go 1.24+ installed:**
   ```bash
   go version  # Should be 1.24 or higher
   ```

3. **Install development tools:**
   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

   # Install revive
   go install github.com/mgechev/revive@latest
   ```

4. **Run tests to verify setup:**
   ```bash
   go test -v ./...
   ```

## Making Changes

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Your Changes

- Write clear, idiomatic Go code
- Follow existing code style and patterns
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v -run TestYourTest

# Run linters
golangci-lint run
revive -config revive.toml ./...
```

### 4. Commit Your Changes

**IMPORTANT:** We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning and changelog generation.

#### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Types

- `feat:` - A new feature (triggers minor version bump)
- `fix:` - A bug fix (triggers patch version bump)
- `perf:` - Performance improvement (triggers patch version bump)
- `docs:` - Documentation only changes
- `style:` - Code style changes (formatting, missing semi colons, etc)
- `refactor:` - Code change that neither fixes a bug nor adds a feature
- `test:` - Adding missing tests or correcting existing tests
- `chore:` - Changes to build process or auxiliary tools

#### Breaking Changes

For breaking changes, add `!` after the type or include `BREAKING CHANGE:` in the footer:

```bash
feat!: rename IntValidator to IsIntegerValidator

BREAKING CHANGE: IntValidator has been renamed to IsIntegerValidator
for consistency with other type validators.
```

#### Examples

**Good commit messages:**
```bash
git commit -m "feat: add JSONPresenter for API error formatting"

git commit -m "fix: correct Unicode character counting in MinLengthValidator"

git commit -m "docs: add examples for new error presenters"

git commit -m "refactor: improve type switch statements in validators"

git commit -m "test: add test cases for edge cases in nullable validator"

git commit -m "feat!: rename validators for consistency

BREAKING CHANGE: The following validators have been renamed:
- IntValidator -> IsIntegerValidator
- StringValidator -> IsStringValidator"
```

**Bad commit messages:**
```bash
git commit -m "update stuff"           # Too vague
git commit -m "WIP"                    # Not descriptive
git commit -m "Fixed bug"              # Should be "fix: <description>"
git commit -m "Added new feature"      # Should be "feat: <description>"
```

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Pull Request Guidelines

### Before Submitting

- [ ] Tests pass locally (`go test -v ./...`)
- [ ] Linters pass (`golangci-lint run`)
- [ ] Code follows existing style
- [ ] Documentation is updated
- [ ] Commit messages follow conventional commits format

### PR Title

Use the same conventional commit format for your PR title:

```
feat: add JSONPresenter for structured error output
fix: correct rune counting in MinLengthValidator
```

### PR Description

Include:
- What changes were made
- Why the changes were made
- How to test the changes
- Any breaking changes
- Related issues (if applicable)

Example:
```markdown
## Description
Adds a new JSONPresenter that formats validation errors as JSON strings,
making it easier to integrate with REST APIs.

## Changes
- Added `JSONPresenter()` function
- Added `JSONDetailedPresenter()` with structured field extraction
- Added comprehensive tests
- Updated documentation

## Testing
Run: `go test -v -run TestJSONPresenter`

## Breaking Changes
None

Fixes #123
```

## Code Style Guidelines

### General Go Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting (done automatically by golangci-lint)

### Project-Specific Guidelines

1. **Naming Conventions:**
   - Validators: `IsTypeValidator` or `ConstraintValidator` pattern
   - Error types: Descriptive names ending with `Error` (e.g., `NotAStringError`, `StringTooShortError`)
   - Presenters: Descriptive names ending with `Presenter` (e.g., `JSONPresenter`)

2. **File Organization:**
   - One validator per file: `validator_name.go`
   - Corresponding test file: `validator_name_test.go`
   - Related types in the same file as the validator

3. **Documentation:**
   - All exported functions, types, and constants must have doc comments
   - Doc comments start with the name of the item
   - Include usage examples for complex functionality

4. **Testing:**
   - Use table-driven tests
   - Test both success and failure cases
   - Test edge cases (nil, empty, boundary values)
   - Aim for >80% code coverage

## Validator Implementation Pattern

When creating a new validator, follow this pattern:

```go
package govalidator

import "context"

// YourError describes what went wrong.
type YourError struct {
    Field string
    Value any
}

func (e YourError) Error() string {
    return "descriptive error message"
}

// YourValidator checks if the value meets your criteria.
// It returns twigBlock=true to stop validation of the current branch.
func YourValidator(param string) ContextValidator {
    return func(_ context.Context, value any) (twigBlock bool, errs []error) {
        // Type assertion
        v, ok := value.(expectedType)
        if !ok {
            return true, []error{NotATypeError{}}
        }

        // Validation logic
        if !isValid(v) {
            return false, []error{YourError{Value: v}}
        }

        return false, nil
    }
}
```

## Presenter Implementation Pattern

When creating a new presenter:

```go
package govalidator

import "context"

// YourPresenter creates a presenter that formats errors in your custom way.
// Explain the output format and provide examples.
func YourPresenter(params) PresenterFunc {
    return func(ctx context.Context, path []string, err error) string {
        // Format the path
        pathStr := formatPath(path)

        // Format the error
        errStr := formatError(err)

        // Combine and return
        return pathStr + errStr
    }
}
```

## Getting Help

- Check existing issues and PRs
- Read the documentation: [README.md](../README.md), [CLAUDE.md](../CLAUDE.md)
- Ask questions in issue comments
- Join discussions in pull requests

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

Thank you for contributing to govalidator! 🎉
