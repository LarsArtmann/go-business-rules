package businessrules_test

import (
	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Generic Builders", func() {
	Describe("OneOf", func() {
		Describe("with strings", func() {
			DescribeTable("validation",
				func(value string, allowed []string, shouldPass bool) {
					expectRuleResult(
						businessrules.OneOf("val", value, allowed, businessrules.SeverityError).Check(),
						shouldPass,
					)
				},
				Entry("value in set", "a", []string{"a", "b"}, true),
				Entry("value not in set", "c", []string{"a", "b"}, false),
			)
		})

		Describe("with integers", func() {
			DescribeTable("validation",
				func(value int, allowed []int, shouldPass bool) {
					expectRuleResult(
						businessrules.OneOf("val", value, allowed, businessrules.SeverityError).Check(),
						shouldPass,
					)
				},
				Entry("value in set", 1, []int{1, 2, 3}, true),
				Entry("value not in set", 4, []int{1, 2, 3}, false),
			)
		})
	})

	Describe("Custom", func() {
		DescribeTable("validation",
			func(fn func() error, shouldPass bool) {
				expectRuleResult(
					businessrules.Custom("val", fn, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("passing function", func() error { return nil }, true),
			Entry("failing function", func() error { return assertError("failed") }, false),
		)
	})

	Describe("Equals", func() {
		Describe("with strings", func() {
			DescribeTable("validation",
				func(value, expected string, shouldPass bool) {
					expectRuleResult(
						businessrules.Equals("val", value, expected, businessrules.SeverityError).Check(),
						shouldPass,
					)
				},
				Entry("equal strings", "active", "active", true),
				Entry("different strings", "inactive", "active", false),
			)
		})

		Describe("with integers", func() {
			DescribeTable("validation",
				func(value, expected int, shouldPass bool) {
					expectRuleResult(
						businessrules.Equals("val", value, expected, businessrules.SeverityError).Check(),
						shouldPass,
					)
				},
				Entry("equal integers", 42, 42, true),
				Entry("different integers", 43, 42, false),
			)
		})
	})
})
