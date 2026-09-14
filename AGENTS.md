# businessrules - Project Guide

Severity-aware validation for Go with multiple outcome levels (Info, Warning, Error, Critical).

## Project Overview

This is a standalone Go library for validation with severity levels. Unlike standard validators that only support pass/fail semantics, businessrules enables nuanced validation outcomes.

## Build Commands

### Nix

Hermetic build/test checks are not included because the project depends on a private Go module (`github.com/larsartmann/go-finding`) which the Nix sandbox cannot access. Use `nix develop --command go test ./...` instead.

## Architecture

### Core Types

- **`Rule`** - Interface for validation rules with `Name()`, `Check()`, `Severity()`, `Message()`
- **`ViolationError`** - Represents a failed rule check with context and timestamp
- **`ValidationResultError`** - Contains validation outcome with methods to filter by severity

### Severity Levels

```go
SeverityInfo     // Advisory: just so you know
SeverityWarning  // Non-blocking: should fix
SeverityError    // Blocking: must fix
SeverityCritical // Blocking: critical failure
```

### File Structure

| File                     | Purpose                                                   |
| ------------------------ | --------------------------------------------------------- |
| `rule.go`                | Rule interface and base implementation                    |
| `severity.go`            | Severity enum and helpers                                 |
| `errors.go`              | ViolationError type and constructors                      |
| `validation_result.go`   | ValidationResultError type with filtering methods         |
| `validator.go`           | Validator builder pattern + Stream(ctx) concurrency       |
| `events.go`              | Event/Listener types (RuleEvaluated, ValidationCompleted) |
| `builders.go`            | Pre-built rule constructors (numeric, string, generic)    |
| `builders_collection.go` | Collection rules + extended numeric rules                 |
| `builders_format.go`     | Format-specific rules (email, URL, UUID)                  |
| `builders_composite.go`  | Composite rules (All, Any, When)                          |

### Nested Go Modules (adapters, examples & listeners)

The repo is multi-module. Nested modules keep optional dependencies out of the
root `go.mod`:

| Module              | Depends on                                    | Purpose                                                              |
| ------------------- | --------------------------------------------- | -------------------------------------------------------------------- |
| `adapters/cqrslite` | `go-cqrs-lite/event/v4`, `id/v4`, `eventtest` | Publishes validation events onto a CQRS event bus (`NewBusListener`) |
| `examples/sse`      | `go-sse`, `datastar-go`                       | Datastar reactive browser feed of validation events                  |
| `listeners/otel`    | `go.opentelemetry.io/otel{,/sdk,/sdk/metric}` | Maps validation events onto OTel spans + metrics (`otel.New`)        |

All use `replace github.com/LarsArtmann/go-business-rules/v2 => ../..`. Run their
tests from inside each directory: `cd adapters/cqrslite && nix develop --command
go test ./...`. Root `golangci-lint run` does NOT lint nested modules (separate
modules); use the `check-all` flake app instead:

	nix run .#check-all

build + vet + test + lint across ALL four modules. It exists because BuildFlow's
gate only covers the module it detects config in, while the nested private-dep
modules need the flake devshell env (GOPRIVATE/GONOSUMDB/GONOPROXY, GOEXPERIMENT).

**Module path & version (fixed 2026-09-14):** the module path is now
`github.com/LarsArtmann/go-business-rules/v2` — the `/v2` suffix makes v2+ tags
consumable. The first such tag is `v2.1.0` (`doc.go` reports `Version = "2.1.0"`),
verified end-to-end by consuming it from a local file proxy. The legacy `v2.0.0`
tag (2026-07-26, suffix-less path) stays untouched — it only ever resolves as a
`+incompatible` version under the old path, which is exactly the mechanism we do
NOT build on going forward. Until `v2.1.0` is pushed, Polish-Customs consumes the
local tree via a temporary `replace`.

## Validation Events & Streaming

- `events.go` defines a sealed `Event` interface with exactly two
  implementations: `RuleEvaluated` (per rule check, passes included, with
  duration/start time) and `ValidationCompleted` (terminal, carries the result).
  `Passed()` is a derived method (`Err == nil`) — there is deliberately NO bool
  field, so an event cannot claim success while carrying an error (and the
  branching-flow PHANTOM count unchanged by the events code).
- `Listener` = `func(Event)`. Delivery is synchronous, registration order,
  before `Build` returns. Listeners must not panic (no recover in the library).
- Zero-cost default path: with no listeners `Build` performs no `time.Now`
  calls and constructs no events. Benchmarked: 189 ns/op without listener,
  468 ns/op with one, ~476 ns/op with three (2 rules).
- `Stream(ctx)` runs rules concurrently, streams events in completion order on
  a channel, and emits `ValidationCompleted` with violations re-sorted into
  rule order. Channel events go ONLY to the channel, never to `WithListener`
  listeners (one delivery mechanism per API). Cancellation skips unstarted
  rules; the results channel is buffered to `len(rules)` so abandoned streams
  never leak goroutines (guarded by `runtime.NumGoroutine` regression specs).
  Ginkgo `-count` >1 is rejected by Ginkgo itself —
  loop the whole `go test` command instead when hunting flakes.
- `WithConcurrency(n)` bounds in-flight checks via a semaphore acquired in the
  scheduler loop (cancellation-aware). Values below 1 keep the one-goroutine-
  per-rule default. `Build` is never affected.
- `ContextRule` is the additive cancellation interface: rules implementing
  `CheckContext(ctx)` get the stream's context from `Stream` (via
  `NewContextRule` or any custom implementation); everything else falls back
  to `Check`. `Build` always calls `Check` (no ctx parameter by design). A
  cancellation error returned by a check is reported as that rule's outcome
  like any other error — it is not special-cased.
- Stream benchmarks: `BenchmarkStream2Rules`, `BenchmarkStream10Rules`,
  `BenchmarkStream10RulesConcurrency4` (channel drain dominates: ~4.4 us/op
  for 2 cheap rules).

## Code Patterns

### Rule Builders

Rule builders follow a consistent pattern:

```go
func RuleName(name string, value T, severity Severity) Rule {
    return NewRule(name, func() error {
        if /* condition fails */ {
            return fmt.Errorf("descriptive message")
        }
        return nil
    }, severity, "template message")
}
```

### Parameter Naming

- Use `minimum`/`maximum` instead of `min`/`max` to avoid shadowing Go 1.21+ builtins
- Use descriptive names that explain the purpose

## Testing

- Uses Ginkgo/Gomega for BDD-style testing
- Test files use dot-imports for Ginkgo/Gomega (standard BDD pattern)
- Run specific test: `go test -run "TestName" ./...`

## Dependencies

**Runtime**:

- `github.com/larsartmann/go-finding` - Provides the shared `Severity` type (re-exported as `finding.Severity`). The `Severity` constants are type aliases to this package.

**Dev**:

- `onsi/ginkgo/v2` - BDD test framework
- `onsi/gomega` - Matcher library

> **Note:** This library previously advertised "zero runtime dependencies." That is
> no longer true since `Severity` was migrated to `finding.Severity` (commit
> `e423de4`). All other runtime code is standard library.

### Private modules & the GOPRIVATE/GONOSUMDB flake override

Always run Go commands inside `nix develop`. Both devShells explicitly set **all
three** of `GOPRIVATE`, `GONOSUMDB`, and `GONOPROXY` to
`github.com/larsartmann/*,github.com/LarsArtmann/*` (both case variants).

Why: Home Manager sets `GONOSUMDB` as an OS env var to an explicit, incomplete
list. Go only auto-derives `GONOSUMDB`/`GONOPROXY` from `GOPRIVATE` when they are
**unset**, so a bare `GOPRIVATE` wildcard loses — and a force-pushed private tag
(like `go-error-family@v0.10.0`) then fails `sum.golang.org` verification with a
checksum-mismatch SECURITY ERROR. Fix recipe (2026-07-26, commits `0bcd851`,
`e50a4c8`, `566bf3e`): set all three vars explicitly, then clear the stale
proxy-cached download with `trash "$(go env GOMODCACHE)/cache/download/github.com/larsartmann/<module>"`.

## CI & Publishing Reality (discovered 2026-09-14)

- **The CI workflow is `disabled_manually` on GitHub since 2026-07-17.** No runs in
  2+ months; the README CI badge reflects nothing. Every run since 2026-06 failed
  within 3-5s (setup-level). Verify CI claims with `gh workflow list --all`, never
  by assuming the badge. Re-enable decision is a user call (private-repo Actions
  minutes) — tracked in TODO_LIST.
- **The GitHub repo is PRIVATE.** `proxy.golang.org` has zero cached versions and
  pkg.go.dev 404s; the module can never be indexed while private. Any old report
  saying "verify on pkg.go.dev" was unachievable. (`go-finding`, the runtime dep,
  IS public — CI needs no `GOPRIVATE` for it.)
- Local gates are the real quality bar: `nix develop --command go test ./...`
  (156/158, 2 known-red branching-flow pins), `go test -cover` (95.9%),
  `golangci-lint run` (0 issues), `nix build .#checks.x86_64-linux.format`.

## Historical docs & archive layout (2026-09-14)

All point-in-time docs are archived with inline strikethrough resolutions — never
edit archived files, and never treat them as backlog:

| Location                  | Contents                                                                                                                      |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `docs/status/`            | Current status reports only (2026-09-14 onward)                                                                               |
| `docs/status/archived/`   | All 2026-03 / 2026-07 status reports, every item annotated with its verdict                                                   |
| `docs/planning/archived/` | Executed/deferred plans: implementation plans, event-driven plan, FINDING-SDK, nix-migration proposal, go-composable analysis |
| `docs/reviews/`           | Point-in-time reviews (BDD tests review)                                                                                      |
| `docs/adr/`               | Standing decision records (go-output non-integration)                                                                         |

## Integration with sivchari/govalid

This library complements structural validators:

1. **Structural validation** (govalid): required, format, type - zero allocations, compile-time safe
2. **Business validation** (businessrules): domain rules, severity levels

## Known Patterns

### Structural Similarity in Builders

Functions like `NonNegative` and `Positive` share similar structure. This is intentional - each builder is self-contained and explicit. Refactoring to reduce "duplication" would add complexity without meaningful benefit.

## Linting

Uses golangci-lint v2 with the following key settings:

- `govet.fieldalignment` disabled (micro-optimization for small structs)
- Test files excluded from `revive` rules (dot-imports for Ginkgo/Gomega)
- Test files excluded from `makezero` (index assignment after `make([]T, n)` is intentional in tests)
- `godot` scope: toplevel (comments should end in period)

## Branching-Flow Analysis

The branching-flow multi-linter may report PHANTOM and DUPE violations. These are **false positives** for this validation library pattern:

### PHANTOM Violations (14)

**False positive for validation libraries.** The linter flags using primitive types (string, int, bool) instead of branded types. However, this library is a validation library where:

- Users pass raw primitives to validate them
- The primitives ARE the domain concept being validated
- Forcing branded types would defeat the library's purpose

The count is 14 (re-pinned from 12 in commit `0453a78`) because the nested
`examples/sse` module also validates raw primitives on purpose. **Pin fragility:**
the analyzer scans the whole directory and has no path-exclude flag, so ANY new
module/example/test changes the counts. The `stats` total pin (36) is separately
stale from analyzer binary drift (actual 15) — policy decision tracked in TODO_LIST.

### DUPE Violations

**Reduced through refactoring.** Functions like `All`/`Any` have been refactored to use a strategy pattern that reduces structural similarity while maintaining clarity.

## Hierarchical-Errors Analyzer

The `hierarchical-errors` analyzer may report violations about functions returning generic `error` instead of specific error types. These are **false positives** for this validation library pattern:

### generic_return Violations

**False positive for standard library interface implementations.** The analyzer flags functions that return the generic `error` interface. However, this library implements standard Go interfaces where the signature is fixed by the standard library:

1. **`MarshalJSON` methods** (`errors.go:70`, `validation_result.go:160`):
   - Implements `json.Marshaler` interface from `encoding/json` (stdlib)
   - The signature `func MarshalJSON() ([]byte, error)` is fixed by Go
   - Cannot return a custom error type without breaking interface compatibility

2. **`Check` method** (`rule.go:36`):
   - Core `Rule` interface method designed to return `error`
   - Uses Go's idiomatic error handling pattern
   - Custom error types would force all implementations to use the same error type, reducing flexibility

3. **Internal helper functions** (`builders_format.go:17`, `builders_composite.go:34`, `builders_composite.go:51`):
   - `checkNonEmpty`, `collectAllViolations`, `anyRulePasses` aggregate `Rule.Check()` results
   - They inherit the `error` return from the `Rule` interface; narrowing the return type would couple them to a single error implementation

**Resolution**: These violations are intentional design decisions that follow Go conventions and cannot be changed without breaking compatibility.

## go-auto-upgrade Analyzer

The `go-auto-upgrade` linter may suggest replacing the manual slice-to-map loop in
`validation_result.go` (`BySeverity`, around line 39) with `samber/lo.SliceToMap`.

**False positive.** Adding `samber/lo` for a 3-line loop would introduce a new
runtime dependency for negligible benefit. The manual loop is idiomatic Go.

## go-structure-linter: root-package-files

The `go-structure-linter` may report that the Go source files live at the project
root rather than under `/internal/` or `/pkg/`.

**False positive for a public Go library.** The root-level package files ARE the
public API; moving them to `/internal/` would make them unexportable, and moving
them to `/pkg/` would change the import path for all consumers.

## encoding/json/v2 Migration

This library uses `encoding/json/v2`, which is experimental and requires `GOEXPERIMENT=jsonv2`.

### How it is wired

- **flake.nix**: Both devShells (`default` and `ci`) set `GOEXPERIMENT = "jsonv2"`, so all `go build`/`go test`/`go vet` commands work inside `nix develop`.
- **.golangci.yml**: `goexperiment.jsonv2` is in `build-tags`.

### Downstream consumer impact

This is a **hard breaking change for downstream consumers**. Anyone who `go get`s this library must set `GOEXPERIMENT=jsonv2` (and use Go 1.26+) or compilation fails with a cryptic "build constraints exclude all Go files in encoding/json/v2" error. There is no way to self-contain this requirement in `go.mod`.

**Re-evaluate when**: `encoding/json/v2` graduates from experimental status (no longer requires `GOEXPERIMENT=jsonv2`). At that point the downstream constraint disappears.

**Affected files**: `errors.go`, `validation_result.go` (MarshalJSON), `bdd_branching_flow_test.go` (Unmarshal)

## gomod-check False Positive

The `gomod-check` tool may report: `go.mod:12: direct and indirect requires are mixed (should be separate blocks since Go 1.17+)`.

**False positive.** The go.mod already has properly separated `require` blocks for direct and indirect dependencies. Running `go mod tidy` confirms no changes needed. The warning cannot be auto-fixed because there is nothing to fix.

## art-dupl Analysis

The `art-dupl` tool finds code clones using suffix tree algorithms. Running with threshold 15 tokens:

```bash
art-dupl --semantic --sort total-tokens -t 15
```

**Status: ZERO clones achieved** ✅

All previously reported clone groups have been eliminated through refactoring:

# 
