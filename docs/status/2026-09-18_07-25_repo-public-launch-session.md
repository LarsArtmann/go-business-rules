# Status Report — Repo Went Public: Launch, CI Repair, Publishing Verification

**Date:** 2026-09-18 07:25 CEST
**Session window:** 2026-09-17 ~19:20 CEST → 2026-09-18 07:25 CEST (interleaved with an independent overnight session, attributed separately below)
**Scope:** PRO/CONTRA decision executed — `go-business-rules` flipped from PRIVATE to PUBLIC, documentation brought to public state, CI repaired to green, and the full public-consumption chain (proxy → sumdb → pkg.go.dev) verified end-to-end.
**Format note:** written as Markdown per explicit user request. The `status-report` skill's canonical output is a styled HTML dashboard in `docs/status/<YYYY-MM-DD_HH-MM_name>.html`; this is a deliberate one-off override, same as the v2.2.0 session's report.

---

## Session Verdict

The public flip is **complete, verified, and healthy**: master CI is green (13/13), `go get` works with a default environment through `proxy.golang.org`, sumdb verification passes for fresh consumers, and pkg.go.dev serves full documentation for `v2.2.0`. The publish path that motivated the flip (PRO arguments #1–#3 from the analysis) is now realized, not just predicted.

The session's own quality debt is **methodological, not material**: two invalid diagnostic conclusions (cache-hit masquerade, pipeline-masked diff) were self-caught and corrected before any destructive action, but they cost roughly six avoidable round trips and briefly produced a wrong "proxy and direct agree" finding. One process risk materialized: I raced a concurrently active session on `ci.yml` and hit a non-fast-forward push rejection.

**Attribution boundary:** the `v2.2.0` release, the gosec job rewrite (Docker action → fail-closed `go install gosec@v2.29.0`), and the branching-flow `Fail` → `Skip` fix were done by an **overnight parallel session** — this session diagnosed those same two CI failures independently at 17:58Z before discovering the parallel fixes, rebased onto them, and fixed the one follow-up breakage they caused (`adapters/cqrslite` go.sum).

---

## a) FULLY DONE

| # | Item | Evidence |
| --- | --- | --- |
| 1 | Repo PUBLIC with description + topics (`go`, `golang`, `validation`, `business-rules`, `severity`, `go-library`) | `gh repo view` → `"visibility":"PUBLIC"` |
| 2 | README public state: private-`GOPRIVATE` caveat replaced with pkg.go.dev link + public-proxy statement + `GOEXPERIMENT=jsonv2` warning | README lines 9–20; rendered on pkg.go.dev |
| 3 | AGENTS.md current: public-flip bullet, CI reality section re-dated, GOPRIVATE-vars-are-no-ops note, sumdb stale-cache gotcha with trash recipe | AGENTS.md §Private modules, §CI & Publishing Reality |
| 4 | `ci.yml` header comment de-staled (billing note replaced with public-repo note) | commit `81e70fd` (body later superseded by parallel session's job rewrite) |
| 5 | **Master CI green, 13/13** | run `35310049647`, conclusion `success` |
| 6 | `adapters/cqrslite` ginkgo misalignment fixed (go.mod → v2.32.2 + `go mod tidy`); local build+vet+test green before push | commit `b039519` |
| 7 | Public consumption verified: clean-cache `go get ...@v2.1.0` and `@v2.2.0` through `proxy.golang.org` with `GOSUMDB=sum.golang.org` succeed; proxy serves both versions | scratch-module downloads, Sum `qAZFk…` verified |
| 8 | **pkg.go.dev indexed**: full API index, MIT, rendered README at `pkg.go.dev/github.com/LarsArtmann/go-business-rules/v2` (v2.2.0) | fetched page |
| 9 | Secrets scan: working tree + full 237-commit history clean (token/AWS/PEM/password patterns) | `rg` + `git grep` over all revs |
| 10 | Dependency-graph audit: every LarsArtmann dep public (go-finding, go-error-family, go-sse, go-branded-id, go-cqrs-lite); `starfederation/datastar-go` third-party public | `gh repo view` per repo |
| 11 | v2.1.0 health proven: the checksum-mismatch SECURITY ERROR was a **local stale-cache artifact**, not a poisoned release — no re-tag needed, sumdb record intact | sumdb lookup `h1:qAZF…` matches fresh proxy zip |
| 12 | Local doc commits from the parallel session (status report, CHANGELOG tweak) rebased onto advanced origin and pushed | `git rebase origin/master` → pushed |

## b) PARTIALLY DONE

| # | Item | Gap |
| --- | --- | --- |
| 1 | Security-gate trust | The new fail-closed gosec job is *green* in run `35305107107` per the parallel session's report, but **I never read its logs myself** to assert `Stats.files > 0`. The go-paperless lesson says a scanner must prove it scanned; the green badge alone doesn't. |
| 2 | Internal narration (CONTRA #2 from the analysis) | ~100 Polish-Customs mentions + session reports in `docs/status/`, `docs/planning/`, AGENTS.md are now public. **No user decision recorded**: keep, prune, or relocate. |
| 3 | pkg.go.dev oddity | The `@v2.2.0` versioned page renders fully but says "not in the latest version of its module / Go to latest" — cosmetic inconsistency I chose not to chase. |
| 4 | Polish-Customs consumer state | AGENTS.md still says PC consumes `v2.1.0` (line ~232). v2.2.0 exists; PC's bump is neither done nor explicitly deferred by anyone. |
| 5 | AGENTS.md consumer-section freshness | Mixed: release facts updated to v2.2.0, consumer paragraph still references v2.1.0 verification prose. |

## c) NOT STARTED

| # | Item |
| --- | --- |
| 1 | Branch protection on `master` (require green CI before merge) — repo is now public and mergeable by rule |
| 2 | GitHub social preview image |
| 3 | Repo "Website" field → pkg.go.dev URL |
| 4 | GitHub Discussions: enable or explicitly decline |
| 5 | CODEOWNERS + issue/PR templates review for a public audience |
| 6 | Obsolete Dependabot PRs: the `adapters/cqrslite` ginkgo PR branch was force-updated overnight and is now redundant with master fix `b039519` — close or merge deliberately |
| 7 | v2.1.0 GitHub Release page existence check (parallel session created the v2.2.0 one) |
| 8 | ROADMAP entry: watch `encoding/json/v2` graduation (downstream `GOEXPERIMENT` constraint disappears then) |
| 9 | `docs-health` HARVEST of this report's section (f) into TODO_LIST/ROADMAP |
| 10 | ANNOTATE archived status reports claiming "repo private / pkg.go.dev unachievable" with the 2026-09-17 flip |
| 11 | mailmap for the 5 `unknown@example.com` commits now publicly visible |
| 12 | GitHub account billing fix (still relevant for the ecosystem's *private* repos; CI here no longer depends on it) |

## d) TOTALLY FUCKED UP

Nothing destructive shipped — no poisoned tags, no force-pushes, no lost work, history linear. But three diagnostics produced **wrong intermediate conclusions**, and one process failure is worth naming:

1. **Cache-hit masquerade (worst).** My "proxy and direct agree on `Lm0US…`" claim was false verification: both runs were cache hits on the same Sep-14 zip in `GOMODCACHE`. I presented it as a finding before catching it with a fresh `GOMODCACHE`. A wrong conclusion at that point would have led to "burn v2.1.0 and re-tag" — an unnecessary, permanently recorded release action.
2. **Repeated a documented anti-pattern.** The zip-diff step piped `diff` through `head` and ran `unzip … 2>/dev/null`, masking failures — the exact pipeline-masking lesson already in the global AGENTS.md from 2026-09-11/13. I then had to redo the diff cleanly. Knowing a lesson didn't stop me from repeating it under time pressure.
3. **Unverified env semantics.** `GONOSUMDB=` (empty) silently fell back to Home Manager's OS-level value, so my "sumdb off" runs had sumdb on. Should have run `go env GONOSUMDB GOPRIVATE GOPROXY` first.
4. **Raced a live concurrent session.** I pushed a `ci.yml` edit while the overnight session was rewriting the same file; result: rejected push, rebase, and two CI runs on quickly-superseded commits. I had evidence concurrency existed (daemon commits timestamped after my session start) and didn't act on it.
5. **CI-trigger surprise.** Two pushes produced no CI runs before I read `paths-ignore: "**/*.md"` in the workflow — predictable if I'd read triggers before pushing.

## e) WHAT WE SHOULD IMPROVE

1. **Module-health claims require fresh `GOMODCACHE`.** Cache hits can impersonate "the proxy says". Make the scratch-dir recipe include `GOMODCACHE=$(mktemp -d)` by default.
2. **`go env` before env-dependent commands.** Empty-string overrides of `GOPRIVATE`/`GONOSUMDB`/`GONOPROXY` do not override OS env. Always print the effective values first.
3. **No masked pipelines during diagnosis.** `unzip`/`diff`/`go` failures were hidden by `2>/dev/null` and `| head` exit-masking. Errors first, filtering second.
4. **Check for concurrent sessions before editing shared files.** `git log --format='%ci'` skew and fresh daemon commits are the telltale; coordinate or wait.
5. **Read workflow triggers before predicting CI behavior.** `paths-ignore` cost two dead pushes.
6. **Trust scanner badges only with instrument assertions.** Read the `Stats.files` line from the gosec report in the logs, not just the green check.
7. **Daemon vs explicit commits.** Commit-per-task was mostly followed, but the daemon still interleaved heuristic commits mid-task (cdabe29, 2c238e1, 2e06602). For release-critical work, consider suspending the daemon.

## f) NEXT 50 (brainstorm, impact-sorted; extra items are ROADMAP fuel, not commitments)

| # | Task | Impact |
| --- | --- | --- |
| 1 | Verify gosec logs in latest green run assert `Stats.files > 0` | High |
| 2 | Close/merge the now-redundant `adapters/cqrslite` ginkgo Dependabot PR | High |
| 3 | Bump Polish-Customs to `v2.2.0`, run its full suite | High |
| 4 | Update AGENTS.md §Polish-Customs (still v2.1.0 prose) | High |
| 5 | Decide fate of public internal narration (docs/status, docs/planning) | High |
| 6 | Branch protection on master (require green CI) | High |
| 7 | Harvest section (f) into TODO_LIST/ROADMAP (docs-health HARVEST) | High |
| 8 | Verify CI badge in README renders green publicly | Med |
| 9 | Social preview image | Med |
| 10 | Repo "Website" field → pkg.go.dev | Med |
| 11 | Check/complete v2.1.0 GitHub Release page | Med |
| 12 | Resolve pkg.go.dev "not latest" oddity on @v2.2.0 page | Med |
| 13 | ROADMAP: watch encoding/json/v2 graduation | Med |
| 14 | CONTRIBUTING.md: document GOEXPERIMENT=jsonv2 for contributors | Med |
| 15 | CONTRIBUTING.md: tag discipline — sumdb pins public tags forever, no force-push | Med |
| 16 | SECURITY.md: supported-versions table (v2.2.x) | Med |
| 17 | ANNOTATE archived reports claiming "pkg.go.dev unachievable" with flip date | Med |
| 18 | Enable/decline GitHub Discussions | Med |
| 19 | Review CODEOWNERS + templates for public inbound | Med |
| 20 | Fresh-container verbatim test of README install instructions | Med |
| 21 | Add Go/GOEXPERIMENT compatibility table to README | Med |
| 22 | Coverage badge (97.1% claimed) | Med |
| 23 | Go Report Card badge | Low |
| 24 | CHANGELOG v2.2.0 entry completeness check | Low |
| 25 | FEATURES.md: add public-launch facts (pkg.go.dev indexed) | Low |
| 26 | mailmap for `unknown@example.com` commits | Low |
| 27 | Glance `.config/metadata.yaml` + `library-policy.yaml` for public-safety | Low |
| 28 | Accept/document internal tooling configs now public (git-town.toml, .buildflow.yml) | Low |
| 29 | CI: ubuntu-26 runner migration Oct 19 — accept or pin | Low |
| 30 | CI: cache gosec install step | Low |
| 31 | CI: consider gosec on nested modules (currently root-only) | Low |
| 32 | CI: optional short-duration fuzz job for the 7 fuzz targets | Low |
| 33 | Dependabot: consistent grouping across all four modules | Low |
| 34 | Release workflow: auto-create GitHub Release on tag | Low |
| 35 | Benchmark results doc (README claims 189/468/476 ns) | Low |
| 36 | ROADMAP: evaluate json v1 build-tag fallback to drop GOEXPERIMENT burden | Low |
| 37 | Fix GitHub account billing (ecosystem-private repos still need it) | Low |
| 38 | Traffic insights baseline after first public week | Low |
| 39 | Audit `examples/sse` for anything unintended-public | Low |
| 40 | Confirm `.gitignore` additions from overnight commit `2e06602` are sane | Low |
| 41 | Review whether `docs/status/*.md` `linguist-documentation` attr still right publicly | Low |
| 42 | Consider pinning gosec version upgrade cadence (currently @v2.29.0) | Low |
| 43 | dprint/gofmt check in CI (buildflow covers locally; CI parity optional) | Low |
| 44 | README: link the v2.2.0 release notes | Low |
| 45 | Add topic `businessrules` (product name) to repo topics | Low |
| 46 | Sitemap of nested modules in README (cqrslite/otel/sse usage snippets) | Low |
| 47 | Publish example dashboards/screenshots from the sse example | Low |
| 48 | Review pkg.go.dev "Imported by: 0" — seed awareness (blog/social) | Low |
| 49 | Decide whether archived `docs/planning` FINDING-SDK proposal still reflects intent | Low |
| 50 | Quarterly re-pin of branching-flow analyzer counts (documented fragility) | Low |

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **Internal narration:** should `docs/status/` (incl. archived), `docs/planning/`, and the ~100 Polish-Customs mentions stay public now that the repo is? This is the one CONTRA item left un-decided, and it gets harder to reverse the longer it's indexed.
2. **Concurrency:** the overnight session (v2.2.0, CI fixes) and the daemon are still active in this repo (commits at 07:08 today; a Dependabot branch force-updated). Should I stand down from this repo while that workstream runs, or is parallel work acceptable?
3. **Polish-Customs:** bump it to `v2.2.0` now, or leave it on `v2.1.0` until its next consumer release?

---

*Point-in-time snapshot. Do not edit archived reports; annotate via docs-health ANNOTATE mode.*
