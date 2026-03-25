package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/artmann/businessrules"
)

var _ = Describe("ValidatorBuilder", func() {
	It("should build valid result when all rules pass", func() {
		rule := businessrules.NewRule(
			"pass",
			func() error { return nil },
			businessrules.SeverityInfo,
			"pass",
		)
		result := businessrules.NewValidator().AddRule(rule).Build()
		Expect(result.Valid).To(BeTrue())
	})

	It("should build invalid result when rule fails", func() {
		rule := businessrules.NewRule(
			"fail",
			func() error { return assertError("failed") },
			businessrules.SeverityError,
			"fail",
		)
		result := businessrules.NewValidator().AddRule(rule).Build()
		Expect(result.Valid).To(BeFalse())
		Expect(result.HasErrors()).To(BeTrue())
	})

	It("should support AddRules for multiple rules", func() {
		rule1 := businessrules.NewRule(
			"r1",
			func() error { return nil },
			businessrules.SeverityInfo,
			"r1",
		)
		rule2 := businessrules.NewRule(
			"r2",
			func() error { return nil },
			businessrules.SeverityInfo,
			"r2",
		)
		result := businessrules.NewValidator().AddRules(rule1, rule2).Build()
		Expect(result.Valid).To(BeTrue())
	})
})
