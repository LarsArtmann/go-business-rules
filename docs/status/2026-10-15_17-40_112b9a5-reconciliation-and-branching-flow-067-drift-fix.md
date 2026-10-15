# Status Report — 112b9a5 Reconciliation + branching-flow 0.6.7 Drift Fix

- **Written:** 2026-10-15 17:40 CEST
- **Session scope:** resume from the 2026-10-11 green-recovery session; mandated steps were (1) todo reconciliation, (2) reconcile pin/nolint docs with commit 112b9a5 (native directive form migration) + re-verify pins, (3) then wait. User directive this session: "READ, UNDERSTAND, RESEARCH, REFLECT … Execute and Verify them one step at the time. Repeat until done."
- **Repo at report time:** master, working tree clean, HEAD `e5caf1b` (daemon). Session work committed by the auto-commit daemon as `0ade223` → `968e8a7` → `e5caf1b`.
- **Environment drift discovered:** branching-flow **0.6.4 (192031d) → 0.6.7 (aa1611d)** via the 2026-10-14 flake bump (`1df5c61`); BuildFlow binary now `518739c9-dirty`.

## TL;DR

The mandated reconciliation is **done and verified at every level**: docs now describe the native `//branching-flow:ignore panic` directive form, both directives were re-proven load-bearing by individual removal, the dead `.golangci.yml` nolintlint exclusion is gone, and — beyond the mandate — a **new silent analyzer drift (0.6.7)** was caught, root-caused (exit-code semantics + stderr pollution), fixed in the BDD gate, re-pinned, and documented. All gates green: root suite ok, golangci-lint 0 issues, `check-all` ALL 4 MODULES GREEN, full BuildFlow gate **exit 0**.

## What I forgot / could have done better (self-critique, honest)

1. **Forgot to check analyzer versions at session start.** The flake bump (`1df5c61`, two days earlier) had already moved branching-flow to 0.6.7. One `branching-flow --version` against the pin's version would have surfaced the drift in 5 seconds instead of after a red 11-spec test run.
2. **Repeated last session's two-mode lesson, almost.** My first manual verification piped `2>/dev/null`, discarding exactly the stderr that was breaking the tests (`CombinedOutput` merges it into the JSON). I nearly concluded "tests are stale, analyzer is fine" — my manual jq run was green while the suite was red. I caught it only because I kept digging after the failures persisted. This is the **second occurrence** of the "filtered one output mode, concluded from the other" class; it belongs in `references/lessons.md` (cross-project) and the gotchas doc.
3. **Self-inflicted, immediately-recovered file damage:** a malformed `multiedit` (old_string included context the new_string dropped) deleted the `checkNoViolations` declaration line. Caught on the very next view, repaired in the same batch. Never should have been constructed.
4. **Verification-by-sed almost verified nothing:** the first directive-removal test used 1 tab of indentation against a 3-tab directive, so it tested unmodified code. Caught because `git diff --stat` printed nothing for that file — the habit of demanding diff evidence before trusting a mutation saved the conclusion.
5. **Doc edits ran in parallel with a red suite.** No harm resulted (baseline failures predated my edits), but the disciplined order is: baseline the suite, then edit.
6. **The gotchas doc records the exit-code change but not the stderr/JSON-corruption mechanism** — that detail lives only in the test's comment. The doc is one step behind the code.

## a) FULLY DONE

| Item | Evidence |
| --- | --- |
| Todo list reconciled (items 4–6 marked complete from 2026-10-11 session; reconciliation task added, then completed) | todo tool state: 7/7 completed |
| Commit archaeology of 112b9a5 and the surrounding daemon commits (`011321b` deletion → `ee9552c` restore → `112b9a5` migration → `acc7573` report → `1df5c61` flake bump) fully reconstructed | `git show` outputs, this session's transcript |
| Native directives re-proven load-bearing **individually**: removing `validator.go:186` re-surfaces the Stream channel-send finding; removing `builders_collection.go:109` re-surfaces the DivisibleBy division-by-zero finding; both restored, tree clean | live `branching-flow panic .` runs, 2026-10-15 |
| Pins verified live under 0.6.7: phantom = 20 (all warning), stats = 24 (20+2+1+1), panic = 0 | live `phantom --format finding` / `stats --format json` / `panic` runs |
| BDD gate fixed for 0.6.7: `runBFCommand` now reads stdout only (stderr captured to a `strings.Builder`, surfaced in the failure message) and tolerates the findings-exit (`cmd.ProcessState.ExitCode() == 1`); same tolerance in the `all` spec; pin comments updated with dated 0.6.7 re-verification | commit `968e8a7`/`e5caf1b`; root suite ok (250 ran, 1 skipped — the pre-existing skip) |
| Dead `.golangci.yml` nolintlint exclusion `(validator|builders_collection)\.go` removed — its reason (nolint unknown-linter warnings) no longer exists (zero `//nolint` directives remain anywhere in the repo) | grep sweep; `buildflow -s golangci-lint` → 0 findings |
| AGENTS.md pin bullet rewritten (native directive form, 0.6.7 exit-code note, exclusion removal) — 216/220 lines, under the agent-config gate | commit `0ade223`; `wc -l AGENTS.md` = 216 |
| `docs/agents-gotchas.md` reconciled: suppression-migration paragraph, 0.6.7 version-history entry, stale-binary-hint note | commit `0ade223` |
| Verification chain green: root suite ok → golangci-lint 0 issues → `nix run .#check-all` ALL 4 MODULES GREEN → full BuildFlow gate **EXIT 0** (90/103 steps, 51 detect-only warnings — vulnix CVE closure noise and doctor availability noise, all documented classes) | command outputs, 2026-10-15 17:21–17:34 |
| wsl_v5 whitespace findings from my own edits resolved (3 iterations) | final `buildflow -s golangci-lint` → "BuildFlow passed", no findings |
| Stale-reference sweep: living docs mention the old `//nolint` form only as historical context; CHANGELOG untouched (append-only) | grep sweep at session end |

## b) PARTIALLY DONE

| Item | Works | Open gap | Blocker | Effort |
| --- | --- | --- | --- | --- |
| Doc consistency after 112b9a5 + 0.6.7 | AGENTS.md, gotchas doc, test comments all describe native form + 0.6.7 reality | gotchas doc lacks the stderr-corruption mechanism (lives only in the test comment); re-pin procedure doesn't say "check binary version first" | none | S |
| 2026-10-11 status report | Committed (`acc7573`) as point-in-time snapshot | Not annotated with what 112b9a5/0.6.7 superseded (docs-health ANNOTATE pending) | none | S |
| BuildFlow binary freshness | Full gate exits 0 on `518739c9-dirty` | `-dirty` build unidentified; not investigated per session scope | none (deliberately out of scope) | S |
| Upstream-issue candidates (branching-flow stale hint, exit-code docs; BuildFlow cache invalidation) | Diagnosed, evidence captured | Not filed — awaiting your file-vs-absorb decision | your decision | M |
| CHANGELOG | Accurate history through the 0.6.4 era | No entries for 2026-10-11 session or this session | none | S |

## c) NOT STARTED

| Item | Why not started | Still wanted? |
| --- | --- | --- |
| HARVEST of both status reports' section (f) into TODO_LIST.md / ROADMAP.md | Standing WAIT-for-instructions contract from 2026-10-11; you have not directed it | presumably yes — it's the docs-health loop |
| govulncheck re-verification under a newer Go (open item from 2026-10-11) | Not in this session's mandate | yes (documented with a re-verify instruction) |
| Upstream issues (see b) | Blocked on your decision | your call |
| lychee policy decision (GITHUB_TOKEN vs excluding `github.com/larsartmann/*`) | Blocked on your decision | your call |
| SystemNix Go ≥ 1.27.1 deployment (would kill the devshell-invocation workaround) | Blocked on your deployment timing | your call |
| Version-aware BDD gate (parse `tool.version` from branching-flow JSON and fail/warn on drift) | Discovered this session, not implemented | recommended |
| Property test run (`GBR_PROPERTY=1`) — not run this session or last | Opt-in, out of mandate | yes, cheap |

## d) TOTALLY FUCKED UP

**Nothing. Zero known breakage at report time.** Every gate is green; the working tree is clean; the last full gate exited 0. The honest entries here are the two **recovered** self-inflicted incidents (checkNoViolations deletion; wrong-indentation sed), listed under the self-critique above with their mitigations — neither left a trace. Residual ugliness that is *not* repo breakage, for the record:

- branching-flow 0.6.7's experimental warning still recommends the **deprecated** `//nolint:branching-flow:panic` syntax — a stale upstream hint, actively misleading (it would reintroduce the golangci-lint unknown-linter noise we just removed).
- vulnix reports 30 advisory CVEs against the nixpkgs build-toolchain closure — documented noise, advisory severity, non-gating.
- 9 BuildFlow tools unavailable per health check — documented doctor noise, non-gating.

## e) WHAT WE SHOULD IMPROVE

1. **Version-gate the analyzers, not just the counts.** The BDD gate pins counts but ignores binary identity; 0.6.7 drifted both counts-semantics (exit codes) and output plumbing. The 0.6.7 finding JSON even carries `tool.version` — assert or warn on it in `runPhantomCommand`. Prevents the entire "red suite, mystery cause" class.
2. **Never verify an analyzer with stderr discarded.** Twice now a filtered output mode nearly produced a wrong conclusion. Rule: manual verification reads stdout and stderr together, and the test helpers decide how to separate them.
3. **Kill the three-way pin split brain.** Every re-pin touches `bdd_branching_flow_test.go` + `docs/agents-gotchas.md` + the AGENTS.md bullet. A single pins file (JSON/YAML) read by the test and rendered into docs — or at minimum a re-pin checklist script — removes the drift surface. The gotchas procedure enumerates the three locations, but enumeration is not enforcement.
4. **Pre-baseline the suite before editing anything.** One background `go test` at session start separates pre-existing drift from self-inflicted breakage.
5. **Edit-tool hygiene:** diff old_string vs new_string mentally before submitting a multiedit — the checkNoViolations deletion came from asymmetric context. A "the replacement must preserve everything outside the intended delta" habit would have caught it.
6. **AGENTS.md headroom is 4 lines.** The 220-line agent-config gate will bite the next substantive addition; plan the next restructure proactively instead of under gate pressure.
7. **Record the cross-project lesson.** "Analyzer verification must cover all output modes" has now recurred across two sessions and two analyzers — it qualifies for `references/lessons.md` in crush-config (committed there, not written in-session).
8. **BuildFlow speed claim is stale:** the skill says full mode takes ~5–10 min; this full gate ran in 29.2 s. Either the mode changed under the new binary or steps were skipped/cached — worth one `--verbose` look so we trust what "full" means again.

## f) 50 things we should get done next

*(Brainstorm for docs-health HARVEST — most items are ROADMAP fuel; ranked by impact within tiers. Impact / Effort: S <30 min, M 30 min–2 h, L >2 h.)*

**Tier 1 — protect the gate we just fixed**

| # | Task | Impact | Effort | Category |
| --- | --- | --- | --- | --- |
| 1 | Parse `tool.version` from branching-flow JSON in the BDD gate; warn/fail on version ≠ pinned version | Critical | S | Quality |
| 2 | Add "check `branching-flow --version` / `brandid-lint --version` against the pin" as step 1 of the re-pin procedure in the gotchas doc | Critical | S | Documentation |
| 3 | Record the "never verify with stderr discarded / all output modes" lesson in gotchas + cross-project `lessons.md` (committed in crush-config) | High | S | Documentation |
| 4 | Add the stderr-corruption mechanism to the gotchas 0.6.7 entry (currently only in the test comment) | High | S | Documentation |
| 5 | Unit-test `runBFCommand`'s contract with a fake binary: exit 1 tolerated, exit 2 fails, stderr surfaces in the failure message | High | M | Quality |
| 6 | Version-gate brandid-lint in the BD003 spec the same way (it has no drift protection at all) | High | S | Quality |
| 7 | Annotate the 2026-10-11 status report with what 112b9a5 + 0.6.7 superseded (docs-health ANNOTATE, inline, non-destructive) | Medium | S | Documentation |
| 8 | HARVEST sections (f) of the 2026-10-11 and this report into TODO_LIST.md / ROADMAP.md | High | M | Documentation |
| 9 | Write CHANGELOG entries for 2026-10-11 (SECURITY.md, re-pin, AGENTS restructure, gotchas doc) and 2026-10-15 (0.6.7 drift fix, exclusion removal) | Medium | S | Documentation |
| 10 | Plan the next AGENTS.md restructure before the 220-line gate bites (216/230? — currently 216/220, 4 lines headroom) | Medium | M | Documentation |

**Tier 2 — upstream closure (close the loop per buildflow skill)**

| # | Task | Impact | Effort | Category |
| --- | --- | --- | --- | --- |
| 11 | File upstream (branching-flow): experimental warning recommends deprecated `//nolint:branching-flow:panic` form — verify at source per verify-before-filing first | High | M | Bug |
| 12 | Check branching-flow 0.6.5–0.6.7 changelogs for the exit-code semantics change; document as deliberate or file an issue asking for documented exit codes / a `--no-fail` machine mode | High | M | Bug |
| 13 | File upstream (BuildFlow): result cache doesn't invalidate on tool-version change (2026-10-11 finding; re-verify against current binary first) | High | M | Bug |
| 14 | Re-run full gate WITHOUT `BUILDFLOW_NO_RESULT_CACHE` to confirm the cache no longer replays stale findings after tool-version changes | Medium | S | Quality |
| 15 | Identify the 9 unavailable BuildFlow tools (`--verbose`) and classify each as mattering or not | Medium | S | Quality |
| 16 | Run `buildflow doctor` to baseline environment health post-flake-bump | Medium | S | Quality |
| 17 | Identify why full mode ran in 29.2 s (was test-race actually among the 90 steps? check the log) — re-establish trust in "full" | High | S | Quality |
| 18 | Investigate the `-dirty` BuildFlow build (`518739c9-dirty`); rebuild cleanly or document the dirt | Low | S | Cleanup |

**Tier 3 — CI and release hygiene**

| # | Task | Impact | Effort | Category |
| --- | --- | --- | --- | --- |
| 19 | Confirm the latest CI run on current HEAD is green (local gates green; remote unverified this session) | High | S | Quality |
| 20 | Echo analyzer versions in CI (branching-flow/brandid-lint) so CI-vs-devshell binary skew is visible in logs | Medium | S | Quality |
| 21 | Decide CI posture for the BD003 spec's skip-when-absent behavior (install brandid-lint in CI or accept the skip) | Medium | S | Quality |
| 22 | Run gitleaks + codespell on-demand before the next release (they never run in pipeline modes) | Medium | S | Security |
| 23 | Run markdown-lint on-demand over the docs touched this session | Low | S | Quality |
| 24 | Dry-run `golangci-lint-auto-configure` to confirm it doesn't want to re-add anything we just removed | Medium | S | Quality |
| 25 | Run the opt-in property test (`GBR_PROPERTY=1`) once against the current toolchain | Medium | S | Quality |
| 26 | Re-run event/streaming benchmarks and refresh the numbers in the gotchas doc (last measured 2026-10-11) | Low | S | Quality |
| 27 | Verify the release checklist pairing rule is still discoverable: bump `doc.go` Version AND `suite_test.go` expectation together before the next tag | Medium | S | Documentation |
| 28 | Set a post-flake-bump habit: re-run the analyzer-gate suite after every daemon flake.lock commit (they move analyzer binaries silently) | High | S | Quality |

**Tier 4 — docs and product surface**

| # | Task | Impact | Effort | Category |
| --- | --- | --- | --- | --- |
| 29 | Link the new SECURITY.md Security Practices section from README | Low | S | Documentation |
| 30 | Spot-check lychee links in the new/edited gotchas sections | Low | S | Documentation |
| 31 | Decide lychee policy: GITHUB_TOKEN vs excluding `github.com/larsartmann/*` (needs you) | Medium | S | Documentation |
| 32 | Re-evaluate the GOEXPERIMENT=jsonv2 requirement against current Go (json/v2 graduation status; standing item) | High | M | Feature |
| 33 | Smoke-test `examples/sse` Datastar demo against the current root module | Medium | M | Quality |
| 34 | Smoke-test `listeners/otel` event→span/metric mapping | Medium | M | Quality |
| 35 | Smoke-test `adapters/cqrslite` bus listener | Medium | M | Quality |
| 36 | Consider exporting the findings-exit tolerance as a tiny reusable helper (other LarsArtmann repos shell out to the same analyzers) | Low | M | Feature |
| 37 | Document the BDD gate's binary-skip behavior (specs skip when branching-flow/brandid-lint absent) for contributors | Low | S | Documentation |
| 38 | Add a version matrix to the gotchas doc: 0.2.0 / 0.6.4 / 0.6.7 × (counts, exit codes, stdout/stderr split) | Medium | S | Documentation |
| 39 | Review whether `expectBFOutputContains("stats", "Total")` should assert something more specific than "Total" (weak assertion) | Low | S | Quality |
| 40 | Consolidate pin knowledge into a single source (pins file read by the test; docs rendered or cross-linked) | High | L | Cleanup |

**Tier 5 — carried-over strategic items (unchanged, still awaiting you)**

| # | Task | Impact | Effort | Category |
| --- | --- | --- | --- | --- |
| 41 | Decide upstream-issue filing overall (merges #11–#13 into one go/no-go) | High | S | Decision |
| 42 | SystemNix Go ≥ 1.27.1 deployment timing (deletes vs perpetuates the devshell-only BuildFlow workaround) | High | — | Decision |
| 43 | govulncheck re-verification under a newer Go toolchain | Medium | M | Quality |
| 44 | Re-verify the residual non-gating findings policy list against the current BuildFlow binary (nix-checker, vulnix, doctor classes) | Medium | M | Documentation |
| 45 | Confirm GOPRIVATE/GONOSUMDB no-op overrides are still needed (OS-level `GONOSUMDB` interplay) post-public-repo | Low | S | Cleanup |
| 46 | Re-confirm `adapters/cqrslite` go 1.27.1 floor is still enforced by `go mod tidy` after flake bumps | Medium | S | Quality |
| 47 | Consider a dependabot config refresh via `buildflow -s dependabot-auto-configure --dry-run` | Low | S | Cleanup |
| 48 | Run `buildflow --semantic` jscpd pass to re-confirm the art-dupl ZERO-clones claim on current code | Low | M | Quality |
| 49 | Review `flake.nix` checks (`skills-no-store-links` etc.) still pass after daemon flake bumps | Medium | S | Quality |
| 50 | Archive-or-annotate this session's predecessor report chain once HARVEST lands, so `docs/status/` stays a clean timeline | Low | S | Documentation |

## g) Three questions I cannot figure out myself

1. **File-or-absorb the upstream issues?** Three diagnosed candidates are waiting on your call: (a) branching-flow's experimental warning recommending the deprecated `//nolint:branching-flow:panic` form, (b) branching-flow 0.6.7's undocumented exit-code semantics change (exit 1 = findings present), (c) BuildFlow's result cache replaying findings after tool-version changes. My inclination is to file all three (source-verified first, in your voice per github-voice), but (b) in particular could also be "accepted convention" you'd rather absorb and document.
2. **Lychee policy:** use a GITHUB_TOKEN so the scanner can check `github.com/larsartmann/*` links, or exclude that host from scanning entirely? I cannot decide your token-hygiene posture; both work.
3. **SystemNix Go ≥ 1.27.1 deployment timing:** when the system toolchain reaches the dependency floor, the "ambient buildflow fails, always use `nix develop --command`" workaround and its documentation become obsolete. Deploy soon (I delete the workaround docs) or far out (I keep maintaining them)?

---

*Next step after this report: docs-health HARVEST on sections (f) of this and the 2026-10-11 report — on your go-ahead. Then I wait for instructions.*
