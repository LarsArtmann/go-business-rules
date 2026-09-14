# ADR-0002: Type Renames to `ViolationError` / `ValidationResultError`

**Date:** 2026-07 (rename commit `1f2976d`)
**Status:** Accepted — standing

## Context

The library's original outcome types carried generic names that did not say
what they contained. Generic names force every reader to re-learn the type,
and they collide with the vocabulary of consumer domains.

## Decision

- `Violation` → `ViolationError`: a value represents ONE failed rule check,
  carries context (rule, message, timestamp), and participates in Go error
  semantics.
- `ValidationResult` → `ValidationResultError`: the outcome type of a run IS
  an error-carrying aggregate (Valid + []ViolationError); the name states it.

## Consequences

- Names now match Go's `*Error` convention for error values.
- Consumers' code reads unambiguously: `ValidationResultError.ViolationErrors`.
- The rename predates the `/v2` module path; the breaking change shipped
  inside the same major version transition documented in CHANGELOG 2.1.0.
