# Domain Language

A **Ubiquitous Language** for `businessrules` — shared across maintainers, contributors, and AI sessions.
Inspired by Domain-Driven Design (DDD).

Every term below means the **same thing** to everyone who reads it. The code uses
these names verbatim; this file keeps their definitions stable.

## Glossary

| Term                    | Definition                                                              | Context                                        |
| ----------------------- | ----------------------------------------------------------------------- | ---------------------------------------------- |
| businessrules           | This Go library: severity-aware validation with multiple outcome levels | The project/package name                       |
| Severity                | The importance level of a rule: `info`, `warning`, `error`, `critical`  | Re-exported from `finding.Severity` (a string) |
| Rule                    | A single validation check with a name, severity, and message            | The `Rule` interface                           |
| RuleImpl                | The concrete, immutable base implementation of `Rule`                   | Struct with `WithName`/`WithSeverity`/`WithMessage` |
| ViolationError          | A failed rule check, with context and timestamp; implements `error`     | Produced when a `Rule.Check()` returns non-nil |
| ValidationResultError   | The outcome of validating many rules: validity flag + all violations     | Returned by `ValidatorBuilder.Build()`         |
| ValidatorBuilder        | Fluent builder that collects rules and runs them to produce a result    | `NewValidator().AddRule(...).Build()`          |

## Value Objects

Immutable objects defined by their attributes.

| Term         | Definition                                                      | Context                              |
| ------------ | --------------------------------------------------------------- | ------------------------------------ |
| Severity     | A four-level string ranking of importance                      | `info < warning < error < critical`  |
| ViolationError | Immutable snapshot of one failure (rule + context + timestamp) | Safe to share across goroutines      |

## Severity Levels

The library's central concept — standard validators are binary; `businessrules`
returns the _degree_ of failure.

| Level             | Meaning                                  | Typical response                         |
| ----------------- | ---------------------------------------- | ---------------------------------------- |
| `SeverityInfo`    | Advisory; just so you know               | Log only, never blocks                   |
| `SeverityWarning` | Non-blocking; should be reviewed         | Surface to user, allow processing        |
| `SeverityError`   | Blocking; must be fixed                  | Reject the input                         |
| `SeverityCritical`| Blocking; severe failure, immediate attention | Reject and alert                     |

## Operations

| Term          | Definition                                                       | Context                                    |
| ------------- | ---------------------------------------------------------------- | ------------------------------------------ |
| Check         | Execute a rule's condition; returns `nil` (pass) or `error` (fail) | `Rule.Check()`                         |
| Build         | Run all collected rules and assemble the `ValidationResultError` | `ValidatorBuilder.Build()`                 |
| Filter        | Select violations matching a predicate                           | `ValidationResultError.Filter(...)`        |
| BySeverity    | Select violations matching one or more severity levels           | `ValidationResultError.BySeverity(...)`    |
| Merge         | Combine two results into one                                     | `ValidationResultError.Merge(other)`       |

---

> **How to use this file:**
>
> - Keep terms concise — one clear sentence per definition
> - Update when new domain concepts emerge
> - Use these terms consistently in code, docs, and conversations
> - When in doubt about a word's meaning, check here first
