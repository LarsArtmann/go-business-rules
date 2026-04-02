# go-output Integration Decision

**Date:** 2026-03-29
**Status:** Reviewed — No integration planned

## Summary

`go-output` is not a dependency of `go-business-rules`. The libraries serve distinct purposes and should remain separate.

## Library Comparison

|                  | `go-business-rules`             | `go-output`                            |
| ---------------- | ------------------------------- | -------------------------------------- |
| **Purpose**      | Validation with severity levels | Structured data output (12 formats)    |
| **Runtime deps** | Zero (stdlib only)              | `lipgloss`, `go-faster/yaml`           |
| **Focus**        | Pass/fail + degree of failure   | Tables, JSON, CSV, trees, graphs, etc. |

## Decision

**No integration.** `go-business-rules` is a validation library with no display layer.

The `fmt.Sprintf` calls in `errors.go` and `validation_result.go` provide sufficient error reporting for library consumers.

## When Integration Would Make Sense

If a future consumer of `go-business-rules` wants formatted output, they can compose the two libraries at the application layer:

```go
import (
    "github.com/artmann/businessrules"
    "github.com/artmann/go-output"
)

result := validator.Validate(data)
table := output.NewMarkdownTable()
table.SetHeaders([]string{"Rule", "Severity", "Message"})
for _, v := range result.Violations() {
    table.AddRow([]string{v.Rule(), v.Severity().String(), v.Message()})
}
table.Render()
```

This keeps `go-business-rules` dependency-free and lets consumers choose their output format.

## When Integration Would NOT Make Sense

- **Adding `go-output` as a dependency** — Forces `lipgloss` and `yaml` on all consumers, even those who only need validation logic
- **Forking `go-output` code into `go-business-rules`** — Duplication without benefit; the code already exists in its canonical location
- **Creating a separate `go-business-rules/output` sub-package** — Premature abstraction; no consumer has requested this

## Conclusion

Keep concerns separated. Validation lives in `go-business-rules`. Output formatting lives in `go-output`. Compose at the application layer.
