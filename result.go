package businessrules

import "encoding/json"

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

// Count returns the total number of violations.
func (r Result) Count() int {
	return len(r.Violations)
}

// FirstError returns the first violation with Error or Critical severity.
// Returns an empty Violation if no errors exist.
func (r Result) FirstError() Violation {
	errors := r.Errors()
	if len(errors) == 0 {
		return Violation{}
	}
	return errors[0]
}

// FirstCritical returns the first Critical severity violation.
// Returns an empty Violation if no critical violations exist.
func (r Result) FirstCritical() Violation {
	critical := r.Critical()
	if len(critical) == 0 {
		return Violation{}
	}
	return critical[0]
}

// FirstWarning returns the first Warning severity violation.
// Returns an empty Violation if no warnings exist.
func (r Result) FirstWarning() Violation {
	warnings := r.Warnings()
	if len(warnings) == 0 {
		return Violation{}
	}
	return warnings[0]
}

// FirstInfo returns the first Info severity violation.
// Returns an empty Violation if no info violations exist.
func (r Result) FirstInfo() Violation {
	info := r.Info()
	if len(info) == 0 {
		return Violation{}
	}
	return info[0]
}

// ForEach calls fn for each violation in the result.
func (r Result) ForEach(fn func(Violation)) {
	for _, v := range r.Violations {
		fn(v)
	}
}

// Filter returns violations that match the predicate.
// Use for custom filtering beyond severity-based methods.
func (r Result) Filter(predicate func(Violation) bool) []Violation {
	var result []Violation
	for _, v := range r.Violations {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Merge combines two results into a new result.
// The merged result is valid only if both inputs are valid.
func (r Result) Merge(other Result) Result {
	violations := make([]Violation, 0, len(r.Violations)+len(other.Violations))
	violations = append(violations, r.Violations...)
	violations = append(violations, other.Violations...)

	return Result{
		Valid:      r.Valid && other.Valid,
		Violations: violations,
	}
}

// MarshalJSON implements json.Marshaler for Result.
func (r Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Valid      bool        `json:"valid"`
		Violations []Violation `json:"violations,omitempty"`
	}{
		Valid:      r.Valid,
		Violations: r.Violations,
	})
}
