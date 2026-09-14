package businessrules

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Collection Rules.

// notEmptyCheck creates a validation rule that checks if a collection is not empty.
func notEmptyCheck(name string, length int, severity Severity) RuleImpl {
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
func NotEmptySlice[T any](name string, value []T, severity Severity) RuleImpl {
	return notEmptyCheck(name, len(value), severity)
}

// NotEmptyMap creates a rule that validates a map has at least one entry.
// Use for validating that maps are not empty.
func NotEmptyMap[T any](name string, value map[string]T, severity Severity) RuleImpl {
	return notEmptyCheck(name, len(value), severity)
}

// Additional Numeric Rules.

// GreaterThan creates a rule that validates value > minimum.
// Use for validating numeric values that must exceed a threshold.
func GreaterThan(name string, value, minimum float64, severity Severity) RuleImpl {
	return thresholdCheck(
		name, value, minimum, severity,
		lessThanOrEqual,
		"must be greater than",
		"must be greater than minimum",
	)
}

// LessThan creates a rule that validates value < maximum.
// Use for validating numeric values that must be below a threshold.
func LessThan(name string, value, maximum float64, severity Severity) RuleImpl {
	return thresholdCheck(
		name, value, maximum, severity,
		greaterThanOrEqual,
		"must be less than",
		"must be less than maximum",
	)
}

// Precision Rules.

// MaxDecimalPlaces creates a rule that validates a float carries at most the
// given number of decimal places in its shortest representation. Use for
// monetary or measurement fields where a value like 0.30000000000000004 is a
// caller bug, not an acceptable amount. Non-finite values fail.
func MaxDecimalPlaces(name string, value float64, places int, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return fmt.Errorf("%s must be a finite number, got %v", name, value)
			}

			formatted := strconv.FormatFloat(value, 'f', -1, 64)
			if dot := strings.IndexByte(formatted, '.'); dot >= 0 && len(formatted)-dot-1 > places {
				return fmt.Errorf(
					"%s must have at most %d decimal places, got %v",
					name,
					places,
					value,
				)
			}

			return nil
		},
		severity,
		name+" must not exceed the maximum decimal places",
	)
}

// DivisibleBy creates a rule that validates an integer is divisible by the
// divisor without remainder. Use for quantity or packaging constraints such
// as "order in multiples of 12". A zero divisor fails the check instead of
// panicking.
func DivisibleBy(name string, value, divisor int, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			if divisor == 0 {
				return fmt.Errorf("%s: divisor must not be zero", name)
			}

			// the zero-guard above already returned; the experimental analyzer
			// does not track the flow
			//nolint:branching-flow:panic
			if value%divisor != 0 {
				return fmt.Errorf(
					"%s must be divisible by %d, got %d",
					name,
					divisor,
					value,
				)
			}

			return nil
		},
		severity,
		name+" must be divisible by the divisor",
	)
}
