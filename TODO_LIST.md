# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-07-26 (`master`, 145 specs pass, 94.8% coverage)

---

## Documentation drift (high value, low effort)

These are factual drifts confirmed against code during the 2026-07-26 docs-health audit.

- [ ] **Fix `Version` constant split-brain** — `doc.go:76` says `"1.1.0"` but the only git tag is `v0.1.0`. Reconcile the version scheme (either tag `v1.1.0` or correct the constant).
- [ ] **Document `go-auto-upgrade` `lo.SliceToMap` false positive** in `AGENTS.md` — `validation_result.go:39-43` uses a manual slice-to-map loop; adding `samber/lo` would break the minimal-dependency principle. Source: `docs/status/2026-07-23_10-11_*.md` §c.7.
- [ ] **Document `go-structure-linter` `root-package-files` false positive** in `AGENTS.md` — root-level `.go` files ARE the public API for this library; moving to `/internal/` or `/pkg/` would break consumers. Source: `docs/status/2026-07-23_10-11_*.md` §c.8.

## Linter false-positive documentation (hierarchical-errors)

The `hierarchical-errors` analyzer reports 6 `generic_return` findings. AGENTS.md documents only 3 and with **stale line numbers**. Verified current locations:

- [ ] **Correct stale line numbers** in AGENTS.md hierarchical-errors section: `errors.go:67`→`:70`, `validation_result.go:152`→`:160`, `rule.go:38`→`:36`.
- [ ] **Document the 3 missing false positives**: `checkNonEmpty` (`builders_format.go:17`), `collectAllViolations` (`builders_composite.go:34`), `anyRulePasses` (`builders_composite.go:51`) — internal helpers aggregating `Rule.Check()` results, so they inherit the `error` return.

## Test robustness

- [ ] **Harden the branching-flow stats test** — `bdd_branching_flow_test.go:163` asserts the substring `"36"` appears anywhere in output; a coincidental duration/count would pass falsely. Parse the stats table's `Total` row instead. Source: `docs/status/2026-07-23_10-11_*.md` §e.8.

## Domain language

- [ ] **Fill in `docs/DOMAIN_LANGUAGE.md`** — currently a placeholder template ("The project/product name", "Example Term"). Define the real ubiquitous language: Rule, Severity, Violation, ValidationResult, etc.

## Integration & release

- [ ] **Integrate into Polish-Customs** as a real-world consumer (replace internal `validation.go`); verify compatibility.
- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers. Track the Go release that graduates `json/v2` from experimental and remove the constraint then.

---

_Items above were harvested from the most recent status report (`docs/status/2026-07-23_10-11_*.md`) and verified still-open against `master` on 2026-07-26. Already-done items were routed to `CHANGELOG.md`._
