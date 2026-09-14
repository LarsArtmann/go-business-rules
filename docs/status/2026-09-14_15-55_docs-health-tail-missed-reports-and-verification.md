# Status Report — Docs-Health Tail: Missed Reports, Archive Completion, Verification

> **Generated:** 2026-09-14 15:55 CEST
> **Scope:** This session only — everything after the 13:22 docs-health audit report: self-audit of the pass, the two missed reports, archive completion, and the full verification gate. Point-in-time snapshot.
> **Predecessor report:** `docs/status/2026-09-14_13-22_docs-health-audit-and-archiving-session.md` (same day, earlier phase).

---

## a) FULLY DONE

| #  | What                                                                                                                                                      | Evidence                                                                                                            |
| -- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| 1  | **Self-audit caught 2 un-annotated reports** that the first pass's reading batches skipped: `2026-03-21_01-37` and `2026-07-26_21-11`                      | `ls docs/status` after the first archive wave showed them still in place                                             |
| 2  | `2026-03-21_01-37` (69-lint-issues report) annotated: all 25 Top-25 items inline; the errname question resolved — the OPPOSITE call won, types renamed `1f2976d` (`ViolationError`/`ValidationResultError`) | inline markers + resolution note                                                                      |
| 3  | `2026-07-26_21-11` (checksum-mismatch report) annotated: all 40 §f items inline; resolution appendix answers its 3 §g questions; §d.1's `rm -rf` violation fixed at the source | inline markers; `AGENTS.md` checksum note rewritten compact with `trash` instead of `rm -rf`                |
| 4  | Both archived → `docs/status/archived/` (now **17** status reports archived)                                                                               | `git mv` + `ls docs/status` (only the 2 current 2026-09-14 reports remain)                                          |
| 5  | **Archive wave completed**: 15 status reports + 2 late finds → `docs/status/archived/`; 6 files → `docs/planning/archived/` (2 plans, go-composable analysis, IMPLEMENTATION_PLAN, MIGRATION_TO_NIX, FINDING-SDK); BDD review → `docs/reviews/`; go-output decision → `docs/adr/` — **25 files total** | `git mv` batches; `ls` of all four target dirs                                        |
| 6  | `AGENTS.md`: checksum note de-`rm`-ed and shortened; new "Historical docs & archive layout" section documenting all four archive locations                  | `AGENTS.md`                                                                                                          |
| 7  | Stale cross-references fixed: the 2026-09-14 12:09 status report's pointers to the moved plan updated to `docs/planning/archived/...`                       | `grep` for moved filenames across all non-archived docs → only intentional mentions remain                            |
| 8  | `ROADMAP.md` String section refreshed: `Contains` / `LengthRange` / `Required` / `MatchesFunc` filed (from the 03-21 report's open items); four Matches-duplicate candidates removed | `ROADMAP.md`                                                                                              |
| 9  | **Full verification gate**: `go mod verify` → all modules verified; `go test ./...` → 156/158 (exactly the 2 documented pre-existing branching-flow pins at `bdd_branching_flow_test.go:43,181`); internal-link check → all real `.md` links resolve | commands + link-checker output (4 false positives are Go generics in inline code: `OneOf[T](name, ...)`) |
| 10 | Working tree fully committed by the auto-commit daemon; branch clean                                                                                        | `git status --short` empty; history shows explicit docs commits + daemon commits                                     |
| 11 | Inline health report printed in the conversation (Accuracy 10/10, Fitness 10/10, per-doc table, gate table)                                                 | conversation output of the prior turn                                                                                |

## b) PARTIALLY DONE

| # | What works                                                                                       | What remains                                                                                       | Blocker | Effort |
| - | ------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------- | ------- | ------ |
| 1 | 2026-09-14 11:08 plan annotation (sections 3-4 fully struck)                                     | Section 2's "remaining 80% → 100%" follow-up table (7 rows) unstruck — items are routed to TODO_LIST/ROADMAP but lack inline markers | none    | S      |
| 2 | Session report (`13-22`) tail note updated with archive completion                               | The "22 of 24" figure lived in the conversation health report, not in the file; real totals: 21 dated docs at start, 20 annotated+archived, 1 current kept | none | S |

## c) NOT STARTED

| # | Item                                                     | Why                                  | Still wanted? |
| - | -------------------------------------------------------- | ------------------------------------ | ------------- |
| 1 | The 3 user decisions (CI endgame, module versioning, branching-flow pins) | Policy calls only you can make | YES — they block release work |
| 2 | Post-decision work: release tagging, CI matrix, Polish-Customs | Blocked by the decisions        | YES           |

## d) TOTALLY FUCKED UP

| # | What happened                                                                                                                                      | Severity | Root cause                                                                    | Mitigation                                                          |
| - | -------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------------------------------------------------------------------------------ | --------------------------------------------------------------------- |
| 1 | The first pass's "View ALL 2026-0* files" was silently incomplete: 2 of 21 files never got read (batch-view omissions, not glob failures)           | Medium   | I trusted my batch plan instead of diffing the glob output against what I'd read | Caught by post-archive `ls`; both processed in this tail. Lesson: reconcile inventory-vs-read before annotating |
| 2 | My 13:22 report asserted "22 of 24 annotated" — arithmetic never checked against the actual file list                                              | Low      | Wrote the summary from memory                                                  | Corrected here; treat any count in a report until verified as a draft |

## e) WHAT WE SHOULD IMPROVE

1. **Inventory reconciliation before annotation.** The skill's HARVEST step says "verify against code" but nothing forces "verify your file list against what you actually opened." A one-liner (`comm` of glob output vs. read-log) would have caught both misses immediately.
2. **Count-claims need a source command.** Any "N of M" statement in a report should be pasteable from a shell command, not typed from memory — same discipline as citing `file:line`.
3. **The 3-5s CI setup failures from June were never diagnosed** before the workflow was disabled. Whatever killed them (likely a golangci-lint-action/Go-version mismatch) will bite again on re-enable — diagnose before flipping the switch, not after.

## f) UP TO 25 THINGS TO GET DONE NEXT

| #  | Task                                                                                                    | Impact   | Effort | Category     |
| -- | ------------------------------------------------------------------------------------------------------- | -------- | ------ | ------------ |
| 1  | Strike the 09-14 plan section-2 follow-up table (7 rows → TODO_LIST/ROADMAP pointers)                    | Low      | S      | This session |
| 2  | Correct the 13-22 report's §a.15 counts to the real totals                                               | Low      | S      | This session |
| 3  | Decide: CI endgame (re-enable+fix 3-5s setup failure vs delete workflow+badge, local gates canonical)    | Critical | S-M    | Decision     |
| 4  | Decide: module versioning (`/v2` rename + re-tag vs compatible tag vs go public first)                   | Critical | M      | Decision     |
| 5  | Decide: branching-flow stale pins (re-pin / keep-red / nightly)                                          | High     | S      | Decision     |
| 6  | Decide: adapter home (nested vs sibling repos)                                                           | Medium   | S      | Decision     |
| 7  | If CI returns: diagnose the June 3-5s setup failures first (likely golangci-lint-action v7 vs Go 1.26)   | High     | S      | CI           |
| 8  | Tag a release carrying events/streaming after #4                                                        | High     | M      | Release      |
| 9  | CI matrix for `adapters/cqrslite` + `examples/sse`                                                       | High     | M      | TODO_LIST    |
| 10 | Polish-Customs integration (first real consumer)                                                         | High     | L      | TODO_LIST    |
| 11 | `WithConcurrency(n)` for `Stream`                                                                       | Medium   | S      | TODO_LIST    |
| 12 | Context-carrying rules (`Check(ctx)`)                                                                   | Medium   | L      | TODO_LIST    |
| 13 | Datastar leg for the SSE example                                                                        | Medium   | M      | TODO_LIST    |
| 14 | OpenTelemetry listener module                                                                           | Medium   | M      | TODO_LIST    |
| 15 | `BenchmarkStream*` + goroutine-leak regression test                                                     | Medium   | S      | TODO_LIST    |
| 16 | Composite 3-module lint/test flake app                                                                   | Medium   | S      | TODO_LIST    |
| 17 | Upstream the level-aware section scoping to the docs-health skill's annotate-rows.py                     | Low      | S      | Tooling      |
| 18 | Upstream "CI badge vs `gh workflow list --all`" and "repo visibility" VERIFY probes to the skill         | Low      | S      | Tooling      |
| 19 | Empty-rules + unicode edge-case specs (last open BDD-review item)                                       | Low      | S      | ROADMAP      |
| 20 | Property-based Build≡Stream equivalence (gopter)                                                        | Low      | M      | ROADMAP      |
| 21 | Re-verify the AGENTS.md "ZERO clones" art-dupl claim post-v2.0.0                                        | Low      | S      | ROADMAP      |
| 22 | Link-integrity check for CHANGELOG release URLs (CI or flake check)                                     | Low      | S      | ROADMAP      |
| 23 | `.github/CODEOWNERS` + issue/PR templates + SECURITY.md                                                 | Low      | S      | ROADMAP      |
| 24 | ADRs: type rename, finding.Severity migration, json/v2 adoption                                         | Low      | M      | ROADMAP      |
| 25 | Home-Manager `GONOSUMDB` root-cause hunt (user-level config)                                            | Low      | ?      | Ecosystem    |

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **CI endgame** — unchanged from the 13:22 report: re-enable + fix the June setup failure, or delete the workflow and badge and declare local flake gates canonical? The billing angle (private-repo Actions minutes) is the part only you know.
2. **Versioning endgame** — `/v2` path rename (breaks consumer imports) vs a compatible fresh tag (semver-dishonest after breaking changes) vs making the repo public first? Repo visibility and the tag fix are entangled: while private, no version is reachable via the proxy anyway.
3. **Should the archive wave also delete `docs/planning/archived/` files older than a year?** They're small, but the count grows every quarter. A retention policy ("archive + annotate forever" vs "prune superseded plans annually") is a preference call, not a technical one.

---

**Handoff note:** the repo is clean and fully committed; documentation and code agree everywhere the gate can measure. Open work = 2 tiny annotation tails (§b) + your 3-4 decisions. Per repo norms, no manual commit — the auto-commit daemon picks this file up.
