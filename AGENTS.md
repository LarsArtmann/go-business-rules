# businessrules - Project Guide

Severity-aware validation for Go with multiple outcome levels (Info, Warning, Error, Critical).

## Project Overview

This is a standalone Go library for validation with severity levels. Unlike standard validators that only support pass/fail semantics, businessrules enables nuanced validation outcomes.

## Build Commands

### Nix

Hermetic build/test checks are not included because the project depends on a private Go module (`github.com/larsartmann/go-finding`) which the Nix sandbox cannot access. Use `nix develop --command just test` instead.

## Architecture

### Core Types

- **`Rule`** - Interface for validation rules with `Name()`, `Check()`, `Severity()`, `Message()`
- **`Violation`** - Represents a failed rule check with context and timestamp
- **`ValidationResult`** - Contains validation outcome with methods to filter by severity

### Severity Levels

```go
SeverityInfo     // Advisory: just so you know
SeverityWarning  // Non-blocking: should fix
SeverityError    // Blocking: must fix
SeverityCritical // Blocking: critical failure
```

### File Structure

| File                    | Purpose                                                |
| ----------------------- | ------------------------------------------------------ |
| `rule.go`               | Rule interface and base implementation                 |
| `severity.go`           | Severity enum and helpers                              |
| `errors.go`             | Violation type and constructors                        |
| `validation_result.go`  | Result type with filtering methods                     |
| `validator.go`          | Validator builder pattern                              |
| `builders.go`           | Pre-built rule constructors (numeric, string, generic) |
| `builders_format.go`    | Format-specific rules (email, URL, UUID)               |
| `builders_composite.go` | Composite rules (All, Any, When)                       |

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

**Runtime**: Zero dependencies (stdlib only)

**Dev**:

- `onsi/ginkgo/v2` - BDD test framework
- `onsi/gomega` - Matcher library

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

### PHANTOM Violations (12)

**False positive for validation libraries.** The linter flags using primitive types (string, int, bool) instead of branded types. However, this library is a validation library where:

- Users pass raw primitives to validate them
- The primitives ARE the domain concept being validated
- Forcing branded types would defeat the library's purpose

### DUPE Violations

**Reduced through refactoring.** Functions like `All`/`Any` have been refactored to use a strategy pattern that reduces structural similarity while maintaining clarity.

## Hierarchical-Errors Analyzer

The `hierarchical-errors` analyzer may report violations about functions returning generic `error` instead of specific error types. These are **false positives** for this validation library pattern:

### generic_return Violations

**False positive for standard library interface implementations.** The analyzer flags functions that return the generic `error` interface. However, this library implements standard Go interfaces where the signature is fixed by the standard library:

1. **`MarshalJSON` methods** (`errors.go:67`, `validation_result.go:152`):
   - Implements `json.Marshaler` interface from `encoding/json` (stdlib)
   - The signature `func MarshalJSON() ([]byte, error)` is fixed by Go
   - Cannot return a custom error type without breaking interface compatibility

2. **`Check` method** (`rule.go:38`):
   - Core `Rule` interface method designed to return `error`
   - Uses Go's idiomatic error handling pattern
   - Custom error types would force all implementations to use the same error type, reducing flexibility

**Resolution**: These violations are intentional design decisions that follow Go conventions and cannot be changed without breaking compatibility.

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
