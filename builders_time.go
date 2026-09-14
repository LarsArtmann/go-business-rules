package businessrules

import (
	"fmt"
	"time"
)

func formatTime(value time.Time) string {
	return value.Format(time.RFC3339)
}

// NotPast creates a rule that validates the time is not before now.
// Pass the current time explicitly (injectable clock) so the rule stays
// deterministic and testable. Use for validating expiry or deadline fields.
func NotPast(name string, value, now time.Time, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			if value.Before(now) {
				return fmt.Errorf(
					"%s must not be in the past, got %s (now %s)",
					name,
					formatTime(value),
					formatTime(now),
				)
			}

			return nil
		},
		severity,
		name+" must not be in the past",
	)
}

// NotFuture creates a rule that validates the time is not after now.
// Pass the current time explicitly (injectable clock) so the rule stays
// deterministic and testable. Use for validating birth dates or issued-at fields.
func NotFuture(name string, value, now time.Time, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			if value.After(now) {
				return fmt.Errorf(
					"%s must not be in the future, got %s (now %s)",
					name,
					formatTime(value),
					formatTime(now),
				)
			}

			return nil
		},
		severity,
		name+" must not be in the future",
	)
}

// DateInRange creates a rule that validates start <= value <= end.
// Use for validating booking windows, contract terms, or reporting periods.
func DateInRange(
	name string,
	value, start, end time.Time,
	severity Severity,
) RuleImpl {
	return NewRule(
		name,
		func() error {
			if value.Before(start) || value.After(end) {
				return fmt.Errorf(
					"%s must be between %s and %s, got %s",
					name,
					formatTime(start),
					formatTime(end),
					formatTime(value),
				)
			}

			return nil
		},
		severity,
		name+" must be within the date range",
	)
}
