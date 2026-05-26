package businessrules

import (
	"fmt"

	"github.com/larsartmann/go-finding"
)

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

// severityName returns the human-readable name of the severity level.
func severityName(s Severity) string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return fmt.Sprintf("UNKNOWN(%s)", s)
	}
}
