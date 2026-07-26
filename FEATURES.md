# Features — businessrules

> Honest feature inventory by status. Verified against code, not trust.
>
> **Status legend:** `FULLY_FUNCTIONAL` (works, tested) · `PARTIALLY_FUNCTIONAL` (ships with known gaps) · `BROKEN` (exists but fails) · `PLANNED` (no code yet).
>
> **Verified:** 2026-07-26 against `master` (145 Ginkgo specs pass, 94.8% coverage, `go test` green inside `nix develop`).

---

## Core Framework

| Feature                                                                                | Status               | Evidence                                                                                      |
| -------------------------------------------------------------------------------------- | -------------------- | --------------------------------------------------------------------------------------------- |
| `Severity` type with 4 levels (Info / Warning / Error / Critical)                      | FULLY_FUNCTIONAL     | `severity.go` — type alias to `finding.Severity` (`string`); constants re-exported            |
| `Rule` interface (Name / Check / Severity / Message)                                   | FULLY_FUNCTIONAL     | `rule.go:6`                                                                                   |
| `RuleImpl` base implementation + immutable `WithName` / `WithSeverity` / `WithMessage` | FULLY_FUNCTIONAL     | `rule.go:24-57`                                                                               |
| `NewRule` constructor                                                                  | FULLY_FUNCTIONAL     | `rule.go:61`                                                                                  |
| `ViolationError` (failed check + context + timestamp), implements `error`              | FULLY_FUNCTIONAL     | `errors.go:12`                                                                                |
| `NewViolation` / `NewViolationFromError` / `WithContext`                               | FULLY_FUNCTIONAL     | `errors.go:41,51,61`                                                                          |
| `ValidationResultError` with severity filtering & first-violation access               | FULLY_FUNCTIONAL     | `validation_result.go:10`                                                                     |
| `ValidatorBuilder` fluent API (`NewValidator` / `AddRule` / `AddRules` / `Build`)      | FULLY_FUNCTIONAL     | `validator.go:5-48`                                                                           |
| `Version` constant                                                                     | PARTIALLY_FUNCTIONAL | `doc.go:76` says `"1.1.0"`, but the only git tag is `v0.1.0` — **split-brain, see TODO_LIST** |

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

## Quality & Testing

| Feature                             | Status           | Evidence                                   |
| ----------------------------------- | ---------------- | ------------------------------------------ |
| BDD test suite (Ginkgo/Gomega)      | FULLY_FUNCTIONAL | 145 specs pass (`go test -v`), `*_test.go` |
| Code coverage                       | FULLY_FUNCTIONAL | 94.8% (`go test -cover`)                   |
| Example tests (godoc-rendered)      | FULLY_FUNCTIONAL | 15 `Example*` funcs, `example_test.go`     |
| Fuzz tests                          | FULLY_FUNCTIONAL | 7 `Fuzz*` targets, `fuzz_test.go`          |
| Benchmarks                          | FULLY_FUNCTIONAL | 7 `Benchmark*` funcs, `benchmark_test.go`  |
| Branching-flow BDD regression tests | FULLY_FUNCTIONAL | `bdd_branching_flow_test.go`               |

## Serialization

| Feature                       | Status               | Evidence                                                                                                                                                            |
| ----------------------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `encoding/json/v2` marshaling | PARTIALLY_FUNCTIONAL | `errors.go:70`, `validation_result.go:160` — works but requires `GOEXPERIMENT=jsonv2` (Go 1.26+); **hard breaking change for downstream consumers** (see AGENTS.md) |

## Tooling & Infrastructure

| Feature                                                     | Status           | Evidence                                                               |
| ----------------------------------------------------------- | ---------------- | ---------------------------------------------------------------------- |
| Nix flake devShell + CI shell                               | FULLY_FUNCTIONAL | `flake.nix` (`nix develop`, sets `GOEXPERIMENT=jsonv2`)                |
| treefmt formatting checks                                   | FULLY_FUNCTIONAL | `flake.nix` treefmt config (gofumpt, goimports, nixfmt)                |
| GitHub Actions CI                                           | FULLY_FUNCTIONAL | `.github/workflows/ci.yml` (tests, lint, security; SHA-pinned actions) |
| golangci-lint v2 config                                     | FULLY_FUNCTIONAL | `.golangci.yml` — 0 issues                                             |
| `sivchari/govalid` structural-validator integration pattern | FULLY_FUNCTIONAL | Documented in README; complementary layer                              |

## Known Gaps & Missing Features

| Item                                                                  | Status               | Note                                |
| --------------------------------------------------------------------- | -------------------- | ----------------------------------- |
| Polish-Customs integration (real-world consumer)                      | PLANNED              | Never started; tracked in TODO_LIST |
| Additional rule builders (Time/Date, Network/ID, Precision)           | PLANNED              | See ROADMAP.md                      |
| Advanced composition (`Not`, `Or`, `Xor`, async rules, rule metadata) | PLANNED              | See ROADMAP.md                      |
| `docs/DOMAIN_LANGUAGE.md` filled in                                   | PARTIALLY_FUNCTIONAL | Currently a placeholder template    |

---

_Owned by this file: the honest feature inventory with status. For what's changing per version, see `CHANGELOG.md`. For open work, see `TODO_LIST.md`. For long-term ideas, see `ROADMAP.md`._
