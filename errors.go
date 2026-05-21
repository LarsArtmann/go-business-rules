package businessrules

import (
	"encoding/json"
	"fmt"
	"time"
)

// ViolationError represents a failed rule check with context and metadata.
// It implements the error interface for seamless integration with
// Go's error handling patterns.
type ViolationError struct {
	Timestamp time.Time
	Context   string
	Rule      Rule
}

// Error implements the error interface, returning a formatted violation message.
// The format is: [SEVERITY] rule_name: message (context: details).
func (v ViolationError) Error() string {
	if v.Context != "" {
		return fmt.Sprintf(
			"[%s] %s: %s (context: %s)",
			v.Rule.Severity().String(),
			v.Rule.Name(),
			v.Rule.Message(),
			v.Context,
		)
	}

	return fmt.Sprintf(
		"[%s] %s: %s",
		v.Rule.Severity().String(),
		v.Rule.Name(),
		v.Rule.Message(),
	)
}

// NewViolation creates a ViolationError with the current timestamp.
// The context parameter provides additional details about the failure.
func NewViolation(rule Rule, context string) ViolationError {
	return ViolationError{
		Rule:      rule,
		Context:   context,
		Timestamp: time.Now(),
	}
}

// NewViolationFromError creates a ViolationError from an error.
// The error's message is used as the context.
func NewViolationFromError(rule Rule, err error) ViolationError {
	return ViolationError{
		Rule:      rule,
		Context:   err.Error(),
		Timestamp: time.Now(),
	}
}

// WithContext returns a new ViolationError with updated context.
// Useful for adding context to an existing violation.
func (v ViolationError) WithContext(context string) ViolationError {
	return ViolationError{
		Rule:      v.Rule,
		Context:   context,
		Timestamp: v.Timestamp,
	}
}

// MarshalJSON implements json.Marshaler for ViolationError.
func (v ViolationError) MarshalJSON() ([]byte, error) {
	marshaled, err := json.Marshal(struct {
		RuleName  string    `json:"ruleName"`
		Severity  string    `json:"severity"`
		Message   string    `json:"message"`
		Context   string    `json:"context,omitempty"`
		Timestamp time.Time `json:"timestamp"`
	}{
		RuleName:  v.Rule.Name(),
		Severity:  v.Rule.Severity().String(),
		Message:   v.Rule.Message(),
		Context:   v.Context,
		Timestamp: v.Timestamp,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ViolationError: %w", err)
	}

	return marshaled, nil
}
