package businessrules

import (
	"context"
	"sync"
	"time"
)

// ValidatorBuilder provides a fluent API for building validators.
// Add rules using AddRule or AddRules, optionally observe the evaluation
// with WithListener, then call Build to get the ValidationResultError.
type ValidatorBuilder struct {
	rules     []Rule
	listeners []Listener
}

// NewValidator creates a new ValidatorBuilder with an empty rule set.
func NewValidator() *ValidatorBuilder {
	return &ValidatorBuilder{
		rules:     make([]Rule, 0),
		listeners: nil,
	}
}

// AddRule adds a single rule to the validator.
// Returns the builder for method chaining.
func (b *ValidatorBuilder) AddRule(rule Rule) *ValidatorBuilder {
	b.rules = append(b.rules, rule)

	return b
}

// AddRules adds multiple rules to the validator.
// Returns the builder for method chaining.
func (b *ValidatorBuilder) AddRules(rules ...Rule) *ValidatorBuilder {
	b.rules = append(b.rules, rules...)

	return b
}

// WithListener registers listeners that receive validation events
// (RuleEvaluated for every rule, ValidationCompleted at the end) while
// Build or Stream runs. Returns the builder for method chaining.
func (b *ValidatorBuilder) WithListener(listeners ...Listener) *ValidatorBuilder {
	b.listeners = append(b.listeners, listeners...)

	return b
}

// Build executes all rules and returns the validation result.
// The ValidationResultError contains all violations and a Valid flag indicating success.
// When listeners are registered, events are emitted synchronously during
// execution; without listeners no timing or emission overhead is incurred.
func (b *ValidatorBuilder) Build() ValidationResultError {
	timed := len(b.listeners) > 0

	var runStart time.Time
	if timed {
		runStart = time.Now()
	}

	violations := make([]ViolationError, 0, len(b.rules))

	for _, rule := range b.rules {
		var ruleStart time.Time
		if timed {
			ruleStart = time.Now()
		}

		err := rule.Check()

		if timed {
			b.emit(RuleEvaluated{
				RuleName: rule.Name(),
				Severity: rule.Severity(),
				Err:      err,
				Duration: time.Since(ruleStart),
				At:       ruleStart,
			})
		}

		if err != nil {
			violations = append(violations, NewViolationFromError(rule, err))
		}
	}

	result := ValidationResultError{
		Valid:           len(violations) == 0,
		ViolationErrors: violations,
	}

	if timed {
		b.emit(ValidationCompleted{
			Result:   result,
			Duration: time.Since(runStart),
			At:       runStart,
		})
	}

	return result
}

// emit delivers an event to every listener, in registration order.
func (b *ValidatorBuilder) emit(event Event) {
	for _, listener := range b.listeners {
		listener(event)
	}
}

// streamResult pairs a rule check outcome with its position in the rule
// list, so the final result can be reported in deterministic rule order
// even though events stream in completion order.
type streamResult struct {
	index     int
	rule      Rule
	evaluated RuleEvaluated
}

// Stream executes all rules concurrently and returns a channel of validation
// events: one RuleEvaluated per rule as checks finish (completion order),
// then a terminal ValidationCompleted whose violations follow the original
// rule order. The channel is closed after the terminal event.
//
// Events are delivered on the returned channel instead of to listeners
// registered with WithListener.
//
// When ctx is canceled, unstarted rules are skipped while running checks
// (which cannot be interrupted) are still awaited. ValidationCompleted is
// emitted only if the stream drains before the cancellation; abandoning a
// canceled stream never leaks goroutines.
func (b *ValidatorBuilder) Stream(ctx context.Context) <-chan Event {
	out := make(chan Event)

	go func() {
		defer close(out)

		runStart := time.Now()
		results := make(chan streamResult, len(b.rules))

		var waitGroup sync.WaitGroup

		for index, rule := range b.rules {
			if ctx.Err() != nil {
				break
			}

			waitGroup.Add(1)

			go func(index int, rule Rule) {
				defer waitGroup.Done()

				ruleStart := time.Now()

				err := rule.Check()
				results <- streamResult{
					index: index,
					rule:  rule,
					evaluated: RuleEvaluated{
						RuleName: rule.Name(),
						Severity: rule.Severity(),
						Err:      err,
						Duration: time.Since(ruleStart),
						At:       ruleStart,
					},
				}
			}(index, rule)
		}

		go func() {
			waitGroup.Wait()
			close(results)
		}()

		violationsByIndex := make(map[int]ViolationError, len(b.rules))

		for result := range results {
			select {
			case out <- result.evaluated:
			case <-ctx.Done():
				return
			}

			if result.evaluated.Err != nil {
				violationsByIndex[result.index] = NewViolationFromError(result.rule, result.evaluated.Err)
			}
		}

		violations := make([]ViolationError, 0, len(violationsByIndex))

		for index := range b.rules {
			if violation, failed := violationsByIndex[index]; failed {
				violations = append(violations, violation)
			}
		}

		result := ValidationResultError{
			Valid:           len(violations) == 0,
			ViolationErrors: violations,
		}

		select {
		case out <- ValidationCompleted{
			Result:   result,
			Duration: time.Since(runStart),
			At:       runStart,
		}:
		case <-ctx.Done():
		}
	}()

	return out
}
