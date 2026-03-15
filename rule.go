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

// baseRule is the default implementation of the Rule interface.
type baseRule struct {
	name     string
	check    func() error
	severity Severity
	message  string
}

// Name returns the rule's identifier.
func (r baseRule) Name() string {
	return r.name
}

// Check executes the validation function.
func (r baseRule) Check() error {
	return r.check()
}

// Severity returns the rule's severity level.
func (r baseRule) Severity() Severity {
	return r.severity
}

// Message returns the rule's validation message.
func (r baseRule) Message() string {
	return r.message
}

// NewRule creates a new Rule with the given parameters.
// The check function should return nil on success or an error on failure.
func NewRule(name string, check func() error, severity Severity, message string) Rule {
	return baseRule{
		name:     name,
		check:    check,
		severity: severity,
		message:  message,
	}
}
