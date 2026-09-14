package businessrules_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	businessrules "github.com/LarsArtmann/go-business-rules"
)

var _ = Describe("Validation Stream", func() {
	var boom error

	BeforeEach(func() {
		boom = errors.New("insufficient balance")
	})

	slowRule := func(name string, delay time.Duration, err error) businessrules.Rule {
		return businessrules.NewRule(name, func() error {
			time.Sleep(delay)

			return err
		}, businessrules.SeverityError, "slow check")
	}

	drain := func(events <-chan businessrules.Event) []businessrules.Event {
		var received []businessrules.Event

		for event := range events {
			received = append(received, event)
		}

		return received
	}

	It("yields one event per rule, a terminal completed event, and closes the channel", func() {
		events := businessrules.NewValidator().
			AddRule(passingRule("r1", businessrules.SeverityInfo, "m")).
			AddRule(failingRule("r2", func() error { return boom }, businessrules.SeverityError, "m")).
			Stream(context.Background())

		received := drain(events)

		Expect(received).To(HaveLen(3))

		completed, ok := received[2].(businessrules.ValidationCompleted)
		Expect(ok).To(BeTrue(), "last event must be ValidationCompleted")
		Expect(completed.Result.Valid).To(BeFalse())
		Expect(completed.Result.ViolationErrors).To(HaveLen(1))
		Expect(completed.Result.ViolationErrors[0].Rule.Name()).To(Equal("r2"))
	})

	It("streams events in completion order, not rule order", func() {
		events := businessrules.NewValidator().
			AddRule(slowRule("slow_rule", 60*time.Millisecond, nil)).
			AddRule(passingRule("fast_rule", businessrules.SeverityInfo, "m")).
			Stream(context.Background())

		received := drain(events)

		Expect(received).To(HaveLen(3))
		Expect(received[0].(businessrules.RuleEvaluated).RuleName).To(Equal("fast_rule"))
		Expect(received[1].(businessrules.RuleEvaluated).RuleName).To(Equal("slow_rule"))
	})

	It("reports the same violations as Build, in rule order", func() {
		rules := []businessrules.Rule{
			failingRule("r1", func() error { return boom }, businessrules.SeverityCritical, "m"),
			passingRule("r2", businessrules.SeverityInfo, "m"),
			failingRule("r3", func() error { return boom }, businessrules.SeverityWarning, "m"),
		}

		streamed := businessrules.NewValidator().AddRules(rules...).Stream(context.Background())
		received := drain(streamed)

		completed := received[len(received)-1].(businessrules.ValidationCompleted)
		sequential := businessrules.NewValidator().AddRules(rules...).Build()

		Expect(completed.Result.Valid).To(Equal(sequential.Valid))
		Expect(completed.Result.Count()).To(Equal(sequential.Count()))
		Expect(completed.Result.ViolationErrors[0].Rule.Name()).To(Equal("r1"))
		Expect(completed.Result.ViolationErrors[1].Rule.Name()).To(Equal("r3"))
	})

	It("delivers events on the channel instead of to registered listeners", func() {
		recorder := &eventRecorder{}

		events := businessrules.NewValidator().
			WithListener(recorder.listen).
			AddRule(passingRule("r1", businessrules.SeverityInfo, "m")).
			Stream(context.Background())

		received := drain(events)

		Expect(received).To(HaveLen(2))
		Expect(recorder.events).To(BeEmpty())
	})

	It("runs no rules when the context is already canceled", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		events := businessrules.NewValidator().
			AddRule(passingRule("r1", businessrules.SeverityInfo, "m")).
			Stream(ctx)

		for _, event := range drain(events) {
			_, isRuleEvent := event.(businessrules.RuleEvaluated)
			Expect(isRuleEvent).To(BeFalse(), "no rule check may run under a canceled context")
		}
	})

	It("stops scheduling and closes the stream after mid-run cancellation", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		events := businessrules.NewValidator().
			AddRule(slowRule("slow_rule", 60*time.Millisecond, nil)).
			AddRule(slowRule("never_scheduled", 60*time.Millisecond, nil)).
			Stream(ctx)

		var received []businessrules.Event

		Eventually(func() int {
			select {
			case event, ok := <-events:
				if !ok {
					return len(received)
				}

				received = append(received, event)

				if len(received) == 1 {
					cancel()
				}

				return -1
			case <-time.After(2 * time.Second):
				Fail("stream did not produce an event in time")

				return -1
			}
		}, 3*time.Second).Should(BeNumerically(">=", 0), "stream must close after cancellation")

		Expect(received).ToNot(BeEmpty())
	})
})
