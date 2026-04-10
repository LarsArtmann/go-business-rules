package businessrules

import "fmt"

// Collection Rules.

// notEmptyCheck creates a validation rule that checks if a collection is not empty.
func notEmptyCheck(name string, length int, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			if length == 0 {
				return fmt.Errorf("%s must not be empty", name)
			}

			return nil
		},
		severity,
		name+" must not be empty",
	)
}

// NotEmptySlice creates a rule that validates a slice has at least one element.
// Use for validating that lists, arrays, or slices are not empty.
func NotEmptySlice[T any](name string, value []T, severity Severity) Rule {
	return notEmptyCheck(name, len(value), severity)
}

// NotEmptyMap creates a rule that validates a map has at least one entry.
// Use for validating that maps are not empty.
func NotEmptyMap[T any](name string, value map[string]T, severity Severity) Rule {
	return notEmptyCheck(name, len(value), severity)
}

// Additional Numeric Rules.

// GreaterThan creates a rule that validates value > minimum.
// Use for validating numeric values that must exceed a threshold.
func GreaterThan(name string, value, minimum float64, severity Severity) Rule {
	return thresholdCheck(name, value, minimum, severity,
		lessThanOrEqual,
		"must be greater than",
		"must be greater than minimum",
	)
}

// LessThan creates a rule that validates value < maximum.
// Use for validating numeric values that must be below a threshold.
func LessThan(name string, value, maximum float64, severity Severity) Rule {
	return thresholdCheck(name, value, maximum, severity,
		greaterThanOrEqual,
		"must be less than",
		"must be less than maximum",
	)
}
