package businessrules

import "fmt"

// OneOf creates a rule that validates the value is in the allowed set.
// Use for validating enum-like values or restricted options.
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

// Custom creates a rule with a user-defined validation function.
// Use for complex validations not covered by built-in rules.
func Custom(name string, check func() error, severity Severity) Rule {
	return NewRule(name, check, severity, name+" validation failed")
}

// All creates a rule that passes only when all sub-rules pass.
// Violations from all failed rules are collected.
func All(name string, rules []Rule, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			var violations []string
			for _, rule := range rules {
				if err := rule.Check(); err != nil {
					violations = append(violations, err.Error())
				}
			}
			if len(violations) > 0 {
				return fmt.Errorf("%s failed: %v", name, violations)
			}
			return nil
		},
		severity,
		name+" all rules must pass",
	)
}

// Any creates a rule that passes when at least one sub-rule passes.
// Fails only when all sub-rules fail.
func Any(name string, rules []Rule, severity Severity) Rule {
	return NewRule(
		name,
		func() error {
			for _, rule := range rules {
				if err := rule.Check(); err == nil {
					return nil
			}
			}
			return fmt.Errorf("%s: none of the alternative rules passed", name)
		},
		severity,
		name+" at least one rule must pass",
	)
}

// When creates a conditional rule that only validates when condition is true.
// Use for conditional validation logic.
func When(name string, condition bool, rule Rule) Rule {
	return NewRule(
		name,
		func() error {
			if !condition {
				return nil
			}
			return rule.Check()
		},
		rule.Severity(),
		name+" conditional validation",
	)
}
