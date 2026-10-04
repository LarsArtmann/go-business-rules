# Roadmap — businessrules

> Long-term direction and raw ideas **not yet refined into actionable tasks**.
>
> When an item here becomes bounded and short-term, it graduates into [`TODO_LIST.md`](TODO_LIST.md). When an item ships, it is recorded in [`CHANGELOG.md`](CHANGELOG.md) and reflected in [`FEATURES.md`](FEATURES.md).

**Last reviewed:** 2026-10-04 (harvest of the toolchain-skew remediation session — upstream/tooling asks routed here, bounded work to [`TODO_LIST.md`](TODO_LIST.md), three new open questions below; shipped items continue to land in [`CHANGELOG.md`](CHANGELOG.md)).

---

## Additional Rule Builders (remaining candidates)

### Network / Identifier

- `PostalCode(name, value, country, severity)` — country-specific postal patterns (the shipped `PostalCode` is deliberately generic; add per-country patterns only when a consumer needs one)

### Miscellaneous

- `UUIDv4(name, value, severity)` — version-checked UUID variant (shipped `UUID` accepts any version)
- `SemVer(name, value, severity)` — semantic-version format
- `Latitude`/`Longitude` range rules

## Advanced Features (candidates)

### Rule Composition (remaining)

- `Nor`/`Nand`/`Implies` — only with a demonstrated use case; `Not`+`All`/`Any` cover them compositionally today

### Rule Metadata (remaining)

- `Priority()` — execution-order hints for `Stream` (shipped: `Description()`, `Tags()`)

### Context-Aware Validation (beyond `ContextRule`)

- `WithContext(key, value)` — typed values passed to rules through a validation-session context
- Rules that access shared session state during a run

### Asynchronous Rules

- `Async(name, asyncCheck, severity)` — rules that check asynchronously
- `Await()` method on `ValidationResultError`

### Validation Groups

- Groups of rules enabled/disabled together
- Group inheritance for complex validation scenarios

## Event-Driven Validation (remaining candidates, from 2026-09-14)

- Rule scheduler with priorities / short-circuit on `Critical` during `Stream`
- Listener backpressure / async queue semantics (listeners are synchronous today, by design)
- Signing validation events (go-cqrs-lite recipes) for tamper-evident audit trails
- Run ID (unique per validation run) on `ValidationCompleted` for distributed provenance
- Event JSON marshaling (deliberate YAGNI today — add only when a consumer asks)
- `RuleID` branded type via `go-composable-business-types/id` (analysis in `docs/planning/archived/`; breaking change, deferred since 2026-03 — revisit at the next major version)

## Ecosystem & Developer Experience

- Tighter integration examples with `sivchari/govalid` (structural + business validation pairing)
- Generate rule documentation from `Description()`/`Tags()` metadata
- CLI tool for running validations
- Interactive playground / more real-world examples
- Sitemap of the nested modules in README with usage snippets (`adapters/cqrslite`, `listeners/otel`, `examples/sse`)

## Performance

- Caching for expensive regex compilations (hot-path builders take precompiled `*regexp.Regexp` today; measure before building)
- Publish a benchmark results document — `BenchmarkStream*` and the 189/468/476 ns listener numbers exist in README/AGENTS but no consolidated results doc

## Versioning & Stability

- API stability review once Polish-Customs (live consumer since 2026-09-14) exercises more surface
- Re-evaluate the `encoding/json/v2` / `GOEXPERIMENT=jsonv2` downstream constraint when the package graduates from experimental (ADR-0004)
- Evaluate a `json/v1` build-tag fallback so consumers can opt out of the `GOEXPERIMENT=jsonv2` requirement

## Public presence & distribution (new 2026-09-18)

- Social preview image for the GitHub repo
- Seed awareness for pkg.go.dev "Imported by: 0" (blog post / social) once a real external consumer exists
- Publish example dashboards/screenshots from the `examples/sse` Datastar feed
- Review the now-public internal tooling configs (`git-town.toml`, `.buildflow.yml`, `.config/metadata.yaml`, `library-policy.yaml`) for anything unintended
- Decide whether the archived `docs/planning/` FINDING-SDK proposal still reflects intent

## CI & tooling polish (new 2026-09-18)

- Cache the `go install gosec@v2.29.0` step; pin an upgrade cadence for it
- Extend gosec to the nested modules (currently root-only)
- Optional short-duration fuzz job over the 7 `Fuzz*` targets
- Consistent Dependabot grouping across all four modules
- Auto-create the GitHub Release on tag push (remove the manual `gh release create` step)
- Decide on the ubuntu-26 runner migration (Oct 19) — accept or pin
- dprint/gofmt check in CI for parity with `buildflow` (currently local-only)

## BuildFlow & tooling ecosystem (upstream / raw, new 2026-10-04)

Harvested from the toolchain-skew remediation ([`docs/status/2026-10-04_07-27_buildflow-toolchain-skew-remediation.md`](docs/status/2026-10-04_07-27_buildflow-toolchain-skew-remediation.md) §e/§f) — mostly upstream asks in OTHER repos (BuildFlow, treefmt-nix, SystemNix), so they live here, not in TODO_LIST:

- BuildFlow: uniform per-mode tool→environment resolution (govalid-generate was green in fast mode, red in full mode from the same tree); make the chosen env visible per step in `-v` output (§e1/§f5)
- BuildFlow `go-version-auto-configure`: tolerate patch floors that `go mod tidy` provably force from a dependency (§e2/§f11)
- BuildFlow `go-auto-upgrade`: first-class project-level finding suppressions so documented false positives stop re-flagging (§e3/§f12; gated on the samber/lo policy answer below)
- treefmt-nix/nixpkgs: expose the Go version formatter wrappers were built with, or keep wrapper-go ≥ default-go (the hermetic sandbox cannot download toolchains) (§e4)
- SystemNix: generate pre-deploy-check lib staging from the deploy script's `source` lines — same manual-staging bug shipped twice (§e5/§f13)
- crush-config `references/lessons.md`: pipefail-after-every-piped-gate discipline + "verify conclusions in the gating mode" (§f19, §e7)
- Flake app wrapping PSI-polling + deploy for "quiet-window deploys" (§f38); consider `nh sw` single-package path for toolchain-only updates (§f33)
- Pre-commit size check for the AGENTS.md 377-line budget (§f27); per-module go floors as a single flake-eval source of truth (§f41)
- Policy: pin nixpkgs input vs track `nixos-unstable` (§f35); `nix flake check --all-systems` feasibility (aarch64 builders) (§f21); check-all inside CI probably duplicates the matrix — confirm skip (§f45)
- Watch govalid upstream for releases > 1.9.0 and bump the mirrored rev in BOTH flakes (§f30); review `go-tool-run`/`ginkgo-version-check` step purposes (§f46)

## Open questions (user decisions needed)

1. ~~CI on a private repo~~ — RESOLVED 2026-09-17: the repo went public, Actions are free, and CI is green 13/13 (run `35310049647`). Account billing still affects other, private repos in the ecosystem.
2. ~~Repo visibility~~ — RESOLVED 2026-09-17: made PUBLIC; pkg.go.dev now indexes `v2.2.0` and the public proxy serves the module.
3. **Adapter home** — keep `adapters/cqrslite`, `examples/sse`, and `listeners/otel` as nested modules here, or promote to sibling repos (collector-extraction pattern) once a second external consumer appears?
4. **Should the nested modules ever be published?** They carry local `replace` directives today; publishing them would require removing those and giving each a release train — or they stay internal-only forever.
5. **Concurrency policy** — the auto-commit daemon plus parallel sessions have interleaved heuristic commits mid-release (e.g. `2e06602`). Should the daemon be suspended for release-critical work?
6. ~~Branch protection on `master`~~ — DECLINED 2026-09-18: kept unprotected (the auto-commit daemon pushes directly to `master`; a required-PR/status-check rule would break it).
7. ~~GitHub Discussions~~ — DECLINED 2026-09-18: issues-only.
8. **Deploy authority & timing (2026-10-04)** — the full-mode `govalid-generate` step stays red until the system Go is 1.27.1 (SystemNix deploy). Fire it myself at the next quiet IO window, or do you want to pick the moment? It switches the running system and touches all concurrent sessions.
9. **Pressure-gate override policy (2026-10-04)** — is `DEPLOY_FORCE_PRESSURE=1` ever acceptable (documented kernel-freeze precursor class), or is that gate always human-only? I have not forced it; the repo stays 2-steps-red until the machine quiets down.
10. **samber/lo policy (2026-10-04)** — "no samber/lo anywhere" is the documented AGENTS.md policy today. Make it permanent fleet-wide (→ push the BuildFlow suppression ask above), or allow `lo` in test files / nested modules (→ convert the 8 flagged loops)?

---

_Last Updated: 2026-10-04_
