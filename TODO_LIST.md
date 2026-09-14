# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-09-14 (events & streaming shipped; `adapters/cqrslite` + `examples/sse` nested modules live)

---

## Integration & release

- [ ] **Fix root module version tags** (blocking, discovered 2026-09-14): git tag `v2.0.0` is not consumable by Go module resolution (v2+ tag requires a `/v2` module-path suffix; this module has none). Downstream `go get` resolves only `v0.1.0` (May 2026). Decide: rename module to `github.com/LarsArtmann/go-business-rules/v2` + re-tag, or tag a compatible `v0.x`/`v1.x` carrying current code. Blocks publishing `adapters/cqrslite` and any versioned consumer require.
- [ ] **Integrate into Polish-Customs** as a real-world consumer (replace internal `validation.go`); verify compatibility. Not a blocker for tagging `v2.0.0`, but the first real-world validation of the breaking-change surface.
- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers (documented in `AGENTS.md`). Track the Go release that graduates `json/v2` from experimental and remove the constraint then. Tracked as a long-term item; cannot be resolved until Go ships it.

## Events & streaming follow-ups (from 2026-09-14 event-driven work)

- [ ] **Wire nested modules into CI** — `.github/workflows/ci.yml` tests only the root module; `adapters/cqrslite` and `examples/sse` need matrix jobs (each is a separate module; private deps need the flake devshell).
- [ ] **`WithConcurrency(n)` option for `Stream`** — currently one goroutine per rule; add an errgroup-style bound for callers with many slow rules (rule count is caller-controlled today, so this is a safety valve, not a correctness fix).
- [ ] **Context-carrying rules** — `Rule.Check()` cannot be interrupted; a `Check(ctx)`-style rule interface would let `Stream(ctx)` cancel in-flight slow checks instead of only skipping unstarted ones. Requires an additive `Rule2`-style interface and builder support.
- [ ] **Datastar reactive-UI example** — extend `examples/sse` with go-datastar signals/patches (research done 2026-09-14: patches are values → `.Event()` → same `sse.Broadcaster`; the wiring is proven in go-datastar's `example/main.go`).
- [ ] **OpenTelemetry listener** — a small `listeners/otel` (or nested module) mapping `RuleEvaluated`/`ValidationCompleted` to spans/metrics.

---

_The release-prep items from the 2026-07-26 docs-health audit — `Version` constant split-brain, branching-flow stats test hardening, and `docs/DOMAIN_LANGUAGE.md` — were resolved for the `v2.0.0` release and are recorded in [`CHANGELOG.md`](CHANGELOG.md). The `json/v2` downstream-constraint item remains open by design (depends on a future Go release)._
