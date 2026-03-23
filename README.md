# businessrules

> **Severity-aware validation for Go** — because not all validation failures are equal.

A Go library that adds severity levels to validation, enabling applications to distinguish between critical errors, warnings, and informational issues. Standard validators return pass/fail; businessrules returns the *degree* of failure.

[![GoDoc](https://pkg.go.dev/badge/github.com/artmann/businessrules.svg)](https://pkg.go.dev/github.com/artmann/businessrules)
[![Go Report Card](https://goreportcard.com/badge/github.com/artmann/businessrules)](https://goreportcard.com/report/github.com/artmann/businessrules)
[![CI](https://github.com/artmann/businessrules/workflows/CI/badge.svg)](https://github.com/artmann/businessrules/actions)

## Installation

```bash
go get github.com/artmann/businessrules
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/artmann/businessrules"
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

func (p Package) Validate() businessrules.ValidationResult {
    return businessrules.NewValidator().
        AddRules(p.Rules()...).
        Build()
}

func main() {
    pkg := Package{Weight: -5, TrackingID: ""}

    result := pkg.Validate()

    if result.HasErrors() {
        for _, err := range result.Errors() {
            fmt.Printf("ERROR: %s\n", err.Message())
        }
    }

    for _, warn := range result.Warnings() {
        fmt.Printf("WARNING: %s\n", warn.Message())
    }

    if !result.HasErrors() {
        processPackage(pkg)
    }
}
```

## Why businessrules?

Standard validators return only valid/invalid. Business rules return the *degree* of failure:

| Output              | Use Case                                              |
| ------------------- | ----------------------------------------------------- |
| **Valid / Invalid** | Binary pass/fail — enough for structural validation    |
| **Severity levels** | Nuanced outcomes — critical for business validation    |

```go
// User submits weight = -5 kg
// → ERROR: Invalid value (blocks processing)

// User submits weight = 2000 kg
// → WARNING: Suspicious but allowed (logs, notifies, allows)

// User submits lowercase "kg" instead of "KG"
// → INFO: Style suggestion (advisory only)
```

### Works with structural validators

businessrules complements type/format validators (like `sivchari/govalid`):

| Layer               | Validator           | Validates                        |
| ------------------- | ------------------- | -------------------------------- |
| Structural          | `govalid`           | Type, format, required fields    |
| Business            | `businessrules`     | Domain rules, severity levels     |

```go
type User struct {
    Email string `govalid:"required,email"`
    Age   int    `govalid:"min=0,max=150"`
}

func (u User) ValidateAll() (*businessrules.ValidationResult, error) {
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
const (
    SeverityInfo     Severity = iota // Advisory: just so you know
    SeverityWarning                  // Non-blocking: should fix
    SeverityError                    // Blocking: must fix
    SeverityCritical                 // Blocking: critical failure
)
```

### ValidationResult

```go
type ValidationResult struct {
    Valid      bool
    Violations []Violation
}

func (r ValidationResult) Errors() []Violation
func (r ValidationResult) Warnings() []Violation
func (r ValidationResult) Info() []Violation
func (r ValidationResult) Critical() []Violation
func (r ValidationResult) BySeverity(severity Severity) []Violation
func (r ValidationResult) HasErrors() bool
func (r ValidationResult) HasWarnings() bool
```

### Pre-built Rules

```go
// Numeric
NonNegative(name string, value float64, severity Severity) Rule
Positive(name string, value float64, severity Severity) Rule
GreaterThan(name string, value, minimum float64, severity Severity) Rule
LessThan(name string, value, maximum float64, severity Severity) Rule
InRange(name string, value, min, max float64, severity Severity) Rule
MinInt(name string, value, min int, severity Severity) Rule
MaxInt(name string, value, max int, severity Severity) Rule

// String
NotEmpty(name string, value string, severity Severity) Rule
NotBlank(name string, value string, severity Severity) Rule
MinLength(name string, value string, min int, severity Severity) Rule
MaxLength(name string, value string, max int, severity Severity) Rule
Matches(name string, value string, pattern *regexp.Regexp, severity Severity) Rule

// Collection
NotEmptySlice[T any](name string, value []T, severity Severity) Rule
NotEmptyMap[T any](name string, value map[string]T, severity Severity) Rule

// Format
Email(name string, value string, severity Severity) Rule
URL(name string, value string, severity Severity) Rule
UUID(name string, value string, severity Severity) Rule

// Generic
Equals[T comparable](name string, value, expected T, severity Severity) Rule
OneOf[T comparable](name string, value T, allowed []T, severity Severity) Rule
Custom(name string, check func() error, severity Severity) Rule

// Composite
All(name string, rules []Rule, severity Severity) Rule
Any(name string, rules []Rule, severity Severity) Rule
When(name string, condition bool, rule Rule) Rule
```

### Validator Builder

```go
result := businessrules.NewValidator().
    AddRule(rule1).
    AddRules(rule2, rule3).
    Build()
```

## Philosophy

- **Zero runtime dependencies** — only standard library
- **Type-safe** — no `any` types
- **Small files** — ≤250 lines per file
- **Small functions** — ≤30 lines per function
- **Composable** — integrates with `sivchari/govalid` for structural validation
- **Tested with Ginkgo/Gomega** — BDD-style testing for behavior specification

## Dependencies

| Dependency       | Purpose          | Notes                      |
| ---------------- | ---------------- | -------------------------- |
| `onsi/ginkgo/v2` | Testing (dev)    | BDD-style test framework   |
| `onsi/gomega`    | Assertions (dev) | Matcher library for Ginkgo |

**Zero runtime dependencies** — only standard library.

## License

MIT
