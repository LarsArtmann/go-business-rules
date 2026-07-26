# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-07-26 (`master`, v2.0.0 release prepared — version split-brain resolved, stats test hardened, DOMAIN_LANGUAGE filled in)

---

## Integration & release

- [ ] **Integrate into Polish-Customs** as a real-world consumer (replace internal `validation.go`); verify compatibility. Not a blocker for tagging `v2.0.0`, but the first real-world validation of the breaking-change surface.
- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers (documented in `AGENTS.md`). Track the Go release that graduates `json/v2` from experimental and remove the constraint then. Tracked as a long-term item; cannot be resolved until Go ships it.

---

_The release-prep items from the 2026-07-26 docs-health audit — `Version` constant split-brain, branching-flow stats test hardening, and `docs/DOMAIN_LANGUAGE.md` — were resolved for the `v2.0.0` release and are recorded in [`CHANGELOG.md`](CHANGELOG.md). The `json/v2` downstream-constraint item remains open by design (depends on a future Go release)._
