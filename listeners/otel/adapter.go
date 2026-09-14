// Package otel provides a businessrules.Listener that maps validation
// events onto OpenTelemetry spans and metrics, so validation runs show up
// in whatever tracing and metrics backend the application already uses.
//
// Wire contract:
//
//	RuleEvaluated       -> span "businessrules.rule" (timed from the event)
//	                       counter businessrules.rule.evaluations
//	                       histogram businessrules.rule.duration
//	ValidationCompleted -> span "businessrules.validation"
//	                       counter businessrules.validations
//	                       histogram businessrules.validation.duration
//	                       counter businessrules.violations (one per violation)
package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

// Instrument and span names follow OpenTelemetry naming: lowercase,
// dot-separated namespace.
const (
	evaluationsInstrument = "businessrules.rule.evaluations"
	ruleDurationName      = "businessrules.rule.duration"
	validationsInstrument = "businessrules.validations"
	runDurationName       = "businessrules.validation.duration"
	violationsInstrument  = "businessrules.violations"

	spanRuleName = "businessrules.rule"
	spanRunName  = "businessrules.validation"

	instrumentationName = "github.com/LarsArtmann/go-business-rules/listeners/otel"
)

// Attribute keys used on spans and metrics.
const (
	attrRuleName       = attribute.Key("businessrules.rule.name")
	attrRuleSeverity   = attribute.Key("businessrules.rule.severity")
	attrRulePassed     = attribute.Key("businessrules.rule.passed")
	attrRunValid       = attribute.Key("businessrules.validation.valid")
	attrViolationCount = attribute.Key("businessrules.validation.violations")
)

// Option configures the listener. Without options the listener uses the
// process-global OpenTelemetry providers (no-ops unless the application
// configured them).
type Option func(*options)

type options struct {
	tracerProvider trace.TracerProvider
	meterProvider  metric.MeterProvider
}

// WithTracerProvider overrides the tracer provider the listener creates
// spans from. Pass the application's sdktrace.TracerProvider so validation
// spans land in the same traces as everything else.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return func(o *options) {
		o.tracerProvider = provider
	}
}

// WithMeterProvider overrides the meter provider the listener creates
// instruments from. Pass the application's sdkmetric.MeterProvider so
// validation metrics are actually exported.
func WithMeterProvider(provider metric.MeterProvider) Option {
	return func(o *options) {
		o.meterProvider = provider
	}
}

// instruments bundles every metric instrument the listener writes. They are
// created once up front, so per-event work never allocates instruments and
// creation errors surface at wiring time instead of mid-run.
type instruments struct {
	evaluations  metric.Int64Counter
	ruleDuration metric.Int64Histogram
	validations  metric.Int64Counter
	runDuration  metric.Int64Histogram
	violations   metric.Int64Counter
	tracer       trace.Tracer
}

// New returns a businessrules.Listener that maps validation events onto
// OpenTelemetry spans and metrics. The returned error reports instrument
// creation failures, so applications fail at wiring time instead of
// silently losing validation telemetry.
//
// The listener is safe for concurrent use; businessrules delivers events
// synchronously per validation run.
func New(opts ...Option) (businessrules.Listener, error) {
	config := options{
		tracerProvider: otel.GetTracerProvider(),
		meterProvider:  otel.GetMeterProvider(),
	}

	for _, opt := range opts {
		opt(&config)
	}

	created, err := newInstruments(config)
	if err != nil {
		return nil, err
	}

	return created.handle, nil
}

func newInstruments(config options) (*instruments, error) {
	meter := config.meterProvider.Meter(instrumentationName)

	evaluations, err := meter.Int64Counter(
		evaluationsInstrument,
		metric.WithDescription("Number of rule checks evaluated"),
	)
	if err != nil {
		return nil, fmt.Errorf("create %s instrument: %w", evaluationsInstrument, err)
	}

	ruleDuration, err := meter.Int64Histogram(
		ruleDurationName,
		metric.WithDescription("Duration of single rule checks"),
		metric.WithUnit("ns"),
	)
	if err != nil {
		return nil, fmt.Errorf("create %s instrument: %w", ruleDurationName, err)
	}

	validations, err := meter.Int64Counter(
		validationsInstrument,
		metric.WithDescription("Number of completed validation runs"),
	)
	if err != nil {
		return nil, fmt.Errorf("create %s instrument: %w", validationsInstrument, err)
	}

	runDuration, err := meter.Int64Histogram(
		runDurationName,
		metric.WithDescription("Duration of whole validation runs"),
		metric.WithUnit("ns"),
	)
	if err != nil {
		return nil, fmt.Errorf("create %s instrument: %w", runDurationName, err)
	}

	violations, err := meter.Int64Counter(violationsInstrument, metric.WithDescription("Number of rule violations"))
	if err != nil {
		return nil, fmt.Errorf("create %s instrument: %w", violationsInstrument, err)
	}

	return &instruments{
		evaluations:  evaluations,
		ruleDuration: ruleDuration,
		validations:  validations,
		runDuration:  runDuration,
		violations:   violations,
		tracer:       config.tracerProvider.Tracer(instrumentationName),
	}, nil
}

// handle is the businessrules.Listener entry point: it switches over the
// sealed Event union and ignores unknown implementations.
func (i *instruments) handle(event businessrules.Event) {
	ctx := context.Background()

	switch ev := event.(type) {
	case businessrules.RuleEvaluated:
		i.recordRuleEvaluated(ctx, ev)
	case businessrules.ValidationCompleted:
		i.recordValidationCompleted(ctx, ev)
	}
}

func (i *instruments) recordRuleEvaluated(ctx context.Context, evaluated businessrules.RuleEvaluated) {
	attrs := metric.WithAttributes(
		attrRuleName.String(evaluated.RuleName),
		attrRuleSeverity.String(string(evaluated.Severity)),
		attrRulePassed.Bool(evaluated.Passed()),
	)

	_, span := i.tracer.Start(ctx, spanRuleName, trace.WithTimestamp(evaluated.At), trace.WithAttributes(
		attrRuleName.String(evaluated.RuleName),
		attrRuleSeverity.String(string(evaluated.Severity)),
		attrRulePassed.Bool(evaluated.Passed()),
	))

	if evaluated.Err != nil {
		span.RecordError(evaluated.Err, trace.WithTimestamp(evaluated.At.Add(evaluated.Duration)))
		span.SetStatus(codes.Error, evaluated.Err.Error())
	}

	span.End(trace.WithTimestamp(evaluated.At.Add(evaluated.Duration)))

	i.evaluations.Add(ctx, 1, attrs)
	i.ruleDuration.Record(ctx, int64(evaluated.Duration), attrs)
}

func (i *instruments) recordValidationCompleted(ctx context.Context, completed businessrules.ValidationCompleted) {
	attrs := metric.WithAttributes(
		attrRunValid.Bool(completed.Result.Valid),
		attrViolationCount.Int(completed.Result.Count()),
	)

	_, span := i.tracer.Start(ctx, spanRunName, trace.WithTimestamp(completed.At), trace.WithAttributes(
		attrRunValid.Bool(completed.Result.Valid),
		attrViolationCount.Int(completed.Result.Count()),
	))
	span.End(trace.WithTimestamp(completed.At.Add(completed.Duration)))

	i.validations.Add(ctx, 1, attrs)
	i.runDuration.Record(ctx, int64(completed.Duration), attrs)

	for _, violation := range completed.Result.ViolationErrors {
		i.violations.Add(ctx, 1, metric.WithAttributes(
			attrRuleName.String(violation.Rule.Name()),
			attrRuleSeverity.String(string(violation.Rule.Severity())),
		))
	}
}
