# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-07-26 (`master`, 145 specs pass, 94.8% coverage)

---

## Versioning & release

- [ ] **Fix `Version` constant split-brain** — `doc.go:76` says `"1.1.0"` but the only git tag is `v0.1.0`. Reconcile the version scheme (either tag the current state or correct the constant). The versioning history is documented in `CHANGELOG.md`.

## Test robustness

- [ ] **Harden the branching-flow stats test** — `bdd_branching_flow_test.go:163` asserts the substring `"36"` appears anywhere in output; a coincidental duration/count would pass falsely. Parse the stats table's `Total` row instead. Source: `docs/status/2026-07-23_10-11_*.md` §e.8.

## Domain language

- [ ] **Fill in `docs/DOMAIN_LANGUAGE.md`** — currently a placeholder template ("The project/product name", "Example Term"). Define the real ubiquitous language: Rule, Severity, ViolationError, ValidationResultError, ValidatorBuilder.

## Integration & release

- [ ] **Integrate into Polish-Customs** as a real-world consumer (replace internal `validation.go`); verify compatibility.
- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers. Track the Go release that graduates `json/v2` from experimental and remove the constraint then.

---

_Items were harvested from the most recent status report (`docs/status/2026-07-23_10-11_*.md`) and verified against `master` on 2026-07-26. Items already resolved during the 2026-07-26 docs-health audit (hierarchical-errors line-number fixes, undocumented `go-auto-upgrade` / `root-package-files` / internal-helper false positives) were applied directly to `AGENTS.md` rather than tracked here._
