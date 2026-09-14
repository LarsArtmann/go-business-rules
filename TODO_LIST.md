# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-09-14 (execution session: `v2.1.0` pushed and verified consumable from the real remote, Polish-Customs `replace` dropped and green against the published version, CI re-enabled — billing fix remains the only blocker)

---

## Release follow-ups (user actions only)

- [x] **Push `master` and the `v2.1.0` tag** — DONE 2026-09-14: both pushed; a fresh scratch module resolved `github.com/LarsArtmann/go-business-rules/v2@v2.1.0` from the real remote, compiled, and ran (`2.1.0 true`).
- [x] **Drop the temporary `replace` in Polish-Customs** — DONE 2026-09-14: `replace` line removed, `go mod tidy` resolves the published `v2.1.0`, full Polish-Customs suite green.
- [ ] **Fix GitHub Actions billing** (only remaining blocker): every job is still rejected at start — _"recent account payments have failed or your spending limit needs to be increased"_ (re-confirmed 2026-09-14 on the first post-enable run: all 13 matrix jobs, zero steps executed). The workflow is already **enabled** and its matrix is correct; once billing is fixed, the next push runs CI automatically.

## Open by design

- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers (documented in `AGENTS.md`). **Re-verified 2026-09-14 on Go 1.26.7: still required** (build fails without the flag, also transitively through `go-finding`). Track the Go release that graduates `json/v2` from experimental and remove the constraint then.

---

_The 2026-09-14 v2.1.0 session resolved: the unresolvable `v2.0.0` tag (module path migrated to `/v2`), the stale branching-flow pins (re-pinned to analyzer reality: 18 PHANTOM, 22 stats; panic false positive suppressed), `WithConcurrency(n)`, context-carrying rules (`ContextRule` + `NewContextRule`), the `listeners/otel` module, the Datastar rewrite of `examples/sse`, Stream benchmarks, goroutine-leak regression tests, the all-modules `check-all` flake app, the CI module matrix, and Polish-Customs compatibility verification against `/v2`. All recorded in [`CHANGELOG.md`](CHANGELOG.md)._
