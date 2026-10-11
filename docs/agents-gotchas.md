# Agent Gotchas & Triage Recipes

Long-form context moved out of `AGENTS.md` (2026-10-11, when the
`go-structure-linter` agent-config rule began gating on AGENTS.md length).
Every entry here is referenced from `AGENTS.md`; nothing in this file is
backlog — these are standing explanations, recipes, and policies.

## Go toolchain skew (why ambient buildflow fails)

The devshell uses `pkgs.go_1_27` (currently 1.27.2) because
`go-cqrs-lite/snapshot/v4@v4.5.1` requires go >= 1.27.1 and BuildFlow runs
tools with `GOTOOLCHAIN=local`. The devshell also ships `go-licenses` and a
`govalid` built with `buildGo127Module` (mirroring `SystemNix/pkgs/govalid.nix`
— keep the rev in sync) so BuildFlow resolves them with the right Go.
treefmt's goimports is rebuilt with Go 1.27 (`go127Gotools`): its internal
`go` would otherwise try a toolchain download and die in the hermetic sandbox.
CI pins `go-version: "1.27"` (matches the module floor; setup-go picks the
latest 1.27.x). Until the system Go is >= 1.27.1 (SystemNix deploy —
user-fired, needs sudo), run BuildFlow as
`nix develop --command buildflow --build-mode=full`: ambient go fails the
adapters floor (green this way; one first-run govalid race).

Symptom when violated: `govalid: go: .../snapshot/v4@v4.5.1 requires go >=
1.27.1 (running go 1.27.0)` and `test-race` failures on ambient invocations,
even with `GOTOOLCHAIN=auto` exported (BuildFlow forces `GOTOOLCHAIN=local`
for its tools).

## Private modules & the GOPRIVATE/GONOSUMDB flake override

Both devShells explicitly set all three of `GOPRIVATE`, `GONOSUMDB`, and
`GONOPROXY` to `github.com/larsartmann/*,github.com/LarsArtmann/*` (both case
variants). The repo and all its dependencies are PUBLIC since 2026-09-17, so
these are harmless no-ops today — but they stay because Home Manager sets
`GONOSUMDB` as an OS env var to an explicit list, and Go only auto-derives it
from `GOPRIVATE` when **unset**: a bare `GOPRIVATE` wildcard loses, and a
force-pushed private tag then fails sumdb verification with a
checksum-mismatch SECURITY ERROR (2026-07-26 gotcha).

Related recipe: a stale private-era module-cache download can fail sumdb
against the healthy proxy zip — verify module health ONLY with a fresh
`GOMODCACHE` (`trash "$(go env GOMODCACHE)/cache/download/github.com/!lars!artmann/go-business-rules"`,
then refetch).

## CI & publishing reality

- The repo is PUBLIC (2026-09-17): GitHub Actions are free on public repos and
  all dependencies resolve from the public proxy — no tokens needed. CI is a
  matrix over all four Go modules (root, `adapters/cqrslite`, `examples/sse`,
  `listeners/otel`), test/lint/build each, gosec root-only with
  `GOEXPERIMENT=jsonv2`. Its exact commands are verified locally via
  `nix run .#check-all`.
- The `Security` job runs `go install ...gosec@v2.29.0` (the container
  action's Go predates jsonv2) and asserts `.Stats.files > 0`. The
  branching-flow integration specs `Skip` when the `branching-flow` binary is
  absent (CI), but enforce pins locally.
- pkg.go.dev/proxy fetch on demand since the public flip; older
  "unverifiable" reports predate it.
- Local gates are the real quality bar: `nix run .#check-all` (all 4 modules:
  build, vet, test, lint), `nix develop --command go test ./...` (root,
  251/251 specs incl. branching-flow pins and the opt-in property test),
  `buildflow` (quality gate), `nix build .#checks.x86_64-linux.format`.

## encoding/json/v2 migration detail

This library uses `encoding/json/v2`, which is experimental and requires
`GOEXPERIMENT=jsonv2`.

- **flake.nix**: Both devShells (`default` and `ci`) set
  `GOEXPERIMENT = "jsonv2"`, so all `go build`/`go test`/`go vet` commands
  work inside `nix develop`.
- **.golangci.yml**: `goexperiment.jsonv2` is in `build-tags`.

**Downstream consumer impact:** this is a hard breaking change for downstream
consumers. Anyone who `go get`s this library must set `GOEXPERIMENT=jsonv2`
(and use Go 1.26+) or compilation fails with a cryptic "build constraints
exclude all Go files in encoding/json/v2" error. There is no way to
self-contain this requirement in `go.mod`.

**Re-evaluate when**: `encoding/json/v2` graduates from experimental status
(no longer requires `GOEXPERIMENT=jsonv2`). At that point the downstream
constraint disappears.

**Affected files**: `errors.go`, `validation_result.go` (MarshalJSON),
`bdd_branching_flow_test.go` (Unmarshal).

## Branching-flow pins: history & re-pin procedure

The branching-flow multi-linter reports PHANTOM/DUPE-class violations that are
**false positives for this validation library**: the primitives ARE the domain
concept being validated (users pass raw primitives; builders take raw
primitive parameters by design; forcing branded types would defeat the
library's purpose).

**Policy (applied 2026-09-14): re-pin to current analyzer reality.** A
permanently red suite hides NEW regressions; the pins still catch drift from
code changes.

Current pins (2026-10-11, branching-flow 0.6.4 / 192031d): 20 PHANTOM total
(all PHANTOM-TRANSPOSE same-type parameter pairs, all warning severity, across
`builders.go`, `builders_collection.go`, `builders_format.go`,
`builders_network.go`) and 24 `stats` totalIssues = 20 phantom + 2 flagparam +
1 ifacecomplete + 1 mixins.

Version history:

- **0.2.0 era:** pins were 25 PHANTOM in mixed severities (7 critical / 6
  error / 11 info / 1 warning) and a 31 total; 0.2.0's stats COUNTED the two
  nolint-suppressed panic findings (Stream result send + DivisibleBy
  zero-guard), pushing 29 → 31.
- **0.6.4:** rewrote the phantom rule set — metadata `name` is now the pair
  (e.g. "value, minimum") instead of a single identifier, severities
  collapsed to warning-only, and `stats` no longer counts the
  nolint-suppressed panic findings (total 31 → 24 with zero code change).

The `//nolint:branching-flow:panic` suppressions REMAIN REQUIRED under 0.6.4:
the panic analyzer still flags the `Stream` result send in `validator.go`
(`results` is closed only after `waitGroup.Wait()`; verified 2026-10-11 —
removing the nolint makes panic report 2 findings again) and the
`DivisibleBy` zero-guard in `builders_collection.go` (early return untracked).
`.golangci.yml` excludes both files from `nolintlint`, which otherwise flags
the analyzer's directive syntax as unknown.

**Re-pin procedure (a red pin means):** diff the findings
(`buildflow -s ... --format finding` or `branching-flow phantom --format
finding .`); confirm they are still the same deliberate-primitives false
positive class; re-measure `stats --format json`; update the pins in
`bdd_branching_flow_test.go` AND this doc AND the AGENTS.md pin bullet, with a
dated comment. **Pin fragility:** the analyzer scans the whole directory and
has no path-exclude flag, so ANY new module/example/builder/test changes the
counts. Analyzer binary drift can move the totals without any code change —
diff before assuming new violations.

## Benchmarks (events & streaming)

- `Build` with no listeners: 189 ns/op; one listener: 468 ns/op; three
  listeners: ~476 ns/op (2 rules).
- Stream: `BenchmarkStream2Rules`, `BenchmarkStream10Rules`,
  `BenchmarkStream10RulesConcurrency4` (channel drain dominates: ~4.4 us/op
  for 2 cheap rules).

## Per-analyzer false-positive rationales

### Hierarchical-errors analyzer

Flags functions returning generic `error`. False positives — signatures fixed
by interfaces or Go convention: `MarshalJSON` (`json.Marshaler`, stdlib-fixed
signature), `Check` (the `Rule` interface; a custom error type would couple
every implementation to one error type), and internal aggregators
(`checkNonEmpty`, `collectAllViolations`, `anyRulePasses`) which inherit
`error` from `Rule.Check()`. Intentional design decisions; changing them
breaks compatibility.

### go-auto-upgrade analyzer

Suggests replacing small manual loops with `samber/lo` helpers (SliceToMap in
`validation_result.go` BySeverity, Filter at `validation_result.go:47/137`,
Map in `property_test.go:48`, Reduce in
`listeners/otel/adapter_test.go:114/130`). False positives: adding
`samber/lo` for 3-line loops would introduce new dependencies (a runtime dep
for the root library, a dep for test helpers in the nested modules) for
negligible benefit. The manual loops are idiomatic Go, and the public API
surface (`Filter`, `BySeverity`) must stay dependency-free.

### go-structure-linter: root-package-files

Flags Go source files at the project root rather than under `/internal/` or
`/pkg/`. False positive for a public Go library: the root-level package files
ARE the public API; moving them to `/internal/` would make them unexportable,
and moving them to `/pkg/` would change the import path for all consumers.
`.go-structure-linter.yaml` pins the `flat` preset for this reason.

### gomod-check false positive

`go.mod:12: direct and indirect requires are mixed` — the blocks are already
separated; `go mod tidy` is a no-op. Nothing to fix.

### DUPE violations

Reduced through refactoring: `All`/`Any` use a strategy pattern that reduces
structural similarity while maintaining clarity.

## Residual BuildFlow findings (policy, 2026-10-04)

All below the error-level findings gate; triaged and intentionally left:

- **nix-checker (4)**: govalid's pinned `hash`/`vendorHash` are mandatory for
  `fetchFromGitHub`/`buildGoModule`; the extract-to-`hash.nix` style
  suggestions are declined (single rarely-bumped derivation;
  `buildflow -s nix-hash-fix --fix` owns hash repair).
- **cqrs-lint (3, info)**: A009 stack preset wrong for an adapter (consumer
  owns wiring); B009 cqrs-gen overkill for one function; D013 moot —
  `event.New` defaults `SchemaVersion` to 1 (reconstructed from storage), so
  `WithSchemaVersion(1)` would restate it.
- **vulnix (28, all warning)**: CVEs in the nixpkgs build-toolchain closure.
  The repo ships Go source, not binaries. Advisory; re-check after
  `nix flake update`; act only on runtime-relevant CVEs.
- **doctor "34 failed"**: per-tool availability across ALL providers (pytest,
  cargo-*, ...) — noise for a Go+Nix repo. Real gaps: `govulncheck` missing;
  `go-licenses`/`lychee` are devshell-only.
- **license-check**: green full-mode without result cache; `Unknown` entries
  are nested-module LICENSE mapping (repo-root LICENSE covers) — cosmetic.

## Lychee link-checker gotchas

`lychee.toml` excludes `docs/planning/archived` and `docs/status/archived`
from link checking: strikethrough planner entries like `OneOf[T]()` parse as
empty markdown links, and archived files are never edited to appease a
linter. The filename must be `lychee.toml` — a dot-prefixed `.lychee.toml` is
NOT auto-discovered.
