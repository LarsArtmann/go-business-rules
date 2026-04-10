package businessrules_test

import (
	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Composite Builders", func() {
	expectCompositeResult := func(result error, shouldPass bool) {
		if shouldPass {
			Expect(result).To(Succeed())
		} else {
			Expect(result).ToNot(Succeed())
		}
	}

	Describe("All", func() {
		DescribeTable("validation",
			func(name string, rules []businessrules.Rule, shouldPass bool) {
				expectCompositeResult(
					businessrules.All(name, rules, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("all pass", "all", []businessrules.Rule{
				businessrules.NonNegative("a", 1, businessrules.SeverityError),
				businessrules.Positive("b", 1, businessrules.SeverityError),
			}, true),
			Entry("one fails", "all", []businessrules.Rule{
				businessrules.NonNegative("a", 1, businessrules.SeverityError),
				businessrules.Positive("b", -1, businessrules.SeverityError),
			}, false),
		)
	})

	Describe("Any", func() {
		DescribeTable("validation",
			func(name string, rules []businessrules.Rule, shouldPass bool) {
				expectCompositeResult(
					businessrules.Any(name, rules, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("one passes", "any", []businessrules.Rule{
				businessrules.Positive("a", -1, businessrules.SeverityError),
				businessrules.Positive("b", 1, businessrules.SeverityError),
			}, true),
			Entry("all fail", "any", []businessrules.Rule{
				businessrules.Positive("a", -1, businessrules.SeverityError),
				businessrules.Positive("b", 0, businessrules.SeverityError),
			}, false),
		)
	})

	Describe("When", func() {
		DescribeTable("validation",
			func(condition bool, shouldPass bool) {
				rule := businessrules.NotEmpty("val", "", businessrules.SeverityError)
				conditional := businessrules.When("conditional", condition, rule)
				expectCompositeResult(conditional.Check(), shouldPass)
			},
			Entry("condition true", true, false),
			Entry("condition false", false, true),
		)
	})
})
