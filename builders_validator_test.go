package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/artmann/businessrules"
)

var _ = Describe("ValidatorBuilder", func() {
	It("should build valid result when all rules pass", func() {
		rule := passingRule("pass", businessrules.SeverityInfo, "pass")
		result := businessrules.NewValidator().AddRule(rule).Build()
		expectValid(result)
	})

	It("should build invalid result when rule fails", func() {
		rule := failingRule(
			"fail",
			func() error { return assertError("failed") },
			businessrules.SeverityError,
			"fail",
		)
		result := businessrules.NewValidator().AddRule(rule).Build()
		expectInvalid(result)
		expectHasErrors(result)
	})

	It("should support AddRules for multiple rules", func() {
		rule1 := passingRule("r1", businessrules.SeverityInfo, "r1")
		rule2 := passingRule("r2", businessrules.SeverityInfo, "r2")
		result := businessrules.NewValidator().AddRules(rule1, rule2).Build()
		expectValid(result)
	})
})
