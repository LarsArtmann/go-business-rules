# TODO List — businessrules

> Actionable items for the next 2-4 weeks.

**Last Updated:** 2026-06-14

---

## Code Quality Hardening (Completed 2026-06-14)

- [x] Fix branching-flow BDD test expectations to match current state (12 PHANTOM violations, 33 total issues)
- [x] Remove unused `severityName` function from `severity.go`
- [x] Add `github.com/larsartmann/go-finding` to depguard allow list
- [x] Fix stale `github.com/artmann/businessrules` references in `.golangci.yml` (ireturn allow, goimports local-prefixes)
- [x] Add explicit `gci` section configuration matching goimports local-prefixes
- [x] Run `golangci-lint run ./... --fix` — 0 issues
- [x] Run `go test -race -cover ./...` — 94.8% coverage, all 145 specs pass
- [x] Run `branching-flow all .` — clean except documented false positives
- [x] Run `nix flake check` — all checks pass

## Phase 6: Integration

- [ ] Add as dependency to Polish-Customs
- [ ] Replace internal `validation.go` with import
- [ ] Run Polish-Customs tests to verify compatibility
- [ ] Commit migration

## Phase 7: Publish

- [ ] Tag release: `git tag v0.1.0`
- [ ] Push to remote: `git push origin master --tags`
- [ ] Verify on pkg.go.dev

## Success Criteria (Remaining)

- [ ] Polish-Customs successfully migrated
- [ ] CI/CD pipeline configured
