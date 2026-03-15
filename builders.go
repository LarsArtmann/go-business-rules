package businessrules

import (
	"fmt"
	"regexp"
)

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

func OneOf[T comparable](name string, value T, allowed []T, severity Severity) Rule {
	allowedSet := make(map[T]bool, len(allowed))
	for _, v := range allowed {
		allowedSet[v] = true
	}

	return NewRule(
		name,
		func() error {
			if !allowedSet[value] {
				return fmt.Errorf("%s must be one of the allowed values", name)
			}
			return nil
		},
		severity,
		name+" must be one of allowed values",
	)
}

func Custom(name string, check func() error, severity Severity) Rule {
	return NewRule(name, check, severity, name+" validation failed")
}
