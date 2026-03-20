# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-03-15

### Breaking Changes

- Renamed `Result` type to `ValidationResult` for improved clarity
- Removed `Result` type alias (previously provided for backwards compatibility)
- Update any code referencing `Result` to use `ValidationResult` instead

### Added

- `FirstError()` method on `ValidationResult` - returns first Error/Critical violation
- `FirstCritical()` method on `ValidationResult` - returns first Critical violation
- `FirstWarning()` method on `ValidationResult` - returns first Warning violation
- `FirstInfo()` method on `ValidationResult` - returns first Info violation
- `HasCritical()` method on `ValidationResult` - checks for Critical violations
- `HasInfo()` method on `ValidationResult` - checks for Info violations
- `Filter(predicate func(Violation) bool)` method on `ValidationResult` - custom filtering
- `Count()` method on `ValidationResult` - returns total violation count
- `Merge(other ValidationResult)` method on `ValidationResult` - combines results

### Changed

- CI now tests Go 1.23, 1.24, 1.25 (removed 1.22, added 1.25)

## [1.0.0] - 2026-03-15

### Added

#### Core Types

- `Severity` type with four levels: Info, Warning, Error, Critical
- `Rule` interface for defining validation rules
- `Violation` struct representing failed rule checks
- `Result` struct with severity filtering methods
- `ValidatorBuilder` with fluent API for building validators
- `Version` constant for library version tracking

#### Numeric Rules

- `NonNegative(name, value, severity)` - validates value >= 0
- `Positive(name, value, severity)` - validates value > 0
- `InRange(name, value, min, max, severity)` - validates min <= value <= max
- `MinInt(name, value, min, severity)` - validates value >= min
- `MaxInt(name, value, max, severity)` - validates value <= max

#### String Rules

- `NotEmpty(name, value, severity)` - validates string is not empty
- `MinLength(name, value, min, severity)` - validates minimum string length
- `MaxLength(name, value, max, severity)` - validates maximum string length
- `Matches(name, value, pattern, severity)` - validates regex pattern match

#### Format Rules

- `Email(name, value, severity)` - validates RFC 5322 email addresses
- `URL(name, value, severity)` - validates HTTP/HTTPS URLs
- `UUID(name, value, severity)` - validates UUID format

#### Generic Rules

- `OneOf[T](name, value, allowed, severity)` - validates value is in allowed set
- `Custom(name, check, severity)` - custom validation function

#### Composite Rules

- `All(name, rules, severity)` - all sub-rules must pass
- `Any(name, rules, severity)` - at least one sub-rule must pass
- `When(name, condition, rule)` - conditional rule execution

#### Documentation

- Comprehensive GoDoc comments on all exported types
- Package-level documentation with usage examples
- Example tests for all major functions
- README with API documentation and quick start guide
- MIT LICENSE

#### Testing

- 33 tests with 93.9% code coverage
- Example tests for godoc documentation
- BDD-style test suite using Ginkgo/Gomega

#### CI/CD

- GitHub Actions workflow for CI
- Multi-version Go testing (1.22-1.25)
- golangci-lint integration
- Security scanning with Gosec
- Code coverage reporting with Codecov

### Technical Details

- Zero external runtime dependencies (only testing dependencies)
- Thread-safe immutable rule instances
- Generic rule support via Go 1.18+ generics
- Compatible with Go 1.22+

[1.0.0]: https://github.com/LarsArtmann/go-business-rules/releases/tag/v1.0.0
