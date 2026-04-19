package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/artmann/businessrules"
)

var _ = Describe("Builders", func() {
	Describe("Numeric Builders", func() {
		Describe("NonNegative", func() {
			DescribeTable("validation",
				func(value float64, shouldPass bool) {
					expectRuleResult(
						businessrules.NonNegative("val", value, businessrules.SeverityError).
							Check(),
						shouldPass,
					)
				},
				Entry("zero", 0.0, true),
				Entry("positive", 10.5, true),
				Entry("negative", -1.0, false),
			)
		})

		Describe("Positive", func() {
			DescribeTable("validation",
				func(value float64, shouldPass bool) {
					expectRuleResult(
						businessrules.Positive("val", value, businessrules.SeverityError).Check(),
						shouldPass,
					)
				},
				Entry("positive integer", 1.0, true),
				Entry("positive decimal", 0.1, true),
				Entry("zero", 0.0, false),
				Entry("negative", -1.0, false),
			)
		})

		Describe("InRange", func() {
			DescribeTable("validation",
				func(value, minimum, maximum float64, shouldPass bool) {
					expectRuleResult(
						businessrules.InRange("val", value, minimum, maximum, businessrules.SeverityError).
							Check(),
						shouldPass,
					)
				},
				Entry("within range", 5.0, 0.0, 10.0, true),
				Entry("at minimum", 0.0, 0.0, 10.0, true),
				Entry("at maximum", 10.0, 0.0, 10.0, true),
				Entry("below minimum", -1.0, 0.0, 10.0, false),
				Entry("above maximum", 11.0, 0.0, 10.0, false),
			)
		})

		Describe("MinInt", func() {
			DescribeTable("validation",
				func(value, minimum int, shouldPass bool) {
					expectRuleResult(
						businessrules.MinInt("val", value, minimum, businessrules.SeverityError).
							Check(),
						shouldPass,
					)
				},
				Entry("above minimum", 5, 3, true),
				Entry("at minimum", 3, 3, true),
				Entry("below minimum", 2, 3, false),
			)
		})

		Describe("MaxInt", func() {
			DescribeTable("validation",
				func(value, maximum int, shouldPass bool) {
					expectRuleResult(
						businessrules.MaxInt("val", value, maximum, businessrules.SeverityError).
							Check(),
						shouldPass,
					)
				},
				Entry("below maximum", 5, 10, true),
				Entry("at maximum", 10, 10, true),
				Entry("above maximum", 15, 10, false),
			)
		})
	})

	Describe("Additional Numeric Builders", func() {
		Describe("GreaterThan", func() {
			DescribeTable("validation",
				func(value, minimum float64, shouldPass bool) {
					expectRuleResult(
						businessrules.GreaterThan("val", value, minimum, businessrules.SeverityError).
							Check(),
						shouldPass,
					)
				},
				Entry("above minimum", 10.0, 5.0, true),
				Entry("at minimum", 5.0, 5.0, false),
				Entry("below minimum", 3.0, 5.0, false),
			)
		})

		Describe("LessThan", func() {
			DescribeTable("validation",
				func(value, maximum float64, shouldPass bool) {
					expectRuleResult(
						businessrules.LessThan("val", value, maximum, businessrules.SeverityError).
							Check(),
						shouldPass,
					)
				},
				Entry("below maximum", 3.0, 5.0, true),
				Entry("at maximum", 5.0, 5.0, false),
				Entry("above maximum", 10.0, 5.0, false),
			)
		})
	})
})
