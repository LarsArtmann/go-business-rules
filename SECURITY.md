# Security Policy

## Supported versions

Only the latest tagged release of `github.com/LarsArtmann/go-business-rules/v2`
receives security fixes. The module path carries the `/v2` major-version suffix.

| Version | Supported                  |
| ------- | -------------------------- |
| v2.2.x  | Yes (current release line) |
| v2.1.x  | No                         |
| < v2.1  | No (pre-`/v2` module path) |

## Reporting a vulnerability

The repository is public. Please report suspected vulnerabilities through
GitHub's private vulnerability reporting on the
[Security tab](https://github.com/LarsArtmann/go-business-rules/security/advisories/new)
rather than in a public issue, so a fix can be prepared before disclosure.

Please include:

1. A description of the vulnerability and its impact
2. Reproduction steps or a proof-of-concept
3. Affected versions

You will receive an initial response within 7 days. Please do not disclose the issue publicly until a fix is released.

## Security Practices

- **Minimal dependency surface.** The only runtime dependency is
  `github.com/larsartmann/go-finding`; everything else is the Go standard
  library. Dependencies are kept current through the BuildFlow pipeline
  (`buildflow update`) rather than ad-hoc edits.
- **Automated vulnerability scanning.** `govulncheck` runs against every Go
  module and `gosec` runs on the root module in CI on every push. The Nix
  build closure is audited with `vulnix`, and secrets are scanned with
  `gitleaks`.
- **No secrets, no I/O.** The library executes caller-supplied check functions
  in-process: it performs no network I/O, no filesystem access, and holds no
  credentials, so its attack surface is limited to the validating code itself.
- **Hardened build environment.** Builds and tests run hermetically inside the
  Nix devshell with pinned toolchain versions, so supply-chain substitution of
  build tools is not possible locally or in CI.

## Scope

This library performs validation only: it executes caller-supplied check functions and produces severity-aware results. It performs no I/O, opens no network connections, and holds no secrets. Reports about caller-supplied check functions are out of scope.
