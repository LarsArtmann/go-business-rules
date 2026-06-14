package businessrules

import "github.com/larsartmann/go-finding"

// Severity represents the importance level of a validation rule.
type Severity = finding.Severity

const (
	// SeverityInfo indicates informational issues that don't block processing.
	SeverityInfo = finding.SeverityInfo
	// SeverityWarning indicates issues that should be reviewed but don't block processing.
	SeverityWarning = finding.SeverityWarning
	// SeverityError indicates validation failures that should block processing.
	SeverityError = finding.SeverityError
	// SeverityCritical indicates severe failures requiring immediate attention.
	SeverityCritical = finding.SeverityCritical
)
