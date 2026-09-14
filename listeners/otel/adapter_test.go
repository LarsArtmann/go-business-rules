package otel_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	businessrules "github.com/LarsArtmann/go-business-rules"
	otelrules "github.com/LarsArtmann/go-business-rules/listeners/otel"
)

var _ = Describe("OpenTelemetry listener", func() {
	validatorWith := func(listener businessrules.Listener) businessrules.ValidationResultError {
		return businessrules.NewValidator().
			WithListener(listener).
			AddRule(businessrules.NewRule("passes", func() error { return nil }, businessrules.SeverityInfo, "m")).
			AddRule(businessrules.NewRule("fails", func() error { return errors.New("boom") }, businessrules.SeverityError, "m")).
			Build()
	}

	It("wires up against the global no-op providers when unconfigured", func() {
		listener, err := otelrules.New()

		Expect(err).ToNot(HaveOccurred())
		Expect(listener).ToNot(BeNil())
	})

	It("counts evaluations, violations, and completed runs", func() {
		reader := metricdata.NewManualReader()
		provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

		listener, err := otelrules.New(otelrules.WithMeterProvider(provider))
		Expect(err).ToNot(HaveOccurred())

		result := validatorWith(listener)

		Expect(result.Valid).To(BeFalse())

		var data metricdata.ResourceMetrics
		Expect(reader.Collect(context.Background(), &data)).To(Succeed())

		metrics := map[string]metricdata.Metrics{}
		for _, scope := range data.ScopeMetrics {
			for _, m := range scope.Metrics {
				metrics[m.Name] = m
			}
		}

		evaluations := counterSum(metrics, "businessrules.rule.evaluations")
		violations := counterSum(metrics, "businessrules.violations")
		validations := counterSum(metrics, "businessrules.validations")

		Expect(evaluations).To(Equal(int64(2)), "one evaluation per rule, passes included")
		Expect(violations).To(Equal(int64(1)), "one violation per failing rule")
		Expect(validations).To(Equal(int64(1)))

		Expect(metrics).To(HaveKey("businessrules.rule.duration"))
		Expect(metrics["businessrules.rule.duration"].Data.(metricdata.Histogram[int64]).Count).To(Equal(uint64(2)))
		Expect(metrics).To(HaveKey("businessrules.validation.duration"))
		Expect(metrics["businessrules.validation.duration"].Data.(metricdata.Histogram[int64]).Count).To(Equal(uint64(1)))
	})

	It("records a span per rule and one for the whole run", func() {
		recorder := tracetest.NewSpanRecorder()
		provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

		listener, err := otelrules.New(otelrules.WithTracerProvider(provider))
		Expect(err).ToNot(HaveOccurred())

		validatorWith(listener)

		spans := recorder.Ended()
		Expect(spans).To(HaveLen(3))

		byName := map[string]tracetest.ReadOnlySpan{}
		for _, span := range spans {
			byName[span.Name()] = span
		}

		Expect(byName).To(HaveKey("businessrules.rule"))
		Expect(byName).To(HaveKey("businessrules.validation"))

		failingSpan := spanForRule(spans, "fails")
		Expect(failingSpan).ToNot(BeNil())
		Expect(failingSpan.Status().Code).To(Equal(codes.Error))

		passingSpan := spanForRule(spans, "passes")
		Expect(passingSpan).ToNot(BeNil())
		Expect(passingSpan.Status().Code).ToNot(Equal(codes.Error))

		Expect(spanEndsAfterStart(byName["businessrules.rule"])).To(BeTrue())
		Expect(spanEndsAfterStart(byName["businessrules.validation"])).To(BeTrue())
	})
})

func counterSum(metrics map[string]metricdata.Metrics, name string) int64 {
	m, found := metrics[name]
	Expect(found).To(BeTrue(), "expected metric %s to be recorded", name)

	sum, ok := m.Data.(metricdata.Sum[int64])
	Expect(ok).To(BeTrue(), "metric %s must be an int64 sum", name)

	var total int64

	for _, point := range sum.DataPoints {
		total += point.Value
	}

	return total
}

func spanForRule(spans []tracetest.ReadOnlySpan, ruleName string) tracetest.ReadOnlySpan {
	for _, span := range spans {
		if attrValue(span.Attributes(), "businessrules.rule.name") == ruleName {
			return span
		}
	}

	return nil
}

func attrValue(attrs []attribute.KeyValue, key string) string {
	for _, attr := range attrs {
		if string(attr.Key) == key {
			return attr.Value.Emit()
		}
	}

	return ""
}

func spanEndsAfterStart(span tracetest.ReadOnlySpan) bool {
	return !span.EndTime().Before(span.StartTime())
}
