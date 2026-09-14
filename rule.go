package businessrules

import "context"

// Rule defines the interface for a validation rule.
// Implementations check a specific condition and report failures
// with an associated severity level.
type Rule interface {
	// Name returns the identifier for this rule.
	// Used in error messages and logging.
	Name() string

	// Check validates the rule condition.
	// Returns nil if the rule passes, or an error describing the violation.
	Check() error

	// Severity returns the importance level of this rule.
	// Used to categorize violations in the ValidationResult.
	Severity() Severity

	// Message returns the human-readable description of what this rule validates.
	// Used as a template for error messages.
	Message() string
}

// RuleImpl implements Rule interface.
type RuleImpl struct {
	n    string
	c    func() error
	s    Severity
	m    string
	d    string
	tags []string
}

// Name returns the identifier for this rule.
func (r RuleImpl) Name() string { return r.n }

// Check validates the rule condition.
func (r RuleImpl) Check() error { return r.c() }

// Severity returns the importance level of this rule.
func (r RuleImpl) Severity() Severity { return r.s }

// Message returns the human-readable description of what this rule validates.
func (r RuleImpl) Message() string { return r.m }

// Description returns optional metadata explaining WHY the rule exists.
// Empty when unset. Unlike Message (the failure template), the description
// documents intent for observers such as event listeners.
func (r RuleImpl) Description() string { return r.d }

// Tags returns optional classification labels for the rule. Nil when unset.
// The returned slice is owned by the rule and must not be modified.
func (r RuleImpl) Tags() []string { return r.tags }

// WithName returns a new RuleImpl with the specified name.
func (r RuleImpl) WithName(n string) RuleImpl {
	return RuleImpl{n: n, c: r.c, s: r.s, m: r.m, d: r.d, tags: r.tags}
}

// WithSeverity returns a new RuleImpl with the specified severity.
func (r RuleImpl) WithSeverity(s Severity) RuleImpl {
	return RuleImpl{n: r.n, c: r.c, s: s, m: r.m, d: r.d, tags: r.tags}
}

// WithMessage returns a new RuleImpl with the specified message.
func (r RuleImpl) WithMessage(m string) RuleImpl {
	return RuleImpl{n: r.n, c: r.c, s: r.s, m: m, d: r.d, tags: r.tags}
}

// WithDescription returns a new RuleImpl with the specified description.
func (r RuleImpl) WithDescription(d string) RuleImpl {
	return RuleImpl{n: r.n, c: r.c, s: r.s, m: r.m, d: d, tags: r.tags}
}

// WithTags returns a new RuleImpl with the specified tags.
func (r RuleImpl) WithTags(tags ...string) RuleImpl {
	return RuleImpl{n: r.n, c: r.c, s: r.s, m: r.m, d: r.d, tags: tags}
}

// NewRule creates a new rule with the given parameters.
// The check function should return nil on success or an error on failure.
func NewRule(name string, check func() error, severity Severity, message string) RuleImpl {
	return RuleImpl{n: name, c: check, s: severity, m: message, d: "", tags: nil}
}

// ContextRule is implemented by rules whose check can observe context
// cancellation. It is optional and additive: Stream calls CheckContext when
// a rule implements it, so slow (e.g. I/O-bound) checks can return early on
// cancellation, and falls back to Check for all other rules. Build keeps
// calling Check, because it takes no context.
type ContextRule interface {
	Rule

	// CheckContext validates the rule condition, returning promptly when ctx
	// is canceled. A context error must mean the check was canceled, not that
	// the validated value is a violation.
	CheckContext(ctx context.Context) error
}

// ContextRuleImpl implements the ContextRule interface with a
// context-aware check function.
type ContextRuleImpl struct {
	n    string
	c    func(ctx context.Context) error
	s    Severity
	m    string
	d    string
	tags []string
}

// Name returns the identifier for this rule.
func (r ContextRuleImpl) Name() string { return r.n }

// Check validates the rule condition with a background context, which never
// cancels. Stream calls CheckContext instead when it can.
func (r ContextRuleImpl) Check() error { return r.c(context.Background()) }

// CheckContext validates the rule condition, returning promptly when ctx is canceled.
func (r ContextRuleImpl) CheckContext(ctx context.Context) error { return r.c(ctx) }

// Severity returns the importance level of this rule.
func (r ContextRuleImpl) Severity() Severity { return r.s }

// Message returns the human-readable description of what this rule validates.
func (r ContextRuleImpl) Message() string { return r.m }

// Description returns optional metadata explaining WHY the rule exists.
// Empty when unset. Unlike Message (the failure template), the description
// documents intent for observers such as event listeners.
func (r ContextRuleImpl) Description() string { return r.d }

// Tags returns optional classification labels for the rule. Nil when unset.
// The returned slice is owned by the rule and must not be modified.
func (r ContextRuleImpl) Tags() []string { return r.tags }

// WithName returns a new ContextRuleImpl with the specified name.
func (r ContextRuleImpl) WithName(n string) ContextRuleImpl {
	return ContextRuleImpl{n: n, c: r.c, s: r.s, m: r.m, d: r.d, tags: r.tags}
}

// WithSeverity returns a new ContextRuleImpl with the specified severity.
func (r ContextRuleImpl) WithSeverity(s Severity) ContextRuleImpl {
	return ContextRuleImpl{n: r.n, c: r.c, s: s, m: r.m, d: r.d, tags: r.tags}
}

// WithMessage returns a new ContextRuleImpl with the specified message.
func (r ContextRuleImpl) WithMessage(m string) ContextRuleImpl {
	return ContextRuleImpl{n: r.n, c: r.c, s: r.s, m: m, d: r.d, tags: r.tags}
}

// WithDescription returns a new ContextRuleImpl with the specified description.
func (r ContextRuleImpl) WithDescription(d string) ContextRuleImpl {
	return ContextRuleImpl{n: r.n, c: r.c, s: r.s, m: r.m, d: d, tags: r.tags}
}

// WithTags returns a new ContextRuleImpl with the specified tags.
func (r ContextRuleImpl) WithTags(tags ...string) ContextRuleImpl {
	return ContextRuleImpl{n: r.n, c: r.c, s: r.s, m: r.m, d: r.d, tags: tags}
}

// NewContextRule creates a new context-aware rule with the given parameters.
// The check function should return nil on success or an error on failure;
// when ctx is canceled mid-check it should return promptly, typically with
// the context error.
func NewContextRule(
	name string,
	check func(ctx context.Context) error,
	severity Severity,
	message string,
) ContextRuleImpl {
	return ContextRuleImpl{n: name, c: check, s: severity, m: message, d: "", tags: nil}
}
