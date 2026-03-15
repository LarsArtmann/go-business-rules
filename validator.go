package businessrules

type ValidatorBuilder struct {
	rules []Rule
}

func NewValidator() *ValidatorBuilder {
	return &ValidatorBuilder{
		rules: make([]Rule, 0),
	}
}

func (b *ValidatorBuilder) AddRule(rule Rule) *ValidatorBuilder {
	b.rules = append(b.rules, rule)
	return b
}

func (b *ValidatorBuilder) AddRules(rules ...Rule) *ValidatorBuilder {
	b.rules = append(b.rules, rules...)
	return b
}

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
