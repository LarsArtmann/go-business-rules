# Security Policy

## Supported versions

Only the latest tagged release of `github.com/LarsArtmann/go-business-rules/v2` receives security fixes.

## Reporting a vulnerability

The repository is private; report suspected vulnerabilities directly to the owner via GitHub (mention `@LarsArtmann` in a private maintainer-facing issue or contact via the email on the GitHub profile).

Please include:

1. A description of the vulnerability and its impact
2. Reproduction steps or a proof-of-concept
3. Affected versions

You will receive an initial response within 7 days. Please do not disclose the issue publicly until a fix is released.

## Scope

This library performs validation only: it executes caller-supplied check functions and produces severity-aware results. It performs no I/O, opens no network connections, and holds no secrets. Reports about caller-supplied check functions are out of scope.
