# TODO List — businessrules

> Actionable, bounded work for the next 2-4 weeks. Open items only.
>
> Completed work lives in [`CHANGELOG.md`](CHANGELOG.md), not here. Long-term ideas live in [`ROADMAP.md`](ROADMAP.md).

**Last verified:** 2026-09-18 (HARVEST of [`docs/status/2026-09-18_07-25_repo-public-launch-session.md`](docs/status/2026-09-18_07-25_repo-public-launch-session.md) §f, plus the two v2.2.0 session reports; every item re-checked against `master` at `7f3d9ed`)

_Source shorthand below: **PubLaunch** = `docs/status/2026-09-18_07-25_repo-public-launch-session.md` §f (its item numbers cited as `PubLaunch §fN`)._

---

## Repo hygiene (public-launch follow-ups)

- [ ] **Decide the fate of internal narration now public** — `docs/status/` (incl. `archived/`), `docs/planning/`, and ~100 Polish-Customs mentions are publicly indexed since the 2026-09-17 flip. Keep, prune, or relocate: pick one. This is the last undecided CONTRA item and gets harder to reverse the longer it is indexed. (PubLaunch §b2, §g1)
- [ ] **Protect `master`** — verified 2026-09-18: `gh api repos/LarsArtmann/go-business-rules/branches/master/protection` → `Branch not protected` (HTTP 404). Require green CI and block force-pushes; the repo is now public and mergeable by rule. (PubLaunch §c1)
- [ ] **Enable GitHub private vulnerability reporting** — verified 2026-09-18: `private-vulnerability-reporting` → `{"enabled":false}`, while [`SECURITY.md`](SECURITY.md) now directs reporters to the Security tab. Enable it (or reword SECURITY.md). (PubLaunch §c, `SECURITY.md`)
- [ ] **Fill in the repo "Website" field** — verified 2026-09-18: `homepageUrl` is empty. Set it to `https://pkg.go.dev/github.com/LarsArtmann/go-business-rules/v2`. Optional: add topic `businessrules`. (PubLaunch §f10, §f45)
- [ ] **Create the missing `v2.1.0` GitHub Release page** — verified 2026-09-18: `gh release list` shows only `v2.2.0` (Latest) and `v2.0.0`; `v2.1.0` was tagged but never given a Release page. (PubLaunch §f11)
- [ ] **Enable or explicitly decline GitHub Discussions** — verified 2026-09-18: `hasDiscussionsEnabled: false`. (PubLaunch §f18)
- [ ] **Add a `.mailmap`** — verified 2026-09-18: history contains `Unknown Author <unknown@example.com>` alongside `Lars Artmann <git@lars.software>`; both are now publicly visible. (PubLaunch §c11, §f26)
- [ ] **Review `.github/CODEOWNERS` + issue/PR templates for public inbound** (PubLaunch §c5, §f19)
- [ ] **Confirm the `.gitignore` additions from commit `2e06602` are sane** (PubLaunch §f40)

## Consumer & publishing follow-ups

- [ ] **Bump Polish-Customs to `v2.2.0` and run its full suite** — verified 2026-09-18: `/home/lars/projects/Polish-Customs/go.mod:11` still pins `v2.1.0`, and its stale comment (lines 9-10) claims a `replace` directive that no longer exists. Timing is a user decision. (PubLaunch §b4, §g3)
- [ ] **Investigate the pkg.go.dev `@v2.2.0` "not the latest version" banner** — the versioned page renders fully but shows a stale "Go to latest" hint. (PubLaunch §b3, §f12)
- [ ] **Verify the README install snippet in a fresh container** (PubLaunch §f20)

## Documentation

- [ ] **Add a Go / `GOEXPERIMENT` compatibility table to README** — the hard `GOEXPERIMENT=jsonv2` requirement deserves a scannable table, not just prose. (PubLaunch §f21)
- [ ] **Link the `v2.2.0` release notes from README** (PubLaunch §f44)
- [ ] **ANNOTATE archived status reports claiming "repo private / pkg.go.dev unachievable"** with the 2026-09-17 flip date (docs-health ANNOTATE mode). (PubLaunch §c10, §f17)
- [ ] **Decide whether `docs/status/*.md` still warrant the `linguist-documentation` attribute** now that the repo is public (PubLaunch §f41)

## CI, badges & release mechanics

- [ ] **Run CI on tag pushes** — verified 2026-09-18: `ci.yml` triggers only on `push`/`pull_request` to branches, so the `v2.2.0` tag ran no CI. Add `tags: ["v*"]`. (v2.2.0 session report)
- [ ] **Add a coverage badge** — 97.1% re-measured 2026-09-18 via `nix develop --command go test -cover ./...`. (PubLaunch §f22)
- [ ] **Fix GitHub account billing** — no longer relevant to this repo (public Actions are free), but the ecosystem's private repos still need it. (PubLaunch §c12, §f37)

## Standing (quarterly)

- [ ] **Quarterly docs-health rerun** (next: ~2026-12): re-verify FEATURES claims against code, re-run `art-dupl -t 15`, and re-check the json/v2 graduation status (ADR-0004).
- [ ] **Re-pin the branching-flow analyzer counts** — verified 2026-09-18: pins are 25 PHANTOM / 29 `stats` per `AGENTS.md` (fragile: any new file moves the totals; re-measure and re-pin with a comment). (PubLaunch §f50)
- [ ] **Re-evaluate `encoding/json/v2`** — the `GOEXPERIMENT=jsonv2` requirement is a hard breaking change for downstream consumers (documented in `AGENTS.md`). **Re-verified 2026-09-18: still required.** Track the Go release that graduates `json/v2` from experimental and drop the constraint then.

---

_Baseline established 2026-09-18 from the v2.2.0 release + repo-public-launch sessions. The v2.2.0 release itself (CHANGELOG cut, `doc.go`+`suite_test.go` version bump, gosec job rewrite, branching-flow `Skip` fix, tag push, GitHub Release, proxy/pkg.go.dev verification) is complete and recorded in [`CHANGELOG.md`](CHANGELOG.md)._
