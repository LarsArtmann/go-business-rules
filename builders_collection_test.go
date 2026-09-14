package businessrules_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

func checkNotEmptyStringSlice(slice []string, shouldPass bool) {
	emptyRule := businessrules.NotEmptySlice("val", slice, businessrules.SeverityError)
	expectRuleResult(emptyRule.Check(), shouldPass)
}

func checkNotEmptyIntSlice(slice []int, shouldPass bool) {
	validationErr := businessrules.NotEmptySlice("val", slice, businessrules.SeverityError).Check()
	expectRuleResult(validationErr, shouldPass)
}

func checkNotEmptyStringMap(m map[string]string, shouldPass bool) {
	sliceRule := businessrules.NotEmptyMap("val", m, businessrules.SeverityError)
	expectRuleResult(sliceRule.Check(), shouldPass)
}

func checkNotEmptyIntMap(m map[string]int, shouldPass bool) {
	mapErr := businessrules.NotEmptyMap("val", m, businessrules.SeverityError).Check()
	expectRuleResult(mapErr, shouldPass)
}

var _ = Describe("Collection Builders", func() {
	Describe("NotEmptySlice", func() {
		DescribeTable(
			"validation with string slice", checkNotEmptyStringSlice,
			Entry("non-empty string slice", []string{"a", "b"}, true),
			Entry("empty string slice", []string{}, false),
		)

		DescribeTable(
			"validation with int slice", checkNotEmptyIntSlice,
			Entry("non-empty int slice", []int{1}, true),
			Entry("empty int slice", []int{}, false),
		)

		It("should fail nil slice", func() {
			var nilSlice []string
			expectRuleResult(
				businessrules.NotEmptySlice("val", nilSlice, businessrules.SeverityError).Check(),
				false,
			)
		})

		It("should handle single-element slice", func() {
			expectRuleResult(
				businessrules.NotEmptySlice("val", []int{42}, businessrules.SeverityError).Check(),
				true,
			)
		})

		It("should handle large slice", func() {
			large := make([]int, 1000)
			large[999] = 1
			expectRuleResult(
				businessrules.NotEmptySlice("val", large, businessrules.SeverityError).Check(),
				true,
			)
		})
	})

	Describe("NotEmptyMap", func() {
		DescribeTable(
			"validation with string value map", checkNotEmptyStringMap,
			Entry("non-empty map", map[string]string{"a": "1"}, true),
			Entry("empty map", map[string]string{}, false),
		)

		DescribeTable(
			"validation with int value map", checkNotEmptyIntMap,
			Entry("non-empty map", map[string]int{"a": 1}, true),
			Entry("empty map", map[string]int{}, false),
		)

		It("should fail nil map", func() {
			var nilMap map[string]string
			expectRuleResult(
				businessrules.NotEmptyMap("val", nilMap, businessrules.SeverityError).Check(),
				false,
			)
		})

		It("should handle single-element map", func() {
			expectRuleResult(
				businessrules.NotEmptyMap("val", map[string]int{"key": 42}, businessrules.SeverityError).
					Check(),
				true,
			)
		})
	})

	It("should return proper error message", func() {
		err := businessrules.NotEmptySlice("items", []string{}, businessrules.SeverityError).Check()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("items"))
	})
})

var _ = Describe("Precision Builders", func() {
	Describe("MaxDecimalPlaces", func() {
		DescribeTable(
			"validation",
			func(value float64, places int, shouldPass bool) {
				expectRuleResult(
					businessrules.MaxDecimalPlaces(
						"val",
						value,
						places,
						businessrules.SeverityError,
					).Check(),
					shouldPass,
				)
			},
			Entry("within places", 3.14, 2, true),
			Entry("at the limit", 3.1, 1, true),
			Entry("integer value", 42.0, 0, true),
			Entry("exceeds places", 3.14159, 2, false),
			Entry("float error accumulation", 0.30000000000000004, 2, false),
			Entry("not a number", math.NaN(), 2, false),
			Entry("infinite", math.Inf(1), 2, false),
		)
	})

	Describe("DivisibleBy", func() {
		DescribeTable(
			"validation",
			func(value, divisor int, shouldPass bool) {
				expectRuleResult(
					businessrules.DivisibleBy(
						"val",
						value,
						divisor,
						businessrules.SeverityError,
					).Check(),
					shouldPass,
				)
			},
			Entry("evenly divisible", 24, 12, true),
			Entry("zero is divisible", 0, 12, true),
			Entry("not divisible", 25, 12, false),
			Entry("negative dividend divisible", -24, 12, true),
			Entry("negative dividend not divisible", -25, 12, false),
			Entry("negative divisor", 24, -12, true),
			Entry("zero divisor fails", 24, 0, false),
		)
	})
})
