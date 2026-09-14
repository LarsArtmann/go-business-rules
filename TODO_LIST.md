# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-09-14 (v2.1.0 session: `/v2` module path, `WithConcurrency`, `ContextRule`, `listeners/otel`, Datastar example, `check-all` flake app, branching-flow re-pin, CI billing root cause, Polish-Customs `/v2` compatibility)

---

## Release follow-ups (user actions only)

- [ ] **Push `master` and the `v2.1.0` tag** (blocking publishing): the annotated `v2.1.0` tag is cut locally on the release commit but is **not pushed** (pushes require explicit approval). Commands: `git push origin master && git push origin v2.1.0`. The module path now ends in `/v2`, so this tag resolves as `github.com/LarsArtmann/go-business-rules/v2@v2.1.0` (verified end-to-end against a local file proxy).
- [ ] **Drop the temporary `replace` in Polish-Customs** once the tag is pushed: `polish-customs/go.mod` currently has `replace github.com/LarsArtmann/go-business-rules/v2 => /home/lars/projects/go-business-rules` so it consumes the local tree. Delete the line and re-run `go mod tidy` after pushing.
- [ ] **Fix GitHub Actions billing, then re-enable CI** (discovered 2026-09-14): every failed run since 2026-06 (e.g. run `29447520877`) was rejected at job start — *"recent account payments have failed or your spending limit needs to be increased"*. The workflow YAML was never the problem. Fix billing / raise the spending limit, then `gh workflow enable CI`. The rewritten `.github/workflows/ci.yml` already runs a matrix over all four modules; its exact commands are verified locally by `nix run .#check-all`.

## Open by design

- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers (documented in `AGENTS.md`). **Re-verified 2026-09-14 on Go 1.26.7: still required** (build fails without the flag, also transitively through `go-finding`). Track the Go release that graduates `json/v2` from experimental and remove the constraint then.

---

_The 2026-09-14 v2.1.0 session resolved: the unresolvable `v2.0.0` tag (module path migrated to `/v2`), the stale branching-flow pins (re-pinned to analyzer reality: 18 PHANTOM, 22 stats; panic false positive suppressed), `WithConcurrency(n)`, context-carrying rules (`ContextRule` + `NewContextRule`), the `listeners/otel` module, the Datastar rewrite of `examples/sse`, Stream benchmarks, goroutine-leak regression tests, the all-modules `check-all` flake app, the CI module matrix, and Polish-Customs compatibility verification against `/v2`. All recorded in [`CHANGELOG.md`](CHANGELOG.md)._
