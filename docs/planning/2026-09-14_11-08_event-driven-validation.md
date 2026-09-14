# Plan: Event-Driven Validation for businessrules

> **Created:** 2026-09-14 11:08 CEST
> **Goal:** Make `go-business-rules` a _fact producer_ that plugs into event-driven systems, without breaking its pure, synchronous core.
> **Constraint:** ZERO breaking changes. ZERO new dependencies in the root `go.mod`. The default `Build()` path must behave byte-identically to today when no listeners are registered.

---

## 1. Context & Research

### Why

`ValidatorBuilder.Build()` (validator.go:34) discards everything except failures: no pass events, no timing, no provenance. An event-driven surface captures all of it and lets consumers wire metrics, audit trails, alerting, and live UIs **without touching the core loop**.

### Research findings (verified 2026-09-14)

| Repo           | Finding                                                                                                                                                                                                                                                     | Relevance                            |
| -------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------ |
| `go-cqrs-lite` | `event/v4`: `event.New(type, streamID, streamType, version, payload, opts...)` auto-marshals payload (CBOR). Bus: `watermill.NewEventBus()` implements `event.Bus` (`Publish(ctx, events...)`, `Subscribe`). Tests: `eventtest.NewFakeBus()` (synchronous). | Adapter target for level 3           |
| `go-sse`       | `Broadcaster[T]` (non-blocking broadcast, subscribe/unsubscribe channels), `Stream` (SSE wire), `Replay`. Requires `GOEXPERIMENT=jsonv2`.                                                                                                                   | Streaming demo (level 2 consumption) |
| `go-datastar`  | Patches are values → `.Event() → sse.Event` → broadcast. Hard dep on `go-sse v0.6.0`. Proven pattern in its `example/main.go`.                                                                                                                              | Follow-up, not in scope today        |

### Design decisions (the "fight for every name" section)

| Decision              | Choice                                                                                                                                                                                                               | Why                                                                                                                                                                          |
| --------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Event model           | Sealed `Event` interface + two concrete structs (`RuleEvaluated`, `ValidationCompleted`)                                                                                                                             | Type-safe, switch-able, impossible to invent untyped garbage. No JSON marshaling of events (YAGNI — events are for programmatic listeners; result already has `MarshalJSON`) |
| Listener type         | `type Listener func(Event)` — plain func                                                                                                                                                                             | Composable (wrap for fan-out), zero interface ceremony. Contract: synchronous, must not panic (documented)                                                                   |
| Default path          | `len(listeners) == 0` → no `time.Now()` calls, identical loop                                                                                                                                                        | Never make the hot path worse. Benchmark proves it                                                                                                                           |
| Streaming             | `Stream(ctx) <-chan Event`: goroutine per rule, events in completion order, `ValidationCompleted` last (violations re-sorted to rule order for determinism), channel closed after; ctx cancellation = partial result | Streaming is the point; deterministic final result is the quality bar                                                                                                        |
| CQRS adapter location | Nested Go module `adapters/cqrslite/` (own `go.mod`)                                                                                                                                                                 | Root `go.mod` stays dependency-free; adapter is opt-in for consumers                                                                                                         |
| Wire payload          | Adapter-owned DTOs (`RuleEvaluatedData`, `ValidationCompletedData`) with plain fields (string/int64/bool)                                                                                                            | Stable wire contract decoupled from library struct evolution; CBOR-safe                                                                                                      |
| SSE demo              | Nested module `examples/sse/` (own `go.mod`)                                                                                                                                                                         | Shows the full story: validation → events → broadcaster → browser                                                                                                            |

### What we deliberately do NOT do

- No goroutines/channels in the default `Build()` path.
- No `go-cqrs-lite` / `go-sse` imports in the root module.
- No command/dispatch machinery — this library is a pure function that _emits_ facts, not a framework that routes them.
- No `recover()` around listeners — listeners are trusted code; document the contract.

---

## 2. Pareto Breakdown

### The 1% → 51% of value

**Core event emission.** `events.go` with `Event` / `RuleEvaluated` / `ValidationCompleted` / `Listener`, `WithListener(...)` on the builder, synchronous emission inside `Build()`, zero-cost when unused. Everything else is built on this.

### The 4% → 64% of value (1% + this)

**Streaming evaluation + proof + docs.** `Stream(ctx)` for concurrent slow rules; benchmarks proving listener overhead is noise; `doc.go`, README, FEATURES.md, CHANGELOG updates so the feature is discoverable and honestly documented.

### The 20% → 80% of value (4% + this)

**Ecosystem bridges.** `adapters/cqrslite` (events → durable CQRS domain events → projections/dashboards) and `examples/sse` (events → live browser feed). These make the feature real for Lars's ecosystem (go-cqrs-lite, go-sse, cqrs-htmx).

### The remaining 80% → 100%

| Item                                                                                         | Where tracked                 |
| -------------------------------------------------------------------------------------------- | ----------------------------- |
| Datastar reactive-UI example (signals/patches from validation events)                        | TODO_LIST → ROADMAP candidate |
| `WithConcurrency(n)` bound for `Stream` (errgroup.SetLimit)                                  | TODO_LIST                     |
| Context-carrying rules (`Rule2` with `Check(ctx)`) so cancellation interrupts running checks | TODO_LIST                     |
| CI matrix wiring for nested modules (adapters + examples)                                    | TODO_LIST                     |
| OpenTelemetry listener implementation                                                        | ROADMAP                       |
| Event JSON marshaling (only if a consumer asks)                                              | YAGNI today                   |
| Polish-Customs integration as first event-driven consumer                                    | existing TODO_LIST item       |

---

## 3. Comprehensive Plan (30–100 min tasks, sorted by impact)

| #  | Task                                                                                                    | Effort | Impact | Value            |
| -- | ------------------------------------------------------------------------------------------------------- | ------ | ------ | ---------------- |
| ~~1~~  | ~~Core events: `events.go` types + `WithListener` + `Build()` emission (zero-cost default path)~~ done at `88ed1c3` | ~~60 min~~ | ~~HIGH~~ | ~~1% tier~~ |
| ~~2~~  | ~~BDD specs for events (`events_test.go`): per-rule events, terminal event, multi-listener, sync contract~~ done at `88ed1c3` | ~~45 min~~ | ~~HIGH~~ | ~~Correctness gate~~ |
| ~~3~~  | ~~`Stream(ctx)`: concurrent evaluation, completion-order events, deterministic final result, cancellation~~ done — stream_test.go 6 specs green (daemon commit) | ~~60 min~~ | ~~HIGH~~ | ~~4% tier~~ |
| ~~4~~  | ~~BDD specs for `Stream` (`stream_test.go`): equivalence with `Build`, cancellation, channel close~~ done — stream specs green, stable across 4 runs | ~~45 min~~ | ~~HIGH~~ | ~~Correctness gate~~ |
| ~~5~~  | ~~Benchmarks: no-listener vs 1 listener vs 3 listeners~~ done at `4221e69` | ~~30 min~~ | ~~MED~~ | ~~Performance gate~~ |
| ~~6~~  | ~~Docs: `doc.go` events section, README "Events" section, FEATURES.md, CHANGELOG `[Unreleased]`~~ done at `85edd99` | ~~45 min~~ | ~~MED~~ | ~~Discoverability~~ |
| ~~7~~  | ~~`adapters/cqrslite` nested module: DTOs + `NewBusListener` + specs with `eventtest.NewFakeBus`~~ done at `3f5f8f8` | ~~90 min~~ | ~~MED~~ | ~~20% tier~~ |
| ~~8~~  | ~~`examples/sse` nested module: broadcaster + `/events` handler + index.html + build check~~ done at `b57a016` | ~~60 min~~ | ~~MED~~ | ~~20% tier~~ |
| ~~9~~  | ~~AGENTS.md + TODO_LIST.md updates (events section, harvested follow-ups)~~ done at `624a95e` | ~~30 min~~ | ~~MED~~ | ~~Memory~~ |
| ~~10~~ | ~~Full verification + detailed commits + push~~ done at `0453a78` | ~~40 min~~ | ~~HIGH~~ | ~~Release gate~~ |

## 4. Detailed Breakdown (≤ 12 min tasks, sorted by impact)

| #    | Task                                                                                                      | Time   | Depends on |
| ---- | --------------------------------------------------------------------------------------------------------- | ------ | ---------- |
| ~~1.1~~ | ~~Write `events.go`: `Event` marker interface, `RuleEvaluated`, `ValidationCompleted`, `Listener`~~ done (see section 3 / execution log) | ~~10 min~~ | ~~—~~ |
| ~~1.2~~ | ~~Add `listeners` field + `WithListener(...)` to `ValidatorBuilder`~~ done (see section 3 / execution log) | ~~5 min~~ | ~~1.1~~ |
| ~~1.3~~ | ~~Rework `Build()`: single loop, timed emission only when listeners exist~~ done (see section 3 / execution log) | ~~10 min~~ | ~~1.2~~ |
| ~~2.1~~ | ~~Spec: emits `RuleEvaluated` per rule (pass + fail) in rule order~~ done (see section 3 / execution log) | ~~10 min~~ | ~~1.3~~ |
| ~~2.2~~ | ~~Spec: failure event carries `Err`, severity, `Passed=false`; pass event `Err=nil`~~ done (see section 3 / execution log) | ~~8 min~~ | ~~2.1~~ |
| ~~2.3~~ | ~~Spec: `ValidationCompleted` last, carries aggregate result + duration; builder returns identical result~~ done (see section 3 / execution log) | ~~8 min~~ | ~~2.2~~ |
| ~~2.4~~ | ~~Spec: multiple listeners all receive full sequence; emission is synchronous (done before `Build` returns)~~ done (see section 3 / execution log) | ~~8 min~~ | ~~2.3~~ |
| ~~2.5~~ | ~~Run suite (`nix develop --command go test ./...`) — must stay green (145+ specs)~~ done (see section 3 / execution log) | ~~5 min~~ | ~~2.4~~ |
| ~~3.1~~ | ~~Implement `Stream(ctx)`: scheduler + collector, completion-order events~~ done (see section 3 / execution log) | ~~12 min~~ | ~~1.3~~ |
| ~~3.2~~ | ~~Violation re-sorting to rule order + terminal `ValidationCompleted` + channel close~~ done (see section 3 / execution log) | ~~10 min~~ | ~~3.1~~ |
| ~~3.3~~ | ~~Context cancellation: stop scheduling, partial result~~ done (see section 3 / execution log) | ~~8 min~~ | ~~3.2~~ |
| ~~4.1~~ | ~~Spec: stream yields N rule events + completed, channel closes~~ done (see section 3 / execution log) | ~~10 min~~ | ~~3.3~~ |
| ~~4.2~~ | ~~Spec: violations in completed event == `Build()` result (same set)~~ done (see section 3 / execution log) | ~~8 min~~ | ~~4.1~~ |
| ~~4.3~~ | ~~Spec: canceled ctx → partial completed, no goroutine leak (events stop)~~ done (see section 3 / execution log) | ~~10 min~~ | ~~4.2~~ |
| ~~5.1~~ | ~~`BenchmarkBuild_NoListener` / `_OneListener` / `_ThreeListeners`~~ done (see section 3 / execution log) | ~~10 min~~ | ~~2.5~~ |
| ~~5.2~~ | ~~Run benchmarks, record ns/op delta in plan doc follow-up note~~ done (see section 3 / execution log) | ~~5 min~~ | ~~5.1~~ |
| ~~6.1~~ | ~~`doc.go`: Events section + thread-safety note for listeners~~ done (see section 3 / execution log) | ~~10 min~~ | ~~2.5~~ |
| ~~6.2~~ | ~~README: "## Events" section with listener + stream examples~~ done (see section 3 / execution log) | ~~12 min~~ | ~~2.5~~ |
| ~~6.3~~ | ~~FEATURES.md: events + streaming rows with evidence; verified date bump~~ done (see section 3 / execution log) | ~~8 min~~ | ~~6.1~~ |
| ~~6.4~~ | ~~CHANGELOG.md: `[Unreleased]` Added entries~~ done (see section 3 / execution log) | ~~5 min~~ | ~~6.3~~ |
| ~~7.1~~ | ~~`adapters/cqrslite/go.mod` + `doc.go` (module doc, wire contract)~~ done (see section 3 / execution log) | ~~10 min~~ | ~~2.5~~ |
| ~~7.2~~ | ~~DTOs: `RuleEvaluatedData`, `ValidationCompletedData` + converters~~ done (see section 3 / execution log) | ~~10 min~~ | ~~7.1~~ |
| ~~7.3~~ | ~~`NewBusListener(bus, streamID, streamType, ...opts)`: publish per event, version counter~~ done (see section 3 / execution log) | ~~12 min~~ | ~~7.2~~ |
| ~~7.4~~ | ~~Specs with `eventtest.NewFakeBus`: subscribe → validate → decode payloads, assert round-trip~~ done (see section 3 / execution log) | ~~12 min~~ | ~~7.3~~ |
| ~~7.5~~ | ~~`go test ./...` inside adapter dir (nix develop, GOPRIVATE)~~ done (see section 3 / execution log) | ~~8 min~~ | ~~7.4~~ |
| ~~8.1~~ | ~~`examples/sse/go.mod` + `main.go`: rules → listener → `sse.Broadcaster`~~ done (see section 3 / execution log) | ~~12 min~~ | ~~2.5~~ |
| ~~8.2~~ | ~~`/events` handler (`sse.NewStream` loop) + minimal `index.html`~~ done (see section 3 / execution log) | ~~10 min~~ | ~~8.1~~ |
| ~~8.3~~ | ~~`go build ./...` inside example dir; manual curl smoke test~~ done (see section 3 / execution log) | ~~10 min~~ | ~~8.2~~ |
| ~~9.1~~ | ~~AGENTS.md: events architecture section, nested-module commands, listener contract~~ done (see section 3 / execution log) | ~~10 min~~ | ~~6.4~~ |
| ~~9.2~~ | ~~TODO_LIST.md: add harvested follow-ups (concurrency bound, ctx rules, CI matrix, datastar)~~ done (see section 3 / execution log) | ~~5 min~~ | ~~9.1~~ |
| ~~10.1~~ | ~~Root: full test suite + `go vet` + golangci-lint~~ done (see section 3 / execution log) | ~~10 min~~ | ~~9.2~~ |
| ~~10.2~~ | ~~Adapters + examples: build + test verification~~ done (see section 3 / execution log) | ~~8 min~~ | ~~10.1~~ |
| ~~10.3~~ | ~~Detailed per-task commits (check `git status` before each — auto-commit daemon races)~~ done (see section 3 / execution log) | ~~12 min~~ | ~~10.2~~ |
| ~~10.4~~ | ~~`git push` + report back with tables~~ done (see section 3 / execution log) | ~~8 min~~ | ~~10.3~~ |

---

## 5. Execution Graph

```mermaid
graph TD
    P[Plan doc + commit] --> E[events.go types + Listener]
    E --> W[WithListener + Build emission]
    W --> T1{Specs: events green?}
    T1 -- no --> F1[Fix, re-run]
    F1 --> T1
    T1 -- yes --> S[Stream ctx implementation]
    S --> T2{Specs: stream green?}
    T2 -- no --> F2[Fix, re-run]
    F2 --> T2
    T2 -- yes --> B[Benchmarks + record]
    B --> D1[doc.go + README]
    D1 --> D2[FEATURES + CHANGELOG]
    D2 --> A1[adapter go.mod + DTOs]
    A1 --> A2[NewBusListener]
    A2 --> A3{FakeBus specs green?}
    A3 -- no --> F3[Fix, re-run]
    F3 --> A3
    A3 -- yes --> X1[example go.mod + main]
    X1 --> X2{Build + smoke ok?}
    X2 -- no --> X3[Simplify to doc-only wiring]
    X3 --> X4[Document fallback]
    X2 -- yes --> M1[AGENTS + TODO_LIST]
    X4 --> M1
    M1 --> V[Full verify: root + nested modules + lint]
    V --> C[Detailed commits]
    C --> PU[git push]
```

---

## 6. Verification Gates

| Gate          | Command                                                       | Pass criteria                                                            |
| ------------- | ------------------------------------------------------------- | ------------------------------------------------------------------------ |
| Root tests    | `nix develop --command go test ./...`                         | All specs green, 0 failures                                              |
| Coverage      | `nix develop --command go test -cover ./...`                  | No drop below current 94.8% (new code included)                          |
| Benchmarks    | `nix develop --command go test -bench . -benchmem`            | Listener delta ≈ ns-level                                                |
| Adapter tests | `cd adapters/cqrslite && nix develop --command go test ./...` | Round-trip via FakeBus green                                             |
| Example build | `cd examples/sse && nix develop --command go build ./...`     | Compiles                                                                 |
| Lint          | `nix develop --command golangci-lint run`                     | 0 new issues                                                             |
| API safety    | `git diff` review                                             | Root `go.mod` untouched; no signature changes to existing exported funcs |

## 7. Rollback Safety

Every tier is an independent, additive commit. If the adapter or example hits external limits (private proxy, GOEXPERIMENT), the plan's fallback (documented wiring, no code) preserves the 1%/4% value. Nothing here modifies existing behavior.

---

## 8. Execution Log (updated as work completes)

| Date       | Result                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| ---------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 2026-09-14 | Plan committed (`099104c`). Baseline verified via worktree at `a2d5938`: suite already red with 2 stale branching-flow count pins (stats expects 36, actual 15; `all` output drift) — pre-existing, not caused by this work.                                                                                                                                                                                                                               |
| 2026-09-14 | Core events shipped (`88ed1c3`). Design hardening: `Passed` field replaced by derived `Passed() = (Err == nil)` method after the linter flagged bool-blindness; removes the representable impossible state (claimed pass + carried error). PHANTOM count stays at documented 12. Suite: 156 pass / 2 pre-existing fails.                                                                                                                                   |
| 2026-09-14 | `Stream(ctx)` shipped (auto-committed by daemon). Suite stable across 4 runs: 156 pass / 2 pre-existing fails.                                                                                                                                                                                                                                                                                                                                             |
| 2026-09-14 | Benchmarks (2 rules, AMD RYZEN AI MAX+ 395): no listener **189 ns/op, 456 B/op, 9 allocs**; one listener **468 ns/op, 688 B/op, 13 allocs**; three listeners **476 ns/op, 704 B/op, 13 allocs**. First listener costs ~280 ns (event structs + timing), additional listeners ~7 ns each. No-listener path unchanged.                                                                                                                                       |
| 2026-09-14 | `adapters/cqrslite` shipped: 4 FakeBus specs green (payload round-trips, stream tagging/versioning, terminal summaries, error routing). Publishing blocked by the root module's unconsumable v2.0.0 tag — documented in AGENTS.md + TODO_LIST.md.                                                                                                                                                                                                          |
| 2026-09-14 | `examples/sse` shipped with a real-HTTP smoke test (SSE subscribe → POST /validate → both event kinds arrive). Discovered SSE headers only flush on first event, so the test triggers validation concurrently.                                                                                                                                                                                                                                             |
| 2026-09-14 | Final state: lint 0 issues, format check green, all three modules' tests green. Suite: 156 pass / 2 fail — both remaining failures are pre-existing branching-flow binary drift (stats pin 36 vs actual, `all` output), red at baseline `a2d5938`. The 12→14 PHANTOM and 1→3 error-severity pins were re-pinned with documented reasoning: the nested example module intentionally validates raw primitives (no path-exclude flag exists in the analyzer). |
