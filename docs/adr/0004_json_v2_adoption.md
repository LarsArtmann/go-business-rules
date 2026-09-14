# ADR-0004: `encoding/json/v2` Adoption and the `GOEXPERIMENT` Constraint

**Date:** 2026-07 (per CHANGELOG 2.0.0)
**Status:** Accepted — re-verified 2026-09-14 on Go 1.26.7

## Context

`ViolationError` and `ValidationResultError` serialize for logs and API
responses. `encoding/json/v2` fixes long-standing v1 design issues and is
the fleet policy target for JSON handling on Go 1.25+. It is still
experimental and requires `GOEXPERIMENT=jsonv2` at build time.

## Decision

Use `encoding/json/v2` (the `encoding/json/v1` façade stays importable for
consumers, but this module's own Marshal/Unmarshal paths target v2).

## Consequences

- **Hard downstream constraint:** every consumer must build with
  `GOEXPERIMENT=jsonv2` (and Go 1.26+) or compilation fails with a cryptic
  "build constraints exclude all Go files in encoding/json/v2" error. This
  cannot be contained in `go.mod`.
- The constraint propagates transitively (e.g. via `go-finding`).
- **Re-evaluate when:** a Go release graduates json/v2 from experimental;
  the requirement then disappears for everyone at once.
- Wired in: `flake.nix` (both devShells set `GOEXPERIMENT=jsonv2`),
  `.golangci.yml` (build-tags), CI (workflow env). Affected files:
  `errors.go`, `validation_result.go`, `bdd_branching_flow_test.go`.
