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

// rule implements Rule interface using functional options.
type rule struct {
	name     string
	check    func() error
	severity Severity
	message  string
}

func (r rule) Name() string     { return r.name }
func (r rule) Check() error     { return r.check() }
func (r rule) Severity() Severity { return r.severity }
func (r rule) Message() string  { return r.message }

func newRule(name string, check func() error, severity Severity, message string) rule {
	return rule{name: name, check: check, severity: severity, message: message}
}

func (r rule) WithName(name string) Rule {
	return newRule(name, r.check, r.severity, r.message)
}

func (r rule) WithSeverity(severity Severity) Rule {
	return newRule(r.name, r.check, severity, r.message)
}

func (r rule) WithMessage(message string) Rule {
	return newRule(r.name, r.check, r.severity, message)
}

// NewRule creates a new Rule with the given parameters.
// The check function should return nil on success or an error on failure.
func NewRule(name string, check func() error, severity Severity, message string) Rule {
	return newRule(name, check, severity, message)
}
