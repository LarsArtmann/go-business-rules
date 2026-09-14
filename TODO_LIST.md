# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-09-14 (events & streaming shipped; `adapters/cqrslite` + `examples/sse` nested modules live; docs-health audit pass)

---

## CI & release pipeline

- [ ] **CI workflow is disabled on GitHub** (discovered 2026-09-14): `.github/workflows/ci.yml` is `disabled_manually` since 2026-07-17 — no runs in 2 months, and every run since 2026-06 failed within 3-5s (setup-level failure). Decide: re-enable and fix the setup failure, or keep disabled (private-repo Actions minutes) and rely on the local flake gates. Until decided, the README CI badge and any "CI green" assumption are meaningless.
- [ ] **Fix root module version tags** (blocking, discovered 2026-09-14): git tag `v2.0.0` is not consumable by Go module resolution (v2+ tag requires a `/v2` module-path suffix; this module has none). The repo is also **private**, so the proxy/pkg.go.dev cannot serve anything today. Decide: rename module to `github.com/LarsArtmann/go-business-rules/v2` + re-tag, or tag a compatible `v0.x`/`v1.x` carrying current code. Then tag a release carrying events/streaming. Blocks publishing `adapters/cqrslite` and any versioned consumer require.
- [ ] **Decide + apply policy for the 2 stale branching-flow specs** (`bdd_branching_flow_test.go:181` and the `stats` total pin): red since before 2026-09-14 due to analyzer binary drift (stats expects 36, actual 15). Options: re-pin to current reality (CI green, weaker detector), keep red as explicit signal, or move to a nightly job.
- [ ] **Integrate into Polish-Customs** as a real-world consumer (replace internal `validation.go`); verify compatibility. Not a blocker for tagging, but the first real-world validation of the breaking-change surface.
- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers (documented in `AGENTS.md`). Track the Go release that graduates `json/v2` from experimental and remove the constraint then. Tracked as a long-term item; cannot be resolved until Go ships it.

## Events & streaming follow-ups (from 2026-09-14 event-driven work)

- [ ] **Wire nested modules into CI** — `.github/workflows/ci.yml` tests only the root module; `adapters/cqrslite` and `examples/sse` need matrix jobs (each is a separate module; private deps need the flake devshell).
- [ ] **`WithConcurrency(n)` option for `Stream`** — currently one goroutine per rule; add an errgroup-style bound for callers with many slow rules (rule count is caller-controlled today, so this is a safety valve, not a correctness fix).
- [ ] **Context-carrying rules** — `Rule.Check()` cannot be interrupted; a `Check(ctx)`-style rule interface would let `Stream(ctx)` cancel in-flight slow checks instead of only skipping unstarted ones. Requires an additive `Rule2`-style interface and builder support.
- [ ] **Datastar reactive-UI example** — extend `examples/sse` with go-datastar signals/patches (research done 2026-09-14: patches are values → `.Event()` → same `sse.Broadcaster`; the wiring is proven in go-datastar's `example/main.go`).
- [ ] **OpenTelemetry listener** — a small `listeners/otel` (or nested module) mapping `RuleEvaluated`/`ValidationCompleted` to spans/metrics.
- [ ] **Stream benchmarks** (`BenchmarkStream*`) — only the `Build` path is measured today (189/468/476 ns/op for 0/1/3 listeners).
- [ ] **Goroutine-leak regression test for `Stream`** — `runtime.NumGoroutine`-based, no new dependency; guards the buffered-channel contract.
- [ ] **Composite lint/test command for all 3 modules** — root `golangci-lint run` and the root test command do not cover `adapters/cqrslite` / `examples/sse`; add a flake app that runs build+test+lint across all three.

---

_The release-prep items from the 2026-07-26 docs-health audit — `Version` constant split-brain, branching-flow stats test hardening, and `docs/DOMAIN_LANGUAGE.md` — were resolved for the `v2.0.0` release and are recorded in [`CHANGELOG.md`](CHANGELOG.md). The `json/v2` downstream-constraint item remains open by design (depends on a future Go release)._
