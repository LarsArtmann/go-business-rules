package businessrules

import (
	"encoding/json"
	"fmt"
)

// ValidationResult contains the outcome of validating multiple rules.
// It provides methods to filter and check violations by severity.
type ValidationResult struct {
	// Valid indicates whether all rules passed (no violations).
	Valid bool

	// Violations contains all rule violations detected during validation.
	Violations []Violation
}

// Errors returns all violations with Error or Critical severity.
func (r ValidationResult) Errors() []Violation {
	return r.BySeverity(SeverityError, SeverityCritical)
}

// Warnings returns all violations with Warning severity.
func (r ValidationResult) Warnings() []Violation {
	return r.BySeverity(SeverityWarning)
}

// Info returns all violations with Info severity.
func (r ValidationResult) Info() []Violation {
	return r.BySeverity(SeverityInfo)
}

// Critical returns all violations with Critical severity.
func (r ValidationResult) Critical() []Violation {
	return r.BySeverity(SeverityCritical)
}

// BySeverity returns violations matching any of the specified severity levels.
func (r ValidationResult) BySeverity(severities ...Severity) []Violation {
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
func (r ValidationResult) HasErrors() bool {
	return len(r.Errors()) > 0
}

// HasWarnings returns true if there are any Warning violations.
func (r ValidationResult) HasWarnings() bool {
	return len(r.Warnings()) > 0
}

// HasCritical returns true if there are any Critical violations.
func (r ValidationResult) HasCritical() bool {
	return len(r.Critical()) > 0
}

// HasInfo returns true if there are any Info violations.
func (r ValidationResult) HasInfo() bool {
	return len(r.Info()) > 0
}

// Count returns the total number of violations.
func (r ValidationResult) Count() int {
	return len(r.Violations)
}

// FirstError returns the first violation with Error or Critical severity.
// Returns an empty Violation if no errors exist.
func (r ValidationResult) FirstError() Violation {
	errors := r.Errors()
	if len(errors) == 0 {
		return Violation{}
	}
	return errors[0]
}

// FirstCritical returns the first Critical severity violation.
// Returns an empty Violation if no critical violations exist.
func (r ValidationResult) FirstCritical() Violation {
	critical := r.Critical()
	if len(critical) == 0 {
		return Violation{}
	}
	return critical[0]
}

// FirstWarning returns the first Warning severity violation.
// Returns an empty Violation if no warnings exist.
func (r ValidationResult) FirstWarning() Violation {
	warnings := r.Warnings()
	if len(warnings) == 0 {
		return Violation{}
	}
	return warnings[0]
}

// FirstInfo returns the first Info severity violation.
// Returns an empty Violation if no info violations exist.
func (r ValidationResult) FirstInfo() Violation {
	info := r.Info()
	if len(info) == 0 {
		return Violation{}
	}
	return info[0]
}

// ForEach calls fn for each violation in the result.
func (r ValidationResult) ForEach(fn func(Violation)) {
	for _, v := range r.Violations {
		fn(v)
	}
}

// Filter returns violations that match the predicate.
// Use for custom filtering beyond severity-based methods.
func (r ValidationResult) Filter(predicate func(Violation) bool) []Violation {
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
func (r ValidationResult) Merge(other ValidationResult) ValidationResult {
	violations := make([]Violation, 0, len(r.Violations)+len(other.Violations))
	violations = append(violations, r.Violations...)
	violations = append(violations, other.Violations...)

	return ValidationResult{
		Valid:      r.Valid && other.Valid,
		Violations: violations,
	}
}

// MarshalJSON implements json.Marshaler for ValidationResult.
func (r ValidationResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Valid      bool        `json:"valid"`
		Violations []Violation `json:"violations,omitempty"`
	}{
		Valid:      r.Valid,
		Violations: r.Violations,
	})
}

// Error implements the error interface for ValidationResult.
// Returns a summary of all violations or nil if valid.
func (r ValidationResult) Error() string {
	if r.Valid {
		return ""
	}
	if len(r.Violations) == 0 {
		return "validation failed"
	}
	if len(r.Violations) == 1 {
		return r.Violations[0].Error()
	}
	return fmt.Sprintf(
		"validation failed with %d violations: %s",
		len(r.Violations),
		r.Violations[0].Error(),
	)
}
