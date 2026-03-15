package businessrules

import (
	"encoding/json"
	"fmt"
	"time"
)

// Violation represents a failed rule check with context and metadata.
// It implements the error interface for seamless integration with
// Go's error handling patterns.
type Violation struct {
	// Rule is the validation rule that failed.
	Rule Rule

	// Context provides additional information about the violation.
	// May include the actual value, comparison details, or other context.
	Context string

	// Timestamp records when the violation was detected.
	Timestamp time.Time
}

// Error implements the error interface, returning a formatted violation message.
// The format is: [SEVERITY] rule_name: message (context: details)
func (v Violation) Error() string {
	if v.Context != "" {
		return fmt.Sprintf("[%s] %s: %s (context: %s)",
			v.Rule.Severity().String(),
			v.Rule.Name(),
			v.Rule.Message(),
			v.Context,
		)
	}
	return fmt.Sprintf("[%s] %s: %s",
		v.Rule.Severity().String(),
		v.Rule.Name(),
		v.Rule.Message(),
	)
}

// NewViolation creates a Violation with the current timestamp.
// The context parameter provides additional details about the failure.
func NewViolation(rule Rule, context string) Violation {
	return Violation{
		Rule:      rule,
		Context:   context,
		Timestamp: time.Now(),
	}
}

// NewViolationFromError creates a Violation from an error.
// The error's message is used as the context.
func NewViolationFromError(rule Rule, err error) Violation {
	return Violation{
		Rule:      rule,
		Context:   err.Error(),
		Timestamp: time.Now(),
	}
}

// WithContext returns a new Violation with updated context.
// Useful for adding context to an existing violation.
func (v Violation) WithContext(context string) Violation {
	return Violation{
		Rule:      v.Rule,
		Context:   context,
		Timestamp: v.Timestamp,
	}
}

// MarshalJSON implements json.Marshaler for Violation.
func (v Violation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		RuleName   string    `json:"rule_name"`
		Severity   string    `json:"severity"`
		Message    string    `json:"message"`
		Context    string    `json:"context,omitempty"`
		Timestamp  time.Time `json:"timestamp"`
	}{
		RuleName:  v.Rule.Name(),
		Severity:  v.Rule.Severity().String(),
		Message:   v.Rule.Message(),
		Context:   v.Context,
		Timestamp: v.Timestamp,
	})
}
