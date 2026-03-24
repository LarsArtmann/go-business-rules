package businessrules

import (
	"encoding/json"
	"fmt"
)

// ValidationResultError contains the outcome of validating multiple rules.
// It provides methods to filter and check violations by severity.
type ValidationResultError struct {
	// Valid indicates whether all rules passed (no violations).
	Valid bool

	// ViolationErrors contains all rule violations detected during validation.
	ViolationErrors []ViolationError
}

// Errors returns all violations with Error or Critical severity.
func (r ValidationResultError) Errors() []ViolationError {
	return r.BySeverity(SeverityError, SeverityCritical)
}

// Warnings returns all violations with Warning severity.
func (r ValidationResultError) Warnings() []ViolationError {
	return r.BySeverity(SeverityWarning)
}

// Info returns all violations with Info severity.
func (r ValidationResultError) Info() []ViolationError {
	return r.BySeverity(SeverityInfo)
}

// Critical returns all violations with Critical severity.
func (r ValidationResultError) Critical() []ViolationError {
	return r.BySeverity(SeverityCritical)
}

// BySeverity returns violations matching any of the specified severity levels.
func (r ValidationResultError) BySeverity(severities ...Severity) []ViolationError {
	severitySet := make(map[Severity]bool, len(severities))
	for _, s := range severities {
		severitySet[s] = true
	}

	var result []ViolationError
	for _, v := range r.ViolationErrors {
		if severitySet[v.Rule.Severity()] {
			result = append(result, v)
		}
	}
	return result
}

// HasErrors returns true if there are any Error or Critical violations.
func (r ValidationResultError) HasErrors() bool {
	return len(r.Errors()) > 0
}

// HasWarnings returns true if there are any Warning violations.
func (r ValidationResultError) HasWarnings() bool {
	return len(r.Warnings()) > 0
}

// HasCritical returns true if there are any Critical violations.
func (r ValidationResultError) HasCritical() bool {
	return len(r.Critical()) > 0
}

// HasInfo returns true if there are any Info violations.
func (r ValidationResultError) HasInfo() bool {
	return len(r.Info()) > 0
}

// Count returns the total number of violations.
func (r ValidationResultError) Count() int {
	return len(r.ViolationErrors)
}

// FirstError returns the first violation with Error or Critical severity.
// Returns an empty ViolationError if no errors exist.
func (r ValidationResultError) FirstError() ViolationError {
	errors := r.Errors()
	if len(errors) == 0 {
		return ViolationError{} //nolint:exhaustruct
	}
	return errors[0]
}

// FirstCritical returns the first Critical severity violation.
// Returns an empty ViolationError if no critical violations exist.
func (r ValidationResultError) FirstCritical() ViolationError {
	critical := r.Critical()
	if len(critical) == 0 {
		return ViolationError{} //nolint:exhaustruct
	}
	return critical[0]
}

// FirstWarning returns the first Warning severity violation.
// Returns an empty ViolationError if no warnings exist.
func (r ValidationResultError) FirstWarning() ViolationError {
	warnings := r.Warnings()
	if len(warnings) == 0 {
		return ViolationError{} //nolint:exhaustruct
	}
	return warnings[0]
}

// FirstInfo returns the first Info severity violation.
// Returns an empty ViolationError if no info violations exist.
func (r ValidationResultError) FirstInfo() ViolationError {
	info := r.Info()
	if len(info) == 0 {
		return ViolationError{} //nolint:exhaustruct
	}
	return info[0]
}

// ForEach calls fn for each violation in the result.
func (r ValidationResultError) ForEach(fn func(ViolationError)) {
	for _, v := range r.ViolationErrors {
		fn(v)
	}
}

// Filter returns violations that match the predicate.
// Use for custom filtering beyond severity-based methods.
func (r ValidationResultError) Filter(predicate func(ViolationError) bool) []ViolationError {
	var result []ViolationError
	for _, v := range r.ViolationErrors {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Merge combines two results into a new result.
// The merged result is valid only if both inputs are valid.
func (r ValidationResultError) Merge(other ValidationResultError) ValidationResultError {
	violations := make([]ViolationError, 0, len(r.ViolationErrors)+len(other.ViolationErrors))
	violations = append(violations, r.ViolationErrors...)
	violations = append(violations, other.ViolationErrors...)

	return ValidationResultError{
		Valid:           r.Valid && other.Valid,
		ViolationErrors: violations,
	}
}

// MarshalJSON implements json.Marshaler for ValidationResultError.
func (r ValidationResultError) MarshalJSON() ([]byte, error) {
	marshaled, err := json.Marshal(struct {
		Valid           bool             `json:"valid"`
		ViolationErrors []ViolationError `json:"violations,omitempty"`
	}{
		Valid:           r.Valid,
		ViolationErrors: r.ViolationErrors,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ValidationResultError: %w", err)
	}
	return marshaled, nil
}

// Error implements the error interface for ValidationResultError.
// Returns a summary of all violations or nil if valid.
func (r ValidationResultError) Error() string {
	if r.Valid {
		return ""
	}
	if len(r.ViolationErrors) == 0 {
		return "validation failed"
	}
	if len(r.ViolationErrors) == 1 {
		return r.ViolationErrors[0].Error()
	}
	return fmt.Sprintf(
		"validation failed with %d violations: %s",
		len(r.ViolationErrors),
		r.ViolationErrors[0].Error(),
	)
}
