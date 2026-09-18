# Roadmap — businessrules

> Long-term direction and raw ideas **not yet refined into actionable tasks**.
>
> When an item here becomes bounded and short-term, it graduates into [`TODO_LIST.md`](TODO_LIST.md). When an item ships, it is recorded in [`CHANGELOG.md`](CHANGELOG.md) and reflected in [`FEATURES.md`](FEATURES.md).

**Last reviewed:** 2026-09-18 (harvest of the v2.2.0 release + repo-public-launch sessions: `v2.2.0` cut and published, repo flipped public on 2026-09-17, pkg.go.dev indexed, CI green 13/13 — shipped items pruned, public-presence and CI-polish candidates added).

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

## Open questions (user decisions needed)

1. ~~CI on a private repo~~ — RESOLVED 2026-09-17: the repo went public, Actions are free, and CI is green 13/13 (run `35310049647`). Account billing still affects other, private repos in the ecosystem.
2. ~~Repo visibility~~ — RESOLVED 2026-09-17: made PUBLIC; pkg.go.dev now indexes `v2.2.0` and the public proxy serves the module.
3. **Adapter home** — keep `adapters/cqrslite`, `examples/sse`, and `listeners/otel` as nested modules here, or promote to sibling repos (collector-extraction pattern) once a second external consumer appears?
4. **Should the nested modules ever be published?** They carry local `replace` directives today; publishing them would require removing those and giving each a release train — or they stay internal-only forever.
5. **Concurrency policy** — the auto-commit daemon plus parallel sessions have interleaved heuristic commits mid-release (e.g. `2e06602`). Should the daemon be suspended for release-critical work?

---

_Last Updated: 2026-09-18_
