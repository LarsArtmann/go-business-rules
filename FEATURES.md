# Features — businessrules

> Honest feature inventory by status. Verified against code, not trust.
>
> **Status legend:** `FULLY_FUNCTIONAL` (works, tested) · `PARTIALLY_FUNCTIONAL` (ships with known gaps) · `BROKEN` (exists but fails) · `PLANNED` (no code yet).
>
> **Verified:** 2026-09-14 against `master` (169/169 Ginkgo specs pass incl. branching-flow pins, coverage 95.9% measured via `go test -cover`).

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
| `Version` constant                                                                     | FULLY_FUNCTIONAL | `doc.go:105` reports `"2.1.0"`, reconciled with the `v2.1.0` release               |

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

| Category           | Rules                                                       | Status           | Evidence                                         |
| ------------------ | ----------------------------------------------------------- | ---------------- | ------------------------------------------------ |
| Numeric            | `NonNegative`, `Positive`, `InRange`, `MinInt`, `MaxInt`    | FULLY_FUNCTIONAL | `builders.go:79-134`                             |
| Numeric (extended) | `GreaterThan`, `LessThan`                                   | FULLY_FUNCTIONAL | `builders_collection.go:39-57`                   |
| String             | `NotEmpty`, `NotBlank`, `MinLength`, `MaxLength`, `Matches` | FULLY_FUNCTIONAL | `builders.go:140-233`                            |
| Collection         | `NotEmptySlice[T]`, `NotEmptyMap[T]`                        | FULLY_FUNCTIONAL | `builders_collection.go:25-33`                   |
| Format             | `Email`, `URL`, `UUID`                                      | FULLY_FUNCTIONAL | `builders_format.go:27,49,84`                    |
| Generic            | `Equals[T]`, `OneOf[T]`, `Custom`                           | FULLY_FUNCTIONAL | `builders.go:239`, `builders_composite.go:10,27` |
| Composite          | `All`, `Any`, `When`                                        | FULLY_FUNCTIONAL | `builders_composite.go:64,76,98`                 |

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
| BDD test suite (Ginkgo/Gomega)      | FULLY_FUNCTIONAL | 169 specs pass (`go test ./...`), `*_test.go`                                                                                           |
| Code coverage                       | FULLY_FUNCTIONAL | 95.9% (`go test -cover`, measured 2026-09-14)                                                                                           |
| Example tests (godoc-rendered)      | FULLY_FUNCTIONAL | 15 `Example*` funcs, `example_test.go`                                                                                                  |
| Fuzz tests                          | FULLY_FUNCTIONAL | 7 `Fuzz*` targets, `fuzz_test.go`                                                                                                       |
| Benchmarks                          | FULLY_FUNCTIONAL | 10 `Benchmark*` funcs incl. `BenchmarkStream*` variants, `benchmark_test.go`                                                            |
| Branching-flow BDD regression tests | FULLY_FUNCTIONAL | `bdd_branching_flow_test.go` — all pins re-pinned 2026-09-14 to analyzer reality (18 PHANTOM / 22 stats); policy: re-pin with a comment |

## Serialization

| Feature                       | Status               | Evidence                                                                                                                                                            |
| ----------------------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `encoding/json/v2` marshaling | PARTIALLY_FUNCTIONAL | `errors.go:70`, `validation_result.go:160` — works but requires `GOEXPERIMENT=jsonv2` (Go 1.26+); **hard breaking change for downstream consumers** (see AGENTS.md) |

## Tooling & Infrastructure

| Feature                                                        | Status           | Evidence                                                                                                                                                                                        |
| -------------------------------------------------------------- | ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Nix flake devShell + CI shell                                  | FULLY_FUNCTIONAL | `flake.nix` (`nix develop`, sets `GOEXPERIMENT=jsonv2`)                                                                                                                                         |
| `check-all` flake app (build+vet+test+lint over all 4 modules) | FULLY_FUNCTIONAL | `nix run .#check-all` (`flake.nix`)                                                                                                                                                             |
| GitHub Actions CI                                              | BROKEN (billing) | `.github/workflows/ci.yml` rewritten as a 4-module matrix, but GitHub rejects every job at start: _account payments failed / spending limit_ (run `29447520877`); enable after billing is fixed |
| golangci-lint v2 config                                        | FULLY_FUNCTIONAL | `.golangci.yml` — 0 issues                                                                                                                                                                      |
| `sivchari/govalid` structural-validator integration pattern    | FULLY_FUNCTIONAL | Documented in README; complementary layer                                                                                                                                                       |

## Known Gaps & Missing Features

| Item                                                                  | Status           | Note                                                                        |
| --------------------------------------------------------------------- | ---------------- | --------------------------------------------------------------------------- |
| Polish-Customs integration (real-world consumer)                      | FULLY_FUNCTIONAL | `pkg/types` consumes `/v2` via local `replace`; full suite green 2026-09-14 |
| Push `v2.1.0` tag + drop Polish-Customs `replace`                     | PLANNED          | User action; see TODO_LIST                                                  |
| CI re-enablement                                                      | PLANNED          | Blocked on GitHub billing; see TODO_LIST                                    |
| Additional rule builders (Time/Date, Network/ID, Precision)           | PLANNED          | See ROADMAP.md                                                              |
| Advanced composition (`Not`, `Or`, `Xor`, async rules, rule metadata) | PLANNED          | See ROADMAP.md                                                              |

---

_Owned by this file: the honest feature inventory with status. For what's changing per version, see `CHANGELOG.md`. For open work, see `TODO_LIST.md`. For long-term ideas, see `ROADMAP.md`._
