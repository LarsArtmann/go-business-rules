# go-composable-business-types/id Integration Analysis

> **Status:** ✅ Analysis Complete — Recommendation: **Use Sparingly**\
> **Date:** 2026-03-18\
> **Analyst:** Crush

---

## Executive Summary

The `go-composable-business-types/id` library provides **branded, strongly-typed identifiers** that prevent mixing different entity IDs at compile time. After comprehensive analysis, this library has **limited but valuable** applications in the businessrules project.

**Verdict:** Integrate for **RuleID** type only. Other identifier types (ViolationID, ContextID) provide marginal value at the cost of API complexity.

---

## 1. Library Overview

### What It Provides

| Feature                    | Benefit                                |
| -------------------------- | -------------------------------------- |
| Branded types via generics | Compile-time prevention of ID mixing   |
| Zero runtime dependencies  | Maintains project philosophy           |
| JSON/SQL serialization     | Seamless integration with storage      |
| Type-safe comparisons      | `ruleID1.Equal(ruleID2)` vs `==`       |
| Ordering support           | Sortable IDs for `int`, `string`, etc. |

### Installation

```bash
go get github.com/larsartmann/go-composable-business-types/id
```

### Core Pattern

```go
// Without branded IDs (current) — runtime bugs possible
func GetRule(name string) Rule { ... }
GetRule(userEmail)  // Compiles! Wrong semantic usage.

// With branded IDs — compile-time safety
type RuleBrand struct{}
type RuleID = id.ID[RuleBrand, string]

func GetRule(id RuleID) Rule { ... }
GetRule(id.NewID[RuleBrand]("required_field"))  // Correct
GetRule(userEmail)  // Compile error: type mismatch
```

---

## 2. Current Codebase Analysis

### Identifier Usage Patterns

| File                   | Current Type | Purpose                                      |
| ---------------------- | ------------ | -------------------------------------------- |
| `rule.go:9`            | `string`     | `Name() string` — rule identifier            |
| `rule.go:54`           | `string`     | `NewRule(name, ...)` — rule constructor      |
| `builders.go:*`        | `string`     | All builder functions take `name string`     |
| `errors.go:25`         | `string`     | Error message interpolation: `v.Rule.Name()` |
| `validation_result.go` | N/A          | No identifiers currently                     |

### Current API Surface

```go
// rule.go
func NewRule(name string, check func() error, severity Severity, message string) Rule

// builders.go (14 functions)
func NonNegative(name string, value float64, severity Severity) Rule
func Positive(name string, value float64, severity Severity) Rule
func NotEmpty(name, value string, severity Severity) Rule
// ... 11 more

// builders_composite.go (4 functions)
func OneOf[T comparable](name string, value T, allowed []T, severity Severity) Rule
func Custom(name string, check func() error, severity Severity) Rule
func All(name string, rules []Rule, severity Severity) Rule
func Any(name string, rules []Rule, severity Severity) Rule
func When(name string, condition bool, rule Rule) Rule
```

**Total:** 19 functions accepting `name string` parameter.

---

## 3. Integration Scenarios

### Scenario A: RuleID Only (Recommended)

Replace rule names with branded `RuleID` type.

**Changes Required:**

1. Add dependency to `go.mod`
2. Create `types.go` with RuleID type
3. Update `Rule.Name()` return type
4. Update `NewRule()` parameter type
5. Update all 19 builder functions
6. Update tests

**Before:**

```go
rule := businessrules.NotEmpty("email", user.Email, businessrules.SeverityError)
```

**After:**

```go
type EmailRule struct{}
var RuleEmail = id.NewID[EmailRule]("email")

rule := businessrules.NotEmpty(RuleEmail, user.Email, businessrules.SeverityError)
```

**Pros:**

- Compile-time safety for rule identification
- Prevents mixing rule names with field names, error messages
- Aligns with domain-driven design principles

**Cons:**

- Breaking API change (major version bump)
- More verbose usage (requires defining brand types)
- All consumers must update

**Impact:** High effort, moderate benefit.

---

### Scenario B: RuleID + ViolationID

Add branded IDs for violations in addition to rules.

**Analysis:**

- Violations are ephemeral (created during validation)
- Typically not stored or referenced by ID
- Current `Violation` struct has no ID field
- **Recommendation:** Not justified. Adds complexity without clear benefit.

---

### Scenario C: RuleID + SeverityID

Replace `Severity` (currently `int` based) with branded ID.

**Analysis:**

- `Severity` is already type-safe via custom type: `type Severity int`
- iota constants provide compile-time safety
- No string mixing risk
- **Recommendation:** Not applicable. Current design is superior.

---

### Scenario D: No Integration

Maintain current `string`-based identifiers.

**Pros:**

- Simple API (strings are familiar)
- No dependencies
- No breaking changes

**Cons:**

- No compile-time prevention of semantic mixing
- Rule names could theoretically be confused with field names

**Assessment:** Current design is defensible. Risk of bugs is low in practice.

---

## 4. Decision Matrix

| Scenario            | Safety | Complexity | Effort | Breaking | Verdict            |
| ------------------- | ------ | ---------- | ------ | -------- | ------------------ |
| A: RuleID only      | ⭐⭐⭐ | ⭐⭐       | ⭐⭐⭐ | Yes      | ✅ **Recommended** |
| B: Rule + Violation | ⭐⭐   | ⭐⭐⭐     | ⭐⭐⭐ | Yes      | ❌ Over-engineered |
| C: SeverityID       | ⭐     | ⭐⭐       | ⭐⭐   | Yes      | ❌ Not applicable  |
| D: No change        | ⭐⭐   | ⭐         | ⭐     | No       | ✅ Acceptable      |

---

## 5. Recommended Implementation

### Phase 1: Type Definition (Minimal)

Create `types.go` with RuleID branded type:

```go
package businessrules

import "github.com/larsartmann/go-composable-business-types/id"

// RuleBrand is the brand type for RuleID.
// Use to create distinct rule identifiers.
type RuleBrand struct{}

// RuleID is a branded identifier for validation rules.
// Prevents mixing rule names with other string values at compile time.
type RuleID = id.ID[RuleBrand, string]

// NewRuleID creates a new RuleID from a string value.
func NewRuleID(value string) RuleID {
    return id.NewID[RuleBrand](value)
}
```

### Phase 2: API Migration (Breaking)

Update `Rule` interface and implementations:

```go
// rule.go

// Rule defines the interface for a validation rule.
type Rule interface {
    // ID returns the branded identifier for this rule.
    ID() RuleID

    // Name returns the string representation of the rule ID.
    // Deprecated: Use ID().String() instead.
    Name() string

    Check() error
    Severity() Severity
    Message() string
}

// Update NewRule to accept RuleID
func NewRule(id RuleID, check func() error, severity Severity, message string) Rule
```

### Phase 3: Builder Updates

Update all builder functions:

```go
// builders.go
func NonNegative(id RuleID, value float64, severity Severity) Rule
func NotEmpty(id RuleID, value string, severity Severity) Rule
// ... etc
```

### Phase 4: Usage Pattern

Consumers define domain-specific rule IDs:

```go
// In consumer's domain package
package domain

import (
    "github.com/artmann/businessrules"
    "github.com/larsartmann/go-composable-business-types/id"
)

// Define rule brands
type UserEmailRule struct{}
type UserAgeRule struct{}

// Create rule IDs
var (
    RuleUserEmail = businessrules.NewRuleID("user_email")
    RuleUserAge   = businessrules.NewRuleID("user_age")
)

// Usage
rules := []businessrules.Rule{
    businessrules.Email(RuleUserEmail, user.Email, businessrules.SeverityError),
    businessrules.NonNegative(RuleUserAge, float64(user.Age), businessrules.SeverityWarning),
}
```

---

## 6. Migration Path

### Option 1: Gradual Deprecation (Preferred)

1. Add `RuleID` type and `ID()` method
2. Keep `Name() string` as deprecated alias
3. Add new builder functions with `RuleID` suffix:
   ```go
   func NotEmptyRule(id RuleID, value string, severity Severity) Rule
   ```
4. Deprecate old functions with `// Deprecated:` comments
5. Remove in v2.0.0

### Option 2: Hard Breaking Change

1. Replace all `string name` parameters with `RuleID`
2. Bump major version: `v2.0.0`
3. Provide migration guide

---

## 7. Risk Assessment

| Risk                   | Likelihood | Impact | Mitigation                               |
| ---------------------- | ---------- | ------ | ---------------------------------------- |
| Consumer resistance    | Medium     | Medium | Deprecation cycle, clear migration guide |
| Increased verbosity    | High       | Low    | Acceptable trade-off for safety          |
| Dependency maintenance | Low        | Low    | Library is stable, zero deps             |
| Breaking existing code | High       | High   | Major version bump, clear changelog      |

---

## 8. Conclusion

### Recommendation

**Integrate `go-composable-business-types/id` for RuleID only.**

**Rationale:**

1. **Aligns with project philosophy:** Zero runtime dependencies (library maintains this)
2. **Improves type safety:** Prevents semantically incorrect string mixing
3. **Domain clarity:** Rule identifiers are distinct from field names, error messages
4. **Future-proof:** Enables richer rule metadata, rule registries, rule relationships

**Implementation Priority:**

| Priority | Task                                   | Effort    |
| -------- | -------------------------------------- | --------- |
| P1       | Add dependency, create RuleID type     | 30 min    |
| P2       | Implement gradual deprecation          | 2-3 hours |
| P3       | Update documentation, examples         | 1 hour    |
| P4       | Release v1.x with deprecation, plan v2 | —         |

**Alternative:** If immediate breaking changes are unacceptable, defer integration until v2.0.0 planning.

---

## 9. Appendix: Library Compatibility

### Dependency Check

| Library                         | Deps      | Compatible |
| ------------------------------- | --------- | ---------- |
| businessrules                   | 0 runtime | ✅         |
| go-composable-business-types/id | 0 runtime | ✅         |
| Combined                        | 0 runtime | ✅         |

### Go Version Requirements

| Project                      | Go Version |
| ---------------------------- | ---------- |
| businessrules                | 1.25.0     |
| go-composable-business-types | 1.23+      |
| Compatibility                | ✅ Full    |

---

## 10. Appendix: Example Integration

### Complete Integration Example

```go
// types_id.go — new file
package businessrules

import (
    "fmt"
    "github.com/larsartmann/go-composable-business-types/id"
)

// RuleBrand is the brand type for RuleID.
type RuleBrand struct{}

// RuleID is a branded, strongly-typed identifier for validation rules.
type RuleID = id.ID[RuleBrand, string]

// NewRuleID creates a RuleID from a string value.
func NewRuleID(value string) RuleID {
    return id.NewID[RuleBrand](value)
}

// MustRuleID panics if the value is empty.
func MustRuleID(value string) RuleID {
    if value == "" {
        panic("RuleID value cannot be empty")
    }
    return NewRuleID(value)
}

// RuleIDFromString is an alias for NewRuleID for explicit conversion.
func RuleIDFromString(value string) RuleID {
    return NewRuleID(value)
}
```

```go
// rule.go — updated
package businessrules

// Rule defines the interface for a validation rule.
type Rule interface {
    // ID returns the branded identifier for this rule.
    ID() RuleID

    // Name returns the string representation (deprecated).
    // Deprecated: Use ID().Get() or ID().String()
    Name() string

    Check() error
    Severity() Severity
    Message() string
}

// baseRule implementation
type baseRule struct {
    id       RuleID
    check    func() error
    severity Severity
    message  string
}

func (r baseRule) ID() RuleID     { return r.id }
func (r baseRule) Name() string   { return r.id.Get() }
func (r baseRule) Check() error   { return r.check() }
func (r baseRule) Severity() Severity { return r.severity }
func (r baseRule) Message() string { return r.message }

// NewRule creates a Rule with the given parameters.
func NewRule(id RuleID, check func() error, severity Severity, message string) Rule {
    return baseRule{
        id:       id,
        check:    check,
        severity: severity,
        message:  message,
    }
}
```

```go
// builders.go — updated signature
func NonNegative(id RuleID, value float64, severity Severity) Rule {
    return NewRule(
        id,
        func() error {
            if value < 0 {
                return fmt.Errorf("%s must be non-negative, got %f", id.Get(), value)
            }
            return nil
        },
        severity,
        id.Get()+" must be non-negative",
    )
}
```

---

_End of Analysis_
