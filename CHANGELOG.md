# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> **Versioning note (2026-07-26 release):** the only git tag on this repository
> was `v0.1.0` (2026-05-05, commit `d9faacb`). The `[1.0.0]` / `[1.1.0]` entries
> below predate that tag and document the pre-release evolution of the library;
> they were never individually tagged. `v2.0.0` consolidates all post-`v0.1.0`
> work — including that evolution and the breaking changes below — into a single
> tagged release. A major bump is warranted by the breaking changes (type renames,
> the `Severity` migration to `finding.Severity`, and the `encoding/json/v2`
> requirement). `doc.go` now reports `Version = "2.0.0"`, resolving the
> split-brain previously tracked in `TODO_LIST.md`.

## [Unreleased]

### Added

- **String builders**: `Contains` (substring), `LengthRange` (bounded length),
  `Required` (visible content — `NotEmpty` + `NotBlank` semantics), and
  `MatchesFunc` (custom predicate).
- **Time builders** (`builders_time.go`): `NotPast`, `NotFuture`, and
  `DateInRange` with an injectable "now" so rules stay deterministic and
  testable.
- **Network/ID builders** (`builders_network.go`): `IPAddress` (IPv4/IPv6 via
  `net/netip`), `CreditCard` (13-19 digits + Luhn checksum), `PhoneNumber`
  (generic E.164-ish, separator-tolerant), and `PostalCode` (generic 3-10
  character pattern; country-specific patterns deliberately deferred).
- **Precision builders**: `MaxDecimalPlaces` (rejects float error
  accumulation like `0.30000000000000004` for money fields; non-finite fails)
  and `DivisibleBy` (zero divisor fails the check instead of panicking).
- **Composite builders**: `Not` (negation), `Or` (variadic alternative to
  `Any`), and `Xor` (exactly-one-passes).
- **Rule metadata**: `WithDescription` / `WithTags` on `RuleImpl` and
  `ContextRuleImpl` (the `Rule` interface is unchanged, so no breaking
  change). Metadata is surfaced on `RuleEvaluated` events (cloned, so
  listeners cannot mutate rule state) and costs nothing on the no-listener
  path.
- **Property-based invariant test** (`property_test.go`): gopter verifies
  `violations(Build) ≡ violations(Stream)` over generated rule sets. Opt-in:
  `GBR_PROPERTY=1 go test ./...` (skipped by default).
- **Hygiene**: `SECURITY.md`, `.github/CODEOWNERS`, issue and PR templates.
- **ADRs**: `docs/adr/0002_type_renames.md`, `0003_finding_severity_alias.md`,
  `0004_json_v2_adoption.md` capture the standing type-rename,
  `finding.Severity` alias, and json/v2 decisions.

### Changed

- Branching-flow analyzer pins re-pinned to post-builder reality (25 PHANTOM;
  29 stats total). The `DivisibleBy` zero-guard carries a
  `//nolint:branching-flow:panic` suppression (experimental analyzer does not
  track the guard); the `.golangci.yml` nolintlint exclusion now covers both
  files carrying that directive class.

## [2.1.0] - 2026-09-14

### Breaking Changes

- **Module path migrated to `/v2`**: the module is now
  `github.com/LarsArtmann/go-business-rules/v2`. This fixes the unresolvable
  `v2.0.0` git tag (a v2+ tag requires a `/v2` module-path suffix; that tag
  only ever worked as a `+incompatible` version of the suffix-less path and
  stays where it is). Consumers: `go get
  github.com/LarsArtmann/go-business-rules/v2@v2.1.0` and update import
  paths. `doc.go` reports `Version = "2.1.0"`.

### Added

- **Validation events**: register listeners via `ValidatorBuilder.WithListener` to
  observe every rule check (`RuleEvaluated`, passing and failing, with duration and
  start time) plus a terminal `ValidationCompleted` carrying the returned result.
  Delivery is synchronous, in registration order, before `Build` returns. Without
  listeners the build path is unchanged: no timing, no allocations for events.
- **Concurrent streaming evaluation**: `ValidatorBuilder.Stream(ctx)` runs all rules
  concurrently and returns a channel of the same events in completion order, closed
  after the terminal event. Violations in the terminal event are re-sorted into the
  original rule order so the aggregated result stays deterministic. Cancellation
  skips unstarted rules and never leaks goroutines.
- **`WithConcurrency(n)`**: caps how many rule checks `Stream` runs at once (one
  goroutine per rule remains the default). Values below 1 keep the unbounded
  default.
- **Context-aware rules**: the additive `ContextRule` interface (`CheckContext(ctx)`)
  and the `NewContextRule` builder let `Stream(ctx)` interrupt slow, I/O-bound
  checks on cancellation instead of only skipping unstarted ones. Plain rules keep
  working unchanged (`Stream` falls back to `Check`; `Build` is untouched).
- `listeners/otel` (nested module): a `businessrules.Listener` mapping validation
  events onto OpenTelemetry spans and metrics (`businessrules.rule.evaluations`,
  `businessrules.rule.duration`, `businessrules.validations`,
  `businessrules.validation.duration`, `businessrules.violations`).
- `examples/sse` now renders events as Datastar patches: validation results stream
  to the browser as `datastar-patch-elements` / `datastar-patch-signals` values
  broadcast through the same `go-sse` Broadcaster.
- Stream benchmarks (`BenchmarkStream2Rules`, `BenchmarkStream10Rules`,
  `BenchmarkStream10RulesConcurrency4`) alongside the existing Build-path
  benchmarks.
- Goroutine-leak regression tests for `Stream` (abandoned-after-cancel and fully
  drained), pinning the buffered-channel contract.
- `check-all` flake app (`nix run .#check-all`): build + vet + test + lint across
  all Go modules (root, `adapters/cqrslite`, `examples/sse`, `listeners/otel`).

### Fixed

- Branching-flow regression pins re-pinned to current analyzer reality (18
  PHANTOM false positives, 22 stats total); the interprocedural
  panic-analyzer false positive on the `Stream` result send is suppressed with
  the analyzer's own `//nolint:branching-flow:panic` mechanism.
- `examples/sse` smoke tests now check `bufio.Scanner.Err()`.
- Listener-overhead benchmarks: 189 ns/op with no listener (unchanged default path),
  468 ns/op with one listener, flat with additional listeners.
- `adapters/cqrslite` (nested module): publishes validation events onto a
  `go-cqrs-lite` event bus as durable domain events with adapter-owned, CBOR-safe
  wire DTOs.
- `examples/sse` (nested module): live browser feed of validation events over
  Server-Sent Events.

## [2.0.0] - 2026-07-26

All changes below are verified against `master` (2026-07-26) and consolidated
into the first properly tagged release since `v0.1.0`. See the versioning note
above.

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

### Fixed

- Hardened the branching-flow stats regression test (`bdd_branching_flow_test.go`):
  it now parses the `totalIssues` field from `stats --format json` instead of
  matching a fragile `"36"` substring that a coincidental duration/count could pass.
- Reconciled a stale `go.sum` checksum for `github.com/larsartmann/go-error-family@v0.10.0`
  after that module was re-released at the same version.
- Resolved the `Version` constant split-brain (`doc.go` now reports `2.0.0`).

### Documentation

- `docs/DOMAIN_LANGUAGE.md` filled in with the real ubiquitous language (Rule,
  Severity, ViolationError, ValidationResultError, ValidatorBuilder, and operations)
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

[Unreleased]: https://github.com/LarsArtmann/go-business-rules/compare/v2.1.0...HEAD
[2.1.0]: https://github.com/LarsArtmann/go-business-rules/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/LarsArtmann/go-business-rules/releases/tag/v2.0.0

<!-- [1.0.0] and [1.1.0] were never tagged; no release URLs exist for them. See the versioning note above. -->
