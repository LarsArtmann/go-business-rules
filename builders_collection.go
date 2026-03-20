package businessrules

import "fmt"

// Collection Rules.

// NotEmptySlice creates a rule that validates a slice has at least one element.
// Use for validating that lists, arrays, or slices are not empty.
func NotEmptySlice[T any](name string, value []T, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if len(value) == 0 {
				return fmt.Errorf("%s must not be empty", name)
			}
			return nil
		},
		severity,
		name+" must not be empty",
	)
}

// NotEmptyMap creates a rule that validates a map has at least one entry.
// Use for validating that maps are not empty.
func NotEmptyMap[T any](name string, value map[string]T, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if len(value) == 0 {
				return fmt.Errorf("%s must not be empty", name)
			}
			return nil
		},
		severity,
		name+" must not be empty",
	)
}

// Additional Numeric Rules.

// GreaterThan creates a rule that validates value > minimum.
// Use for validating numeric values that must exceed a threshold.
func GreaterThan(name string, value, minimum float64, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value <= minimum {
				return fmt.Errorf("%s must be greater than %f, got %f", name, minimum, value)
			}
			return nil
		},
		severity,
		name+" must be greater than minimum",
	)
}

// LessThan creates a rule that validates value < maximum.
// Use for validating numeric values that must be below a threshold.
func LessThan(name string, value, maximum float64, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if value >= maximum {
				return fmt.Errorf("%s must be less than %f, got %f", name, maximum, value)
			}
			return nil
		},
		severity,
		name+" must be less than maximum",
	)
}
