package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
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
		DescribeTable(
			"validation",
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
		DescribeTable(
			"validation",
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
		DescribeTable(
			"validation",
			func(condition, shouldPass bool) {
				rule := businessrules.NotEmpty("val", "", businessrules.SeverityError)
				conditional := businessrules.When("conditional", condition, rule)
				expectRuleResult(conditional.Check(), shouldPass)
			},
			Entry("condition true", true, false),
			Entry("condition false", false, true),
		)
	})
})

var _ = Describe("Extended Composite Builders", func() {
	Describe("Not", func() {
		It("passes when the inner rule fails", func() {
			forbidden := businessrules.Equals("status", "active", "banned", businessrules.SeverityError)
			expectRuleResult(
				businessrules.Not("not-banned", forbidden, businessrules.SeverityError).Check(),
				true,
			)
		})

		It("fails when the inner rule passes", func() {
			forbidden := businessrules.Equals("status", "active", "active", businessrules.SeverityError)
			expectRuleResult(
				businessrules.Not("not-active", forbidden, businessrules.SeverityError).Check(),
				false,
			)
		})
	})

	Describe("Or", func() {
		It("passes when any variadic sub-rule passes", func() {
			expectRuleResult(
				businessrules.Or(
					"contact",
					businessrules.SeverityError,
					businessrules.NotEmpty("email", "", businessrules.SeverityError),
					businessrules.NotEmpty("phone", "+49 123", businessrules.SeverityError),
				).Check(),
				true,
			)
		})

		It("fails when all variadic sub-rules fail", func() {
			expectRuleResult(
				businessrules.Or(
					"contact",
					businessrules.SeverityError,
					businessrules.NotEmpty("email", "", businessrules.SeverityError),
					businessrules.NotEmpty("phone", "", businessrules.SeverityError),
				).Check(),
				false,
			)
		})

		It("fails vacuously with no sub-rules", func() {
			expectRuleResult(
				businessrules.Or("contact", businessrules.SeverityError).Check(),
				false,
			)
		})
	})

	Describe("Xor", func() {
		It("passes when exactly one sub-rule passes", func() {
			expectRuleResult(
				businessrules.Xor(
					"delivery",
					businessrules.NotEmpty("email", "a@b.c", businessrules.SeverityError),
					businessrules.NotEmpty("phone", "", businessrules.SeverityError),
					businessrules.SeverityError,
				).Check(),
				true,
			)
		})

		It("fails when both sub-rules pass", func() {
			expectRuleResult(
				businessrules.Xor(
					"delivery",
					businessrules.NotEmpty("email", "a@b.c", businessrules.SeverityError),
					businessrules.NotEmpty("phone", "+49 123", businessrules.SeverityError),
					businessrules.SeverityError,
				).Check(),
				false,
			)
		})

		It("fails when neither sub-rule passes", func() {
			expectRuleResult(
				businessrules.Xor(
					"delivery",
					businessrules.NotEmpty("email", "", businessrules.SeverityError),
					businessrules.NotEmpty("phone", "", businessrules.SeverityError),
					businessrules.SeverityError,
				).Check(),
				false,
			)
		})
	})
})
