package businessrules_test

import (
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

var _ = Describe("String Builders", func() {
	Describe("NotEmpty", func() {
		DescribeTable(
			"validation",
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
		DescribeTable(
			"validation",
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
		DescribeTable(
			"validation",
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
		DescribeTable(
			"validation",
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

var _ = Describe("Extended String Builders", func() {
	Describe("Contains", func() {
		DescribeTable(
			"validation",
			func(value, substring string, shouldPass bool) {
				expectRuleResult(
					businessrules.Contains("val", value, substring, businessrules.SeverityError).
						Check(),
					shouldPass,
				)
			},
			Entry("contains substring", "hello world", "world", true),
			Entry("substring at start", "hello", "hel", true),
			Entry("substring at end", "hello", "llo", true),
			Entry("missing substring", "hello", "xyz", false),
			Entry("case sensitive", "Hello", "hello", false),
		)
	})

	Describe("LengthRange", func() {
		DescribeTable(
			"validation",
			func(value string, minimum, maximum int, shouldPass bool) {
				expectRuleResult(
					businessrules.LengthRange(
						"val",
						value,
						minimum,
						maximum,
						businessrules.SeverityError,
					).Check(),
					shouldPass,
				)
			},
			Entry("within range", "abc", 2, 5, true),
			Entry("at minimum", "ab", 2, 5, true),
			Entry("at maximum", "abcde", 2, 5, true),
			Entry("below range", "a", 2, 5, false),
			Entry("above range", "abcdef", 2, 5, false),
		)
	})

	Describe("Required", func() {
		DescribeTable(
			"validation",
			func(value string, shouldPass bool) {
				expectRuleResult(
					businessrules.Required("val", value, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("visible content", "hello", true),
			Entry("single character", "x", true),
			Entry("empty", "", false),
			Entry("whitespace only", "  \t\n  ", false),
		)
	})

	Describe("MatchesFunc", func() {
		DescribeTable(
			"validation",
			func(value string, shouldPass bool) {
				isLowercase := func(s string) bool { return s == strings.ToLower(s) }
				expectRuleResult(
					businessrules.MatchesFunc(
						"val",
						value,
						isLowercase,
						businessrules.SeverityError,
					).Check(),
					shouldPass,
				)
			},
			Entry("satisfies predicate", "hello", true),
			Entry("violates predicate", "Hello", false),
		)
	})
})
