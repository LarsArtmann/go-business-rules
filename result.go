package businessrules

// Result contains the outcome of validating multiple rules.
// It provides methods to filter and check violations by severity.
type Result struct {
	// Valid indicates whether all rules passed (no violations).
	Valid bool

	// Violations contains all rule violations detected during validation.
	Violations []Violation
}

// Errors returns all violations with Error or Critical severity.
func (r Result) Errors() []Violation {
	return r.BySeverity(SeverityError, SeverityCritical)
}

// Warnings returns all violations with Warning severity.
func (r Result) Warnings() []Violation {
	return r.BySeverity(SeverityWarning)
}

// Info returns all violations with Info severity.
func (r Result) Info() []Violation {
	return r.BySeverity(SeverityInfo)
}

// Critical returns all violations with Critical severity.
func (r Result) Critical() []Violation {
	return r.BySeverity(SeverityCritical)
}

// BySeverity returns violations matching any of the specified severity levels.
func (r Result) BySeverity(severities ...Severity) []Violation {
	severitySet := make(map[Severity]bool, len(severities))
	for _, s := range severities {
		severitySet[s] = true
	}

	var result []Violation
	for _, v := range r.Violations {
		if severitySet[v.Rule.Severity()] {
			result = append(result, v)
		}
	}
	return result
}

// HasErrors returns true if there are any Error or Critical violations.
func (r Result) HasErrors() bool {
	return len(r.Errors()) > 0
}

// HasWarnings returns true if there are any Warning violations.
func (r Result) HasWarnings() bool {
	return len(r.Warnings()) > 0
}

// HasCritical returns true if there are any Critical violations.
func (r Result) HasCritical() bool {
	return len(r.Critical()) > 0
}

// HasInfo returns true if there are any Info violations.
func (r Result) HasInfo() bool {
	return len(r.Info()) > 0
}
