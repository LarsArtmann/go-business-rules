package businessrules

import (
	"fmt"
	"slices"
)

// OneOf creates a rule that validates the value is in the allowed set.
// Use for validating enum-like values or restricted options.
func OneOf[T comparable](name string, value T, allowed []T, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			if slices.Index(allowed, value) == -1 {
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
func Custom(name string, check func() error, severity Severity) RuleImpl {
	return NewRule(name, check, severity, name+" validation failed")
}

// compositeRule creates a rule that combines multiple sub-rules with a custom strategy.
type ruleStrategy func(name string, rules []Rule) error

func collectAllViolations(name string, rules []Rule) error {
	var violations []string

	for _, rule := range rules {
		err := rule.Check()
		if err != nil {
			violations = append(violations, err.Error())
		}
	}

	if len(violations) > 0 {
		return fmt.Errorf("%s failed: %v", name, violations)
	}

	return nil
}

func anyRulePasses(name string, rules []Rule) error {
	for _, rule := range rules {
		err := rule.Check()
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("%s: none of the alternative rules passed", name)
}

// All creates a rule that passes only when all sub-rules pass.
// Violations from all failed rules are collected.
func All(name string, rules []Rule, severity Severity) RuleImpl {
	return compositeRuleWith(
		name,
		rules,
		severity,
		collectAllViolations,
		name+" all rules must pass",
	)
}

// Any creates a rule that passes when at least one sub-rule passes.
// Fails only when all sub-rules fail.
func Any(name string, alternatives []Rule, severity Severity) RuleImpl {
	return compositeRuleWith(
		name,
		alternatives,
		severity,
		anyRulePasses,
		name+" at least one rule must pass",
	)
}

func compositeRuleWith(
	name string,
	rules []Rule,
	severity Severity,
	strategy ruleStrategy,
	msg string,
) RuleImpl {
	return NewRule(name, func() error { return strategy(name, rules) }, severity, msg)
}

// When creates a conditional rule that only validates when condition is true.
// Use for conditional validation logic.
func When(name string, condition bool, rule Rule) RuleImpl {
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

// Not creates a rule that passes only when the inner rule FAILS. Use for
// expressing forbidden conditions, such as rejecting values an allowlist rule
// would accept.
func Not(name string, rule Rule, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			if err := rule.Check(); err == nil {
				return fmt.Errorf("%s: forbidden rule %s passed", name, rule.Name())
			}

			return nil
		},
		severity,
		name+" must not satisfy the forbidden rule",
	)
}

// Or creates a rule that passes when at least one sub-rule passes, accepting
// the sub-rules variadically (unlike Any, which takes a slice).
func Or(name string, severity Severity, rules ...Rule) RuleImpl {
	return compositeRuleWith(
		name,
		rules,
		severity,
		anyRulePasses,
		name+" at least one rule must pass",
	)
}

// Xor creates a rule that passes when EXACTLY ONE of the two sub-rules passes.
// Use for mutually exclusive alternatives, such as "either email or phone is set".
func Xor(name string, first, second Rule, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			passing := 0

			for _, rule := range []Rule{first, second} {
				if rule.Check() == nil {
					passing++
				}
			}

			if passing != 1 {
				return fmt.Errorf(
					"%s: exactly one of the rules must pass, got %d passing",
					name,
					passing,
				)
			}

			return nil
		},
		severity,
		name+" exactly one rule must pass",
	)
}
