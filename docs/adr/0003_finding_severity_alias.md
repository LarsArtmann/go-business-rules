# ADR-0003: `Severity` Migrated to `finding.Severity` Alias

**Date:** 2026 (commit `e423de4`)
**Status:** Accepted — standing

## Context

This library declared its own `Severity` enum while the fleet-wide
`github.com/larsartmann/go-finding` module defined an identical one. Two
enums for one concept force mapping code at every boundary between
validation results and findings (the downstream common representation).

## Decision

`Severity` is a type ALIAS (`type Severity = finding.Severity`) of the
go-finding enum; the constants re-export it. This keeps:

- zero-diff assignment at boundaries with go-finding consumers;
- the ergonomic unqualified names (`SeverityError`, ...) in this package;
- a runtime dependency on go-finding (the only non-stdlib runtime dep,
  superseding the earlier "zero runtime dependencies" claim).

## Consequences

- A breaking change in go-finding's Severity propagates here by design;
  the alias couples the two modules' release cadence for this one type.
- Reverting to a local enum would reintroduce boundary mapping code; only
  reconsider if go-finding is abandoned.
