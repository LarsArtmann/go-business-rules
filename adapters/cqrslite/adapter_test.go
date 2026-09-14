package cqrslite_test

import (
	"context"
	"errors"

	"github.com/larsartmann/go-cqrs-lite/event/v4"
	"github.com/larsartmann/go-cqrs-lite/event/v4/eventtest"
	"github.com/larsartmann/go-cqrs-lite/id/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
	cqrslite "github.com/LarsArtmann/go-business-rules/adapters/cqrslite"
)

var _ = Describe("Bus Listener", func() {
	var (
		bus        *eventtest.FakeBus
		streamID   id.StreamID
		streamType id.StreamType
	)

	BeforeEach(func() {
		bus = eventtest.NewFakeBus()
		streamID = id.NewStreamID()
		streamType = "Order"
	})

	type publishedEvent struct {
		typ     event.Type
		version event.Version
		event   event.Event
	}

	collect := func(types ...event.Type) (func() []publishedEvent, func()) {
		var collected []publishedEvent

		for _, typ := range types {
			Expect(bus.Subscribe(typ, func(_ context.Context, e event.Event) error {
				collected = append(collected, publishedEvent{typ: typ, version: e.Version(), event: e})

				return nil
			})).To(Succeed())
		}

		return func() []publishedEvent { return collected }, func() { collected = nil }
	}

	runValidation := func(listener businessrules.Listener) businessrules.ValidationResultError {
		return businessrules.NewValidator().
			WithListener(listener).
			AddRule(businessrules.NewRule("balance", func() error {
				return errors.New("insufficient balance")
			}, businessrules.SeverityError, "balance must cover the order")).
			AddRule(businessrules.NewRule("items", func() error {
				return nil
			}, businessrules.SeverityInfo, "at least one item")).
			Build()
	}

	It("publishes per-rule events with round-trippable payloads", func() {
		collected, _ := collect(cqrslite.TypeRuleEvaluated)

		listener := cqrslite.NewBusListener(bus, streamID, streamType)
		runValidation(listener)

		Expect(collected()).To(HaveLen(2))

		first, err := event.DecodePayloadAuto[cqrslite.RuleEvaluatedData](collected()[0].event)
		Expect(err).ToNot(HaveOccurred())
		Expect(first.RuleName).To(Equal("balance"))
		Expect(first.Severity).To(Equal("error"))
		Expect(first.Passed).To(BeFalse())
		Expect(first.ErrorMessage).To(Equal("insufficient balance"))

		second, err := event.DecodePayloadAuto[cqrslite.RuleEvaluatedData](collected()[1].event)
		Expect(err).ToNot(HaveOccurred())
		Expect(second.RuleName).To(Equal("items"))
		Expect(second.Passed).To(BeTrue())
		Expect(second.ErrorMessage).To(BeEmpty())
	})

	It("tags events with the stream identity and increasing versions", func() {
		collected, _ := collect(cqrslite.TypeRuleEvaluated, cqrslite.TypeValidationCompleted)

		runValidation(cqrslite.NewBusListener(bus, streamID, streamType))

		Expect(collected()).To(HaveLen(3))

		for index, published := range collected() {
			Expect(published.event.StreamID()).To(Equal(streamID))
			Expect(published.event.StreamType()).To(Equal(streamType))
			Expect(published.version).To(Equal(event.Version(index + 1)))
		}

		Expect(collected()[2].typ).To(Equal(cqrslite.TypeValidationCompleted))
	})

	It("publishes a terminal completed event summarizing the run", func() {
		collected, _ := collect(cqrslite.TypeValidationCompleted)

		result := runValidation(cqrslite.NewBusListener(bus, streamID, streamType))

		Expect(collected()).To(HaveLen(1))

		data, err := event.DecodePayloadAuto[cqrslite.ValidationCompletedData](collected()[0].event)
		Expect(err).ToNot(HaveOccurred())
		Expect(data.Valid).To(Equal(result.Valid))
		Expect(data.Valid).To(BeFalse())
		Expect(data.ViolationCount).To(Equal(1))
	})

	It("reports publish failures to the configured error handler without aborting validation", func() {
		failing := eventtest.NewFakeBus()
		publishError := errors.New("broker down")

		var observed []error

		listener := cqrslite.NewBusListener(failing, streamID, streamType,
			cqrslite.WithPublishErrorHandler(func(err error) { observed = append(observed, err) }))
		_ = failing.Subscribe(cqrslite.TypeRuleEvaluated, func(_ context.Context, _ event.Event) error {
			return publishError
		})

		result := runValidation(listener)

		Expect(result.Valid).To(BeFalse())
		Expect(observed).To(HaveLen(2), "one publish error per evaluated rule")

		for _, err := range observed {
			Expect(err).To(MatchError(publishError))
		}
	})
})
