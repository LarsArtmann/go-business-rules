# businessrules - Project Guide

Severity-aware validation for Go with multiple outcome levels (Info, Warning, Error, Critical).

## Project Overview

This is a standalone Go library for validation with severity levels. Unlike standard validators that only support pass/fail semantics, businessrules enables nuanced validation outcomes.

## Build Commands

```bash
# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Run linter
golangci-lint run --timeout 5m

# Run full build pipeline
buildflow --semantic --fix
```

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
- `godot` scope: toplevel (comments should end in period)

## Branching-Flow Analysis

The branching-flow multi-linter may report PHANTOM and DUPE violations. These are **false positives** for this validation library pattern:

### PHANTOM Violations (16)

**False positive for validation libraries.** The linter flags using primitive types (string, int, bool) instead of branded types. However, this library is a validation library where:

- Users pass raw primitives to validate them
- The primitives ARE the domain concept being validated
- Forcing branded types would defeat the library's purpose

### DUPE Violations

**False positive for intentionally similar functions.** Functions like `All`/`Any`, `MinInt`/`MaxInt`, `NonNegative`/`Positive` share similar structure but implement different validation logic. Refactoring to reduce structural similarity would:

- Add complexity without meaningful benefit
- Make the code harder to understand
- Remove the explicit, self-contained nature of each builder

## library-policy Scanner

The `library-policy` tool may report `encoding_json_v2_replacement` violations. These are **false positives** that have been disabled via local configuration:

### encoding_json_v2_replacement

**Do NOT migrate to encoding/json/v2.** The recommendation is premature because:

1. **Experimental API**: `encoding/json/v2` is experimental and not subject to the Go 1 compatibility promise
2. **Requires build flag**: Only available with `GOEXPERIMENT=jsonv2` environment variable
3. **Go's own recommendation**: The documentation explicitly states "Most users should use [encoding/json]"
4. **Breaking change**: Migrating would break compatibility for users not using the experimental flag

**Resolution**: A local `library-policy.yaml` config disables this rule by setting `go_version_min: "1.99"` (higher than current Go 1.26).

**Re-evaluate when**: `encoding/json/v2` graduates from experimental status (no longer requires `GOEXPERIMENT=jsonv2`).

**Affected files**: `errors.go`, `validation_result.go` (MarshalJSON implementations)
