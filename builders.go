package businessrules

import (
	"fmt"
	"regexp"
)

// ============================================================================
// Numeric Rules
// ============================================================================

// NonNegative creates a rule that validates value >= 0.
// Use for validating non-negative numeric values like ages, quantities, or prices.
func NonNegative(name string, value float64, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value < 0 {
				return fmt.Errorf("%s must be non-negative, got %f", name, value)
			}
			return nil
		},
		severity,
		name+" must be non-negative",
	)
}

// Positive creates a rule that validates value > 0.
// Use for validating positive numeric values like counts or amounts.
func Positive(name string, value float64, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value <= 0 {
				return fmt.Errorf("%s must be positive, got %f", name, value)
			}
			return nil
		},
		severity,
		name+" must be positive",
	)
}

// InRange creates a rule that validates min <= value <= max.
// Use for validating numeric values within a specific range.
func InRange(name string, value, min, max float64, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value < min || value > max {
				return fmt.Errorf("%s must be between %f and %f, got %f", name, min, max, value)
			}
			return nil
		},
		severity,
		name+" must be within range",
	)
}

// MinInt creates a rule that validates value >= min.
// Use for validating integer values meet a minimum threshold.
func MinInt(name string, value, min int, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value < min {
				return fmt.Errorf("%s must be at least %d, got %d", name, min, value)
			}
			return nil
		},
		severity,
		name+" must meet minimum",
	)
}

// MaxInt creates a rule that validates value <= max.
// Use for validating integer values don't exceed a maximum.
func MaxInt(name string, value, max int, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value > max {
				return fmt.Errorf("%s must be at most %d, got %d", name, max, value)
			}
			return nil
		},
		severity,
		name+" must not exceed maximum",
	)
}

// ============================================================================
// String Rules
// ============================================================================

// NotEmpty creates a rule that validates the string is not empty.
// Use for validating required string fields.
func NotEmpty(name string, value string, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value == "" {
				return fmt.Errorf("%s must not be empty", name)
			}
			return nil
		},
		severity,
		name+" must not be empty",
	)
}

// MinLength creates a rule that validates len(value) >= min.
// Use for validating minimum string length requirements.
func MinLength(name string, value string, min int, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if len(value) < min {
				return fmt.Errorf("%s must be at least %d characters, got %d", name, min, len(value))
			}
			return nil
		},
		severity,
		name+" must meet minimum length",
	)
}

// MaxLength creates a rule that validates len(value) <= max.
// Use for validating maximum string length constraints.
func MaxLength(name string, value string, max int, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if len(value) > max {
				return fmt.Errorf("%s must not exceed %d characters, got %d", name, max, len(value))
			}
			return nil
		},
		severity,
		name+" must not exceed maximum length",
	)
}

// Matches creates a rule that validates the string matches a regex pattern.
// Use for validating strings against custom patterns.
func Matches(name string, value string, pattern *regexp.Regexp, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if !pattern.MatchString(value) {
				return fmt.Errorf("%s must match pattern %s, got %q", name, pattern.String(), value)
			}
			return nil
		},
		severity,
		name+" must match required pattern",
	)
}


