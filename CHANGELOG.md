# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CombinedPresenter for formatting errors as "path: message"
- CombinedBracketPresenter with smart array index handling
- JSONPresenter for API-friendly error formatting
- JSONDetailedPresenter with structured error field extraction
- DetailedErrorPresenter for human-readable error messages
- VerboseErrorPresenter for technical error details
- FlatListCollector for collecting errors as string slices
- FlatListValidator convenience interface
- Comprehensive linting setup with golangci-lint and revive
- GitHub Actions workflow for automated linting
- GitHub Actions workflows for automated releases (release-please and semantic-release)
- RELEASING.md with complete release process documentation

### Changed
- **BREAKING**: Renamed `IntValidator` to `IsIntegerValidator` for consistency
- **BREAKING**: Renamed `StringValidator` to `IsStringValidator` for consistency
- **BREAKING**: Renamed `FloatishValidator` to `FloatValidator` (removed informal naming)
- **BREAKING**: Renamed `FloatIsLesserThanError` to `FloatTooSmallError` (grammar fix)
- **BREAKING**: Renamed `FloatIsGreaterThanError` to `FloatTooLargeError` (consistency)
- Upgraded Go version from 1.19 to 1.24
- Improved type switch statements to use variable binding (is_map.go, is_list.go, nullable.go, non_nullable.go)
- Updated GitHub Actions workflow to use Go 1.24

### Fixed
- **CRITICAL**: MinLengthValidator now correctly uses rune length instead of byte length for Unicode strings
- Fixed typo in README.md: "LowerCasValidator" → "LowerCaseValidator"
- Corrected grammar in error type name: "Lesser" → "Small"

## [0.1.0] - YYYY-MM-DD

### Added
- Initial release of govalidator library
- Core validation engine with Definition-based approach
- Context-aware validator functions
- Path-based error collection
- Built-in validators:
  - Type validators (IsBoolean, IsList, IsMap, etc.)
  - String validators (MinLength, MaxLength, UpperCase, LowerCase, Regexp)
  - Numeric validators (MinFloat, MaxFloat, Float precision)
  - Collection validators (MinSize, MaxSize)
  - Value validators (OneOf, Nullable, NonNullable)
- Error presenter system with PathPresenter and SimpleErrorPresenter
- RegistryPresenter for custom error formatting
- Comprehensive test coverage
- README with usage examples

[Unreleased]: https://github.com/gstachniukrsk/govalidator/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/gstachniukrsk/govalidator/releases/tag/v0.1.0
