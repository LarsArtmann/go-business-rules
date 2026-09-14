// Package cqrslite bridges businessrules validation events onto a
// go-cqrs-lite event bus. See doc.go for the wire contract and usage.
package cqrslite

import (
	"context"
	"sync/atomic"

	businessrules "github.com/LarsArtmann/go-business-rules"
	"github.com/larsartmann/go-cqrs-lite/event/v4"
	"github.com/larsartmann/go-cqrs-lite/id/v4"
)

// TypeRuleEvaluated is the event type published for every evaluated rule.
const TypeRuleEvaluated event.Type = "businessrules.rule.evaluated.v1"

// TypeValidationCompleted is the event type published when a validation run finishes.
const TypeValidationCompleted event.Type = "businessrules.validation.completed.v1"

// RuleEvaluatedData is the wire payload for a single rule evaluation.
type RuleEvaluatedData struct {
	RuleName         string
	Severity         string
	Passed           bool
	ErrorMessage     string
	DurationNanos    int64
	StartedAtUnixNano int64
}

// ValidationCompletedData is the wire payload for a finished validation run.
type ValidationCompletedData struct {
	Valid             bool
	ViolationCount    int
	DurationNanos     int64
	StartedAtUnixNano int64
}

// Option configures the bus listener.
type Option func(*busListener)

// WithPublishErrorHandler sets the handler invoked when publishing to the bus
// fails. The validation run itself never aborts on publish errors. Without
// this option, publish errors are silently discarded.
func WithPublishErrorHandler(handler func(error)) Option {
	return func(l *busListener) {
		l.onPublishError = handler
	}
}

type busListener struct {
	bus            event.Bus
	streamID       id.StreamID
	streamType     id.StreamType
	version        atomic.Uint64
	onPublishError func(error)
}

// NewBusListener returns a businessrules.Listener that publishes each
// validation event onto the bus as a go-cqrs-lite domain event, tagged with
// the given stream identity and a monotonically increasing version.
func NewBusListener(bus event.Bus, streamID id.StreamID, streamType id.StreamType, opts ...Option) businessrules.Listener {
	listener := &busListener{
		bus:            bus,
		streamID:       streamID,
		streamType:     streamType,
		onPublishError: func(error) {},
	}

	for _, opt := range opts {
		opt(listener)
	}

	return listener.handle
}

func (l *busListener) handle(e businessrules.Event) {
	ctx := context.Background()

	switch ev := e.(type) {
	case businessrules.RuleEvaluated:
		l.publish(ctx, TypeRuleEvaluated, newRuleEvaluatedData(ev))
	case businessrules.ValidationCompleted:
		l.publish(ctx, TypeValidationCompleted, newValidationCompletedData(ev))
	}
}

func (l *busListener) publish(ctx context.Context, eventType event.Type, payload any) {
	evt, err := event.New(eventType, l.streamID, l.streamType, event.Version(l.version.Add(1)), payload)
	if err != nil {
		l.onPublishError(err)

		return
	}

	if err := l.bus.Publish(ctx, evt); err != nil {
		l.onPublishError(err)
	}
}

func newRuleEvaluatedData(re businessrules.RuleEvaluated) RuleEvaluatedData {
	message := ""
	if re.Err != nil {
		message = re.Err.Error()
	}

	return RuleEvaluatedData{
		RuleName:          re.RuleName,
		Severity:          string(re.Severity),
		Passed:            re.Passed(),
		ErrorMessage:      message,
		DurationNanos:     int64(re.Duration),
		StartedAtUnixNano: re.At.UnixNano(),
	}
}

func newValidationCompletedData(vc businessrules.ValidationCompleted) ValidationCompletedData {
	return ValidationCompletedData{
		Valid:             vc.Result.Valid,
		ViolationCount:    vc.Result.Count(),
		DurationNanos:     int64(vc.Duration),
		StartedAtUnixNano: vc.At.UnixNano(),
	}
}
