# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-10-04 (HARVEST of [`docs/status/2026-10-04_07-27_buildflow-toolchain-skew-remediation.md`](docs/status/2026-10-04_07-27_buildflow-toolchain-skew-remediation.md) §e–§g on top of the 2026-09-18 PubLaunch harvest; every report item was re-checked against code — resolved ones were dropped, not carried)

_Source shorthands below: **PubLaunch** = `docs/status/2026-09-18_07-25_repo-public-launch-session.md` §f; **BFSkew** = `docs/status/2026-10-04_07-27_buildflow-toolchain-skew-remediation.md` §f._

---

## Deploy-gated (2026-10-04 toolchain skew)

- [ ] **Fire the SystemNix deploy in a quiet IO window** — restores AMBIENT parity so plain `buildflow --fix --build-mode=full` works outside the devshell. Repo-side verification is DONE: full mode ran 0-failed via `nix develop --command buildflow --build-mode=full` (2026-10-04). MUST be user-fired: the deploy script requires `sudo`, which agent sessions cannot run, and the PSI gate must not be forced (`DEPLOY_FORCE_PRESSURE` — ROADMAP open question 9). (BFSkew §f1, §f2, §f16)

## Toolchain-bump follow-ups (2026-10-04)

- [ ] **Confirm CI green on next push** — setup-go `"1.27"` must resolve ≥ 1.27.1; gosec + `GOEXPERIMENT=jsonv2` must hold across the four-module matrix. (BFSkew §f17)
- [ ] **Fleet sweep: stale `go_1_26` flake pins** — the same pin that broke this repo's hermetic `.#checks…format` is likely in other `~/projects` flakes. (BFSkew §f39)
- [ ] **Fleet sweep: treefmt goimports vs `go 1.27` floors** — formatter wrappers silently carry nixpkgs' default Go and die on toolchain download in hermetic sandboxes. (BFSkew §f40)
- [ ] **go-cqrs-lite upstream: relax snapshot v4.5.1's `go 1.27.1` floor to `go 1.27`** — would let `adapters/cqrslite` use a major.minor directive and silence go-version-auto-configure. (BFSkew §f25)
- [ ] **Write the ecosystem Go-bump runbook** — nixpkgs `go_1_XX` attr → repo flake pins → Home Manager → module floors → CI; it broke four different ways on 2026-10-04. (BFSkew §e6, §f14)
- [ ] **Verify Polish-Customs against the 1.27 floor** — it consumes published `v2.2.0`; confirm its toolchain and `go` directive are ≥ 1.27. (BFSkew §f36)
- [ ] **Fresh-`GOMODCACHE` health check of published `v2.2.0`** — stale-cache recipe in `AGENTS.md`. (BFSkew §f24)
- [ ] **`go.sum` hygiene review** after four in-loop go-mod-update cycles (BFSkew §f48)
- [ ] **devShells.ci parity: add `go-licenses` + `govalid`** (BFSkew §f22)
- [ ] **Decide jscpd handling** — keep the intentional-duplication note vs exclude `.github/workflows` from scanning. (BFSkew §f28)
- [ ] **Review `.buildflow.yml` GOEXPERIMENT duplication** — the "caller wins" info line prints every run. (BFSkew §f29)
- [ ] **Re-measure + re-pin branching-flow pins after the next analyzer upgrade** — current pins: 25 PHANTOM / 31 `stats` (re-pinned 2026-10-04). (BFSkew §f20)

## Repo hygiene (public-launch follow-ups)

- [ ] **Decide the fate of internal narration now public** — `docs/status/` (incl. `archived/`), `docs/planning/`, and ~100 Polish-Customs mentions are publicly indexed since the 2026-09-17 flip. Keep, prune, or relocate: pick one. This is the last undecided CONTRA item and gets harder to reverse the longer it is indexed. (PubLaunch §b2, §g1)
- [ ] **Review `.github/CODEOWNERS` + issue/PR templates for public inbound** (PubLaunch §c5, §f19)
- [ ] **Confirm the `.gitignore` additions from commit `2e06602` are sane** (PubLaunch §f40)
- [ ] **Decide whether `docs/status/*.md` still warrant the `linguist-documentation` attribute** now that the repo is public (PubLaunch §f41)

## Consumer & publishing follow-ups

- [ ] **pkg.go.dev `@v2.2.0` "not the latest version" banner** — investigated 2026-09-18: `proxy.golang.org/.../@latest` returns `v2.2.0` (hash `2e06602`), so this is a pkg.go.dev-side index/staleness artifact, not a module or proxy problem. Not locally actionable; re-check after pkg.go.dev's next refresh, else report upstream. (PubLaunch §b3, §f12)

## CI, badges & release mechanics

- [ ] **Run CI on tag pushes** — verified 2026-09-18: `ci.yml` triggers only on `push`/`pull_request` to branches, so the `v2.2.0` tag ran no CI. Add `tags: ["v*"]`. (v2.2.0 session report)
- [ ] **Add a coverage badge** — 97.1% re-measured 2026-09-18 via `nix develop --command go test -cover ./...`. (PubLaunch §f22)
- [ ] **Fix GitHub account billing** — no longer relevant to this repo (public Actions are free), but the ecosystem's private repos still need it. (PubLaunch §c12, §f37)

## Standing (quarterly)

- [ ] **Quarterly docs-health rerun** (next: ~2026-12): re-verify FEATURES claims against code, re-run `art-dupl -t 15`, and re-check the json/v2 graduation status (ADR-0004).
- [ ] **Re-pin the branching-flow analyzer counts** — verified 2026-10-04: pins are 25 PHANTOM / 31 `stats` per `AGENTS.md` (fragile: any new file moves the totals; re-measure and re-pin with a comment). (PubLaunch §f50)
- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers (documented in `AGENTS.md`). **Re-verified 2026-10-04: still required.** Track the Go release that graduates `json/v2` from experimental and drop the constraint then.

---

_Baseline established 2026-09-18 from the v2.2.0 release + repo-public-launch sessions; harvested 2026-10-04 from the toolchain-skew session. Resolved 2026-10-04 items (nix-checker/cqrs-lint/vulnix triage, doctor triage, license-check cache-free verification, lychee + GOTOOLCHAIN + README-version docs, CHANGELOG entry, property test + benchmark smoke on 1.27.1) are recorded in [`CHANGELOG.md`](CHANGELOG.md), `AGENTS.md`, or dropped as verified._
