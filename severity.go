package businessrules

import "fmt"

// Severity represents the importance level of a validation rule.
// Higher values indicate more severe violations.
type Severity int

const (
	// SeverityInfo indicates informational issues that don't block processing.
	// Use for non-critical validations that provide helpful feedback.
	SeverityInfo Severity = iota

	// SeverityWarning indicates issues that should be reviewed but don't block processing.
	// Use for validations where the data might still be acceptable.
	SeverityWarning

	// SeverityError indicates validation failures that should block processing.
	// Use for critical validations where invalid data must not proceed.
	SeverityError

	// SeverityCritical indicates severe failures requiring immediate attention.
	// Use for validations where failure indicates a serious system or data problem.
	SeverityCritical
)

// String returns the human-readable name of the severity level.
func (s Severity) String() string {
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
		return fmt.Sprintf("UNKNOWN(%d)", int(s))
	}
}
