# TODO List — businessrules

> Extraction plan for `github.com/artmann/businessrules`

**Source:** `Polish-Customs/pkg/types/validation.go`
**Effort:** 4-8 hours
**Status:** Ready for extraction

---

## Phase 1: Repository Setup

- [x] Create repository directory
- [x] Initialize git
- [x] Create README.md
- [x] Initial commit

## Phase 2: Go Module & Structure

- [ ] Run `go mod init github.com/artmann/businessrules`
- [ ] Add dependency: `go get github.com/cockroachdb/errors`
- [ ] Create file structure:
  ```
  businessrules/
  ├── severity.go          # Severity type and constants (20 lines)
  ├── severity_test.go     # Severity tests (30 lines)
  ├── rule.go              # Rule interface and base implementation (50 lines)
  ├── rule_test.go         # Rule tests (80 lines)
  ├── builders.go          # Pre-built rule constructors (100 lines)
  ├── builders_test.go     # Builder tests (150 lines)
  ├── result.go            # Result type and methods (60 lines)
  ├── result_test.go       # Result tests (100 lines)
  ├── validator.go         # Validator interface and builder (40 lines)
  ├── validator_test.go    # Validator tests (80 lines)
  ├── errors.go            # Violation type and error handling (30 lines)
  ├── errors_test.go       # Error tests (40 lines)
  ├── go.mod
  ├── go.sum
  ├── README.md
  ├── LICENSE
  └── examples/
      └── basic_test.go
  ```

## Phase 3: Core Implementation

### 3.1 Severity Type
- [ ] Create `severity.go` with `Severity` type
- [ ] Define constants: `SeverityInfo`, `SeverityWarning`, `SeverityError`, `SeverityCritical`
- [ ] Implement `String()` method
- [ ] Write tests in `severity_test.go`

### 3.2 Rule Interface
- [ ] Create `rule.go` with `Rule` interface
- [ ] Implement base rule struct
- [ ] Implement `NewRule()` constructor
- [ ] Write tests in `rule_test.go`

### 3.3 Violation & Errors
- [ ] Create `errors.go` with `Violation` struct
- [ ] Implement `Error()` method on `Violation`
- [ ] Write tests in `errors_test.go`

### 3.4 Result Type
- [ ] Create `result.go` with `Result` struct
- [ ] Implement grouped accessors: `Errors()`, `Warnings()`, `Info()`, `Critical()`
- [ ] Implement `BySeverity()` and `Has*()` methods
- [ ] Write tests in `result_test.go`

### 3.5 Validator Builder
- [ ] Create `validator.go` with `Validator` interface
- [ ] Implement `ValidatorBuilder` with fluent API
- [ ] Implement `AddRule()`, `AddRules()`, `Build()` methods
- [ ] Write tests in `validator_test.go`

### 3.6 Pre-built Rules
- [ ] Create `builders.go` with rule constructors:
  - [ ] `NonNegative(name, value, severity)`
  - [ ] `Positive(name, value, severity)`
  - [ ] `InRange(name, value, min, max, severity)`
  - [ ] `MinInt(name, value, min, severity)`
  - [ ] `MaxInt(name, value, max, severity)`
  - [ ] `NotEmpty(name, value, severity)`
  - [ ] `MaxLength(name, value, max, severity)`
  - [ ] `Matches(name, value, pattern, severity)`
  - [ ] `OneOf[T comparable](name, value, allowed, severity)`
  - [ ] `Custom(name, check, severity)`
- [ ] Write comprehensive tests in `builders_test.go`

## Phase 4: Examples & Documentation

- [ ] Create `examples/basic_test.go` with usage example
- [ ] Add MIT LICENSE file
- [ ] Update README.md with any API changes

## Phase 5: Quality Gates

- [ ] Run `go build` — must compile
- [ ] Run `go test ./...` — all tests pass
- [ ] Run `go test -cover ./...` — achieve 95%+ coverage
- [ ] Verify no `any` types in codebase
- [ ] Verify all files ≤250 lines
- [ ] Verify all functions ≤30 lines
- [ ] Run `go vet ./...`
- [ ] Run `staticcheck ./...` (if available)

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

- [ ] Library compiles with `go build`
- [ ] All tests pass with `go test ./...`
- [ ] 95%+ test coverage
- [ ] No external dependencies except `cockroachdb/errors`
- [ ] Each file ≤250 lines
- [ ] Each function ≤30 lines
- [ ] No `any` types
- [ ] README with examples
- [ ] Polish-Customs successfully migrated
- [ ] CI/CD pipeline configured

---

*Created: 2026-03-15*
