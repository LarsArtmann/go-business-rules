# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> **Versioning note (2026-07-26 audit):** the only git tag on this repository is
> `v0.1.0` (2026-05-05, commit `d9faacb`). The `[1.0.0]` / `[1.1.0]` entries below
> predate that tag and document the pre-release evolution of the library. The
> changes below `[Unreleased]` are everything that landed after `v0.1.0` and have
> not yet been given a release tag. `doc.go` still reports `Version = "1.1.0"`,
> which is tracked as a split-brain in `TODO_LIST.md`.

## [Unreleased]

All changes below are verified against `master` (2026-07-26) and are not yet
released under a git tag.

### Breaking Changes

- **Type renames** (commit `1f2976d`): `Violation` → `ViolationError`,
  `ValidationResult` → `ValidationResultError`. The new names reflect that these
  types represent validation failures and implement the `error` interface.
  Update all references.
- **`Severity` migrated to `finding.Severity`** (commit `e423de4`): `Severity` is
  now a type alias for `github.com/larsartmann/go-finding.Severity` (a `string`,
  not an `int`/`iota`). Severity constants are re-exported from `go-finding`.
  This makes `go-finding` a **direct runtime dependency** (the library is no
  longer zero-runtime-dependency). Severity comparisons must use
  `.GreaterThanOrEqual()` etc. instead of numeric operators.
- **`encoding/json/v2`** (commit `4b93f12`): JSON marshaling in
  `ViolationError.MarshalJSON` and `ValidationResultError.MarshalJSON` now uses
  `encoding/json/v2`. Building or consuming this library requires
  `GOEXPERIMENT=jsonv2` with Go 1.26+. This is a hard constraint on downstream
  consumers until `json/v2` graduates from experimental.

### Added

- Extended numeric rules: `GreaterThan`, `LessThan` (`builders_collection.go`)
- Collection rules: `NotEmptySlice[T]`, `NotEmptyMap[T]` (`builders_collection.go`)
- String rule: `NotBlank` (`builders.go`)
- Generic rule: `Equals[T]` (`builders.go`)
- Generic threshold engine `thresholdCheck` and `numericCheck` helpers reducing builder duplication
- BDD branching-flow regression suite (`bdd_branching_flow_test.go`)
- Nix flake with `default` and `ci` devShells (`flake.nix`), setting `GOEXPERIMENT=jsonv2`

### Changed

- Bumped Go to 1.26.4 (`go.mod`)
- GitHub Actions hardened: all actions pinned to commit SHAs, `GOEXPERIMENT=jsonv2`
  set in CI, cancel-in-progress, paths-ignore, timeout-minutes (`.github/workflows/ci.yml`)
- golangci-lint v2 configuration with documented false positives

### Documentation

- `FEATURES.md` created (honest feature inventory by status)
- `ROADMAP.md` pruned of now-shipped items

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

[Unreleased]: https://github.com/LarsArtmann/go-business-rules/compare/v0.1.0...HEAD
[1.0.0]: https://github.com/LarsArtmann/go-business-rules/releases/tag/v1.0.0
[1.1.0]: https://github.com/LarsArtmann/go-business-rules/releases/tag/v1.1.0
