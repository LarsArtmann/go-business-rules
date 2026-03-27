package businessrules

import (
	"fmt"
	"regexp"
)

// Numeric Rules.

// thresholdCheck creates a threshold validation rule with a custom comparison.
// The shouldError function returns true when the value violates the threshold.
func thresholdCheck[T int | float64](
	name string,
	value, threshold T,
	severity Severity,
	shouldError func(value, threshold T) bool,
	errMsg string,
	templateMsg string,
) Rule {
	return NewRule(
		name,
		func() error {
			if shouldError(value, threshold) {
				return fmt.Errorf("%s %s %v, got %v", name, errMsg, threshold, value)
			}
			return nil
		},
		severity,
		name+" "+templateMsg,
	)
}

// numericCheck creates a numeric validation rule with a custom condition.
func numericCheck(
	name string,
	value float64,
	severity Severity,
	condition bool,
	errMsg string,
) Rule {
	return NewRule(
		name,
		func() error {
			if condition {
				return fmt.Errorf("%s %s, got %f", name, errMsg, value)
			}
			return nil
		},
		severity,
		name+" "+errMsg,
	)
}

// NonNegative creates a rule that validates value >= 0.
// Use for validating non-negative numeric values like ages, quantities, or prices.
func NonNegative(name string, value float64, severity Severity) Rule {
	return numericCheck(name, value, severity, value < 0, "must be non-negative")
}

// Positive creates a rule that validates value > 0.
// Use for validating positive numeric values like counts or amounts.
func Positive(name string, value float64, severity Severity) Rule {
	return numericCheck(name, value, severity, value <= 0, "must be positive")
}

// InRange creates a rule that validates minimum <= value <= maximum.
// Use for validating numeric values within a specific range.
func InRange(name string, value, minimum, maximum float64, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value < minimum || value > maximum {
				return fmt.Errorf(
					"%s must be between %f and %f, got %f",
					name,
					minimum,
					maximum,
					value,
				)
			}
			return nil
		},
		severity,
		name+" must be within range",
	)
}

// MinInt creates a rule that validates value >= minimum.
// Use for validating integer values meet a minimum threshold.
func MinInt(name string, value, minimum int, severity Severity) Rule {
	return thresholdCheck(name, value, minimum, severity,
		func(v, t int) bool { return v < t },
		"must be at least",
		"must meet minimum",
	)
}

// MaxInt creates a rule that validates value <= maximum.
// Use for validating integer values don't exceed a maximum.
func MaxInt(name string, value, maximum int, severity Severity) Rule {
	return thresholdCheck(name, value, maximum, severity,
		func(v, t int) bool { return v > t },
		"must be at most",
		"must not exceed maximum",
	)
}

// String Rules.

// NotEmpty creates a rule that validates the string is not empty.
// Use for validating required string fields.
func NotEmpty(name, value string, severity Severity) Rule {
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

// NotBlank creates a rule that validates the string is not blank (empty or whitespace-only).
// Use for validating required string fields that should have visible content.
func NotBlank(name, value string, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if len(value) == 0 {
				return fmt.Errorf("%s must not be blank", name)
			}
			for _, r := range value {
				if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
					return nil
				}
			}
			return fmt.Errorf("%s must not be blank (whitespace-only)", name)
		},
		severity,
		name+" must not be blank",
	)
}

// MinLength creates a rule that validates len(value) >= minimum.
// Use for validating minimum string length requirements.
func MinLength(name, value string, minimum int, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if len(value) < minimum {
				return fmt.Errorf(
					"%s must be at least %d characters, got %d",
					name,
					minimum,
					len(value),
				)
			}
			return nil
		},
		severity,
		name+" must meet minimum length",
	)
}

// MaxLength creates a rule that validates len(value) <= maximum.
// Use for validating maximum string length constraints.
func MaxLength(name, value string, maximum int, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if len(value) > maximum {
				return fmt.Errorf(
					"%s must not exceed %d characters, got %d",
					name,
					maximum,
					len(value),
				)
			}
			return nil
		},
		severity,
		name+" must not exceed maximum length",
	)
}

// Matches creates a rule that validates the string matches a regex pattern.
// Use for validating strings against custom patterns.
func Matches(name, value string, pattern *regexp.Regexp, severity Severity) Rule {
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

// Generic Rules.

// Equals creates a rule that validates value == expected.
// Use for validating equality of any comparable type.
func Equals[T comparable](name string, value, expected T, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value != expected {
				return fmt.Errorf("%s must equal %v, got %v", name, expected, value)
			}
			return nil
		},
		severity,
		name+" must equal expected value",
	)
}
