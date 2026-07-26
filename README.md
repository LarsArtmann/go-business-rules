# businessrules

> **Severity-aware validation for Go** — because not all validation failures are equal.

A Go library that adds severity levels to validation, enabling applications to distinguish between critical errors, warnings, and informational issues. Standard validators return pass/fail. `businessrules` returns the _degree_ of failure.

[![GoDoc](https://pkg.go.dev/badge/github.com/LarsArtmann/go-business-rules.svg)](https://pkg.go.dev/github.com/LarsArtmann/go-business-rules)
[![CI](https://github.com/LarsArtmann/go-business-rules/actions/workflows/ci.yml/badge.svg)](https://github.com/LarsArtmann/go-business-rules/actions)

## Installation

```bash
go get github.com/LarsArtmann/go-business-rules
```

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/LarsArtmann/go-business-rules"
)

type Package struct {
    Weight     float64
    TrackingID string
}

func (p Package) Rules() []businessrules.Rule {
    return []businessrules.Rule{
        businessrules.NonNegative("weight", p.Weight, businessrules.SeverityError),
        businessrules.InRange("weight", p.Weight, 0.001, 100, businessrules.SeverityWarning),
        businessrules.NotEmpty("tracking_id", p.TrackingID, businessrules.SeverityError),
    }
}

func (p Package) Validate() businessrules.ValidationResultError {
    return businessrules.NewValidator().
        AddRules(p.Rules()...).
        Build()
}

func main() {
    pkg := Package{Weight: -5, TrackingID: ""}

    result := pkg.Validate()

    if result.HasErrors() {
        for _, err := range result.Errors() {
            fmt.Printf("ERROR: %s\n", err.Error())
        }
    }

    for _, warn := range result.Warnings() {
        fmt.Printf("WARNING: %s\n", warn.Error())
    }

    if !result.HasErrors() {
        processPackage(pkg)
    }
}
```

## Why businessrules?

Standard validators return valid/invalid. Business rules return the _degree_ of failure:

| Output              | Use case                                            |
| ------------------- | --------------------------------------------------- |
| **Valid / Invalid** | Binary pass/fail — enough for structural validation |
| **Severity levels** | Nuanced outcomes — critical for business validation |

```go
// User submits weight = -5 kg
// → ERROR: Invalid value (blocks processing)

// User submits weight = 2000 kg
// → WARNING: Suspicious but allowed (logs, notifies, allows)

// User submits lowercase "kg" instead of "KG"
// → INFO: Style suggestion (advisory only)
```

### Works with structural validators

`businessrules` complements type/format validators like [`sivchari/govalid`](https://github.com/sivchari/govalid):

| Layer      | Validator       | Validates                     |
| ---------- | --------------- | ----------------------------- |
| Structural | `govalid`       | Type, format, required fields |
| Business   | `businessrules` | Domain rules, severity levels |

```go
type User struct {
    Email string `govalid:"required,email"`
    Age   int    `govalid:"min=0,max=150"`
}

func (u User) ValidateAll() (*businessrules.ValidationResultError, error) {
    if err := govalid.Validate(u); err != nil {
        return nil, err
    }

    result := businessrules.NewValidator().
        AddRule(businessrules.Custom("email_domain", func() error {
            if !strings.HasSuffix(u.Email, "@company.com") {
                return errors.New("must be company email")
            }
            return nil
        }, businessrules.SeverityWarning)).
        Build()

    return &result, nil
}
```

## API

### Severity Levels

```go
// Severity is a type alias for github.com/larsartmann/go-finding.Severity (a string).
const (
    SeverityInfo     Severity = "info"     // Advisory: just so you know
    SeverityWarning  Severity = "warning"  // Non-blocking: should fix
    SeverityError    Severity = "error"    // Blocking: must fix
    SeverityCritical Severity = "critical" // Blocking: critical failure
)
```

Severity values are re-exported from [`go-finding`](https://github.com/LarsArtmann/go-finding),
so they are interchangeable with that library's `finding.Severity` type.

### Rule Interface

```go
type Rule interface {
    Name() string
    Check() error
    Severity() Severity
    Message() string
}
```

`RuleImpl` implements `Rule` and provides immutable mutation methods:

```go
func NewRule(name string, check func() error, severity Severity, message string) RuleImpl

func (r RuleImpl) WithName(name string) RuleImpl
func (r RuleImpl) WithSeverity(severity Severity) RuleImpl
func (r RuleImpl) WithMessage(message string) RuleImpl
```

### ViolationError

```go
type ViolationError struct {
    Timestamp time.Time
    Context   string
    Rule      Rule
}

func (v ViolationError) Error() string
func (v ViolationError) WithContext(context string) ViolationError
func (v ViolationError) MarshalJSON() ([]byte, error)

func NewViolation(rule Rule, context string) ViolationError
func NewViolationFromError(rule Rule, err error) ViolationError
```

### ValidationResultError

```go
type ValidationResultError struct {
    Valid           bool
    ViolationErrors []ViolationError
}
```

**Filtering:**

```go
func (r ValidationResultError) Errors() []ViolationError
func (r ValidationResultError) Warnings() []ViolationError
func (r ValidationResultError) Info() []ViolationError
func (r ValidationResultError) Critical() []ViolationError
func (r ValidationResultError) BySeverity(severities ...Severity) []ViolationError
func (r ValidationResultError) Filter(predicate func(ViolationError) bool) []ViolationError
```

**Checks:**

```go
func (r ValidationResultError) HasErrors() bool
func (r ValidationResultError) HasWarnings() bool
func (r ValidationResultError) HasCritical() bool
func (r ValidationResultError) HasInfo() bool
func (r ValidationResultError) Count() int
```

**First violation:**

```go
func (r ValidationResultError) FirstError() ViolationError
func (r ValidationResultError) FirstCritical() ViolationError
func (r ValidationResultError) FirstWarning() ViolationError
func (r ValidationResultError) FirstInfo() ViolationError
```

**Combining and iterating:**

```go
func (r ValidationResultError) Merge(other ValidationResultError) ValidationResultError
func (r ValidationResultError) ForEach(fn func(ViolationError))
```

**Serialization:**

```go
func (r ValidationResultError) MarshalJSON() ([]byte, error)
func (r ValidationResultError) Error() string
```

### Pre-built Rules

#### Numeric

```go
NonNegative(name string, value float64, severity Severity) RuleImpl
Positive(name string, value float64, severity Severity) RuleImpl
GreaterThan(name string, value, minimum float64, severity Severity) RuleImpl
LessThan(name string, value, maximum float64, severity Severity) RuleImpl
InRange(name string, value, minimum, maximum float64, severity Severity) RuleImpl
MinInt(name string, value, minimum int, severity Severity) RuleImpl
MaxInt(name string, value, maximum int, severity Severity) RuleImpl
```

#### String

```go
NotEmpty(name, value string, severity Severity) RuleImpl
NotBlank(name, value string, severity Severity) RuleImpl
MinLength(name, value string, minimum int, severity Severity) RuleImpl
MaxLength(name, value string, maximum int, severity Severity) RuleImpl
Matches(name, value string, pattern *regexp.Regexp, severity Severity) RuleImpl
```

#### Collection

```go
NotEmptySlice[T any](name string, value []T, severity Severity) RuleImpl
NotEmptyMap[T any](name string, value map[string]T, severity Severity) RuleImpl
```

#### Format

```go
Email(name, value string, severity Severity) RuleImpl
URL(name, value string, severity Severity) RuleImpl
UUID(name, value string, severity Severity) RuleImpl
```

#### Generic

```go
Equals[T comparable](name string, value, expected T, severity Severity) RuleImpl
OneOf[T comparable](name string, value T, allowed []T, severity Severity) RuleImpl
Custom(name string, check func() error, severity Severity) RuleImpl
```

#### Composite

```go
All(name string, rules []Rule, severity Severity) RuleImpl
Any(name string, alternatives []Rule, severity Severity) RuleImpl
When(name string, condition bool, rule Rule) RuleImpl
```

### Validator Builder

```go
result := businessrules.NewValidator().
    AddRule(rule1).
    AddRules(rule2, rule3).
    Build()
```

## Philosophy

- **Minimal dependencies** — one runtime dependency (`go-finding`) for the shared `Severity` type; everything else is the standard library
- **Type-safe** — generics, no `any` types
- **Immutable rules** — safe for concurrent use after creation
- **Composable** — integrates with structural validators like `sivchari/govalid`
- **Tested with Ginkgo/Gomega** — BDD-style testing for behavior specification
- **94.8% test coverage** — 145 specs, 15 examples, 7 fuzz targets, 7 benchmarks

## Dependencies

| Dependency                     | Purpose          | Notes                                  |
| ------------------------------ | ---------------- | -------------------------------------- |
| `github.com/larsartmann/go-finding` | Runtime          | Provides the shared `Severity` type    |
| `onsi/ginkgo/v2`               | Testing (dev)    | BDD-style test framework               |
| `onsi/gomega`                  | Assertions (dev) | Matcher library for Ginkgo             |

> **Building note:** JSON marshaling uses `encoding/json/v2`, which requires
> `GOEXPERIMENT=jsonv2` (Go 1.26+). The Nix devShell sets this automatically;
> downstream consumers must set it until `json/v2` graduates from experimental.
> See [AGENTS.md](AGENTS.md) for details.

## License

[MIT](LICENSE) © 2026 Lars Artmann
