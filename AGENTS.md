# businessrules - Project Guide

Severity-aware validation for Go with multiple outcome levels (Info, Warning, Error, Critical).

Long-form narratives, triage recipes, and pin history live in
[docs/agents-gotchas.md](docs/agents-gotchas.md) — consult it when a tool
finding, pin, or environment check goes red.

## Project Overview

Standalone Go library for validation with severity levels. Unlike standard validators that only support pass/fail semantics, businessrules enables nuanced validation outcomes.

## Build Commands

Run ALL Go commands inside `nix develop` (network, GOPRIVATE, and
`GOEXPERIMENT=jsonv2` live there):

    nix develop --command go test ./...              # root module
    nix run .#check-all                              # build+vet+test+lint, ALL 4 modules
    nix develop --command buildflow --build-mode=full  # quality gate

**Toolchain skew gotcha:** ambient `buildflow` FAILS (`govalid-generate`,
`test-race`): BuildFlow runs tools with `GOTOOLCHAIN=local`, but dependencies
require go >= 1.27.1 while the ambient system Go is older. The devshell ships a
matching Go plus correctly built `govalid`/`go-licenses`. Always invoke
BuildFlow through `nix develop --command` (full story in the gotchas doc).

**Go directive rule:** `adapters/cqrslite` MUST keep `go 1.27.1` (a dependency
floor `go mod tidy` enforces); the go-version-auto-configure "patch component"
warning there is a Go-rule necessity, not a preference (other modules stay
`go 1.27`).

## Architecture

### Core Types

- **`Rule`** - Interface for validation rules with `Name()`, `Check()`, `Severity()`, `Message()`
- **`ContextRule`** - Optional additive interface: `Rule` + `CheckContext(ctx)` for cancellation-aware checks
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

| File                     | Purpose                                                     |
| ------------------------ | ----------------------------------------------------------- |
| `rule.go`                | Rule interface and base implementation                      |
| `severity.go`            | Severity enum and helpers                                   |
| `errors.go`              | ViolationError type and constructors                        |
| `validation_result.go`   | ValidationResultError type with filtering methods           |
| `validator.go`           | Validator builder pattern + Stream(ctx) concurrency         |
| `events.go`              | Event/Listener types (RuleEvaluated, ValidationCompleted)   |
| `builders.go`            | Pre-built rule constructors (numeric, string, generic)      |
| `builders_collection.go` | Collection rules + extended numeric + precision rules       |
| `builders_format.go`     | Format-specific rules (email, URL, UUID)                    |
| `builders_composite.go`  | Composite rules (All, Any, When, Not, Or, Xor)              |
| `builders_time.go`       | Time/date rules (NotPast, NotFuture, DateInRange)           |
| `builders_network.go`    | Network/ID rules (IPAddress, CreditCard, Phone, Postal)     |
| `property_test.go`       | gopter Build≡Stream property test (opt-in `GBR_PROPERTY=1`) |

### Nested Go Modules (adapters, examples & listeners)

Nested modules keep optional dependencies out of the root `go.mod`; all use
`replace github.com/LarsArtmann/go-business-rules/v2 => ../..`:

| Module              | Depends on                                    | Purpose                                                              |
| ------------------- | --------------------------------------------- | -------------------------------------------------------------------- |
| `adapters/cqrslite` | `go-cqrs-lite/event/v4`, `id/v4`, `eventtest` | Publishes validation events onto a CQRS event bus (`NewBusListener`) |
| `examples/sse`      | `go-sse`, `datastar-go`                       | Datastar reactive browser feed of validation events                  |
| `listeners/otel`    | `go.opentelemetry.io/otel{,/sdk,/sdk/metric}` | Maps validation events onto OTel spans + metrics (`otel.New`)        |

Root `golangci-lint run` does NOT lint nested modules (separate modules) — use
`nix run .#check-all` (exists because BuildFlow's gate only covers the module it
detects config in, and nested modules need the devshell env).

**Module path & version:** `github.com/LarsArtmann/go-business-rules/v2` (the
`/v2` suffix makes v2+ tags consumable). Current release: `v2.2.0`.
**Release checklist gotcha:** bump `doc.go` `Version` AND the `suite_test.go`
version expectation TOGETHER. Real-world consumer: Polish-Customs consumes the
published v2.2.0 (no `replace`).

## Validation Events & Streaming

- `events.go` defines a sealed `Event` interface with exactly two
  implementations: `RuleEvaluated` (per rule check, passes included) and
  `ValidationCompleted` (terminal, carries the result). `Passed()` is a derived
  method (`Err == nil`) — deliberately NO bool field, so an event cannot claim
  success while carrying an error.
- `Listener` = `func(Event)`: synchronous, registration order, before `Build`
  returns; listeners must not panic (no recover in the library). With no
  listeners `Build` constructs no events and calls no `time.Now` (zero-cost).
- `Stream(ctx)` runs rules concurrently, streams events in completion order on
  a channel, emits `ValidationCompleted` with violations re-sorted into rule
  order. Channel events go ONLY to the channel, never to `WithListener`
  listeners (one delivery mechanism per API). The results channel is buffered
  to `len(rules)` so abandoned streams never leak goroutines (guarded by
  `runtime.NumGoroutine` specs). Cancellation skips unstarted rules; a
  cancellation error from a check is reported as that rule's outcome like any
  other error.
- `WithConcurrency(n)` bounds in-flight checks (cancellation-aware; values
  below 1 keep the one-goroutine-per-rule default; `Build` unaffected).
- `ContextRule`: rules implementing `CheckContext(ctx)` get the stream's
  context from `Stream`; everything else falls back to `Check`. `Build` always
  calls `Check` (no ctx parameter by design).
- Ginkgo `-count` >1 is rejected by Ginkgo itself — loop the whole `go test`
  command instead when hunting flakes. Benchmarks live in the gotchas doc.

## Code Patterns

### Rule Builders

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

- Use `minimum`/`maximum` (not `min`/`max`, which shadow Go 1.21+ builtins) and descriptive names that explain the purpose

## Testing

- Uses Ginkgo/Gomega for BDD-style testing; test files use dot-imports (standard BDD pattern)
- Run specific test: `go test -run "TestName" ./...`

## Dependencies

**Runtime**: only `github.com/larsartmann/go-finding` (shared `Severity` type,
re-exported; constants are type aliases). **Dev**: `onsi/ginkgo/v2`,
`onsi/gomega`. Everything else is standard library.

## Environment & Publishing Essentials

- **GOEXPERIMENT=jsonv2 is REQUIRED** (`encoding/json/v2` is experimental):
  both devShells set it; `.golangci.yml` carries `goexperiment.jsonv2` in
  build-tags. Hard breaking constraint for downstream consumers (they must set
  it too; see gotchas doc). Affected files: `errors.go`,
  `validation_result.go`, `bdd_branching_flow_test.go`. Re-evaluate when
  json/v2 graduates.
- Repo is PUBLIC (2026-09-17): CI runs a matrix over all four modules (test,
  lint, build each; gosec root-only; setup-go pinned `"1.27"`); no tokens
  needed. CI commands are verified locally via `nix run .#check-all`.
- GOPRIVATE/GONOSUMDB/GONOPROXY flake overrides are deliberate no-ops today
  (public since 2026-09-17) but must stay — interplay with the OS-level
  `GONOSUMDB` env var; stale-module-cache sumdb recipe in the gotchas doc.
- `lychee.toml` (NOT `.lychee.toml` — not auto-discovered) excludes the
  archived docs dirs from link checking; archived files are never edited.

## Linting

golangci-lint v2 key settings: `govet.fieldalignment` disabled; test files
excluded from `revive`/`makezero`; `godot` scope toplevel;
`exhaustruct_v5` nolint gotcha — the v5 variant requires
`//nolint:exhaustruct_v5` (plain `exhaustruct` is silently ignored).

## Analyzer Pins & Documented False Positives

- **branching-flow pins (0.6.4/192031d, 2026-10-11; re-verified unchanged
  under 0.6.7/aa1611d, 2026-10-15):** 20 PHANTOM (all PHANTOM-TRANSPOSE
  same-type param pairs, all warning, four builder files); `stats` total 24 =
  20 phantom + 2 flagparam + 1 ifacecomplete + 1 mixins. 0.6.7 made
  phantom/stats/all exit 1 on findings (golangci-lint convention); the BDD
  gate tolerates that findings exit. Panic findings in `validator.go` (Stream
  result send) and `builders_collection.go` (DivisibleBy zero-guard) are
  suppressed with the native `//branching-flow:ignore panic` directive
  (commit 112b9a5; the old `//nolint:branching-flow:panic` form triggered
  golangci nolintlint unknown-linter warnings, whose `.golangci.yml`
  exclusions are since removed). A red pin
  means: diff the findings; if still the deliberate-primitives class, re-pin
  (procedure + version history in the gotchas doc).
- **hierarchical-errors / go-auto-upgrade / go-structure-linter
  root-package-files / gomod-check mixed-requires:** documented false
  positives (rationale per analyzer in the gotchas doc). Root-level Go files
  ARE the public API (`.go-structure-linter.yaml` uses the `flat` preset).
- **Residual non-gating BuildFlow findings (policy):** nix-checker
  hash-extraction suggestions declined (`nix-hash-fix` owns hashes); vulnix
  CVEs are nixpkgs build-toolchain closure noise (advisory); doctor per-tool
  availability noise. Detail in the gotchas doc.

## Known Patterns

- **Structural similarity in builders is intentional** — each builder is
  self-contained; deduplicating them adds complexity without benefit.
- **ci.yml per-job setup duplication (jscpd ~23 lines) is intentional** —
  Actions requires per-job setup; a composite action adds indirection.
- **art-dupl:** ZERO semantic clones (suffix-tree, -t 15) — all eliminated.
- **govalid integration:** complementary layers — govalid covers structural
  validation (required/format/type); businessrules covers business domain
  rules with severity levels.

## Historical docs & archive layout (2026-09-14)

All point-in-time docs are archived with inline strikethrough resolutions —
never edit archived files, and never treat them as backlog:

| Location                  | Contents                                                                                                                      |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `docs/status/`            | Current status reports only (2026-09-14 onward)                                                                               |
| `docs/status/archived/`   | All 2026-03 / 2026-07 status reports, every item annotated with its verdict                                                   |
| `docs/planning/archived/` | Executed/deferred plans: implementation plans, event-driven plan, FINDING-SDK, nix-migration proposal, go-composable analysis |
| `docs/reviews/`           | Point-in-time reviews (BDD tests review)                                                                                      |
| `docs/adr/`               | Standing decision records (go-output non-integration)                                                                         |
