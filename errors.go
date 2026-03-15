package businessrules

import (
	"fmt"
	"time"
)

type Violation struct {
	Rule      Rule
	Context   string
	Timestamp time.Time
}

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

func NewViolation(rule Rule, context string) Violation {
	return Violation{
		Rule:      rule,
		Context:   context,
		Timestamp: time.Now(),
	}
}

func NewViolationFromError(rule Rule, err error) Violation {
	return Violation{
		Rule:      rule,
		Context:   err.Error(),
		Timestamp: time.Now(),
	}
}
