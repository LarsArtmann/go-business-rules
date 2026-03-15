package businessrules

// ValidatorBuilder provides a fluent API for building validators.
// Add rules using AddRule or AddRules, then call Build to get the Result.
type ValidatorBuilder struct {
	rules []Rule
}

// NewValidator creates a new ValidatorBuilder with an empty rule set.
func NewValidator() *ValidatorBuilder {
	return &ValidatorBuilder{
		rules: make([]Rule, 0),
	}
}

// AddRule adds a single rule to the validator.
// Returns the builder for method chaining.
func (b *ValidatorBuilder) AddRule(rule Rule) *ValidatorBuilder {
	b.rules = append(b.rules, rule)
	return b
}

// AddRules adds multiple rules to the validator.
// Returns the builder for method chaining.
func (b *ValidatorBuilder) AddRules(rules ...Rule) *ValidatorBuilder {
	b.rules = append(b.rules, rules...)
	return b
}

// Build executes all rules and returns the validation result.
// The Result contains all violations and a Valid flag indicating success.
func (b *ValidatorBuilder) Build() Result {
	violations := make([]Violation, 0, len(b.rules))

	for _, rule := range b.rules {
		if err := rule.Check(); err != nil {
			violations = append(violations, NewViolationFromError(rule, err))
		}
	}

	return Result{
		Valid:      len(violations) == 0,
		Violations: violations,
	}
}
