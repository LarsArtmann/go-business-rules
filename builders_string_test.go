package businessrules_test

import (
	"regexp"

	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("String Builders", func() {
	Describe("NotEmpty", func() {
		DescribeTable("validation",
			func(value string, shouldPass bool) {
				expectRuleResult(
					businessrules.NotEmpty("val", value, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("non-empty", "hello", true),
			Entry("empty", "", false),
		)
	})

	Describe("NotBlank", func() {
		DescribeTable("validation",
			func(value string, shouldPass bool) {
				expectRuleResult(
					businessrules.NotBlank("val", value, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("non-blank", "hello", true),
			Entry("whitespace only", "  \t\n  ", false),
			Entry("empty", "", false),
			Entry("single character", "x", true),
		)
	})

	Describe("MinLength", func() {
		DescribeTable("validation",
			func(value string, minimum int, shouldPass bool) {
				expectRuleResult(
					businessrules.MinLength("val", value, minimum, businessrules.SeverityError).
						Check(),
					shouldPass,
				)
			},
			Entry("meets minimum", "hello", 3, true),
			Entry("below minimum", "hi", 3, false),
		)
	})

	Describe("MaxLength", func() {
		DescribeTable("validation",
			func(value string, maximum int, shouldPass bool) {
				expectRuleResult(
					businessrules.MaxLength("val", value, maximum, businessrules.SeverityError).
						Check(),
					shouldPass,
				)
			},
			Entry("within limit", "hi", 5, true),
			Entry("exceeds limit", "hello world", 5, false),
		)
	})

	Describe("Matches", func() {
		It("should validate regex pattern", func() {
			pattern := regexp.MustCompile(`^[a-z]+$`)
			testMatches := func(value string, shouldPass bool) {
				expectRuleResult(
					businessrules.Matches("val", value, pattern, businessrules.SeverityError).
						Check(),
					shouldPass,
				)
			}
			testMatches("hello", true)
			testMatches("Hello123", false)
		})
	})
})
