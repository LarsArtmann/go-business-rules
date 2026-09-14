package businessrules

// ValidatorBuilder provides a fluent API for building validators.
// Add rules using AddRule or AddRules, optionally observe the evaluation
// with WithListener, then call Build to get the ValidationResultError.
type ValidatorBuilder struct {
	rules     []Rule
	listeners []Listener
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

// WithListener registers listeners that receive validation events
// (RuleEvaluated for every rule, ValidationCompleted at the end) while
// Build or Stream runs. Returns the builder for method chaining.
func (b *ValidatorBuilder) WithListener(listeners ...Listener) *ValidatorBuilder {
	b.listeners = append(b.listeners, listeners...)

	return b
}

// Build executes all rules and returns the validation result.
// The ValidationResultError contains all violations and a Valid flag indicating success.
// When listeners are registered, events are emitted synchronously during
// execution; without listeners no timing or emission overhead is incurred.
func (b *ValidatorBuilder) Build() ValidationResultError {
	timed := len(b.listeners) > 0

	var runStart time.Time
	if timed {
		runStart = time.Now()
	}

	violations := make([]ViolationError, 0, len(b.rules))

	for _, rule := range b.rules {
		var ruleStart time.Time
		if timed {
			ruleStart = time.Now()
		}

		err := rule.Check()

		if timed {
			b.emit(RuleEvaluated{
				RuleName: rule.Name(),
				Severity: rule.Severity(),
				Passed:   err == nil,
				Err:      err,
				Duration: time.Since(ruleStart),
				At:       ruleStart,
			})
		}

		if err != nil {
			violations = append(violations, NewViolationFromError(rule, err))
		}
	}

	result := ValidationResultError{
		Valid:           len(violations) == 0,
		ViolationErrors: violations,
	}

	if timed {
		b.emit(ValidationCompleted{
			Result:   result,
			Duration: time.Since(runStart),
			At:       runStart,
		})
	}

	return result
}

// emit delivers an event to every listener, in registration order.
func (b *ValidatorBuilder) emit(event Event) {
	for _, listener := range b.listeners {
		listener(event)
	}
}
