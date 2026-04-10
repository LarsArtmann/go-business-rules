package businessrules_test

import (
	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
)

func passingAndFailingRules() []businessrules.Rule {
	return []businessrules.Rule{
		businessrules.NonNegative("a", 1, businessrules.SeverityError),
		businessrules.Positive("b", 1, businessrules.SeverityError),
	}
}

func twoFailingRules() []businessrules.Rule {
	firstFail := businessrules.Positive("a", -1, businessrules.SeverityError)
	secondFail := businessrules.Positive("b", 0, businessrules.SeverityError)
	return []businessrules.Rule{firstFail, secondFail}
}

func anyPassingRules() []businessrules.Rule {
	rules := make([]businessrules.Rule, 2)
	rules[0] = businessrules.Positive("a", -1, businessrules.SeverityError)
	rules[1] = businessrules.Positive("b", 1, businessrules.SeverityError)
	return rules
}

var _ = Describe("Composite Builders", func() {
	Describe("All", func() {
		DescribeTable("validation",
			func(name string, rules []businessrules.Rule, shouldPass bool) {
				expectRuleResult(
					businessrules.All(name, rules, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("all pass", "all", passingAndFailingRules(), true),
			Entry("one fails", "all", twoFailingRules(), false),
		)
	})

	Describe("Any", func() {
		DescribeTable("validation",
			func(name string, rules []businessrules.Rule, shouldPass bool) {
				expectRuleResult(
					businessrules.Any(name, rules, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("one passes", "any", anyPassingRules(), true),
			Entry("all fail", "any", twoFailingRules(), false),
		)
	})

	Describe("When", func() {
		DescribeTable("validation",
			func(condition bool, shouldPass bool) {
				rule := businessrules.NotEmpty("val", "", businessrules.SeverityError)
				conditional := businessrules.When("conditional", condition, rule)
				expectRuleResult(conditional.Check(), shouldPass)
			},
			Entry("condition true", true, false),
			Entry("condition false", false, true),
		)
	})
})
