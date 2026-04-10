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

// rule implements Rule interface.
type rule struct {
	n string
	c func() error
	s Severity
	m string
}

func (r rule) Name() string     { return r.n }
func (r rule) Check() error    { return r.c() }
func (r rule) Severity() Severity { return r.s }
func (r rule) Message() string { return r.m }

func (r rule) WithName(n string) Rule {
	return rule{n: n, c: r.c, s: r.s, m: r.m}
}

func (r rule) WithSeverity(s Severity) Rule {
	return rule{n: r.n, c: r.c, s: s, m: r.m}
}

func (r rule) WithMessage(m string) Rule {
	return rule{n: r.n, c: r.c, s: r.s, m: m}
}

// NewRule creates a new Rule with the given parameters.
// The check function should return nil on success or an error on failure.
func NewRule(name string, check func() error, severity Severity, message string) Rule {
	return rule{n: name, c: check, s: severity, m: message}
}
