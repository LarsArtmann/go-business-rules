# Features — businessrules

> Honest feature inventory by status. Verified against code, not trust.
>
> **Status legend:** `FULLY_FUNCTIONAL` (works, tested) · `PARTIALLY_FUNCTIONAL` (ships with known gaps) · `BROKEN` (exists but fails) · `PLANNED` (no code yet).
>
> **Verified:** 2026-09-18 against `master` (251/251 Ginkgo specs pass incl. branching-flow pins and the opt-in gopter property test, coverage 97.1% re-measured via `go test -cover`).

---

## Core Framework

| Feature                                                                                | Status           | Evidence                                                                           |
| -------------------------------------------------------------------------------------- | ---------------- | ---------------------------------------------------------------------------------- |
| `Severity` type with 4 levels (Info / Warning / Error / Critical)                      | FULLY_FUNCTIONAL | `severity.go` — type alias to `finding.Severity` (`string`); constants re-exported |
| `Rule` interface (Name / Check / Severity / Message)                                   | FULLY_FUNCTIONAL | `rule.go:6`                                                                        |
| `RuleImpl` base implementation + immutable `WithName` / `WithSeverity` / `WithMessage` | FULLY_FUNCTIONAL | `rule.go:24-57`                                                                    |
| `NewRule` constructor                                                                  | FULLY_FUNCTIONAL | `rule.go:61`                                                                       |
| `ViolationError` (failed check + context + timestamp), implements `error`              | FULLY_FUNCTIONAL | `errors.go:12`                                                                     |
| `NewViolation` / `NewViolationFromError` / `WithContext`                               | FULLY_FUNCTIONAL | `errors.go:41,51,61`                                                               |
| `ValidationResultError` with severity filtering & first-violation access               | FULLY_FUNCTIONAL | `validation_result.go:10`                                                          |
| `ValidatorBuilder` fluent API (`NewValidator` / `AddRule` / `AddRules` / `Build`)      | FULLY_FUNCTIONAL | `validator.go:5-48`                                                                |
| `Version` constant                                                                     | FULLY_FUNCTIONAL | `doc.go:135` reports `"2.2.0"`, reconciled with the `v2.2.0` release               |

### `ValidationResultError` methods

| Method group                                                                        | Status           | Evidence                                             |
| ----------------------------------------------------------------------------------- | ---------------- | ---------------------------------------------------- |
| Filter by severity: `Errors` / `Warnings` / `Info` / `Critical` / `BySeverity`      | FULLY_FUNCTIONAL | `validation_result.go:19-54`                         |
| Existence checks: `HasErrors` / `HasWarnings` / `HasCritical` / `HasInfo` / `Count` | FULLY_FUNCTIONAL | `validation_result.go:57-79`                         |
| First violation: `FirstError` / `FirstCritical` / `FirstWarning` / `FirstInfo`      | FULLY_FUNCTIONAL | `validation_result.go:83-123`                        |
| `Filter(predicate)`, `ForEach(fn)`, `Merge(other)`                                  | FULLY_FUNCTIONAL | `validation_result.go:126-157`                       |
| JSON marshaling (`MarshalJSON`)                                                     | FULLY_FUNCTIONAL | `validation_result.go:160` (uses `encoding/json/v2`) |
| `Error()` (implements `error`)                                                      | FULLY_FUNCTIONAL | `validation_result.go:177`                           |

## Pre-built Rule Builders

| Category           | Rules                                                                                                             | Status           | Evidence                                         |
| ------------------ | ----------------------------------------------------------------------------------------------------------------- | ---------------- | ------------------------------------------------ |
| Numeric            | `NonNegative`, `Positive`, `InRange`, `MinInt`, `MaxInt`                                                          | FULLY_FUNCTIONAL | `builders.go:79-134`                             |
| Numeric (extended) | `GreaterThan`, `LessThan`                                                                                         | FULLY_FUNCTIONAL | `builders_collection.go:39-57`                   |
| String             | `NotEmpty`, `NotBlank`, `Required`, `MinLength`, `MaxLength`, `LengthRange`, `Matches`, `Contains`, `MatchesFunc` | FULLY_FUNCTIONAL | `builders.go:140-233,241-310`                    |
| Collection         | `NotEmptySlice[T]`, `NotEmptyMap[T]`                                                                              | FULLY_FUNCTIONAL | `builders_collection.go:25-33`                   |
| Format             | `Email`, `URL`, `UUID`                                                                                            | FULLY_FUNCTIONAL | `builders_format.go:27,49,84`                    |
| Time / Date        | `NotPast`, `NotFuture`, `DateInRange` (injectable clock)                                                          | FULLY_FUNCTIONAL | `builders_time.go:23,45,67`                      |
| Network / ID       | `IPAddress`, `CreditCard` (Luhn), `PhoneNumber`, `PostalCode`                                                     | FULLY_FUNCTIONAL | `builders_network.go:22,37,86,117`               |
| Precision          | `MaxDecimalPlaces`, `DivisibleBy`                                                                                 | FULLY_FUNCTIONAL | `builders_collection.go:100,133`                 |
| Generic            | `Equals[T]`, `OneOf[T]`, `Custom`                                                                                 | FULLY_FUNCTIONAL | `builders.go:239`, `builders_composite.go:10,27` |
| Composite          | `All`, `Any`, `When`, `Not`, `Or`, `Xor`                                                                          | FULLY_FUNCTIONAL | `builders_composite.go:64,76,98,114,130,143`     |
| Rule metadata      | `WithDescription` / `WithTags` on `RuleImpl` + `ContextRuleImpl`, surfaced on `RuleEvaluated`                     | FULLY_FUNCTIONAL | `rule.go`, `events.go`, `validator.go`           |

## Observability & Events

| Feature                                                                                                  | Status           | Evidence                                                                                             |
| -------------------------------------------------------------------------------------------------------- | ---------------- | ---------------------------------------------------------------------------------------------------- |
| `Event` sealed interface with `RuleEvaluated` (passes included, timing, error) and `ValidationCompleted` | FULLY_FUNCTIONAL | `events.go:10`, `events.go:21`, `events.go:49`                                                       |
| `Listener` type + `WithListener(...)` builder option (synchronous, registration-order delivery)          | FULLY_FUNCTIONAL | `events.go:66`, `validator.go:28`; specs in `events_test.go`                                         |
| Zero-cost default path (no listeners → no timing, no events, unchanged result)                           | FULLY_FUNCTIONAL | `validator.go:40-45`; `BenchmarkValidatorNoListener` 189 ns/op vs one listener 468 ns/op             |
| Derived pass/fail (`Passed()` = `Err == nil`, impossible states unrepresentable)                         | FULLY_FUNCTIONAL | `events.go:41`                                                                                       |
| `Stream(ctx)` concurrent evaluation, completion-order events, deterministic final result                 | FULLY_FUNCTIONAL | `validator.go`; specs in `stream_test.go`                                                            |
| `WithConcurrency(n)` bounded in-flight checks (default unbounded)                                        | FULLY_FUNCTIONAL | `validator.go` (`scheduleRules` semaphore); specs in `stream_test.go`                                |
| `ContextRule` additive cancellation interface + `NewContextRule` builder                                 | FULLY_FUNCTIONAL | `rule.go`; specs in `context_rule_test.go`                                                           |
| Goroutine-leak regression guards for `Stream` (abandoned + drained)                                      | FULLY_FUNCTIONAL | `stream_test.go` (`runtime.NumGoroutine` pins)                                                       |
| go-cqrs-lite event-bus bridge (`adapters/cqrslite`, nested module)                                       | FULLY_FUNCTIONAL | `adapters/cqrslite/` — 4 specs green via `eventtest.NewFakeBus`; root `go.mod` stays dependency-free |
| OpenTelemetry listener (`listeners/otel`, nested module): spans + 5 metric instruments                   | FULLY_FUNCTIONAL | `listeners/otel/` — specs assert metrics (manual reader) and spans (span recorder)                   |
| Datastar reactive feed example (`examples/sse`, nested module)                                           | FULLY_FUNCTIONAL | `examples/sse/` — real-HTTP smoke tests over the `@post` SSE stream and the `/events` fan-out        |

## Quality & Testing

| Feature                             | Status           | Evidence                                                                                                                                |
| ----------------------------------- | ---------------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| BDD test suite (Ginkgo/Gomega)      | FULLY_FUNCTIONAL | 251 specs pass (`go test ./...`), `*_test.go`                                                                                           |
| Property-based invariant test       | FULLY_FUNCTIONAL | `property_test.go` — gopter `violations(Build) ≡ violations(Stream)`; opt-in via `GBR_PROPERTY=1`                                       |
| Code coverage                       | FULLY_FUNCTIONAL | 97.1% (`go test -cover`, re-measured 2026-09-18)                                                                                        |
| Example tests (godoc-rendered)      | FULLY_FUNCTIONAL | 16 `Example*` funcs, `example_test.go`                                                                                                  |
| Fuzz tests                          | FULLY_FUNCTIONAL | 7 `Fuzz*` targets, `fuzz_test.go`                                                                                                       |
| Benchmarks                          | FULLY_FUNCTIONAL | 13 `Benchmark*` funcs incl. `BenchmarkStream*` variants, `benchmark_test.go`                                                            |
| Branching-flow BDD regression tests | FULLY_FUNCTIONAL | `bdd_branching_flow_test.go` — all pins re-pinned 2026-09-14 to analyzer reality (25 PHANTOM / 29 stats); policy: re-pin with a comment |

## Serialization

| Feature                       | Status               | Evidence                                                                                                                                                            |
| ----------------------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `encoding/json/v2` marshaling | PARTIALLY_FUNCTIONAL | `errors.go:70`, `validation_result.go:160` — works but requires `GOEXPERIMENT=jsonv2` (Go 1.26+); **hard breaking change for downstream consumers** (see AGENTS.md) |

## Tooling & Infrastructure

| Feature                                                        | Status           | Evidence                                                                                                                                                                                                     |
| -------------------------------------------------------------- | ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Nix flake devShell + CI shell                                  | FULLY_FUNCTIONAL | `flake.nix` (`nix develop`, sets `GOEXPERIMENT=jsonv2`)                                                                                                                                                      |
| `check-all` flake app (build+vet+test+lint over all 4 modules) | FULLY_FUNCTIONAL | `nix run .#check-all` (`flake.nix`)                                                                                                                                                                          |
| GitHub Actions CI                                              | FULLY_FUNCTIONAL | `.github/workflows/ci.yml` — 4-module matrix (test/lint/build per module + root gosec with a `Stats.files > 0` assertion); green 13/13 on run `35310049647` (2026-09-18; public repo, no billing dependency) |
| Public module publication                                      | FULLY_FUNCTIONAL | Repo public since 2026-09-17; `proxy.golang.org` serves `v2.1.0`/`v2.2.0`; pkg.go.dev renders the API at `pkg.go.dev/github.com/LarsArtmann/go-business-rules/v2`                                            |
| golangci-lint v2 config                                        | FULLY_FUNCTIONAL | `.golangci.yml` — 0 issues                                                                                                                                                                                   |
| `sivchari/govalid` structural-validator integration pattern    | FULLY_FUNCTIONAL | Documented in README; complementary layer                                                                                                                                                                    |

## Known Gaps & Missing Features

| Item                                                                       | Status           | Note                                                                                      |
| -------------------------------------------------------------------------- | ---------------- | ----------------------------------------------------------------------------------------- |
| Polish-Customs integration (real-world consumer)                           | FULLY_FUNCTIONAL | consumes the PUBLISHED `.../v2@v2.2.0` (bumped 2026-09-18; `go build` + full suite green) |
| Country-specific `PostalCode` patterns, async rules, `Priority()` metadata | PLANNED          | shipped builders/composition/metadata are FULLY_FUNCTIONAL; remainder in ROADMAP.md       |

---

_Owned by this file: the honest feature inventory with status. For what's changing per version, see `CHANGELOG.md`. For open work, see `TODO_LIST.md`. For long-term ideas, see `ROADMAP.md`._
