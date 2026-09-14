# Roadmap — businessrules

> Long-term direction and raw ideas **not yet refined into actionable tasks**.
>
> When an item here becomes bounded and short-term, it graduates into [`TODO_LIST.md`](TODO_LIST.md). When an item ships, it is recorded in [`CHANGELOG.md`](CHANGELOG.md) and reflected in [`FEATURES.md`](FEATURES.md).

**Last reviewed:** 2026-09-14 (post-v2.1.0 execution session: builder batches, composition, metadata, property-based equivalence, hygiene, and ADRs shipped — pruned from this file).

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

## Performance

- Caching for expensive regex compilations (hot-path builders take precompiled `*regexp.Regexp` today; measure before building)
- Publish benchmark numbers (`BenchmarkStream*` exist; add `BenchmarkBuild*` sweeps)

## Versioning & Stability

- API stability review once Polish-Customs (live consumer since 2026-09-14) exercises more surface
- Re-evaluate the `encoding/json/v2` / `GOEXPERIMENT=jsonv2` downstream constraint when the package graduates from experimental (ADR-0004)
- Cut `v2.2.0` from `[Unreleased]` when the 2026-09-14 builder batch needs a tag

## Open questions (user decisions needed)

1. ~~CI on a private repo~~ — RESOLVED 2026-09-14: workflow re-enabled; runs remain blocked only by GitHub Actions billing (user action, tracked in TODO_LIST).
2. **Adapter home** — keep `adapters/cqrslite`, `examples/sse`, and `listeners/otel` as nested modules here, or promote to sibling repos (collector-extraction pattern) once a second external consumer appears?
3. **Repo visibility** — staying private keeps pkg.go.dev indexing off; making it public changes consumption docs but not code.

---

_Last Updated: 2026-09-14_
