package businessrules_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	businessrules "github.com/LarsArtmann/go-business-rules"
)

type eventRecorder struct {
	events []businessrules.Event
}

func (r *eventRecorder) listen(event businessrules.Event) {
	r.events = append(r.events, event)
}

func (r *eventRecorder) ruleEvents() []businessrules.RuleEvaluated {
	var evaluated []businessrules.RuleEvaluated

	for _, event := range r.events {
		if re, ok := event.(businessrules.RuleEvaluated); ok {
			evaluated = append(evaluated, re)
		}
	}

	return evaluated
}

func (r *eventRecorder) completed() (businessrules.ValidationCompleted, bool) {
	for _, event := range r.events {
		if completed, ok := event.(businessrules.ValidationCompleted); ok {
			return completed, true
		}
	}

	return businessrules.ValidationCompleted{}, false
}

var _ = Describe("Validation Events", func() {
	var (
		recorder *eventRecorder
		boom     error
	)

	BeforeEach(func() {
		recorder = &eventRecorder{}
		boom = errors.New("weight must be positive")
	})

	Describe("Build with listeners", func() {
		It("reports every evaluated rule in rule order, passes and failures alike", func() {
			result := businessrules.NewValidator().
				WithListener(recorder.listen).
				AddRule(passingRule("weight_positive", businessrules.SeverityError, "weight must be positive")).
				AddRule(failingRule("weight_range", func() error { return boom }, businessrules.SeverityWarning, "weight out of range")).
				Build()

			Expect(result.Valid).To(BeFalse())

			evaluated := recorder.ruleEvents()
			Expect(evaluated).To(HaveLen(2))

			Expect(evaluated[0].RuleName).To(Equal("weight_positive"))
			Expect(evaluated[0].Passed()).To(BeTrue())
			Expect(evaluated[0].Err).ToNot(HaveOccurred())
			Expect(evaluated[0].Severity).To(Equal(businessrules.SeverityError))

			Expect(evaluated[1].RuleName).To(Equal("weight_range"))
			Expect(evaluated[1].Passed()).To(BeFalse())
			Expect(evaluated[1].Err).To(MatchError(boom))
			Expect(evaluated[1].Severity).To(Equal(businessrules.SeverityWarning))
		})

		It("timestamps every rule event at check start with a non-negative duration", func() {
			businessrules.NewValidator().
				WithListener(recorder.listen).
				AddRule(passingRule("any_rule", businessrules.SeverityInfo, "always fine")).
				Build()

			evaluated := recorder.ruleEvents()
			Expect(evaluated).To(HaveLen(1))
			Expect(evaluated[0].At).ToNot(BeZero())
			Expect(evaluated[0].Duration).To(BeNumerically(">=", 0))
		})

		It("closes the run with a ValidationCompleted event carrying the returned result", func() {
			result := businessrules.NewValidator().
				WithListener(recorder.listen).
				AddRule(failingRule("weight_range", func() error { return boom }, businessrules.SeverityError, "weight out of range")).
				Build()

			completed, found := recorder.completed()
			Expect(found).To(BeTrue(), "expected a terminal ValidationCompleted event")
			Expect(completed.Result.Valid).To(BeFalse())
			Expect(completed.Result.ViolationErrors).To(HaveLen(1))
			Expect(completed.Result.Valid).To(Equal(result.Valid))
			Expect(completed.Duration).To(BeNumerically(">=", 0))
			Expect(completed.At).ToNot(BeZero())
		})

		It("emits exactly one RuleEvaluated per rule plus one ValidationCompleted", func() {
			businessrules.NewValidator().
				WithListener(recorder.listen).
				AddRule(passingRule("r1", businessrules.SeverityInfo, "m")).
				AddRule(passingRule("r2", businessrules.SeverityWarning, "m")).
				AddRule(passingRule("r3", businessrules.SeverityError, "m")).
				Build()

			Expect(recorder.events).To(HaveLen(4))
			_, isCompleted := recorder.events[3].(businessrules.ValidationCompleted)
			Expect(isCompleted).To(BeTrue(), "last event must be ValidationCompleted")
		})

		It("delivers the full sequence to every registered listener", func() {
			second := &eventRecorder{}

			businessrules.NewValidator().
				WithListener(recorder.listen, second.listen).
				AddRule(passingRule("r1", businessrules.SeverityInfo, "m")).
				Build()

			Expect(recorder.events).To(HaveLen(2))
			Expect(second.events).To(HaveLen(2))
			Expect(second.events[0].(businessrules.RuleEvaluated).RuleName).To(Equal("r1"))
		})

		It("emits events synchronously before Build returns", func() {
			sawEventDuringBuild := false

			probing := func(event businessrules.Event) {
				if _, ok := event.(businessrules.RuleEvaluated); ok {
					sawEventDuringBuild = true
				}
			}

			businessrules.NewValidator().
				WithListener(probing).
				AddRule(passingRule("r1", businessrules.SeverityInfo, "m")).
				Build()

			Expect(sawEventDuringBuild).To(BeTrue())
		})

		It("returns the same result as an unobserved build", func() {
			rules := []businessrules.Rule{
				passingRule("r1", businessrules.SeverityInfo, "m"),
				failingRule("r2", func() error { return boom }, businessrules.SeverityCritical, "m"),
			}

			observed := businessrules.NewValidator().WithListener(recorder.listen).AddRules(rules...).Build()
			unobserved := businessrules.NewValidator().AddRules(rules...).Build()

			Expect(observed.Valid).To(Equal(unobserved.Valid))
			Expect(observed.Count()).To(Equal(unobserved.Count()))
			Expect(observed.ViolationErrors).To(HaveLen(len(unobserved.ViolationErrors)))
		})
	})
})
