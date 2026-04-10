package businessrules_test

import (
	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidatorBuilder", func() {
	It("should build valid result when all rules pass", func() {
		rule := passingRule("pass", businessrules.SeverityInfo, "pass")
		result := businessrules.NewValidator().AddRule(rule).Build()
		Expect(result.Valid).To(BeTrue())
	})

	It("should build invalid result when rule fails", func() {
		rule := failingRule("fail", func() error { return assertError("failed") }, businessrules.SeverityError, "fail")
		result := businessrules.NewValidator().AddRule(rule).Build()
		Expect(result.Valid).To(BeFalse())
		Expect(result.HasErrors()).To(BeTrue())
	})

	It("should support AddRules for multiple rules", func() {
		rule1 := passingRule("r1", businessrules.SeverityInfo, "r1")
		rule2 := passingRule("r2", businessrules.SeverityInfo, "r2")
		result := businessrules.NewValidator().AddRules(rule1, rule2).Build()
		Expect(result.Valid).To(BeTrue())
	})
})
