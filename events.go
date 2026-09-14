package businessrules

import (
	"slices"
	"time"
)

// Event is a validation lifecycle event emitted to registered Listeners.
// The interface is sealed: the only implementations are RuleEvaluated and
// ValidationCompleted, so consumers can switch over events exhaustively.
type Event interface {
	event()
}

// ruleMetadata is implemented by rules carrying optional metadata
// (RuleImpl and ContextRuleImpl). Custom Rule implementations without
// metadata simply do not implement it.
type ruleMetadata interface {
	Description() string
	Tags() []string
}

// metadataOf extracts optional rule metadata, returning zero values for
// rules that do not carry any.
func metadataOf(rule Rule) (string, []string) {
	if m, ok := rule.(ruleMetadata); ok {
		return m.Description(), m.Tags()
	}

	return "", nil
}

// RuleEvaluated reports the outcome of a single rule check.
// It is emitted for passing and failing rules alike, so listeners observe
// the complete evaluation, not only the violations.
type RuleEvaluated struct {
	// RuleName is the name of the evaluated rule.
	RuleName string

	// Severity is the severity of the evaluated rule.
	Severity Severity

	// Err is the error returned by the rule check; nil when the rule passed.
	Err error

	// Duration is how long the rule check took.
	Duration time.Duration

	// At is the time the rule check started.
	At time.Time

	// Description is the rule's optional intent metadata; empty when unset.
	Description string

	// Tags are the rule's optional classification labels; nil when unset.
	Tags []string
}

// Passed reports whether the rule check succeeded.
// It is derived from Err, so an event can never claim success while
// carrying an error.
func (e RuleEvaluated) Passed() bool { return e.Err == nil }

func (RuleEvaluated) event() {}

// ValidationCompleted is the terminal event of a validation run.
// It carries the same result value that the run returns, plus the total
// duration of the evaluation.
type ValidationCompleted struct {
	// Result is the outcome of the validation run.
	Result ValidationResultError

	// Duration is the total time spent evaluating all rules.
	Duration time.Duration

	// At is the time the validation run started.
	At time.Time
}

func (ValidationCompleted) event() {}

// Listener receives validation events synchronously, in evaluation order,
// before the validation run returns. Listeners are trusted code and must
// not panic: a panicking listener crashes the caller of Build.
type Listener func(Event)
