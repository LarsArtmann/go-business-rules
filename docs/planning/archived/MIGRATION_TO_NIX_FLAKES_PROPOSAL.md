# Migration to Nix Flakes — Proposal

> **RESOLVED 2026-09-14:** migration complete — `flake.nix` (flake-parts + treefmt) is the canonical build/task automation, the justfile was removed, and AGENTS.md documents the devShell workflow. Archived from the repo root.
>
> ~~**Status:** Draft~~ | **Date:** 2026-04-09 | **Target:** `github.com/artmann/businessrules`

---

## Table of Contents

- [Executive Summary](#executive-summary)
- [Current State Analysis](#current-state-analysis)
- [Why Nix Flakes](#why-nix-flakes)
- [Proposed Architecture](#proposed-architecture)
- [Migration Steps](#migration-steps)
- [File Inventory](#file-inventory)
- [Reference: flake.nix Blueprint](#reference-flakenix-blueprint)
- [CI Integration Plan](#ci-integration-plan)
- [Risk Assessment](#risk-assessment)
- [Open Questions](#open-questions)
- [Success Criteria](#success-criteria)

---

## Executive Summary

This proposal recommends adopting Nix Flakes as a **supplemental** (not replacement) development environment definition for the businessrules Go library. The flake provides reproducible dev shells, hermetic checks, and pinned tool versions — while the existing justfile and GitHub Actions CI remain the primary build system.

**TL;DR:** Add a `flake.nix` + `flake.lock` that gives contributors a single-command reproducible environment (`nix develop`) with all tools pinned, and enable `nix flake check` for hermetic CI validation.

---

## Current State Analysis

### Build & Dev Tooling Inventory

| Tool            | Version Pinning             | Source                         | Used By               |
| --------------- | --------------------------- | ------------------------------ | --------------------- |
| Go              | 1.26.1 (go.mod) / 1.25 (CI) | Direct install                 | Build, test, vet, fmt |
| golangci-lint   | `latest` (CI), v2 config    | Direct install / GitHub Action | Lint                  |
| gosec           | `master` (CI)               | GitHub Action                  | Security scan         |
| just            | Unpinned                    | Direct install                 | Task runner           |
| go-mod-outdated | Unpinned                    | Direct install                 | Dependency audit      |
| godoc           | Unpinned (stdlib)           | Go toolchain                   | Documentation         |
| Ginkgo/Gomega   | v2.28.1 / v1.39.1           | go.mod (dev dep)               | Test framework        |
| codecov         | v4 action                   | GitHub Action                  | Coverage upload       |

### Configuration Files

| File                       | Purpose                                                    |
| -------------------------- | ---------------------------------------------------------- |
| `justfile`                 | 17 recipes: test, lint, bench, fuzz, build, security, etc. |
| `.golangci.yml`            | 100+ linters, formatters, exclusions                       |
| `.github/workflows/ci.yml` | 4 jobs: test, lint, build, security                        |
| `library-policy.yaml`      | Scanner config (encoding/json/v2 disable)                  |
| `.editorconfig`            | Formatting rules (tabs for Go, spaces for YAML/Markdown)   |
| `go.mod` / `go.sum`        | Module definition and lock file                            |

### Current Pain Points

1. **No tool version pinning.** golangci-lint runs as `latest` in CI. gosec runs from `master`. Contributors may have different versions locally.
2. **Go version mismatch.** go.mod specifies 1.26.1 but CI uses 1.25. This works today but risks future breakage.
3. **Onboarding friction.** New contributors must install Go, just, golangci-lint, gosec independently. No single command sets up the environment.
4. **Platform assumptions.** CI runs on ubuntu-latest only. macOS contributors have no CI signal for their platform.
5. **No hermetic verification.** Tests depend on whatever Go toolchain is on PATH. Results are not bit-for-bit reproducible.

---

## Why Nix Flakes

| Benefit                            | Impact for This Project                                             |
| ---------------------------------- | ------------------------------------------------------------------- |
| **Reproducible environments**      | Every contributor gets identical tool versions via `nix develop`    |
| **Single-command onboarding**      | `nix develop` installs Go, just, golangci-lint, gosec, gopls, delve |
| **Pinned dependencies**            | `flake.lock` ensures deterministic tool versions across time        |
| **Multi-platform CI**              | Same flake works on Linux and macOS runners                         |
| **Hermetic checks**                | `nix flake check` runs tests in isolated sandbox                    |
| **Complementary, not replacement** | Existing justfile and CI continue to work unchanged                 |
| **Cachix integration**             | Binary cache speeds up CI and local builds                          |

### Why Not Alternatives

| Alternative                   | Why Not                                                                     |
| ----------------------------- | --------------------------------------------------------------------------- |
| **Docker/DevContainer**       | Heavyweight for a pure Go library. No container runtime needed.             |
| **asdf/mise**                 | Only pins runtime versions, not linters and tools.                          |
| **Nix classic (default.nix)** | No lock file, no standard CLI, no registry. Flakes are the modern standard. |
| **Makefile**                  | Already have justfile; doesn't solve tool pinning.                          |

---

## Proposed Architecture

### Layer Diagram

```
┌─────────────────────────────────────────────────┐
│                  CI (GitHub Actions)              │
│  ┌──────────────┐  ┌───────────────────────────┐ │
│  │ Existing jobs │  │ Nix flake check job (new) │ │
│  │ (unchanged)   │  │ nix flake check           │ │
│  └──────────────┘  │ nix build                  │ │
│                     └───────────────────────────┘ │
├─────────────────────────────────────────────────┤
│               flake.nix (new)                     │
│  ┌────────────┐ ┌─────────┐ ┌────────────────┐  │
│  │ devShells  │ │ checks  │ │ formatter      │  │
│  │ Go + tools │ │ test    │ │ nixfmt         │  │
│  │ just       │ │ lint    │ │                │  │
│  │ gopls      │ │ vet     │ │                │  │
│  │ delve      │ │ fmt     │ │                │  │
│  └────────────┘ └─────────┘ └────────────────┘  │
├─────────────────────────────────────────────────┤
│               justfile (unchanged)                │
│         Existing recipes, still work              │
├─────────────────────────────────────────────────┤
│            go.mod / .golangci.yml (unchanged)     │
└─────────────────────────────────────────────────┘
```

### Design Principles

1. **Additive, not replacing.** The flake is a new entry point. Existing workflows remain primary.
2. **Library-aware.** Since businessrules is a library (no `main` package), the flake focuses on `devShells` and `checks` rather than binary packages.
3. **Multi-system.** Support `x86_64-linux`, `aarch64-linux`, `x86_64-darwin`, `aarch64-darwin`.
4. **Minimal inputs.** Only `nixpkgs` as a flake input. No flake-parts, no utility flakes. Keep it simple.

---

## Migration Steps

### Phase 1: Foundation (Estimated: 1 session)

#### Step 1.1 — Create `flake.nix`

Create the flake with three outputs: `devShells`, `checks`, and `formatter`.

**Inputs:**

- `nixpkgs` — pinned to `nixos-unstable` for latest Go toolchain

**Outputs:**

- `devShells.<system>.default` — Development environment
- `checks.<system>.test` — Run `go test -race ./...`
- `checks.<system>.lint` — Run `golangci-lint run ./...`
- `checks.<system>.vet` — Run `go vet ./...`
- `checks.<system>.fmt` — Verify `gofmt` compliance
- `formatter.<system>` — `nixfmt-classic` for `.nix` files

**Tools in devShell:**

| Package           | Purpose                       |
| ----------------- | ----------------------------- |
| `go`              | Go toolchain                  |
| `golangci-lint`   | Linting                       |
| `gopls`           | Language server (IDE support) |
| `delve`           | Debugger                      |
| `just`            | Task runner                   |
| `gosec`           | Security scanner              |
| `go-mod-outdated` | Dependency audit              |
| `nixfmt-classic`  | Nix file formatting           |
| `gofumpt`         | Go formatting (strict)        |
| `goimports`       | Import management             |

#### Step 1.2 — Generate `flake.lock`

Run `nix flake lock` to pin all input versions. Commit both `flake.nix` and `flake.lock`.

#### Step 1.3 — Update `.gitignore`

Add Nix-related entries:

```
# Nix
result
result-*
```

#### Step 1.4 — Verify Locally

- [ ] `nix develop` — enters dev shell with all tools
- [ ] `just test` inside dev shell — passes
- [ ] `just lint` inside dev shell — passes
- [ ] `nix flake check` — all checks pass
- [ ] `nix fmt` — formats flake.nix (no-op if already formatted)

### Phase 2: CI Integration (Estimated: 1 session)

#### Step 2.1 — Add Nix Check Job to CI

Add a new `nix` job to `.github/workflows/ci.yml`:

```yaml
nix:
  name: Nix
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: cachix/install-nix-action@v30
      with:
        github_access_token: ${{ secrets.GITHUB_TOKEN }}
    - uses: cachix/cachix-action@v15
      with:
        name: artmann-businessrules
        authToken: ${{ secrets.CACHIX_AUTH_TOKEN }}
    - run: nix flake check
    - run: nix develop --command just test
```

This job runs **in parallel** with existing jobs, not replacing them.

#### Step 2.2 — Set Up Cachix (Optional)

1. Create a Cachix cache: `artmann-businessrules`
2. Add `CACHIX_AUTH_TOKEN` secret to GitHub repo
3. Configure `cachix/cachix-action` in CI workflow

**Benefits:** Cached Nix derivations speed up CI by 2-10x for subsequent runs.

#### Step 2.3 — Verify CI

- [ ] Nix job passes on GitHub Actions
- [ ] Existing jobs still pass (unchanged)
- [ ] Cachix cache populated (if configured)

### Phase 3: Documentation (Estimated: 0.5 session)

#### Step 3.1 — Update README.md

Add a "Development with Nix" section:

```markdown
## Development with Nix

If you use Nix, get a fully-pinned dev environment with one command:

    nix develop

This provides Go, golangci-lint, just, gosec, gopls, and all other tools.

Run all checks:

    nix flake check
```

#### Step 3.2 — Update AGENTS.md

Add Nix-specific build commands to the Build Commands section.

### Phase 4: Polish (Estimated: 0.5 session)

#### Step 4.1 — Fix Go Version Consistency

Align CI Go version with go.mod (both should be 1.26.1). This is an existing issue independent of Nix, but the flake will make it visible.

#### Step 4.2 — Consider direnv Integration

Add `.envrc` with:

```
use flake
```

This auto-activates the dev shell when entering the project directory (for direnv users).

#### Step 4.3 — Periodic flake.lock Updates

Add Dependabot or Renovate config to auto-update `flake.lock` monthly (if those tools are already in use), or document a manual `nix flake update` cadence.

---

## File Inventory

### New Files

| File         | Purpose                                     | Priority           |
| ------------ | ------------------------------------------- | ------------------ |
| `flake.nix`  | Nix Flake definition                        | Required (Phase 1) |
| `flake.lock` | Pinned dependency versions (auto-generated) | Required (Phase 1) |
| `.envrc`     | direnv integration (optional)               | Optional (Phase 4) |

### Modified Files

| File                       | Change                      | Priority              |
| -------------------------- | --------------------------- | --------------------- |
| `.gitignore`               | Add `result`, `result-*`    | Required (Phase 1)    |
| `.github/workflows/ci.yml` | Add Nix CI job              | Required (Phase 2)    |
| `README.md`                | Add Nix development section | Recommended (Phase 3) |
| `AGENTS.md`                | Add Nix build commands      | Recommended (Phase 3) |

### Unchanged Files

| File                  | Reason                                         |
| --------------------- | ---------------------------------------------- |
| `justfile`            | Continues to work as-is inside and outside Nix |
| `.golangci.yml`       | Consumed by golangci-lint from any environment |
| `go.mod` / `go.sum`   | Not affected by Nix                            |
| `library-policy.yaml` | Project policy, not Nix concern                |
| `.editorconfig`       | Editor settings, orthogonal to Nix             |

---

## Reference: flake.nix Blueprint

This is a **starting point** for review, not the final implementation. Adjust versions and packages after testing.

```nix
{
  description = "Severity-aware validation for Go — github.com/artmann/businessrules";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      pkgsFor = system: nixpkgs.legacyPackages.${system};
    in
    {
      # Development environments
      devShells = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
        in
        {
          default = pkgs.mkShell {
            name = "businessrules-dev";

            packages = with pkgs; [
              go
              golangci-lint
              gopls
              delve
              just
              gosec
              gotools
              gofumpt
              nixfmt-classic
            ];

            shellHook = ''
              echo "businessrules dev shell"
              echo "Go: $(go version)"
              echo "golangci-lint: $(golangci-lint version --format short 2>/dev/null || echo 'available')"
              echo ""
              echo "Quick commands:"
              echo "  just test          - Run tests"
              echo "  just lint          - Run linter"
              echo "  just check         - Run all quality checks"
              echo "  nix flake check    - Run hermetic checks"
              echo ""
            '';
          };
        }
      );

      # Hermetic checks (run via `nix flake check`)
      checks = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
        in
        {
          test = pkgs.runCommand "businessrules-test" { nativeBuildInputs = [ pkgs.go ]; } ''
            cd ${self}
            go test -v -race ./...
            touch $out
          '';

          vet = pkgs.runCommand "businessrules-vet" { nativeBuildInputs = [ pkgs.go ]; } ''
            cd ${self}
            go vet ./...
            touch $out
          '';

          fmt = pkgs.runCommand "businessrules-fmt-check" { nativeBuildInputs = [ pkgs.go ]; } ''
            cd ${self}
            test -z "$(gofmt -l .)" || (echo "Files need formatting:"; gofmt -l .; exit 1)
            touch $out
          '';
        }
      );

      # Nix file formatter (`nix fmt`)
      formatter = forAllSystems (system: (pkgsFor system).nixfmt-classic);
    };
}
```

### Notes on the Blueprint

1. **No `packages` output.** This is a Go library with no `main` package. `buildGoModule` would build nothing useful. The value is in `devShells` and `checks`.
2. **`checks.test` uses `runCommand`.** Simpler than `buildGoModule` for a library. Downloads deps fresh each time but provides true hermeticity.
3. **No `golangci-lint` in checks.** The project has an extensive `.golangci.yml` with 100+ linters. Running this in a Nix check would be slow and would couple the flake to the linter config. The CI lint job handles this better.
4. **`nixfmt-classic`** is used over `nixfmt-rfc-style` to match the project's existing formatting conventions. Switch to `nixfmt` (RFC 166 style) if preferred.

### Alternative: buildGoModule for Checks

If `buildGoModule` is desired for caching benefits (via Cachix), the check phase can be integrated:

```nix
checks = forAllSystems (
  system:
  let
    pkgs = pkgsFor system;
    pkg = pkgs.buildGoModule {
      pname = "businessrules";
      version = "1.1.0";
      src = self;
      vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
      doCheck = true;
      checkFlags = [ "-v" "-race" ];
    };
  in
  {
    inherit (pkg) passthru.tests;
    # Or: test = pkg;
  }
);
```

**Trade-off:** `buildGoModule` provides Cachix caching but requires updating `vendorHash` on every `go.sum` change. For a library, the simpler `runCommand` approach is recommended initially.

---

## CI Integration Plan

### Proposed CI Job Addition

```yaml
# Add to .github/workflows/ci.yml
nix:
  name: Nix
  runs-on: ubuntu-latest
  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Install Nix
      uses: cachix/install-nix-action@v30
      with:
        github_access_token: ${{ secrets.GITHUB_TOKEN }}
        extra_nix_config: |
          experimental-features = nix-command flakes

    - name: Setup Cachix
      uses: cachix/cachix-action@v15
      with:
        name: artmann-businessrules
        authToken: ${{ secrets.CACHIX_AUTH_TOKEN }}

    - name: Check Nix flake
      run: nix flake check --print-build-logs

    - name: Verify dev shell
      run: nix develop --command echo "Dev shell works"
```

### CI Job Comparison

| Aspect            | Existing Jobs                      | Nix Job                  |
| ----------------- | ---------------------------------- | ------------------------ |
| Tool installation | `setup-go` action + direct         | Nix flake                |
| Version pinning   | Partial (go.mod, CI yaml)          | `flake.lock`             |
| Reproducibility   | Moderate                           | High                     |
| Cache mechanism   | Go module cache + actions/cache    | Cachix                   |
| Platforms         | Linux only                         | Linux + macOS            |
| Coverage          | Full (test, lint, build, security) | Partial (test, vet, fmt) |

### Migration Path for CI

1. **Phase 2:** Add Nix job alongside existing jobs (no removal).
2. **Phase 3+ (optional):** If Nix job proves reliable, consider replacing the `build` and `test` jobs with the Nix equivalents. The `lint` job stays separate due to the complex `.golangci.yml` configuration.
3. **Not recommended:** Replacing the security job. `gosec` is better run via its dedicated GitHub Action.

---

## Risk Assessment

### Low Risk

| Risk                                     | Mitigation                                                                          |
| ---------------------------------------- | ----------------------------------------------------------------------------------- |
| Flake doesn't work on some system        | `supportedSystems` limits exposure. Fallback to non-Nix workflow.                   |
| Tool version mismatch between Nix and CI | Both use `nixos-unstable` / `latest`. Document known-good versions in `flake.lock`. |
| Contributors don't have Nix installed    | Nix is entirely optional. All tools work without it.                                |
| `flake.lock` drift                       | Document `nix flake update` cadence.                                                |

### Medium Risk

| Risk                                          | Mitigation                                           |
| --------------------------------------------- | ---------------------------------------------------- |
| `nix flake check` slow on first run           | Cachix caching. `runCommand` checks are lightweight. |
| `vendorHash` updates if using `buildGoModule` | Use `runCommand` approach instead (recommended).     |
| Cachix costs for public cache                 | Free for open-source projects.                       |

### Negligible Risk

| Risk                      | Mitigation                                                |
| ------------------------- | --------------------------------------------------------- |
| Existing workflows break  | No modifications to existing files (except `.gitignore`). |
| go.sum conflicts          | Nix does not modify go.sum.                               |
| `.editorconfig` conflicts | Nix respects project formatting rules.                    |

---

## Open Questions

These questions should be resolved before implementation:

1. **Nixfmt style:** Use `nixfmt-classic` (traditional) or `nixfmt` (RFC 166 style)? RFC 166 is the newer standard but may look unfamiliar.

2. **Cachix cache name:** Is `artmann-businessrules` acceptable, or is there a preferred Cachix org/cache name?

3. **CI Go version alignment:** CI uses Go 1.25, go.mod specifies 1.26.1. Should these be aligned? The flake will use whatever `nixpkgs` provides (currently 1.24.x on stable, 1.25+ on unstable).

4. **Direnv support:** Should `.envrc` be included, or left to individual contributors?

5. **golangci-lint in checks:** Should `nix flake check` include a full lint run? This would be slow (~5 min with 100+ linters) but would provide complete hermetic verification.

6. **Flake inputs scope:** Only `nixpkgs` is proposed. Should we also consider:
   - `flake-utils` for cleaner system iteration? (Avoidable with `lib.genAttrs`)
   - `golangci-lint` as a separate input for version pinning? (Overkill for now)

7. **`buildGoModule` vs `runCommand`:** The blueprint uses `runCommand` for simplicity. Should we instead use `buildGoModule` for Cachix caching benefits? This requires maintaining a `vendorHash` that updates on every `go.sum` change.

---

## Success Criteria

The migration is complete when:

- [ ] `nix develop` provides a working dev environment with all tools
- [ ] `nix flake check` passes on all supported systems
- [ ] `just test` works inside the Nix dev shell
- [ ] CI Nix job passes alongside existing jobs
- [ ] `flake.lock` is committed and provides reproducible builds
- [ ] README documents the Nix workflow
- [ ] No existing workflows are broken
- [ ] At least one contributor has verified the dev shell on macOS

### Post-Migration Metrics

| Metric                      | Target                              |
| --------------------------- | ----------------------------------- |
| Time to `nix develop` ready | < 60 seconds (with warm cache)      |
| `nix flake check` duration  | < 120 seconds (with warm cache)     |
| First-time setup            | < 5 minutes (including Nix install) |
| Existing CI job duration    | Unchanged                           |

---

## Appendix A: Tool Version Pinning Strategy

The `flake.lock` file pins the `nixpkgs` commit, which determines all tool versions:

| Tool          | Version Source | Update Mechanism   |
| ------------- | -------------- | ------------------ |
| Go            | nixpkgs commit | `nix flake update` |
| golangci-lint | nixpkgs commit | `nix flake update` |
| just          | nixpkgs commit | `nix flake update` |
| gosec         | nixpkgs commit | `nix flake update` |
| gopls         | nixpkgs commit | `nix flake update` |
| delve         | nixpkgs commit | `nix flake update` |

All tools update together when `nix flake update` is run. This is simpler than pinning individual versions and ensures compatibility.

### Update Cadence

- **Recommended:** Monthly `nix flake update` to pull latest nixpkgs.
- **Maximum:** Quarterly to avoid security vulnerabilities in pinned tools.

## Appendix B: Nix Install Quick Reference

```bash
# macOS ( Determinate Systems installer — recommended)
curl --proto '=https' --tlsv1.2 -sSf -L \
  https://install.determinate.systems/nix | sh -s -- install

# Linux (same installer)
curl --proto '=https' --tlsv1.2 -sSf -L \
  https://install.determinate.systems/nix | sh -s -- install

# Verify
nix --version
nix develop  # Enter dev shell
```

## Appendix C: Mapping — Justfile Recipes to Nix

| Justfile Recipe | Nix Equivalent                                                   | Notes                                                |
| --------------- | ---------------------------------------------------------------- | ---------------------------------------------------- |
| `just test`     | `nix flake check` (partial) or `nix develop --command just test` | Flake check runs test + vet + fmt                    |
| `just lint`     | No direct Nix equivalent                                         | Run via dev shell: `nix develop --command just lint` |
| `just bench`    | No Nix check (intentional)                                       | Benchmarks are ad-hoc, not hermetic                  |
| `just fuzz`     | No Nix check (intentional)                                       | Fuzzing is interactive, not hermetic                 |
| `just build`    | Implicit in `nix flake check`                                    | `go build` is part of the check derivation           |
| `just check`    | `nix flake check` + `nix develop --command just lint`            | Closest full equivalent                              |
| `just security` | No Nix check (intentional)                                       | gosec better handled by dedicated CI job             |
| `just fmt`      | `nix fmt` (Nix files only)                                       | Go formatting still via `gofmt`/`gofumpt`            |

**Design choice:** Not every justfile recipe needs a Nix check. Hermetic checks are best for deterministic, fast operations. Linting (100+ linters), fuzzing, and security scanning remain better suited for CI or manual invocation.
