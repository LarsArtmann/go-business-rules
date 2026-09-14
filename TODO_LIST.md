# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-09-14 16:10 (context rules, `WithConcurrency`, Stream benchmarks, and goroutine-leak tests shipped by the concurrent 16:0x session; 164/169 specs pass, 5 = branching-flow count pins)

---

## CI & release pipeline

- [ ] **CI workflow is disabled on GitHub** (discovered 2026-09-14): `.github/workflows/ci.yml` is `disabled_manually` since 2026-07-17 — no runs in 2 months, and every run since 2026-06 failed within 3-5s (setup-level failure). Decide: re-enable and fix the setup failure, or keep disabled (private-repo Actions minutes) and rely on the local flake gates. Until decided, the README CI badge and any "CI green" assumption are meaningless.
- [ ] **Fix root module version tags** (blocking, discovered 2026-09-14): git tag `v2.0.0` is not consumable by Go module resolution (v2+ tag requires a `/v2` module-path suffix; this module has none). The repo is also **private**, so the proxy/pkg.go.dev cannot serve anything today. Decide: rename module to `github.com/LarsArtmann/go-business-rules/v2` + re-tag, or tag a compatible `v0.x`/`v1.x` carrying current code. Then tag a release carrying events/streaming. Blocks publishing `adapters/cqrslite` and any versioned consumer require.
- [ ] **Decide + apply policy for the 5 stale branching-flow specs** (`bdd_branching_flow_test.go:102,109,181` and the `all`-output spec): red since 2026-09-14, count drifted 2→5 when the 16:0x session's new code (context rules, concurrency) added PHANTOM hits — the analyzer scans the whole directory and has no path-exclude flag. Options: re-pin to current reality (green suite, weaker detector), keep red as explicit signal, or move to a nightly job.
- [ ] **Integrate into Polish-Customs** as a real-world consumer (replace internal `validation.go`); verify compatibility. Not a blocker for tagging, but the first real-world validation of the breaking-change surface.
- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers (documented in `AGENTS.md`). Track the Go release that graduates `json/v2` from experimental and remove the constraint then. Tracked as a long-term item; cannot be resolved until Go ships it.

## Events & streaming follow-ups (from 2026-09-14 event-driven work)

- [ ] **Wire nested modules into CI** — `.github/workflows/ci.yml` tests only the root module; `adapters/cqrslite` and `examples/sse` need matrix jobs (each is a separate module; private deps need the flake devshell).
- [ ] **Datastar reactive-UI example** — extend `examples/sse` with go-datastar signals/patches (research done 2026-09-14: patches are values → `.Event()` → same `sse.Broadcaster`; the wiring is proven in go-datastar's `example/main.go`).
- [ ] **OpenTelemetry listener** — a small `listeners/otel` (or nested module) mapping `RuleEvaluated`/`ValidationCompleted` to spans/metrics.
- [ ] **Composite lint/test command for all 3 modules** — root `golangci-lint run` and the root test command do not cover `adapters/cqrslite` / `examples/sse`; add a flake app that runs build+test+lint across all three.
- [ ] **Document the new context/concurrency APIs** — `CheckContext`, `ContextRuleImpl`, `WithConcurrency` shipped 2026-09-14 16:0x but are absent from README's API section, FEATURES.md, and `docs/DOMAIN_LANGUAGE.md`.

---

_The release-prep items from the 2026-07-26 docs-health audit — `Version` constant split-brain, branching-flow stats test hardening, and `docs/DOMAIN_LANGUAGE.md` — were resolved for the `v2.0.0` release and are recorded in [`CHANGELOG.md`](CHANGELOG.md). The `json/v2` downstream-constraint item remains open by design (depends on a future Go release)._
