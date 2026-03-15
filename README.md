# businessrules

> Severity-aware validation for Go — because not all validation failures are equal.

A standalone Go library for validation with severity levels (Info, Warning, Error, Critical). Unlike standard validators that only support pass/fail semantics, businessrules enables nuanced validation outcomes.

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

func (p Package) Validate() businessrules.Result {
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

| Library | Outcome | Missing |
|---------|---------|---------|
| `go-playground/validator` | Valid / Invalid | No severity |
| `ozzo-validation` | Valid / Invalid | No severity |
| `asaskevich/govalidator` | Valid / Invalid | No severity |
| **businessrules** | Valid / Errors / Warnings / Info | **Severity levels** |

### Real-world Example

```go
// User submits weight = -5 kg
// → ERROR (blocking): Invalid value

// User submits weight = 2000 kg (2 tonnes for a package)
// → WARNING (non-blocking): Suspicious but allowed

// User submits weight with lowercase unit "kg" instead of "KG"
// → INFO (advisory): Style suggestion
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

### Core Types

```go
type Rule interface {
    Name() string
    Check() error
    Severity() Severity
    Message() string
}

type Violation struct {
    Rule      Rule
    Context   string
    Timestamp time.Time
}

type Result struct {
    Valid      bool
    Violations []Violation
}

func (r Result) Errors() []Violation
func (r Result) Warnings() []Violation
func (r Result) Info() []Violation
func (r Result) Critical() []Violation
func (r Result) BySeverity(severity Severity) []Violation
func (r Result) HasErrors() bool
func (r Result) HasWarnings() bool
```

### Pre-built Rules

```go
// Numeric
NonNegative(name string, value float64, severity Severity) Rule
Positive(name string, value float64, severity Severity) Rule
InRange(name string, value, min, max float64, severity Severity) Rule
MinInt(name string, value, min int, severity Severity) Rule
MaxInt(name string, value, max int, severity Severity) Rule

// String
NotEmpty(name string, value string, severity Severity) Rule
MaxLength(name string, value string, max int, severity Severity) Rule
Matches(name string, value string, pattern *regexp.Regexp, severity Severity) Rule

// Generic
OneOf[T comparable](name string, value T, allowed []T, severity Severity) Rule
Custom(name string, check func() error, severity Severity) Rule
```

### Validator Builder

```go
result := businessrules.NewValidator().
    AddRule(rule1).
    AddRules(rule2, rule3).
    Build()
```

## Integration with go-playground/validator

businessrules complements, not replaces, tag-based validators:

```go
import (
    "github.com/go-playground/validator/v10"
    "github.com/artmann/businessrules"
)

type User struct {
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=0,lte=150"`
}

func (u User) ValidateAll() (*businessrules.Result, error) {
    // 1. Structural validation (tags)
    if err := validator.New().Struct(u); err != nil {
        return nil, err
    }

    // 2. Business validation (severity-aware)
    result := businessrules.NewValidator().
        AddRule(businessrules.Custom("email_domain", func() error {
            if !strings.HasSuffix(u.Email, "@company.com") {
                return errors.New("must be company email")
            }
            return nil
        }, businessrules.SeverityWarning)). // Non-blocking!
        Build()

    return &result, nil
}
```

| Validator | Use Case |
|-----------|----------|
| `go-playground/validator` | Structural validation (required, format, type) |
| `businessrules` | Business validation (domain rules, severity levels) |

## Philosophy

- **Zero dependencies** (except `cockroachdb/errors`)
- **Type-safe** — no `any` types
- **Small files** — ≤250 lines per file
- **Small functions** — ≤30 lines per function
- **Composable** — use with existing validators

## License

MIT
