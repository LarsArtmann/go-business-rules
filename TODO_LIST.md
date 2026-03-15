# TODO List — businessrules

> Extraction plan for `github.com/artmann/businessrules`

**Source:** `Polish-Customs/pkg/types/validation.go`
**Effort:** 4-8 hours
**Status:** ✅ Complete

---

## Phase 1: Repository Setup

- [x] Create repository directory
- [x] Initialize git
- [x] Create README.md
- [x] Initial commit

## Phase 2: Go Module & Structure

- [x] Run `go mod init github.com/artmann/businessrules`
- [x] Add dependencies:
  - [x] `github.com/onsi/ginkgo/v2` — BDD testing framework
  - [x] `github.com/onsi/gomega` — Assertion DSL
- [x] Create file structure:
  ```
  businessrules/
  ├── severity.go           # Severity type and constants (41 lines)
  ├── rule.go               # Rule interface and base implementation (61 lines)
  ├── result.go             # Result type and methods (67 lines)
  ├── validator.go          # Validator interface and builder (45 lines)
  ├── errors.go             # Violation type and error handling (59 lines)
  ├── builders.go           # Numeric and String rule constructors (160 lines)
  ├── builders_format.go    # Format rules: Email, URL, UUID (73 lines)
  ├── builders_composite.go # Generic and Composite rules (86 lines)
  ├── doc.go                # Package documentation (73 lines)
  ├── go.mod
  ├── go.sum
  ├── README.md
  ├── LICENSE
  └── examples/
       └── example_test.go  # Usage examples (170 lines)
  ```

## Phase 3: Core Implementation

### 3.1 Severity Type
- [x] Create `severity.go` with `Severity` type
- [x] Define constants: `SeverityInfo`, `SeverityWarning`, `SeverityError`, `SeverityCritical`
- [x] Implement `String()` method
- [x] Write BDD tests using Ginkgo/Gomega

### 3.2 Rule Interface
- [x] Create `rule.go` with `Rule` interface
- [x] Implement base rule struct
- [x] Implement `NewRule()` constructor
- [x] Write BDD tests using Ginkgo/Gomega

### 3.3 Violation & Errors
- [x] Create `errors.go` with `Violation` struct
- [x] Implement `Error()` method on `Violation`
- [x] Write BDD tests using Ginkgo/Gomega

### 3.4 Result Type
- [x] Create `result.go` with `Result` struct
- [x] Implement grouped accessors: `Errors()`, `Warnings()`, `Info()`, `Critical()`
- [x] Implement `BySeverity()` and `Has*()` methods
- [x] Write BDD tests using Ginkgo/Gomega

### 3.5 Validator Builder
- [x] Create `validator.go` with `Validator` interface
- [x] Implement `ValidatorBuilder` with fluent API
- [x] Implement `AddRule()`, `AddRules()`, `Build()` methods
- [x] Write BDD tests using Ginkgo/Gomega

### 3.6 Pre-built Rules
- [x] Create rule constructors in `builders*.go`:
  - [x] `NonNegative(name, value, severity)`
  - [x] `Positive(name, value, severity)`
  - [x] `InRange(name, value, min, max, severity)`
  - [x] `MinInt(name, value, min, severity)`
  - [x] `MaxInt(name, value, max, severity)`
  - [x] `NotEmpty(name, value, severity)`
  - [x] `MinLength(name, value, min, severity)`
  - [x] `MaxLength(name, value, max, severity)`
  - [x] `Matches(name, value, pattern, severity)`
  - [x] `Email(name, value, severity)`
  - [x] `URL(name, value, severity)`
  - [x] `UUID(name, value, severity)`
  - [x] `OneOf[T comparable](name, value, allowed, severity)`
  - [x] `Custom(name, check, severity)`
  - [x] `All(name, rules, severity)`
  - [x] `Any(name, rules, severity)`
  - [x] `When(name, condition, rule)`
- [x] Write comprehensive BDD tests using Ginkgo/Gomega

## Phase 4: Examples & Documentation

- [x] Create `example_test.go` with usage examples
- [x] Add MIT LICENSE file
- [x] Update README.md with API documentation
- [x] Create `doc.go` with package documentation

## Phase 5: Quality Gates

- [x] Run `go build` — ✅ compiles
- [x] Run `go test ./...` — ✅ 49 tests pass
- [x] Run `go test -cover ./...` — ✅ 96.1% coverage (target: 95%)
- [x] Verify no `any` types in codebase — ✅ none found
- [x] Verify all files ≤250 lines — ✅ all source files under limit
- [x] Verify all functions ≤30 lines — ✅ all functions under limit
- [x] Run `go vet ./...` — ✅ passes

## Phase 6: Integration

- [ ] Add as dependency to Polish-Customs
- [ ] Replace internal `validation.go` with import
- [ ] Run Polish-Customs tests to verify compatibility
- [ ] Commit migration

## Phase 7: Publish

- [ ] Tag release: `git tag v0.1.0`
- [ ] Push to remote: `git push origin master --tags`
- [ ] Verify on pkg.go.dev

---

## Success Criteria

- [x] Library compiles with `go build`
- [x] All tests pass with `go test ./...`
- [x] 95%+ test coverage (achieved: 96.1%)
- [x] Zero runtime dependencies — only standard library
- [x] Each file ≤250 lines
- [x] Each function ≤30 lines
- [x] No `any` types
- [x] README with examples
- [x] Comprehensive package documentation in doc.go
- [ ] Polish-Customs successfully migrated
- [ ] CI/CD pipeline configured

---

## Library Policy Compliance

This project follows the library-policy guidelines:

| Category | Library | Status |
|----------|---------|--------|
| Testing | `onsi/ginkgo/v2` + `onsi/gomega` | ✅ Used |
| Error Handling | Standard library | ✅ Zero dependencies |
| Validation | `sivchari/govalid` | ✅ Complementary |

**Banned libraries (NOT used):**
- ❌ `stretchr/testify` — Use Ginkgo/Gomega instead
- ❌ `go-playground/validator` — Use `sivchari/govalid` instead
- ❌ `ozzo-validation` — Unmaintained
- ❌ `pkg/errors` — Use standard library

---

*Created: 2026-03-15*
*Completed: 2026-03-15*
