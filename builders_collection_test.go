package businessrules_test

import (
	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Collection Builders", func() {
	expectResult := func(result error, shouldPass bool) {
		if shouldPass {
			Expect(result).To(Succeed())
		} else {
			Expect(result).ToNot(Succeed())
		}
	}

	Describe("NotEmptySlice", func() {
		DescribeTable("validation with string slice",
			func(slice []string, shouldPass bool) {
				expectResult(
					businessrules.NotEmptySlice("val", slice, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("non-empty string slice", []string{"a", "b"}, true),
			Entry("empty string slice", []string{}, false),
		)

		DescribeTable("validation with int slice",
			func(slice []int, shouldPass bool) {
				expectResult(
					businessrules.NotEmptySlice("val", slice, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("non-empty int slice", []int{1}, true),
			Entry("empty int slice", []int{}, false),
		)

		It("should fail nil slice", func() {
			var nilSlice []string
			expectResult(
				businessrules.NotEmptySlice("val", nilSlice, businessrules.SeverityError).Check(),
				false,
			)
		})

		It("should handle single-element slice", func() {
			expectResult(
				businessrules.NotEmptySlice("val", []int{42}, businessrules.SeverityError).Check(),
				true,
			)
		})

		It("should handle large slice", func() {
			large := make([]int, 1000)
			large[999] = 1
			expectResult(
				businessrules.NotEmptySlice("val", large, businessrules.SeverityError).Check(),
				true,
			)
		})
	})

	Describe("NotEmptyMap", func() {
		DescribeTable("validation with string value map",
			func(m map[string]string, shouldPass bool) {
				expectResult(
					businessrules.NotEmptyMap("val", m, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("non-empty map", map[string]string{"a": "1"}, true),
			Entry("empty map", map[string]string{}, false),
		)

		DescribeTable("validation with int value map",
			func(m map[string]int, shouldPass bool) {
				expectResult(
					businessrules.NotEmptyMap("val", m, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("non-empty map", map[string]int{"a": 1}, true),
			Entry("empty map", map[string]int{}, false),
		)

		It("should fail nil map", func() {
			var nilMap map[string]string
			expectResult(
				businessrules.NotEmptyMap("val", nilMap, businessrules.SeverityError).Check(),
				false,
			)
		})

		It("should handle single-element map", func() {
			expectResult(
				businessrules.NotEmptyMap("val", map[string]int{"key": 42}, businessrules.SeverityError).
					Check(),
				true,
			)
		})
	})

	It("should return proper error message", func() {
		err := businessrules.NotEmptySlice("items", []string{}, businessrules.SeverityError).Check()
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("items"))
	})
})
