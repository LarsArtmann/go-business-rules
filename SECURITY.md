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

## Scope

This library performs validation only: it executes caller-supplied check functions and produces severity-aware results. It performs no I/O, opens no network connections, and holds no secrets. Reports about caller-supplied check functions are out of scope.
