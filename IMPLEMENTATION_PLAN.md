# Implementation Plan — businessrules

> Extraction plan for `github.com/artmann/businessrules`
> **Status:** Ready for extraction
> **Created:** 2026-03-15

---

## 1. Repository Setup

- [x] Create repository directory
- [x] Initialize git
- [x] Create README.md
- [x] Initial commit

- [x] Create TODO_LIST.md
- [x] Update docs for library-policy compliance

- [x] Commit documentation changes

- [ ] Push to remote

---

## 2. Core Types (Priority: HIGH - Foundation)

Priority: Critical

Effort: 10-15 min each

Impact: Must-have for all features

---

| ID  | Task                                              | Category   | Effort   | Impact   | Dependencies                         |
| --- | ------------------------------------------------- | ---------- | -------- | -------- | ------------------------------------ | --------------------------------------------- | ---------- | -------- | -------- | -------------------------------------- | --------------------------------------------- | ---------- | -------- | -------- | ---------------------------------- | ---------- | ----- | -------- | -------- | -------------------------------------- | ---------- | ----- | -------- | ---- |
| 1   | Initialize Go module (`go mod init`)              | Setup      | 5min     | Critical | -                                    |
| 2   | Add `cockroachdb/errors` dependency               | Setup      | 5min     | Critical | 1                                    |
| 3   | Add `onsi/ginkgo/v2` + `onsi/gomega` dependencies | Setup      | 10min    | Critical | 1, 2                                 |
| 4   | Create `severity.go` - Core Types                 | 15min      | Critical | 2        |
| 5   | Implement `Severity.String()` method              | Core Types | 5min     | Critical | 4                                    |
| 6   | Create `rule.go` with `Rule` interface            | Core Types | 20min    | Critical | 4, 5                                 |
| 7   | Implement `NewRule()` constructor                 | Core Types | 10min    | Critical | 6                                    |
| 8   | Implement base rule struct                        | Core Types | 15min    | Critical | 6, 7, 8                              |
| 9   | Create `errors.go` with `Violation` struct        | Core Types | 10min    | Critical | 6-8                                  |
| 10  | Implement `Violation.Error()` method              | Core Types | 5min     | Critical | 6-8, 11                              | Create `result.go` with `Result` struct       | Core Types | 15min    | Critical | 6-10                                   |
| 12  | Implement `Result.Valid` field and Core Types     | 5min       | Critical | 6-8, 13  | Implement `Result.Errors()` accessor | Core Types                                    | 10min      | Critical | 6-10, 14 | Implement `Result.Warnings()` accessor | Core Types                                    | 10min      | Critical | 6-10, 16 | Implement `Result.Info()` accessor | Core Types | 10min | Critical | 6-10, 18 | Implement `Result.Critical()` accessor | Core Types | 10min | Critical | 6-10 |
| 20  | Implement `Result.BySeverity()` method            | Core Types | 10min    | Critical | 6-10                                 |
| 22  | Implement `Result.HasErrors()` method             | Core Types | 5min     | Critical | 6-10                                 |
| 24  | Implement `Result.HasWarnings()` method           | Core Types | 5min     | Critical | 6-10                                 |
| 26  | Implement `Result.HasCritical()` method           | Core Types | 5min     | Critical | 6-10, 28                             | Create `validator.go` with `ValidatorBuilder` | Core Types | 15min    | Critical | 1-10, 26- 29                           | Implement `ValidatorBuilder.AddRule()` method | Core Types | 10min    | Critical | 28                                 |
| 30  | Implement `ValidatorBuilder.AddRules()` method    | Core Types | 10min    | Critical | 28                                   |
| 31  | Implement `ValidatorBuilder.Build()` method       | Core Types | 10min    | Critical | 28-32                                |
| 32  | Commit core types with tests                      | Core Types | 5min     | Critical | 1-10                                 |

| **Subtotal: Core Types Phase** | **~3.5 hours** | **~3.5 hours** | **15 tasks × 15 min each** |

---

## 3. Pre-built Rules (Priority: HIGH - Immediate value)

Priority: High
Effort: 10-20 min each
Impact: Quick wins, immediate productivity

Dependencies: Rule, Severity types
| ID | Task | Category | Effort | Impact | Dependencies |
|----|------|----------|----------|---------|--------------|
| 33 | Implement `NonNegative()` rule | Builders | 15min | High | severity.go |
| 34 | Implement `Positive()` rule | Builders | 10min | High | severity.go |
| 35 | Implement `InRange()` rule | Builders | 15min | High | severity.go |
| 36 | Implement `MinInt()` rule | Builders | 10min | High | severity.go |
| 37 | Implement `MaxInt()` rule | Builders | 10min | High | severity.go |
| 38 | Implement `NotEmpty()` rule | Builders | 10min | High | severity.go |
| 39 | Implement `MaxLength()` rule | Builders | 10min | High | severity.go |
| 40 | Implement `Matches()` rule | Builders | 15min | High | severity.go |
| 41 | Implement `OneOf[T comparable]()` rule | Builders | 15min | High | severity.go |
| 42 | Implement `Custom()` rule | Builders | 10min | High | severity.go |
| 43 | Commit pre-built rules with tests | Builders | 5min | High | 33-42 |

| **Subtotal: Builders Phase** | **~3 hours** | **~2.5 hours** | **10 tasks × 15 min each** |

---

## 4. Testing (Priority: HIGH - Quality Gates)

Priority: High
Effort: 15-30 min each
Impact: Ensure correctness, prevent regressions
Dependencies: Core implementation |
| ID | Task | Category | Effort | Impact | Dependencies |
|----|------|----------|----------|---------|--------------|
| 44 | Create `severity_test.go` with Ginkgo suite | Testing | 15min | High | severity.go |
| 45 | Write tests for `SeverityInfo` constant | Testing | 10min | High | 44 |
| 46 | Write tests for `SeverityWarning` constant | Testing | 10min | High | 44 |
| 47 | Write tests for `SeverityError` constant | Testing | 10min | High | 44 |
| 48 | Write tests for `SeverityCritical` constant | Testing | 10min | High | 44 |
| 49 | Write tests for `Severity.String()` method | Testing | 10min | High | 44 |
| 50 | Create `rule_test.go` with Ginkgo suite | Testing | 15min | High | rule.go |
| 51 | Write tests for `NewRule()` constructor | Testing | 10min | High | 50 |
| 52 | Write tests for `Rule` interface methods | Testing | 15min | High | 50 |
| 53 | Write tests for base rule struct | Testing | 10min | High | 50-54 |
| 54 | Create `errors_test.go` with Ginkgo suite | Testing | 15min | High | errors.go |
| 55 | Write tests for `Violation` struct | Testing | 10min | High | 55 |
| 56 | Write tests for `Violation.Error()` method | Testing | 10min | High | 55-57 |
| 57 | Create `result_test.go` with Ginkgo suite | Testing | 15min | High | result.go |
| 58 | Write tests for `Result.Valid` field | Testing | 10min | High | 57-59 |
| 59 | Write tests for `Result.Errors()` accessor | Testing | 10min | High | 57-59 |
| 60 | Write tests for `Result.Warnings()` accessor | Testing | 10min | High | 57-59 |
| 61 | Write tests for `Result.Info()` accessor | Testing | 10min | High | 57-59 |
| 62 | Write tests for `Result.Critical()` accessor | Testing | 10min | High | 57-59 |
| 63 | Write tests for `Result.BySeverity()` method | Testing | 10min | High | 57-59 |
| 64 | Write tests for `Result.Has*()` methods | Testing | 10min | High | 57-64 |
| 65 | Create `validator_test.go` with Ginkgo suite | Testing | 15min | High | validator.go |
| 66 | Write tests for `ValidatorBuilder` struct | Testing | 10min | High | 65 |
| 67 | Write tests for `ValidatorBuilder.AddRule()` method | Testing | 10min | High | 66 |
| 68 | Write tests for `ValidatorBuilder.addRules()` method | Testing | 10min | High | 66-68 |
| 69 | Write tests for `ValidatorBuilder.Build()` method | Testing | 15min | High | 65-68 |
| 70 | Create `builders_test.go` with Ginkgo suite | Testing | 20min | High | builders.go |
| 71 | Write tests for `NonNegative()` rule | Testing | 10min | High | 71 |
| 72 | Write tests for `Positive()` rule | Testing | 10min | High | 71 |
| 73 | Write tests for `InRange()` rule | Testing | 10min | High | 71-73 |
| 74 | Write tests for `MinInt()` rule | Testing | 10min | High | 74 |
| 75 | Write tests for `MaxInt()` rule | Testing | 10min | High | 74-76 |
| 76 | Write tests for `NotEmpty()` rule | Testing | 10min | High | 76 |
| 77 | Write tests for `MaxLength()` rule | Testing | 10min | High | 77-78 |
| 78 | Write tests for `Matches()` rule | Testing | 10min | High | 78-79 |
| 79 | Write tests for `OneOf()` rule | Testing | 10min | High | 79-80 |
| 80 | Write tests for `Custom()` rule | Testing | 10min | High | 80 |
| 81 | Run all tests with `go test ./...` | Testing | 5min | Critical | 44-80 |
| **Subtotal: Testing Phase** | **~3 hours** | **~3 hours** | **38 tasks × 5-15 min each** |

---

## 5. Examples (Priority: MEDIUM - Documentation)

Priority: Medium
Effort: 15-30 min each
Impact: User onboarding, customer confidence
Dependencies: All implementation |
| ID | Task | Category | Effort | Impact | Dependencies |
|----|------|----------|----------|---------|--------------|
| 82 | Create `examples/` directory | Examples | 5min | Medium | - |
| 83 | Create `examples/basic_test.go` | Examples | 15min | Medium | Core implementation |
| 84 | Write package example with full usage | Examples | 15min | Medium | 83 |
| 85 | Write result grouping example | Examples | 10min | Medium | 83-84 |
| 86 | Write severity filtering example | Examples | 10min | Medium | 83-86 |
| 87 | Commit examples with tests | Examples | 5min | Medium | 82-86 |
| **Subtotal: Examples Phase** | **~1 hour** | **~1 hour** | **6 tasks × 5-15 min each** |

---

## 6. Quality Gates (Priority: MEDIUM - Validation)

Priority: Medium
Effort: 5-15 min each
Impact: Production readiness
Dependencies: All implementation |
| ID | Task | Category | Effort | Impact | Dependencies |
|----|------|----------|----------|---------|--------------|
| 88 | Verify all files ≤250 lines | Quality | 5min | Critical | - |
| 89 | Verify all functions ≤30 lines | Quality | 5min | Critical | - |
| 90 | Verify no `any` types in codebase | Quality | 5min | Critical | - |
| 91 | Run `go build ./...` | Quality | 5min | Critical | - |
| 92 | Run `go vet ./...` | Quality | 5min | Critical | - |
| 93 | Run `go test -cover ./...` | Quality | 10min | High | - |
| 94 | Verify 95%+ test coverage | Quality | 5min | High | 88-93 |
| **Subtotal: Quality Gates Phase** | **~45 min** | **~45 min** | **6 tasks x 5-10 min each** |

---

## 7. Final Integration (Priority: LOW - Future)

Priority: Low
Effort: 30-60 min each
Impact: Integration, publication
Dependencies: Quality gates pass |
| ID | Task | Category | Effort | Impact | Dependencies |
|----|------|----------|----------|---------|--------------|
| 95 | Create MIT LICENSE file | Publish | 5min | Low | - |
| 96 | Tag release: `git tag v0.1.0` | Publish | 5min | Low | 95 |
| 97 | Push to remote repository | Publish | 10min | Low | 96 |
| **Subtotal: Publish Phase** | **~20 min** | **~20 min** | **3 tasks × 5-10 min each** |

---

## Summary by Category

| Category      | Tasks  | Total Effort  |
| ------------- | ------ | ------------- |
| Core Types    | 32     | ~3.5 hours    |
| Builders      | 10     | ~2.5 hours    |
| Testing       | 38     | ~3 hours      |
| Examples      | 6      | ~1 hour       |
| Quality Gates | 6      | ~45 min       |
| Publish       | 3      | ~20 min       |
| **TOTAL**     | **95** | **~11 hours** |

---

## Execution Strategy

- Work in order: Core Types → Builders → Testing → Examples → Quality Gates → Publish
- Commit after each major milestone
- Run tests frequently to ensure no regressions
- Keep files small and focused

---

## Next Steps

1. ✅ Setup (Tasks 1-3) - READY
2. ⏳ Core Types (Tasks 4-32) - NEXT
3. ⏳ Builders (Tasks 33-42) - NEXT
4. ⏳ Testing (Tasks 43-80) - NEXT
5. ⏳ Examples (Tasks 81-86) - NEXT
6. ⏳ Quality Gates (Tasks 87-92) - NEXT
7. ⏳ Publish (Tasks 93-95) - LAST
