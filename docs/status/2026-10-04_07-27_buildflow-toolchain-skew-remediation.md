# Status Report: BuildFlow Toolchain-Skew Remediation (go-business-rules)

- **Written:** 2026-10-04 07:27 CEST
- **Session scope:** triage + fix of `buildflow --fix --build-mode=full` exit 69
  (10 failed steps, 4 tools with remaining findings, 5 detect-only reporters)
- **Repo:** `github.com/LarsArtmann/go-business-rules` (+ one cross-repo fix in
  `SystemNix`)
- **Final state at write time:** buildflow full run down from **10 failed steps →
  2 failed steps** (both `govalid-generate`, blocked by ONE external cause); all
  repo-local fixes verified green; working trees clean in both repos (auto-commit
  daemon).

---

## Root cause chain (what actually went wrong)

Everything cascaded from **Go toolchain fragmentation** across the ecosystem:

| Layer                                    | Go version | Problem                                                         |
| ---------------------------------------- | ---------- | --------------------------------------------------------------- |
| nixpkgs default `go`                     | 1.26.8     | older than this repo's module floors                            |
| flake.nix pin (before fix)               | `go_1_26`  | 1.26.8 in devshells + check-all                                 |
| treefmt goimports wrapper's internal go  | 1.26.8     | tried to DOWNLOAD toolchain go1.27.0 → died in hermetic sandbox |
| Home Manager system go (`/etc/profiles`) | 1.27.0     | what BuildFlow's ambient env actually runs                      |
| `go-cqrs-lite/snapshot/v4@v4.5.1` (dep)  | requires   | **go >= 1.27.1** — newer than every installed toolchain         |
| nixpkgs `go_1_27` (locked rev)           | **1.27.1** | available all along; the flake just pinned the wrong attr       |

With `GOTOOLCHAIN=local` (BuildFlow's policy), nothing could satisfy the dep
floor; without it, toolchain downloads fail in the nix sandbox.

---

## a) FULLY DONE (verified green)

1. **flake.nix: `go_1_26` → `go_1_27` (1.27.1)** in all three places
   (`check-all` runtimeInputs, `devShells.default`, `devShells.ci`). Verified:
   `nix develop --command go version` → go1.27.1, native, no download.
2. **treefmt goimports hermetic fix** — `go127Gotools` override
   (`gotools.override { buildGoModule = buildGo127Module; go = go_1_27; }`).
   Root cause proven by reading the wrapper: nixpkgs' gotools appends its
   callPackage `go` (1.26.8) to PATH; goimports shells out to `go`; the old go
   tried a toolchain download inside the sandbox. Verified:
   `nix build .#checks.x86_64-linux.format` GREEN and `nix flake check` →
   "all checks passed".
3. **exhaustruct_v5 nolint fix** — the enabled linter is the v5 variant, so the
   four `//nolint:exhaustruct` directives in `validation_result.go` were
   silently ignored. Renamed to `//nolint:exhaustruct_v5` (4 sites). Verified:
   `buildflow -s "golangci-lint [root]"` → **0 issues**.
4. **CI go-version bump** — 4× `go-version: "1.26"` → `"1.27"` in
   `.github/workflows/ci.yml` (matches module floor; setup-go picks latest
   1.27.x ≥ the 1.27.1 dep floor).
5. **branching-flow stats re-pin 29 → 31** in `bdd_branching_flow_test.go` —
   branching-flow 0.2.0's `stats` now COUNTS the two nolint-suppressed panic
   findings (Stream result send + DivisibleBy zero-guard) while
   `panic --format finding` still reports them suppressed. Pure analyzer drift:
   PHANTOM pins unchanged (25 = 7 critical / 6 error / 11 info / 1 warning).
   Verified: root `go test ./...` → ok, 251/251 specs.
6. **lychee empty-URL fix** — `OneOf[T]()` in an
   archived planning doc parses as an empty markdown link. Created
   **`lychee.toml`** excluding `docs/planning/archived` +
   `docs/status/archived` (archived files are never edited per repo policy).
   Verified: lychee → 0 errors. (Gotcha: `.lychee.toml` is NOT auto-discovered;
   the plain name is.)
7. **license-check green** — `go-licenses` added to `devShells.default`; the
   step previously ran via `nix run nixpkgs#go-licenses` WITHOUT the project's
   Go env (wrong Go → "Package context does not have module info" spam). In the
   devshell it exits 0; the step went ✔.
8. **flake meta completed** — `homepage`/`maintainers`/`platforms` on the
   `check-all` package AND the new `govalid` package; `meta.description` on the
   `check-all` app (nix-flake-check warning cleared). `flake-meta-checker`
   findings cleared.
9. **AGENTS.md within line budget** — 390 → **376 lines** (max 377) while
   ADDING: Go toolchain 2026-10-04 section, the go-directive rule (below),
   branching-flow 31-pin rationale, the `exhaustruct_v5` nolint gotcha, the
   expanded go-auto-upgrade false-positive list, the lychee exclusion note, the
   intentional ci.yml duplication note; removed a dead trailing `#` heading;
   compressed the GOPRIVATE history, CI history, and module-version paragraphs.
10. **Go directive rule established empirically** — `go mod tidy` ENFORCES
    main-module go ≥ every dependency floor: adapters/cqrslite must stay
    `go 1.27.1` while go-cqrs-lite/snapshot v4.5.1 declares that floor. The
    go-version-auto-configure "patch component" warning on that line is a
    Go-rule necessity, not a preference. Root/examples/listeners stay `go 1.27`.
    Documented in AGENTS.md.
11. **SystemNix cross-repo fix (root cause of the blocked deploy):**
    `pre-deploy-check` aborted EVERY deploy because the app derivation didn't
    stage `scripts/lib/vendor-freshness.sh` that the script sources (same bug
    class as the documented 2026-09-02 metrics-gate incident). Added the
    staging line; audited ALL scripts' `source`d libs vs staged libs (one false
    alarm: `zz-smart-nix.sh` is a conditional HOME check, not a lib). Verified:
    pre-deploy-check → **75 passed, 8 warnings, 0 failed**.
12. **Everything committed** — auto-commit daemon picked up all changes in both
    repos; `git status` clean in both.

## b) PARTIALLY DONE

1. **govalid-generate (the last failing gate step)** — `govalid` is now built
   from source in the devshell (mirrors `SystemNix/pkgs/govalid.nix`, rev pinned
   with a keep-in-sync comment). Status is MODE-DEPENDENT:
   - fast mode: ✔ (9.0 s real run)
   - full mode: ✗ — the fan-out run for `adapters/cqrslite` executes with the
     AMBIENT environment where `go` = 1.27.0 (Home Manager) < dep floor 1.27.1.
     The mechanism (which env each BuildFlow mode uses for tool resolution) is
     only partially diagnosed — see Improvements.
2. **SystemNix deploy (the actual fix for #1)** — prepared and validated
   (`nix run .#pre-deploy-check` → 0 failed), but the deploy's own
   **IO-pressure gate aborted it**: PSI some avg10 40–80% for 40+ minutes,
   caused by OTHER concurrent agent sessions (cmdguard `go test`, a 13 GB RSS
   `nix flake check`, timesheets lint, another `crush -y`). The gate cites
   kernel-freeze precursors (2026-08-22 crash class); I did NOT use
   `DEPLOY_FORCE_PRESSURE=1`. Command ready:
   `cd ~/projects/SystemNix && nix run .#deploy`.
3. **Final buildflow verdict** — full run now fails **2 steps** (was 10):
   govalid-generate + govalid-generate [adapters/cqrslite], both blocked by the
   single external cause above. Remaining findings lists: go-auto-upgrade (8,
   documented false positives), go-version-auto-configure (1, Go-rule
   necessity), nix-checker (4, **not yet investigated**).
4. **Run-to-run inconsistency** — an earlier full run (06:52) recorded ZERO
   failures; the 07:08 full run failed govalid again. Fast/single-step/full
   modes evidently resolve tool environments differently; not fully explained.

## c) NOT STARTED

1. SystemNix **deploy execution** (blocked by pressure gate; everything prepared).
2. **BuildFlow binary rebuild** — stale (13d32f7 vs f88bdb5; advisory
   preflight warning every run). Deferred because that repo has a concurrent
   session.
3. **nix-checker: 4 findings** — never opened.
4. **cqrs-lint: 3 findings** (fast mode) — never opened.
5. **vulnix: 28 CVE findings** — nixpkgs toolchain CVEs (binutils, cargo,
   bison, async); no policy decision made (ignore / threshold / exclude).
6. **9 unavailable tools** health-check advisory (interrogate etc.) — never
   ran `buildflow doctor` to enumerate/resolve.
7. **docs-health HARVEST** of this report's next-task list into TODO_LIST.md /
   ROADMAP.md.
8. aarch64 flake coverage (`nix flake check --all-systems` currently omitted).
9. Fresh-GOMODCACHE module-health verification of v2.2.0 (AGENTS recipe;
   periodic hygiene).
10. GBR_PROPERTY=1 opt-in property test run (not exercised this session).

## d) TOTALLY FUCKED UP (honest near-misses — nothing corrupted, all caught)

1. **False "FORMAT CHECK GREEN" claim** — I wrote
   `nix build … | tail -3 && echo GREEN`; the `&&` tests `tail`'s exit status,
   not nix's (pipeline masks failure). Caught it one command later by rerunning
   with `echo $?` — it was RED (nixfmt style issue). Lesson: `set -o pipefail`
   or explicit `$?` after every piped verification.
2. **Premature success conclusion on govalid** — declared it fixed based on a
   fast-mode green + a `--failed-only` "no failures recorded" (which referred to
   an earlier run). A fresh FULL run proved the fan-out still fails. Conclusions
   about a multi-mode pipeline must be verified in the SAME mode that gated.
3. **Incomplete experiment verification** — the `go 1.27` directive experiment
   passed build+vet, so I called it good; `go test`/check-all then failed with
   "updates to go.mod needed" (tidy consistency) and I had to revert to
   `go 1.27.1`. An experiment must run the full gate, not a prefix of it.
4. **Wrong lychee config filename** — wrote `.lychee.toml` first (not
   auto-discovered), verified red, then renamed to `lychee.toml`. One avoidable
   cycle.
5. **Incomplete nix override** — first `gotools.override` patched only
   `buildGoModule`, not the `go` wrapper arg; cost one full nix build cycle to
   discover via the rebuilt wrapper. Should have read the wrapper script BEFORE
   writing the override.
6. **AGENTS.md budget miss** — trimmed to 372, then appended two more sections
   → 383 (> 377) → second trim pass to 376. Should have budgeted the additions
   in the first edit.
7. **Self-inflicted IO storm contribution** — ran `check-all` and the SystemNix
   deploy concurrently; both are heavy IO. The deploy then aborted on the
   pressure gate that my own parallel load helped trip. Heavy IO work on this
   workstation must be serialized.

## e) WHAT WE SHOULD IMPROVE

1. **BuildFlow tool→environment resolution is mode-inconsistent** — the same
   step (govalid-generate) is green in fast mode and red in full mode;
   license-check only worked once its binary was IN the devshell; single-step
   `-s` runs behaved differently from pipeline runs. The chosen env per tool
   should be uniform across modes and visible in `-v` output
   (e.g. `exec: [devshell] govalid ./...`).
2. **go-version-auto-configure vs `go mod tidy` contradiction** — the analyzer
   demands major.minor-only go directives, but `go mod tidy` FORCE-feeds patch
   floors into the main go.mod when any dependency declares one. The analyzer
   should accept a patch floor it can prove comes from a dependency (or suggest
   the real remedy).
3. **go-auto-upgrade re-flags documented false positives every run** — the
   samber/lo suggestions (8 findings) are already rejected in AGENTS.md policy.
   Projects need a first-class suppression channel for analyzer findings
   (like .golangci exclusions), or the analyzer should read the project's
   documented policy.
4. **treefmt-nix formatter Go coupling** — formatter wrappers silently carry
   whatever Go nixpkgs built them with; a module-floor bump breaks the hermetic
   check with an opaque download error. treefmt-nix could expose the formatter
   Go version or nixpkgs could keep wrapper-go ≥ default-go floors.
5. **SystemNix pre-deploy-check staging is manual** — the same staging bug
   (script sources a lib the derivation doesn't ship) has now happened twice.
   Generate the staging list from the script's `source` lines at build time.
6. **Fleet Go-bump runbook missing** — the update order (nixpkgs go_1_XX attr →
   repo flake pins → Home Manager → module floors → CI) is implicit tribal
   knowledge; it broke 4 different ways today. Deserves a gotcha in AGENTS.md
   / crush-config lessons.
7. **My own verification discipline** — pipefail after every piped gate; full
   (not partial) gate for experiments; conclusions re-proven in the gating mode.

## f) TOP 50 THINGS WE SHOULD GET DONE NEXT

_Brainstorm list — most items beyond #10 are ROADMAP fuel; route through
docs-health HARVEST before committing them to TODO_LIST.md._

| #  | Task                                                                                                                           | Impact |
| -- | ------------------------------------------------------------------------------------------------------------------------------ | ------ |
| 1  | Run SystemNix `nix run .#deploy` in a quiet window (unblocks full-mode govalid fleet-wide)                                     | HIGH   |
| 2  | Re-run full `buildflow --fix --build-mode=full` after deploy → expect 0 failed steps                                           | HIGH   |
| 3  | Investigate nix-checker's 4 remaining findings                                                                                 | MED    |
| 4  | Investigate cqrs-lint's 3 findings                                                                                             | MED    |
| 5  | Diagnose BuildFlow's per-mode tool-env resolution (fast ✔ vs full ✗ for govalid); file upstream                                | HIGH   |
| 6  | Rebuild stale BuildFlow binary (`cd ~/projects/BuildFlow && nix build . && nix run .#reinstall`) when that repo is quiet       | MED    |
| 7  | Decide vulnix policy (thresholds/exclusions for nixpkgs toolchain CVEs)                                                        | MED    |
| 8  | `buildflow doctor` — enumerate + resolve the 9 unavailable tools (interrogate …)                                               | LOW    |
| 9  | Verify go-fix repair's diff introduced nothing unwanted (`git log -p`)                                                         | LOW    |
| 10 | HARVEST this list into TODO_LIST.md / ROADMAP.md (docs-health)                                                                 | MED    |
| 11 | Upstream: go-version-auto-configure should tolerate dependency-forced patch floors                                             | MED    |
| 12 | Upstream: go-auto-upgrade needs project-level finding suppressions                                                             | MED    |
| 13 | SystemNix: generate pre-deploy-check lib staging from the script's source-lines                                                | MED    |
| 14 | Write the ecosystem Go-bump runbook (nixpkgs attr → flake pins → HM → floors → CI)                                             | MED    |
| 15 | Confirm license-check full-mode green is real, not a result-cache artifact (BUILDFLOW_NO_RESULT_CACHE=1)                       | LOW    |
| 16 | Re-verify test-race step explicitly in a full run after the deploy                                                             | LOW    |
| 17 | Confirm CI workflow goes green on next push (setup-go "1.27" resolves ≥ 1.27.1)                                                | MED    |
| 18 | Add an AGENTS.md gotcha: `.lychee.toml` is not auto-discovered; use `lychee.toml`                                              | LOW    |
| 19 | Add pipefail discipline note to crush-config lessons.md (cross-project lesson)                                                 | LOW    |
| 20 | Re-measure + re-pin branching-flow pins after the next analyzer upgrade                                                        | LOW    |
| 21 | Evaluate `nix flake check --all-systems` (aarch64 builders needed?)                                                            | LOW    |
| 22 | Consider adding go-licenses + govalid to `devShells.ci` for env parity                                                         | LOW    |
| 23 | Opt-in property test run: `GBR_PROPERTY=1 nix develop --command go test ./...`                                                 | LOW    |
| 24 | Fresh-GOMODCACHE health check of published v2.2.0 (AGENTS stale-cache recipe)                                                  | LOW    |
| 25 | go-cqrs-lite upstream: consider relaxing snapshot v4.5.1's `go 1.27.1` floor to `go 1.27` (would let adapters use major.minor) | MED    |
| 26 | Restart LSP to clear the stale exhaustruct editor warnings (cosmetic)                                                          | LOW    |
| 27 | Pre-commit size check for AGENTS.md's 377-line budget (catch early)                                                            | LOW    |
| 28 | Decide jscpd: keep intentional-duplication doc vs exclude `.github/workflows`                                                  | LOW    |
| 29 | Review `.buildflow.yml` GOEXPERIMENT env duplication ("caller wins" info noise)                                                | LOW    |
| 30 | Watch govalid upstream for releases > 1.9.0; bump the mirrored rev in BOTH flakes                                              | LOW    |
| 31 | Re-verify flake-meta-checker fully clean (0 findings) in the next full run                                                     | LOW    |
| 32 | Benchmarks sanity after toolchain bump (Stream 2/10-rules numbers)                                                             | LOW    |
| 33 | Consider `nh sw` single-package path vs full `.#deploy` for toolchain-only updates                                             | LOW    |
| 34 | Periodic: `nix flake update` + check-all cadence (go-mod-update ran 4× in-loop today)                                          | LOW    |
| 35 | Consider pinning nixpkgs input vs tracking `nixos-unstable` branch (policy call)                                               | LOW    |
| 36 | Verify Polish-Customs consumer unaffected (its own repo; consumes published v2.2.0)                                            | LOW    |
| 37 | AGENTS.md: document the fast-vs-full govalid asymmetry gotcha                                                                  | LOW    |
| 38 | Consider a `just`-free flake app for "quiet-window deploy" (wraps pressure polling)                                            | LOW    |
| 39 | Sweep other fleet repos for `go_1_26` pins (same bug likely elsewhere)                                                         | MED    |
| 40 | Sweep fleet for treefmt goimports + `go 1.27` floors (same hermetic-download bug)                                              | MED    |
| 41 | Consider exposing per-module go floors as a single source of truth (flake eval)                                                | LOW    |
| 42 | Check whether BuildFlow's govalid full-mode fan-out respects devshell GOWORK=off                                               | LOW    |
| 43 | Add `docs/status/` index or staleness annotations for old reports (docs-health ANNOTATE)                                       | LOW    |
| 44 | Revisit `GOTOOLCHAIN=local` policy: keep (hermetic) — document explicitly in AGENTS.md                                         | LOW    |
| 45 | Consider running check-all inside CI (it duplicates the matrix; probably skip)                                                 | LOW    |
| 46 | Look at `go-tool-run` / `ginkgo-version-check` steps' purpose (unreviewed, green)                                              | LOW    |
| 47 | Consider vendoring nothing / verify no vendor/ crept in (nix-hash-fix ran in-loop)                                             | LOW    |
| 48 | Confirm `go.sum` hygiene after 4 go-mod-update cycles (diff review)                                                            | LOW    |
| 49 | CHANGELOG entry for the toolchain bump (user-visible: requires Go ≥ 1.27.1 toolchain)                                          | MED    |
| 50 | Decide whether README needs a "requires Go 1.27.1 toolchain to build" note                                                     | LOW    |

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **Deploy authority & timing:** `govalid-generate` in full mode stays red until
   the system Go becomes 1.27.1. Shall I run `cd ~/projects/SystemNix &&
   nix run .#deploy` myself as soon as the IO-pressure gate allows (it switches
   your RUNNING system, restarts hermes/agent sessions, and affects the other
   concurrent Crush sessions), or do you want to fire it yourself at a time you
   choose?
2. **Pressure-gate override policy:** is `DEPLOY_FORCE_PRESSURE=1` ever an
   acceptable override for you (documented kernel-freeze precursor class), or is
   that gate always human-only? I chose not to force; if waiting is the answer,
   the repo stays 2-steps-red until the machine quiets down.
3. **samber/lo policy resolution:** is "no samber/lo anywhere" permanent
   fleet-wide — in which case I should push a finding-suppression mechanism for
   go-auto-upgrade upstream into BuildFlow — or is `lo` acceptable in test
   files / nested modules, in which case I should convert the 8 flagged loops?

---

_Point-in-time snapshot; goes stale. Route section (f) through docs-health
HARVEST rather than treating this file as backlog._
