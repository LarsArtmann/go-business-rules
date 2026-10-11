# Status Report: BuildFlow Green Recovery & Analyzer Re-Pin

**Date:** 2026-10-11 10:31 CEST
**Session scope:** one BuildFlow full-mode run handed in red (exit 69) → driven to green (exit 0), plus everything the fixes dragged into the light.
**Binary note:** this session ran against TWO BuildFlow binaries (ambient `25b07c4`, then `200ab99` after the devshell rebuild picked up the toolsdk-provider-switchover branch). Findings below are labeled where the binary identity matters.

| Metric                | Session start                         | Session end                                     |
| --------------------- | ------------------------------------- | ----------------------------------------------- |
| BuildFlow full gate   | exit 69 (3 steps failed)              | **exit 0** (89/102 ran, gate satisfied)         |
| Root suite            | 242 passed / **8 failed** / 1 skipped | all green (`ok`)                                |
| check-all (4 modules) | not run in-session start state        | **ALL MODULES GREEN** (build+vet+test+lint)     |
| AGENTS.md             | 376 lines (new error-level rule)      | **211 lines** (< 220 limit)                     |
| Error-level findings  | securitymd 1 + agent-config 1         | **0**                                           |
| Detect-only findings  | ~8 tools                              | 12 tools / 106 findings, all documented classes |

## What the red run actually was (diagnosis)

1. **test-race, 8 failed specs** — branching-flow pin drift. The pins in `bdd_branching_flow_test.go` were made against branching-flow **0.2.0**; the installed binary is **0.6.4 (192031d)**. 0.6.4 rewrote the PHANTOM rule set (metadata `name` is now the pair, e.g. `"value, minimum"`; severities collapsed to warning-only) and `stats` no longer counts the nolint-suppressed panic findings. Real counts: **20 PHANTOM** (all PHANTOM-TRANSPOSE same-type param pairs, all warning, four builder files), **stats total 24** = 20 + 2 flagparam + 1 ifacecomplete + 1 mixins. Zero new violations — pure analyzer drift.
2. **govalid-generate (×2)** — toolchain skew, environment not code: BuildFlow forces `GOTOOLCHAIN=local` for tools; `go-cqrs-lite/snapshot/v4@v4.5.1` requires go ≥ 1.27.1; ambient system Go is older. Even the user's `GOTOOLCHAIN=auto` export cannot help — BuildFlow overrides it. The devshell ships go 1.27.2 + a matching govalid. Canonical invocation: `nix develop --command buildflow --build-mode=full`.
3. **securitymd** — SECURITY.md lacked a "Security Practices" section (error-level → gate-blocking).
4. **New this session:** the updated BuildFlow binary shipped a **gating `agent-config` rule** (go-structure-linter): AGENTS.md > 220 lines is an _error_, previously a non-gating preflight warning.

## a) FULLY DONE

1. **Branching-flow re-pin, verified twice.** Diffed actual findings vs pins _before_ touching anything (20 findings enumerated by file:line:name), rewrote the pins in `bdd_branching_flow_test.go` (total 20; warning=20; no critical/error/info; DescribeTable entries on current metadata names `value, minimum` / `value, divisor` / `name, value`; a new "flags every builder file" spec replacing the dead bool-param one; stats 31→24), re-measured live after the edit: panic 0, stats 24 — exact match. Suite green.
2. **SECURITY.md "Security Practices" section** added (minimal dependency surface; automated scanning: govulncheck/gosec/vulnix/gitleaks; no-secrets no-I/O surface; hermetic Nix builds). Single-step verified: `buildflow -s securitymd` → passed.
3. **AGENTS.md restructured 376 → 211 lines** with zero concept loss; every removed narrative carried into the new `docs/agents-gotchas.md` (toolchain-skew story, GOPRIVATE/GONOSUMDB recipe, CI reality, json/v2 migration detail, branching-flow pin history + re-pin procedure, per-analyzer false-positive rationales, residual-findings policy, benchmark numbers, lychee gotchas). `agent-config` finding cleared (cache-replay trap diagnosed, see d/e).
4. **Toolchain-skew root cause confirmed and documented as the canonical rule** (ambient buildflow cannot work, full stop — not a "try GOTOOLCHAIN=auto" situation).
5. **Full verification chain green:** `nix develop --command go test ./...` ok · `nix run .#check-all` ALL MODULES GREEN · `nix develop --command buildflow --build-mode=full` **exit 0**.
6. **Residual-findings policy extended** in the gotchas doc: govulncheck "requires newer Go version" warnings identified as detection-environment noise (embedded go1.27.0 toolchain inside govulncheck vs source requiring 1.27+), doc-files-age advisories classified.
7. **Wrong intermediate state fully reverted** (see d) — final tree is correct, and the nolint requirement is now documented with _empirical_ evidence (with/without-nolint measurements, both output modes).

## b) PARTIALLY DONE

1. **Panic-suppression semantics:** final state correct, understanding arrived at via a wrong turn (d1). The gotchas doc records the verified behavior, but the _upstream reporting quirk_ (below) is documented locally only, not filed.
2. **Branching-flow reporting quirk (candidate upstream issue):** with the nolints present, `panic --format finding` reports `total: 0, suppressed: 0` — the suppression is invisible _and mislabeled_; plain-text mode shows the truth; 0.2.0's stats counted suppressed panics, 0.6.4's doesn't. I have the complete with/without evidence collected. Documented in `docs/agents-gotchas.md`; **not filed upstream**.
3. **govulncheck skew:** classified as noise with a written re-verification instruction, but **not verified** with a govulncheck built on the current toolchain.
4. **Gotchas doc spec-count claim:** inherited "251/251 specs"; actual recent runs show **250 ran + 1 skipped** (the skip pre-existed this session). Not reconciled.
5. **Git history hygiene:** the auto-commit daemon captured the wrong intermediate state (nolint deletion) in history before my revert. Final state correct; history contains a 2-commit wrong-then-right pair.

## c) NOT STARTED (noticed in-session, deliberately untouched)

1. `doc-files-age-check` flagged **3 files** (FEATURES.md 552h stale is one) — not even enumerated, let alone triaged.
2. **CHANGELOG.md** has no entries for: pin re-pin, SECURITY.md section, AGENTS.md restructure.
3. BuildFlow **binary-freshness advisory** (binary 200ab99 vs HEAD b4539d2 in the BuildFlow repo) — rebuild/reinstall not run; results reflect a branch a few commits behind.
4. **lychee policy decision** (fleet-level): `GITHUB_TOKEN` for authenticated checks vs `exclude` for `github.com/larsartmann/*` — warns on every run, needs a human policy pick.
5. **SystemNix system-Go ≥ 1.27.1 deploy** (user-fired, sudo): would obsolete the devshell-invocation workaround entirely.
6. The identity of the **1 skipped spec** — unknown.
7. On-demand tools never run this session: **gitleaks, codespell, markdown-lint**.
8. `nix build .#checks.x86_64-linux.format` (a listed local gate in the gotchas doc) — not run.
9. **GBR_PROPERTY=1** opt-in property test — not re-run after this session's edits (low risk: no touched file overlaps, but unverified).
10. **Push to origin** — branch carries this session's daemon commits; user call.

## d) TOTALLY FUCKED UP

1. **The nolint-deletion misdiagnosis.** I saw `panic --format finding` report `total: 0, suppressed: 0` with the nolints present and concluded 0.6.4 no longer flagged those spots → deleted both `//nolint:branching-flow:panic` directives + the `.golangci.yml` nolintlint exclusions + wrote an AGENTS.md paragraph claiming they were "removed as dead". **Wrong.** The test run then failed: plain-mode `panic` showed _Detections: 2_, stats jumped 24 → 26. The suppressions were load-bearing; `--format finding` just hides suppressed findings _and reports the suppressed counter as 0_. I reverted everything (code, config, docs) and re-verified empirically in both modes. **Cost:** two extra test cycles, one wrong daemon-committed intermediate state in history, and a temporarily-lying AGENTS.md. **Root cause:** verified in exactly ONE output mode, and resolved an ambiguous sentence in the old docs ("reports them suppressed") the convenient way instead of the empirical way.
2. **Evidence-context mismatch in my own new test.** I ran `branching-flow phantom --format finding .` from the module root and saw bare relative filenames (`builders.go`), so I wrote my new spec with exact-name map keys. The test invokes the binary with the **absolute modulePath**, so `Position.File` is an absolute path → immediate failure → fixed with `HaveKey(ContainSubstring(file))`. I had the evidence file on disk and still reasoned from the wrong invocation context.

Neither shipped broken — the suite caught both within the session. But both were _preventable_ by replicating the consumer's exact invocation before concluding.

## e) WHAT WE SHOULD IMPROVE

1. **Suppression-deletion rule:** before deleting ANY suppression, verify across ALL output modes (finding, plain text, stats) AND pin the binary version in the same breath. One mode is not evidence.
2. **Replicate the test's exact invocation** (args, cwd, env) when pinning analyzer output — never a shell approximation from a different directory.
3. **Binary identity discipline:** the BuildFlow binary changed mid-session (25b07c4 → 200ab99 via devshell rebuild) and _changed the rule set under analysis_. Preflight freshness warnings deserve an immediate decision (rebuild or pin), not deferral — and the binary version should be recorded in any report that draws conclusions from findings.
4. **Document after green, not during:** I wrote documentation describing in-flight (wrong) states twice. Corrections to durable docs should land only once the suite confirms them.
5. **Result-cache reflex:** after fixing a finding's root cause, the cached finding replays (agent-config kept claiming 376 lines on a 211-line file). `BUILDFLOW_NO_RESULT_CACHE=1` is the documented bypass — it should be the default reflex for verification runs after fixes, not discovered mid-debug.
6. **Close the loop upstream:** the branching-flow suppressed-findings invisibility (evidence in hand) should become an upstream issue, not just a local gotcha paragraph.
7. **Make pin drift loud early:** add a spec asserting the branching-flow binary version (or major.minor) so binary drift announces itself before the count pins break.

## f) Up to 50 things to get done next (brainstorm — HARVEST should route; most are ROADMAP fuel)

**Upstream / BuildFlow (high leverage):**

1. File branching-flow issue: suppressed findings invisible in `--format finding` (`suppressed: 0` while suppressions active) — evidence collected.
2. File BuildFlow issue: result cache does not invalidate when a flagged file's _content_ changes (agent-config replayed a 376-line claim on a 211-line file).
3. Rebuild + reinstall BuildFlow (`nix build . && nix run .#reinstall`) to clear the freshness advisory; re-run the gate on the fresh binary.
4. Decide fleet lychee policy: GITHUB_TOKEN vs exclude `github.com/larsartmann/*`; apply to this repo's `lychee.toml`.
5. Consider a BuildFlow provider flag/feature: scope go-auto-upgrade away from test files (8 findings, all documented FPs).
6. Ask upstream for a branching-flow `--exclude` path flag — it would dissolve the pin-fragility class entirely.
7. Evaluate making the agent-config line limit configurable per repo (220 hard limit vs repo reality).

**Verification debt (quick wins):**
8. Identify the 1 skipped spec (250 ran / 1 skipped) — fix or document why.
9. Run `GBR_PROPERTY=1 go test` to re-verify the property test after this session's edits.
10. Run gitleaks, codespell, markdown-lint (on-demand tools, unexercised this session).
11. Run `nix build .#checks.x86_64-linux.format` (unrun local gate).
12. Verify govulncheck findings with a govulncheck built on go1.27.2; upgrade the gotchas note from "noise (unverified)" to a verified verdict.
13. Loop `go test` ~10× per the gotchas flake recipe to reconfirm stability after the pin rewrite.
14. Reconcile the gotchas "251/251 specs" claim with the actual 250+1skipped.
15. Reconcile the gotchas doctor note ("34 failed" metric) with the current "9 tools unavailable" health-check phrasing — likely two different metrics; document both.

**Docs & memory:**
16. Add CHANGELOG.md entries: pin re-pin (0.6.4), SECURITY.md section, AGENTS.md restructure, docs/agents-gotchas.md creation.
17. Refresh FEATURES.md content (552h stale — content update, not a timestamp touch).
18. Enumerate + triage the other 2 doc-files-age flagged files.
19. Add a dated "verified 2026-10-11" line to the pin docs (both AGENTS bullet and gotchas) — done for pins, missing for the nolint-verification evidence chain.
20. Link docs/agents-gotchas.md from TODO_LIST or docs index so HARVEST/findability don't depend on AGENTS.md alone.
21. Update the SystemNix govalid rev-mirror check ("keep the rev in sync") — verify it is actually in sync today.

**Test quality:**
22. Add a branching-flow binary-version assertion spec (loud drift detection).
23. Consider collapsing the 8 phantom-invoking specs to share one parsed result per run (8 analyzer shells → 1) — faster suite, fewer drift surfaces.
24. Give `examples/sse` + `listeners/otel` first-class thought in the pin procedure: any new module file shifts the whole-tree counts (fragility is documented; a mitigation is not).
25. Add a regression spec that the panic nolints are _present_ (grep-based guard) so nobody "cleans them up" without hitting the suite first — cheap insurance against a repeat of d1.

**Environment / infra:**
26. Fire the SystemNix Go ≥ 1.27.1 deploy (user, sudo); then delete the devshell-invocation workaround and simplify the gotchas story.
27. Update the "currently 1.27.2" date-stamps in docs when the devshell Go bumps again (fragile by design — consider phrasing without patch numbers).
28. Re-run the full gate after the next `nix flake update` (vulnix policy + check-all).
29. Audit `.buildflow.yml` skips (4 via config) — each still justified?
30. Review the 8 no-op + 59 not-applicable tool lists once with `--verbose` to confirm nothing that should run is silently NA.
31. CI: assert the branching-flow binary version (or absence-skip) explicitly so CI and local never compare different analyzers.

**Repo state / process:**
32. Push origin (daemon commits pending locally).
33. Confirm the daemon committed the final 7-line gotchas addition (pending at session end).
34. Run `docs-health` HARVEST on section (f) → route into TODO_LIST.md / ROADMAP.md properly.
35. Decide: should status reports from BuildFlow-triage sessions keep using `.md` (this one did, per operator instruction) — record as standing preference or one-off.

**Longer-term (ROADMAP fuel):**
36. Re-evaluate GOEXPERIMENT=jsonv2 constraint when json/v2 graduates (standing trigger, documented).
37. Re-check cqrs-lint 3 info findings after the next go-cqrs-lite bump.
38. Re-check vulnix list after next nixpkgs bump; act only on runtime-relevant CVEs (standing policy).
39. Consider upstreaming the "verify suppression in all modes" lesson to crush-config references/lessons.md (cross-project class).
40. Polish-Customs compatibility spot-check on the next release (documented consumer, no replace).

## g) Questions I can NOT figure out myself

1. **Lychee policy:** authenticate (`GITHUB_TOKEN` secret in CI + local env) or exclude the `github.com/larsartmann/*` namespace from link checking? This is a fleet-level policy the BuildFlow warning explicitly leaves to you — whichever you pick, I'll apply it here and it should propagate fleet-wide.
2. **File the two upstream issues?** (a) branching-flow: suppressed findings invisible in `--format finding`; (b) BuildFlow: result cache not invalidating on flagged-file content change. Both have session evidence attached. Yes to both, one, or neither?
3. **SystemNix Go deploy timing:** are you deploying system Go ≥ 1.27.1 soon (I then delete the devshell-invocation workaround and simplify the docs), or is the devshell workaround the long-term reality (I then harden the docs around it as permanent)?

---

_Point-in-time snapshot; will go stale. Section (f) is HARVEST input for TODO_LIST.md / ROADMAP.md. Prior related report: `docs/status/2026-10-10_*_buildflow-verification-triage-and-harvest.md`._
