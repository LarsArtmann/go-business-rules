package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/artmann/businessrules"
)

var _ = Describe("Composite Builders", func() {
	It("should validate All - all pass", func() {
		rules := []businessrules.Rule{
			businessrules.NonNegative("a", 1, businessrules.SeverityError),
			businessrules.Positive("b", 1, businessrules.SeverityError),
		}
		Expect(
			businessrules.All("all", rules, businessrules.SeverityError).Check(),
		).To(Succeed())
	})

	It("should validate All - one fails", func() {
		rules := []businessrules.Rule{
			businessrules.NonNegative("a", 1, businessrules.SeverityError),
			businessrules.Positive("b", -1, businessrules.SeverityError),
		}
		Expect(
			businessrules.All("all", rules, businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should validate Any - one passes", func() {
		rules := []businessrules.Rule{
			businessrules.Positive("a", -1, businessrules.SeverityError),
			businessrules.Positive("b", 1, businessrules.SeverityError),
		}
		Expect(
			businessrules.Any("any", rules, businessrules.SeverityError).Check(),
		).To(Succeed())
	})

	It("should validate Any - all fail", func() {
		rules := []businessrules.Rule{
			businessrules.Positive("a", -1, businessrules.SeverityError),
			businessrules.Positive("b", 0, businessrules.SeverityError),
		}
		Expect(
			businessrules.Any("any", rules, businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should validate When - condition true", func() {
		rule := businessrules.NotEmpty("val", "", businessrules.SeverityError)
		conditional := businessrules.When("conditional", true, rule)
		Expect(conditional.Check()).ToNot(Succeed())
	})

	It("should validate When - condition false", func() {
		rule := businessrules.NotEmpty("val", "", businessrules.SeverityError)
		conditional := businessrules.When("conditional", false, rule)
		Expect(conditional.Check()).To(Succeed())
	})
})
