# Contributing to go-business-rules

> **Thank you for contributing!** This guide covers how to contribute effectively to `businessrules`.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Development Setup](#development-setup)
- [Building & Testing](#building--testing)
- [Code Standards](#code-standards)
- [Linting & Quality](#linting--quality)
- [Pull Request Process](#pull-request-process)
- [Commit Messages](#commit-messages)

---

## Code of Conduct

We are committed to providing a welcoming and respectful environment. All contributors are expected to:

- **Be respectful** — Treat others with kindness and professionalism
- **Be inclusive** — Welcome diverse perspectives and experiences
- **Be constructive** — Provide feedback that helps improve the project
- **Be collaborative** — Work together to achieve the best outcomes

---

## Development Setup

This project uses [Nix flakes](https://nixos.wiki/wiki/Flakes) for its development environment.

### Prerequisites

| Tool | Version | Purpose |
| --- | --- | --- |
| Nix | 2.18+ (with flakes enabled) | Reproducible dev environment |
| Go | 1.26.4 (provided by the Nix shell) | Language runtime |

### Enter the development shell

```bash
nix develop
```

The shell provides `go`, `golangci-lint`, `gopls`, `delve`, `gosec`, `gofumpt`, and the `gotools`.
It also sets `GOEXPERIMENT=jsonv2`, which is **required** because this library uses `encoding/json/v2`.

> Without the Nix shell, `go build` / `go test` fail with a "build constraints exclude all Go files in encoding/json/v2" error.

---

## Building & Testing

All commands assume you are inside `nix develop`:

```bash
# Build
go build ./...

# Run the full test suite
go test ./...

# Run with the race detector
go test -race ./...

# Run with coverage
go test -cover ./...

# Run a specific test
go test -run "TestName" ./...

# Verify the Nix flake evaluates
nix flake check --no-build
```

### Test conventions

- Tests use [Ginkgo](https://onsi.github.io/ginkgo/) / [Gomega](BDD style) with dot-imports (standard BDD pattern).
- Example tests (`Example*` functions) are rendered in godoc and double as executable documentation.
- Fuzz tests (`Fuzz*`) and benchmarks (`Benchmark*`) live in `fuzz_test.go` and `benchmark_test.go`.

---

## Code Standards

This is a **single-package Go library**: all source files live at the repository root (they ARE the public API). There is no `cmd/`, `internal/`, or `pkg/` layout.

### Mandatory rules

1. **Zero `any` in public APIs** — use generics (`OneOf[T]`, `Equals[T]`, `NotEmptySlice[T]`).
2. **Immutable rules** — `RuleImpl` is immutable; mutation returns a new value (`WithName`, `WithSeverity`, `WithMessage`).
3. **Parameter naming** — use `minimum`/`maximum` (not `min`/`max`) to avoid shadowing Go 1.21+ builtins.
4. **Early returns** — guard clauses over deep nesting.
5. **`Severity`** is a type alias for `finding.Severity` (a string). Do not introduce a parallel local severity type.

### Naming conventions

| Type | Convention | Example |
| --- | --- | --- |
| Packages | lowercase, single word | `businessrules` |
| Interfaces | PascalCase | `Rule` |
| Rule builders | PascalCase verb/noun | `NonNegative`, `NotEmptySlice` |
| Functions | PascalCase | `NewValidator` |
| Variables | camelCase | `severity` |

---

## Linting & Quality

```bash
# Run golangci-lint (config in .golangci.yml)
golangci-lint run --timeout 5m

# Format check (treefmt, configured in flake.nix)
nix build .#checks.x86_64-linux.format
```

### Quality gate (before merging)

- [ ] `go test -race ./...` passes
- [ ] `golangci-lint run` reports 0 issues
- [ ] `go vet ./...` passes
- [ ] `nix flake check --no-build` passes
- [ ] Code is formatted (`gofumpt` / `goimports`)

### Known false positives

Some linters (`branching-flow`, `hierarchical-errors`, `go-auto-upgrade`, `go-structure-linter`) report false positives specific to a validation library. These are documented in [AGENTS.md](AGENTS.md) — read that file before "fixing" a reported finding.

---

## Pull Request Process

### Branch naming

```
feat/description
fix/description
docs/description
refactor/description
test/description
```

### Review checklist

1. **Self-review first** — run the quality gate locally.
2. **Small PRs** — keep changes focused and digestible.
3. **Explain "why"** — not just "what", but rationale.
4. **Update docs** — if you add a rule builder, update `FEATURES.md`, `README.md`, `ROADMAP.md`, and `docs/DOMAIN_LANGUAGE.md`.

---

## Commit Messages

### Format

```
<type>(<scope>): <subject>

<body>
```

### Types

| Type | Description |
| --- | --- |
| feat | New feature (e.g. a new rule builder) |
| fix | Bug fix |
| docs | Documentation changes |
| style | Formatting, whitespace |
| refactor | Code restructuring |
| test | Adding/updating tests |
| chore | Build, tooling, CI |
| perf | Performance improvements |
| ci | CI/CD changes |

### Examples

```bash
# Good
feat(builders): add NotEmptySlice and NotEmptyMap collection rules

# Bad
added rules

# Good
docs(readme): correct Severity type documentation after finding.Severity migration
```

---

## Getting Help

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- Check [AGENTS.md](AGENTS.md) for non-obvious project context

---

_Thank you for contributing to go-business-rules!_
