# Finding SDK — Unified Data Model Proposal

**Status:** Draft | **Date:** 2026-04-10

---

## Problem

Seven tools in the ecosystem each define their own issue/diagnostic/finding types with overlapping but incompatible models:

| Project | Finding Type | Severity | Position | Fix Model | Output |
|---------|-------------|----------|----------|-----------|--------|
| **golangci-lint-auto-configure** | `LinterRecommendation`, `ValidationError` | `LinterPriority` (Critical/High/Medium/Optional) | Config path only (no source lines) | `AutoFix bool` on linter metadata | JSON, HTML |
| **BuildFlow** | `BinaryViolation`, `TODOViolation`, `HierarchicalErrorsViolation`, etc. | `ViolationSeverity` (Critical/High/Medium/Low/Info) | `GetFile()`, `GetLine()`, `GetColumn()` | `AutoFix bool` in config; `Suggestion string` on violations | JSON, SARIF, HTML |
| **rules** | `analysis.Analyzer` diagnostics (no custom type) | None (delegated to host framework) | `token.Pos` via `go/analysis` | None (report-only) | Text (via `go vet`) |
| **art-dupl** | `Clone`, `CloneGroup` | `CloneSeverity` (low/medium/high/critical) | `LineNumber`, `BytePosition`, `Filename` (no column) | None; `Suggestion` text only | JSON, SARIF, HTML, CSV, text |
| **branching-flow** | `StrongIDViolation`, `BoolBlindnessViolation`, `PrimitiveTypeViolation`, `DuplicateGroup`, `Detection` | `Severity` (critical/high/medium/low) | `SourceLocation` (file, line, column) | Rich: `StrongIDSuggestion{BeforeCode, AfterCode}`, `BitFlagSuggestion`, `EnumSuggestion`, `CompositionSuggestion`; `--fix` flag | JSON, SARIF, HTML, Markdown, text |
| **go-auto-upgrade** | `Change`, `Warning` | None (implicit: change vs warning vs error) | `PathString`, `LineInt` (no column) | All changes are auto-fixable via AST rewrite; `Result.Content` has new code | Text (slog) |
| **hierarchical-errors** | `ErrorViolation`, `ErrorFlow`, `ErrorHierarchy` | `Severity` (low/medium/high) | `token.Position` (file, line, column, offset) | `Suggestion string`; SARIF `Fix` structs (descriptive only) | JSON, SARIF, HTML, DOT, Mermaid, agent, text |

**Result:** No tool can consume another tool's findings. No unified reporting. No way to correlate findings across tools. Each tool reinvents serialization, severity mapping, and SARIF generation.

---

## Proposal: `finding` — A Common Finding SDK

A zero-dependency Go library defining a **universal finding data model** that every tool can produce and consume.

### Design Principles

1. **SARIF-aligned** — SARIF 2.1.0 is the most comprehensive existing standard. Our model maps 1:1 to SARIF but is simpler and Go-native
2. **Minimal core, extensible perimeter** — Every field except `Rule`, `Message`, and `Position` is optional
3. **Fix-strategy aware** — Explicitly models whether a fix is available, and whether it's deterministic or AI-assisted
4. **Composable** — Findings can be grouped, correlated, and merged across tools
5. **Zero dependencies** — stdlib only, consistent with all projects
6. **Round-trippable** — Go struct → JSON → Go struct without loss

---

## Core Types

### Severity

```go
type Severity string

const (
    SeverityInfo     Severity = "info"
    SeverityWarning  Severity = "warning"
    SeverityError    Severity = "error"
    SeverityCritical Severity = "critical"
)
```

**Mapping from existing projects:**

| Project | Native | Maps To |
|---------|--------|---------|
| BuildFlow `ViolationSeverityCritical` | Critical | `critical` |
| BuildFlow `ViolationSeverityHigh` | High | `error` |
| BuildFlow `ViolationSeverityMedium` | Medium | `warning` |
| BuildFlow `ViolationSeverityLow` | Low | `info` |
| BuildFlow `ViolationSeverityInfo` | Info | `info` |
| art-dupl `CloneSeverityCritical` | critical | `critical` |
| art-dupl `CloneSeverityHigh` | high | `error` |
| art-dupl `CloneSeverityMedium` | medium | `warning` |
| art-dupl `CloneSeverityLow` | low | `info` |
| branching-flow `SeverityCritical` | critical | `critical` |
| branching-flow `SeverityHigh` | high | `error` |
| branching-flow `SeverityMedium` | medium | `warning` |
| branching-flow `SeverityLow` | low | `info` |
| hierarchical-errors `SeverityHigh` | high | `error` |
| hierarchical-errors `SeverityMedium` | medium | `warning` |
| hierarchical-errors `SeverityLow` | low | `info` |
| go-auto-upgrade `Warning` | (implicit) | `warning` |
| go-auto-upgrade `Change` | (implicit) | `info` |
| go-auto-upgrade `Error` | (implicit) | `error` |

### FixStrategy

```go
type FixStrategy string

const (
    FixStrategyNone       FixStrategy = "none"       // No fix available
    FixStrategySuggest    FixStrategy = "suggest"    // Human-readable suggestion, not machine-applicable
    FixStrategyDirect     FixStrategy = "direct"     // Deterministic code transformation (AST rewrite, formatter, etc.)
    FixStrategyAI         FixStrategy = "ai"         // Requires AI/LLM to generate context-aware fix
)
```

**Mapping from existing projects:**

| Project | Scenario | FixStrategy |
|---------|----------|-------------|
| art-dupl | Clone detected | `none` |
| art-dupl | `Suggestion` text present | `suggest` |
| branching-flow | `StrongIDSuggestion{BeforeCode, AfterCode}` | `suggest` |
| branching-flow | `--fix` flag + deterministic rewrite | `direct` |
| BuildFlow | `AutoFix=true` + formatter step | `direct` |
| BuildFlow | `Suggestion string` on violation | `suggest` |
| go-auto-upgrade | `Result.Content` with AST rewrite | `direct` |
| go-auto-upgrade | `Warning` (manual review needed) | `suggest` |
| hierarchical-errors | `Suggestion string` on ErrorViolation | `suggest` |
| golangci-lint-auto-configure | `AutoFix bool` on linter | `direct` (via `golangci-lint --fix`) |
| rules | No fixes | `none` |

### Position

```go
type Position struct {
    File   string `json:"file"`
    Line   int    `json:"line,omitempty"`   // 1-based; 0 = not set
    Column int    `json:"column,omitempty"` // 1-based; 0 = not set
    Offset int    `json:"offset,omitempty"` // byte offset; 0 = not set
}

type Range struct {
    Start Position `json:"start"`
    End   Position `json:"end,omitempty"`
}
```

**Design notes:**
- `Line` and `Column` are 1-based (consistent with Go's `token.Position`, editors, SARIF)
- `Offset` is 0-based byte offset (for tools like art-dupl that work with byte positions)
- `Range` supports tools that report spans (art-dupl clones, branching-flow patterns)
- When only `File` is set → file-level finding (golangci-lint-auto-configure)
- When `File` + `Line` → line-level finding (go-auto-upgrade, art-dupl)
- When `File` + `Line` + `Column` → precise finding (hierarchical-errors, branching-flow)

### Finding

The central type:

```go
type Finding struct {
    // Identity
    ID       string `json:"id"`                // Stable, unique identifier (e.g., "branching-flow:STRONG_ID:pkg/types.go:42:5")
    Rule     string `json:"rule"`              // Rule/check name (e.g., "STRONG_ID", "clone-detected", "silent-swallow")
    ToolName string `json:"toolName"`          // Source tool name (e.g., "branching-flow", "art-dupl")

    // Core
    Message  string   `json:"message"`          // Human-readable description
    Severity Severity `json:"severity"`         // info, warning, error, critical
    Position Position `json:"position"`         // Where the issue is

    // Classification
    Category  string `json:"category,omitempty"`  // Domain: "security", "style", "performance", "duplication", "error-handling", "migration", etc.
    Tag       string `json:"tag,omitempty"`       // Sub-classification: "phantom-type", "bool-blindness", "clone", "panic", etc.

    // Fix
    FixStrategy FixStrategy `json:"fixStrategy"`              // none, suggest, direct, ai
    Suggestion  string      `json:"suggestion,omitempty"`     // Human-readable fix description
    BeforeCode  string      `json:"beforeCode,omitempty"`     // Code before the fix
    AfterCode   string      `json:"afterCode,omitempty"`      // Code after the fix (for direct fixes, this IS the fix)

    // Context
    Range      Range           `json:"range,omitempty"`      // For span-based findings (clones, selections)
    Snippet    string          `json:"snippet,omitempty"`    // Surrounding code context
    Confidence float64         `json:"confidence,omitempty"` // 0.0-1.0 (art-dupl, branching-flow composition)
    Related    []RelatedRef    `json:"related,omitempty"`    // Related findings (clone groups, error flows)

    // Extensibility
    Metadata map[string]string `json:"metadata,omitempty"` // Tool-specific key-value pairs
}
```

### RelatedRef

For linking findings (clone groups, error flows, duplicate types):

```go
type RelatedRef struct {
    FindingID string `json:"findingId"`           // ID of the related finding
    Relation  string `json:"relation"`            // "clone-of", "wraps", "duplicates", "causes", "fixes"
    Position  Position `json:"position,omitempty"` // Quick access to the related location
}
```

### Report

Top-level container for a tool run:

```go
type Report struct {
    Tool     ToolInfo  `json:"tool"`
    Findings []Finding `json:"findings"`
    Summary  Summary   `json:"summary"`
}

type ToolInfo struct {
    Name    string `json:"name"`
    Version string `json:"version,omitempty"`
}

type Summary struct {
    Total         int               `json:"total"`
    BySeverity    map[Severity]int  `json:"bySeverity"`
    ByCategory    map[string]int    `json:"byCategory,omitempty"`
    ByFixStrategy map[FixStrategy]int `json:"byFixStrategy,omitempty"`
    FilesAffected int               `json:"filesAffected,omitempty"`
    DurationMs    int64             `json:"durationMs,omitempty"`
}
```

---

## File Structure

```
finding/
├── finding.go          # Finding, Position, Range, RelatedRef types
├── severity.go         # Severity enum, validation, ordering
├── fix_strategy.go     # FixStrategy enum
├── report.go           # Report, ToolInfo, Summary types
├── id.go               # Finding ID generation (tool:rule:file:line:col)
├── json.go             # JSON marshal/unmarshal helpers
├── sarif.go            # Report → SARIF 2.1.0 conversion
├── filter.go           # Query/filter functions (BySeverity, ByCategory, ByFixStrategy, etc.)
├── merge.go            # Merge multiple Reports (dedup, correlate)
├── category.go         # Standard category constants
├── converters/
│   ├── artdupl.go      # art-dupl → Finding converter
│   ├── buildflow.go    # BuildFlow violations → Finding converter
│   ├── branching.go    # branching-flow → Finding converter
│   ├── hiererrors.go   # hierarchical-errors → Finding converter
│   ├── goupgrade.go    # go-auto-upgrade → Finding converter
│   └── golintlintercfg.go # golangci-lint-auto-configure → Finding converter
└── finding_test.go
```

---

## Converter Examples

### art-dupl → Finding

```go
func FromArtDuplClone(group *artdupl.CloneGroup, clone *artdupl.Clone, toolVer string) Finding {
    return Finding{
        ID:       fmt.Sprintf("art-dupl:clone:%s:%d", clone.Filename, clone.StartLine),
        Rule:     "clone-detected",
        ToolName: "art-dupl",
        Message:  fmt.Sprintf("Duplicate code (%d tokens)", group.Size),
        Severity: cloneSeverityToCommon(group.Severity),
        Position: Position{File: clone.Filename, Line: clone.StartLine},
        Category: "duplication",
        Tag:      "clone",
        Range: Range{
            Start: Position{File: clone.Filename, Line: clone.StartLine},
            End:   Position{File: clone.Filename, Line: clone.EndLine},
        },
        FixStrategy: FixStrategySuggest,
        Confidence:  cloneComplexityToConfidence(clone.Complexity),
        Related:     cloneGroupToRelated(group, clone),
        Metadata:    map[string]string{"hash": clone.Hash, "tokens": fmt.Sprint(group.Size)},
    }
}
```

### branching-flow → Finding

```go
func FromStrongIDViolation(v core.StrongIDViolation, toolVer string) Finding {
    finding := Finding{
        ID:          fmt.Sprintf("branching-flow:STRONG_ID:%s:%d:%d", v.Location.FilePath(), v.Location.Line(), v.Location.Column()),
        Rule:        "STRONG_ID",
        ToolName:    "branching-flow",
        Message:     v.Message,
        Severity:    bfSeverityToCommon(v.Severity),
        Position:    Position{File: v.Location.FilePath(), Line: v.Location.Line(), Column: v.Location.Column()},
        Category:    "type-safety",
        Tag:         "phantom-type",
        FixStrategy: FixStrategySuggest,
        Suggestion:  v.Suggestion.ImportPath,
        BeforeCode:  v.Suggestion.BeforeCode,
        AfterCode:   v.Suggestion.AfterCode,
    }
    return finding
}
```

### go-auto-upgrade → Finding

```go
func FromMigrationResult(r migrator.Result, change migrator.Change, toolVer string) Finding {
    return Finding{
        ID:          fmt.Sprintf("go-auto-upgrade:%s:%s:%d", change.Type, r.Path, change.Line),
        Rule:        string(change.Type),
        ToolName:    "go-auto-upgrade",
        Message:     string(change.Description),
        Severity:    SeverityInfo,
        Position:    Position{File: string(r.Path), Line: int(change.Line)},
        Category:    "migration",
        FixStrategy: FixStrategyDirect,
        AfterCode:   string(r.Content), // The full modified file content
    }
}

func FromMigrationWarning(r migrator.Result, w migrator.Warning, toolVer string) Finding {
    return Finding{
        ID:          fmt.Sprintf("go-auto-upgrade:warning:%s:%d", r.Path, w.Line),
        Rule:        "manual-review",
        ToolName:    "go-auto-upgrade",
        Message:     string(w.Message),
        Severity:    SeverityWarning,
        Position:    Position{File: string(r.Path), Line: int(w.Line)},
        Category:    "migration",
        FixStrategy: FixStrategySuggest,
    }
}
```

### hierarchical-errors → Finding

```go
func FromErrorViolation(v types.ErrorViolation, toolVer string) Finding {
    return Finding{
        ID:          fmt.Sprintf("hierarchical-errors:%s:%s:%d:%d", v.Type, v.Position.Filename, v.Position.Line, v.Position.Column),
        Rule:        string(v.Type),
        ToolName:    "hierarchical-errors",
        Message:     v.Message,
        Severity:    heSeverityToCommon(v.Severity),
        Position:    Position{File: v.Position.Filename, Line: v.Position.Line, Column: v.Position.Column, Offset: v.Position.Offset},
        Category:    "error-handling",
        Tag:         string(v.Type),
        FixStrategy: FixStrategySuggest,
        Suggestion:  v.Suggestion,
        Snippet:     v.CodeSnippet,
    }
}
```

---

## Standard Categories

```go
const (
    CategorySecurity     = "security"
    CategoryStyle        = "style"
    CategoryPerformance  = "performance"
    CategoryCorrectness  = "correctness"
    CategoryComplexity   = "complexity"
    CategoryDuplication  = "duplication"
    CategoryErrorHandling = "error-handling"
    CategoryMigration    = "migration"
    CategoryTypeSafety   = "type-safety"
    CategoryStructure    = "structure"
    CategoryConfiguration = "configuration"
    CategoryDocumentation = "documentation"
    CategoryTesting      = "testing"
)
```

---

## SARIF Mapping

Every `Finding` maps directly to SARIF 2.1.0:

| Finding field | SARIF path |
|--------------|-----------|
| `Rule` | `result.ruleId` |
| `Message` | `result.message.text` |
| `Severity` → SARIF level | `result.level` (info→note, warning→warning, error→error, critical→error) |
| `Position.File` | `result.locations[0].physicalLocation.artifactLocation.uri` |
| `Position.Line` | `result.locations[0].physicalLocation.region.startLine` |
| `Position.Column` | `result.locations[0].physicalLocation.region.startColumn` |
| `Range.End.Line` | `result.locations[0].physicalLocation.region.endLine` |
| `Range.End.Column` | `result.locations[0].physicalLocation.region.endColumn` |
| `Suggestion` | `result.fixes[0].description.text` |
| `BeforeCode`/`AfterCode` | `result.fixes[0].artifactChanges[0].replacements[0]` |
| `Related` | `result.relatedLocations` |
| `Metadata` | `result.properties` |
| `Confidence` | `result.rank` (0.0-100.0) |
| `ToolName`/`Version` | `run.tool.driver.name`/`version` |
| `Category` | `rule.properties.category` |

---

## Integration Points

### How Each Tool Would Use This

| Tool | Integration |
|------|------------|
| **art-dupl** | Emit `Finding` per clone, linked via `Related` for clone groups |
| **branching-flow** | Emit `Finding` per violation with `BeforeCode`/`AfterCode` suggestions |
| **BuildFlow** | Replace `PrioritizedViolation` interface with `Finding`; adapters convert external tool output to `Finding` |
| **go-auto-upgrade** | Emit `Finding` per change (direct fix) and per warning (suggest fix) |
| **hierarchical-errors** | Emit `Finding` per violation; `ErrorFlow` maps to `Related` chains |
| **golangci-lint-auto-configure** | Emit `Finding` per recommendation at config level |
| **rules** | Future: custom analyzers could emit `Finding` alongside standard diagnostics |

### Consumer Use Cases

1. **Unified report** — Run all tools, merge `Report`s, output single SARIF/HTML/JSON
2. **Auto-fix pipeline** — Filter `FixStrategy == "direct"`, apply all direct fixes, re-run
3. **AI fix pipeline** — Filter `FixStrategy == "suggest"`, send to AI with context, apply
4. **Quality dashboard** — Aggregate `Summary` across tools, track trends over time
5. **IDE integration** — Consume `Report` via LSP, show diagnostics from all tools in one view

---

## Open Questions

1. **Repository name**: `finding`? `finding-sdk`? `go-finding`?
2. **Stable ID format**: Should IDs be deterministic hashes or readable strings? Proposal uses readable `"tool:rule:file:line:col"` but hash-based IDs might be more robust for deduplication
3. **Converter placement**: Should converters live in the SDK (requiring dependency on all tools) or in each tool (requiring dependency on the SDK)? Recommendation: **in each tool** — tools depend on the SDK, not the other way around
4. **Breaking change from existing types**: Should BuildFlow replace `PrioritizedViolation` or wrap it? Recommendation: **migrate incrementally** — `Finding` becomes the canonical type, old types get converter methods
5. **AI fix metadata**: Should `FixStrategyAI` have additional fields (model, prompt template, context window)? Recommendation: start simple with `Metadata` key-value, add structured fields if needed

---

## Implementation Roadmap

### Phase 1: Core SDK (1-2 days)
- `finding.go`, `severity.go`, `fix_strategy.go`, `report.go`, `id.go`
- `filter.go` (query helpers)
- `json.go` (serialization)
- Tests

### Phase 2: SARIF Output (1 day)
- `sarif.go` (Report → SARIF 2.1.0)
- Tests with SARIF schema validation

### Phase 3: Converters (1-2 days)
- One converter per tool
- Integration tests using real tool output

### Phase 4: Integration (2-3 days)
- Integrate into BuildFlow as primary output type
- Add `--format finding` flag to other tools
- Unified report generation

### Phase 5: Advanced (optional)
- `merge.go` (deduplication, correlation)
- AI fix pipeline integration
- LSP server for unified diagnostics
