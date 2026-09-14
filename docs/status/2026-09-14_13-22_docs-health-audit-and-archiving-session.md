# Status Report — Docs-Health Audit & Historical-Docs Archiving Session

> **Generated:** 2026-09-14 13:22 CEST
> **Scope:** This session only — full `docs-health` AUDIT (read all `**/2026-0*` files → verify living docs against code → annotate historical docs inline → archive fully-resolved files). Point-in-time snapshot.
> **Repo state at writing:** `master` @ `0453a78` + auto-commits; working tree has this session's uncommitted doc edits (auto-commit daemon picks them up).
> **Skill followed:** `docs-health` SKILL.md (AUDIT = BUILD + HARVEST + VERIFY, plus ANNOTATE for historical docs).

---

## a) FULLY DONE

| #  | What                                                                                                                                                                    | Evidence                                                                                                                                                            | Scope                                                        |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------ |
| 1  | Read ALL 19 `**/2026-0*` files (17 status reports + 2 planning docs) plus the 5 root-level proposal/review docs and all 6 living docs                                    | this report's analysis                                                                                                                                              | `docs/status/*`, `docs/planning/*`, root `*.md`              |
| 2  | Established ground truth independently: ran test suite (156/158, 2 pre-existing branching-flow failures), coverage (**95.9%**), golangci-lint (**0 issues**)            | `nix develop --command go test ./...`, `go test -cover`, `golangci-lint run`                                                                                        | root module                                                  |
| 3  | **Discovered CI is dead: workflow `disabled_manually` since 2026-07-17** — no runs in 2 months; every run since 2026-06 failed in 3-5s (setup-level)                    | `gh workflow list --all` → `CI disabled_manually 246463642`; `gh run list --workflow ci.yml` (latest run 2026-07-15)                                                | `.github/workflows/ci.yml`                                   |
| 4  | **Discovered repo is PRIVATE** — `proxy.golang.org` has zero cached versions, pkg.go.dev 404s; the module can never be indexed while private                            | `gh repo view --json visibility` → `PRIVATE`; proxy `@v/list` and `@latest` both 404                                                                                | repo settings                                                |
| 5  | Confirmed `go-finding` (the one runtime dep) is PUBLIC — CI needs no `GOPRIVATE`                                                                                         | `gh repo view LarsArtmann/go-finding --json visibility` → `PUBLIC`                                                                                                  | dependency reality                                           |
| 6  | CHANGELOG: fixed the broken `v1.0.0`/`v1.1.0` release-link definitions (404s live since the v2.0.0 release; flagged in the 2026-07-26 review, never fixed until now)    | `CHANGELOG.md` link defs removed + HTML comment pointing at the versioning note; verified GitHub release body never contained the links (`gh release view v2.0.0`) | `CHANGELOG.md`                                               |
| 7  | FEATURES.md: fixed 6 stale claims — DOMAIN_LANGUAGE "placeholder" (false, it's filled in), dangling "See below" adapter row, 145→156 specs, 94.8%→95.9% coverage, branching-flow row (2 red pins → PARTIALLY_FUNCTIONAL), **CI row FULLY_FUNCTIONAL → BROKEN (disabled)** | `FEATURES.md`                                                                                  | `FEATURES.md`                                                |
| 8  | README.md: removed dead GoDoc badge (404, private repo), added private-repo install note (`GOPRIVATE` requirement + v0.1.0-only reality), added Ecosystem section (adapter + example modules), fixed pointer/value inconsistency in the govalid example, 94.8%/145 → 95.9%/156 stats | `README.md`                                                              | `README.md`                                                  |
| 9  | TODO_LIST.md: harvested 2026-09-14 report items — added CI-disabled decision, branching-flow stale-pins policy, Stream benchmarks, goroutine-leak test, composite 3-module lint/test command; kept tag fix / Polish-Customs / json/v2 / CI matrix / WithConcurrency / ctx-rules / Datastar / OTel | `TODO_LIST.md`                                                            | `TODO_LIST.md`                                               |
| 10 | ROADMAP.md: removed the "Already shipped" trophy-case section (anti-pattern flagged in the 2026-07-26 self-review), added Event-Driven Validation candidates (scheduler/short-circuit, backpressure, signing, Run ID, tags, gopter equivalence, RuleID), added "Open questions" section (versioning endgame, branching-flow policy, private-CI, adapter home) | `ROADMAP.md`                        | `ROADMAP.md`                                                 |
| 11 | AGENTS.md: PHANTOM count 12→14 with pin-fragility documentation (analyzer has no path-exclude; any new module shifts counts); new "CI & Publishing Reality" section (CI disabled, repo private, local gates are the real quality bar) | `AGENTS.md`                                                                                 | `AGENTS.md`                                                  |
| 12 | `docs/DOMAIN_LANGUAGE.md`: added event vocabulary (Event, RuleEvaluated, ValidationCompleted, Listener, Stream, fact producer) + Stream operation                        | `docs/DOMAIN_LANGUAGE.md`                                                                                                                                           | domain glossary                                              |
| 13 | ANNOTATE pass over ALL 18 historical docs: every numbered item resolved inline with `~~item~~ done at <hash>` / `done (<evidence>)` / `done (docs-health pass <where it moved>)` / `Won't implement` — no appendix-only files | `docs/status/*` (13 March + 3 July reports), `docs/planning/2026-03-15_*`, `docs/planning/2026-09-14_*` (sections 3+4 fully struck)                        | 18 files                                                     |
| 14 | Resolution banners/appendices on the 6 non-numbered historical docs (IMPLEMENTATION_PLAN, MIGRATION_TO_NIX_FLAKES, FINDING-SDK, BDD_TESTS_REVIEW, INTEGRATION_DECISION, go-composable analysis) + 2 resolution appendices on the 2026-07-26 reports (their §g questions answered) | file headers/endings                                                         | 8 files                                                      |
| 15 | Level-aware annotation tooling: skill's `annotate-rows.py` only scopes to `## ` boundaries, so sub-section (`###`) tables with colliding row numbers could not be annotated; wrote a shape-checked level-aware variant and used it for all colliding tables | `/tmp/annotate-rows2.py` (session tool, not committed)                                                                                      | tooling gap in skill                                         |

## b) PARTIALLY DONE

| # | What works                                                                      | What remains                                                                                                                       | Blocker                          | Effort |
| - | ------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | -------------------------------- | ------ |
| 1 | All historical docs annotated and ready to archive                              | The `git mv` into `archived/` subdirectories not yet executed (stopped for this report)                                             | none                             | S      |
| 2 | Living docs verified against code                                               | Final post-edit quality gate (re-run tests/lint after doc edits) not yet executed                                                   | none                             | S      |
| 3 | Health report (Accuracy/Fitness scores per docs-health AUDIT format)            | Not yet printed (was mid-flight when this report was requested)                                                                     | none                             | S      |
| 4 | 2026-09-14 status report's harvest items                                        | Items f.29-f.31 (HARVEST/VERIFY/ANNOTATE follow-ups) are now moot — this session IS that harvest; the report itself stays in `docs/status/` as the current snapshot | none | —      |

## c) NOT STARTED

| # | Item                                                                                | Why not started                                             | Still wanted? |
| - | ----------------------------------------------------------------------------------- | ----------------------------------------------------------- | ------------- |
| 1 | `git mv` of 13 status reports + 4 planning docs + 5 root docs into archived locations | Report requested before the archive step                    | YES — next    |
| 2 | Post-edit quality gate re-run (tests + lint)                                        | Same                                                        | YES           |
| 3 | Inline health report (Accuracy/Fitness with visible math)                           | Same                                                        | YES           |
| 4 | Re-enabling/fixing the disabled CI workflow                                         | User decision (private-repo Actions minutes); in TODO_LIST  | Decision      |
| 5 | Root module re-versioning (`/v2` path + re-tag)                                     | User decision (blast radius); in TODO_LIST                  | Decision      |
| 6 | Branching-flow stale-pins policy (re-pin vs nightly vs keep-red)                    | User decision; in TODO_LIST                                 | Decision      |

## d) TOTALLY FUCKED UP

| # | What happened                                                                                                                                                          | Severity | Root cause                                                                                              | Mitigation/Status                                                                                     |
| - | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| 1 | First B-table strike script on `2026-03-20_07-29` missed 3 of 9 rows (regex only matched backtick-led cells, not `**`-led ones)                                        | Low      | Wrote the match pattern against the rows I had visible, not all rows in the table                       | Caught on read-back (`grep -c` returned 6, not 9); struck the remaining 3 with exact-match replacements |
| 2 | First attempt to annotate `2026-03-20_06-57` failed loudly (`row 1: expected 1 match, found 3`) — I assumed the skill script scoped to `###` headings; it only scopes to `## ` | None (script refused atomically) | Did not read the script's section logic before assuming                                                 | Wrote the level-aware variant; dry-ran before live use per the skill's own 2026-08-18 lesson            |
| 3 | Edit-tool rejection on the 09-14 plan banner ("file modified since last read") — my own Python strike-all had rewritten the file between view and edit                    | Low      | Mixed edit tools on the same file without re-viewing                                                    | Re-viewed, applied cleanly                                                                            |

Nothing destructive: no `rm`, no `git reset`, no force-push, no unrequested commits. All mutations are documentation edits; zero source-code changes this session.

## e) WHAT WE SHOULD IMPROVE

1. **The `disabled_manually` CI state was invisible to every local gate.** Two months of "CI badge in README" meaning nothing. Candidate: a docs-health VERIFY checklist item — `gh workflow list --all` when a README shows a CI badge. (Proposed to the docs-health skill upstream.)
2. **Private-repo reality should be a first-class VERIFY probe.** Multiple historical reports carry "verify on pkg.go.dev" tasks that were never achievable. One `gh repo view --json visibility` call kills a whole class of ghost tasks.
3. **The annotate-rows.py section scoping should respect heading level**, not just `## `. Three of today's files needed the workaround. Candidate PR to the docs-health skill.
4. **Status reports should stop proposing `justfile`/`Makefile` items.** Repo policy (flake.nix canonical) was re-proposed in 3 separate reports and rejected 3 times. AGENTS.md already records it; reports should check before listing.
5. **"Verify on pkg.go.dev" appeared in 6+ reports across 6 months** without anyone checking repo visibility — the harvest step now explicitly drops items whose premise is impossible, but the status-report skill could gate this at writing time.

## f) UP TO 50 THINGS TO GET DONE NEXT

Ranked by impact (S<30min M<2h L>2h). Items 1-3 are this session's unfinished tail; 4-10 are the user-decision cluster; the rest are TODO_LIST/ROADMAP fuel already filed there.

| #  | Task                                                                                                              | Impact   | Effort | Category   |
| -- | ----------------------------------------------------------------------------------------------------------------- | -------- | ------ | ---------- |
| 1  | Execute the archive moves: 13 status reports → `docs/status/archived/`, 4 planning docs → `docs/planning/archived/`, root docs → `docs/{reviews,adr,planning/archived}/` | High | S | This session |
| 2  | Re-run quality gate after doc edits (`go test`, `golangci-lint`, link check)                                      | High     | S      | This session |
| 3  | Print the docs-health health report (Accuracy/Fitness, per-doc table)                                             | Medium   | S      | This session |
| 4  | Decide: re-enable + fix CI (setup fails in 3-5s since June) or formally rely on local flake gates                 | Critical | S-M    | Decision   |
| 5  | Decide: module re-versioning (`/v2` path + re-tag vs compatible fresh tag)                                        | Critical | M      | Decision   |
| 6  | Decide: branching-flow stale pins (re-pin / keep-red / nightly)                                                   | High     | S      | Decision   |
| 7  | Decide: adapter home (nested here vs sibling repos)                                                               | Medium   | S      | Decision   |
| 8  | If CI re-enabled: add matrix jobs for `adapters/cqrslite` + `examples/sse`                                        | High     | M      | TODO_LIST  |
| 9  | Tag a release carrying events/streaming once #5 lands                                                             | High     | M      | TODO_LIST  |
| 10 | Polish-Customs integration (first real-world consumer)                                                            | High     | L      | TODO_LIST  |
| 11 | `WithConcurrency(n)` bound for `Stream`                                                                           | Medium   | S      | TODO_LIST  |
| 12 | Context-carrying rules (`Check(ctx)`)                                                                             | Medium   | L      | TODO_LIST  |
| 13 | Datastar leg for the SSE example                                                                                  | Medium   | M      | TODO_LIST  |
| 14 | OpenTelemetry listener module                                                                                     | Medium   | M      | TODO_LIST  |
| 15 | `BenchmarkStream*` + goroutine-leak regression test                                                               | Medium   | S      | TODO_LIST  |
| 16 | Composite lint/test command across all 3 modules (flake app)                                                      | Medium   | S      | TODO_LIST  |
| 17 | Propose annotate-rows.py level-scoping fix upstream to docs-health skill                                          | Low      | S      | Tooling    |
| 18 | Propose "CI badge vs `gh workflow list --all`" verify-checklist item upstream                                     | Low      | S      | Tooling    |
| 19 | Empty-rules + unicode edge-case specs (last gap from the 2026-03 BDD review)                                      | Low      | S      | ROADMAP    |
| 20 | Property-based Build≡Stream equivalence (gopter)                                                                  | Low      | M      | ROADMAP    |

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **CI endgame:** the workflow was disabled *manually* on 2026-07-17 — possibly deliberately (private-repo Actions minutes cost money). Re-enable and fix the 3-5s setup failure, or declare local flake gates the canonical CI and delete the workflow + badge? I tried: run history (fails since June, all setup-level), workflow state API, badge audit. The billing tradeoff is yours.
2. **Versioning endgame:** rename module to `.../v2` + re-tag (breaks every consumer import), or tag a compatible `v0.2.0`/`v1.0.0` on the current path (semver-misleading given the breaking changes already shipped)? Also entangled: the repo is private — is going public part of the plan? If yes, pkg.go.dev becomes achievable and the README badge comes back.
3. **Branching-flow pins:** re-pin `stats` (expects 36, actual 15) and the `all`-output spec to current reality (green but weaker), keep them red as a standing signal, or move branching-flow specs to a nightly job? They were red before today's session and survive every re-pin because the analyzer binary keeps drifting.

---

**Handoff note:** the archive `git mv` step, the final quality-gate re-run, and the inline health report are the only unfinished session tasks (§c.1-3). Everything else open is a user decision (§g) or already filed in TODO_LIST/ROADMAP. Per repo norms, no manual commit — the auto-commit daemon picks this file up.
