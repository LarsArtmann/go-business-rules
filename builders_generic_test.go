package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"

	businessrules "github.com/LarsArtmann/go-business-rules"
)

func checkOneOfString(value string, allowed []string, shouldPass bool) {
	rule := businessrules.OneOf("val", value, allowed, businessrules.SeverityError)
	expectRuleResult(rule.Check(), shouldPass)
}

func checkOneOfInt(value int, allowed []int, shouldPass bool) {
	checkResult := businessrules.OneOf("val", value, allowed, businessrules.SeverityError).Check()
	expectRuleResult(checkResult, shouldPass)
}

func checkEqualsString(value, expected string, shouldPass bool) {
	r := businessrules.Equals("val", value, expected, businessrules.SeverityError)
	expectRuleResult(r.Check(), shouldPass)
}

func checkEqualsInt(value, expected int, shouldPass bool) {
	equalsRule := businessrules.Equals("val", value, expected, businessrules.SeverityError)
	equalsResult := equalsRule.Check()
	expectRuleResult(equalsResult, shouldPass)
}

func checkCustom(fn func() error, shouldPass bool) {
	result := businessrules.Custom("val", fn, businessrules.SeverityError).Check()
	expectRuleResult(result, shouldPass)
}

var _ = Describe("Generic Builders", func() {
	Describe("OneOf", func() {
		Describe("with strings", func() {
			DescribeTable(
				"validation", checkOneOfString,
				Entry("value in set", "a", []string{"a", "b"}, true),
				Entry("value not in set", "c", []string{"a", "b"}, false),
			)
		})

		Describe("with integers", func() {
			DescribeTable(
				"validation", checkOneOfInt,
				Entry("value in set", 1, []int{1, 2, 3}, true),
				Entry("value not in set", 4, []int{1, 2, 3}, false),
			)
		})
	})

	Describe("Custom", func() {
		DescribeTable(
			"validation", checkCustom,
			Entry("passing function", func() error { return nil }, true),
			Entry("failing function", func() error { return assertError("failed") }, false),
		)
	})

	Describe("Equals", func() {
		Describe("with strings", func() {
			DescribeTable(
				"validation", checkEqualsString,
				Entry("equal strings", "active", "active", true),
				Entry("different strings", "inactive", "active", false),
			)
		})

		Describe("with integers", func() {
			DescribeTable(
				"validation", checkEqualsInt,
				Entry("equal integers", 42, 42, true),
				Entry("different integers", 43, 42, false),
			)
		})
	})
})
