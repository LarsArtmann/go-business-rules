package businessrules

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
	n string
	c func() error
	s Severity
	m string
}

// Name returns the identifier for this rule.
func (r RuleImpl) Name() string { return r.n }

// Check validates the rule condition.
func (r RuleImpl) Check() error { return r.c() }

// Severity returns the importance level of this rule.
func (r RuleImpl) Severity() Severity { return r.s }

// Message returns the human-readable description of what this rule validates.
func (r RuleImpl) Message() string { return r.m }

// WithName returns a new RuleImpl with the specified name.
func (r RuleImpl) WithName(n string) RuleImpl {
	return RuleImpl{n: n, c: r.c, s: r.s, m: r.m}
}

// WithSeverity returns a new RuleImpl with the specified severity.
func (r RuleImpl) WithSeverity(s Severity) RuleImpl {
	return RuleImpl{n: r.n, c: r.c, s: s, m: r.m}
}

// WithMessage returns a new RuleImpl with the specified message.
func (r RuleImpl) WithMessage(m string) RuleImpl {
	return RuleImpl{n: r.n, c: r.c, s: r.s, m: m}
}

// NewRule creates a new rule with the given parameters.
// The check function should return nil on success or an error on failure.
func NewRule(name string, check func() error, severity Severity, message string) RuleImpl {
	return RuleImpl{n: name, c: check, s: severity, m: message}
}
