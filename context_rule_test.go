package businessrules_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

type markerKey struct{}

var _ = Describe("Context-aware rules", func() {
	It("exposes the same metadata as a plain rule", func() {
		rule := businessrules.NewContextRule(
			"email_mx",
			func(context.Context) error { return nil },
			businessrules.SeverityError,
			"mail exchanger must resolve",
		)

		Expect(rule.Name()).To(Equal("email_mx"))
		Expect(rule.Severity()).To(Equal(businessrules.SeverityError))
		Expect(rule.Message()).To(Equal("mail exchanger must resolve"))
	})

	It("runs Check with a background context that never cancels", func() {
		var sawBackground bool

		rule := businessrules.NewContextRule("probe", func(ctx context.Context) error {
			sawBackground = ctx == context.Background()

			return nil
		}, businessrules.SeverityInfo, "m")

		Expect(rule.Check()).To(Succeed())
		Expect(sawBackground).To(BeTrue(), "Check must run with a non-canceling context")
	})

	It("works anywhere a plain Rule is expected", func() {
		var rule businessrules.Rule = businessrules.NewContextRule(
			"always_passes",
			func(context.Context) error { return nil },
			businessrules.SeverityError,
			"m",
		)

		result := businessrules.NewValidator().AddRule(rule).Build()

		Expect(result.Valid).To(BeTrue())
	})

	It("keeps the context-aware check when derived via With helpers", func() {
		base := businessrules.NewContextRule("base", func(context.Context) error {
			return errors.New("boom")
		}, businessrules.SeverityWarning, "m")

		derived := base.WithName("renamed").
			WithSeverity(businessrules.SeverityCritical).
			WithMessage("other")

		Expect(derived.Name()).To(Equal("renamed"))
		Expect(derived.Severity()).To(Equal(businessrules.SeverityCritical))
		Expect(derived.Message()).To(Equal("other"))
		Expect(derived.Check()).ToNot(Succeed())

		_, stillContextAware := businessrules.Rule(derived).(businessrules.ContextRule)
		Expect(stillContextAware).To(BeTrue(), "derived rules must stay context-aware")
	})

	Describe("Stream integration", func() {
		It("passes the streaming context to CheckContext", func() {
			rule := businessrules.NewContextRule("probe", func(ctx context.Context) error {
				if ctx.Value(markerKey{}) != "stream" {
					return errors.New("stream context not visible to the check")
				}

				return nil
			}, businessrules.SeverityError, "m")

			ctx := context.WithValue(context.Background(), markerKey{}, "stream")

			var received []businessrules.Event

			for event := range businessrules.NewValidator().AddRule(rule).Stream(ctx) {
				received = append(received, event)
			}

			completed, ok := received[1].(businessrules.ValidationCompleted)
			Expect(ok).To(BeTrue(), "last event must be ValidationCompleted")
			Expect(completed.Result.Valid).To(BeTrue())
		})

		It("interrupts a check through a context derived from the streaming context", func() {
			rule := businessrules.NewContextRule("cancellable", func(ctx context.Context) error {
				child, childCancel := context.WithCancel(ctx)
				childCancel()

				<-child.Done()

				return child.Err()
			}, businessrules.SeverityError, "slow check")

			var received []businessrules.Event

			for event := range businessrules.NewValidator().AddRule(rule).Stream(context.Background()) {
				received = append(received, event)
			}

			Expect(received).To(HaveLen(2))

			ruleEvent, ok := received[0].(businessrules.RuleEvaluated)
			Expect(ok).To(BeTrue())
			Expect(errors.Is(ruleEvent.Err, context.Canceled)).To(BeTrue(),
				"a canceled check must return promptly with the context error")

			completed := received[1].(businessrules.ValidationCompleted)
			Expect(completed.Result.Valid).To(BeFalse())
			Expect(completed.Result.Count()).To(Equal(1))
		})
	})
})
