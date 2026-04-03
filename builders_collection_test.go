package businessrules_test

import (
	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Collection Builders", func() {
	It("should validate NotEmptySlice with string slice", func() {
		slice := []string{"a", "b"}
		Expect(
			businessrules.NotEmptySlice("val", slice, businessrules.SeverityError).Check(),
		).To(Succeed())
	})

	It("should validate NotEmptySlice with int slice", func() {
		Expect(
			businessrules.NotEmptySlice("val", []int{1}, businessrules.SeverityError).Check(),
		).To(Succeed())
	})

	It("should fail NotEmptySlice with empty slice", func() {
		emptySlice := []string{}
		Expect(
			businessrules.NotEmptySlice("val", emptySlice, businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should fail NotEmptySlice with nil slice", func() {
		var nilSlice []string
		Expect(
			businessrules.NotEmptySlice("val", nilSlice, businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should validate NotEmptyMap", func() {
		m := map[string]int{"a": 1}
		Expect(
			businessrules.NotEmptyMap("val", m, businessrules.SeverityError).Check(),
		).To(Succeed())
	})

	It("should fail NotEmptyMap with empty map", func() {
		emptyMap := map[string]string{}
		Expect(
			businessrules.NotEmptyMap("val", emptyMap, businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should fail NotEmptyMap with nil map", func() {
		var nilMap map[string]string
		Expect(
			businessrules.NotEmptyMap("val", nilMap, businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should return proper error message", func() {
		err := businessrules.NotEmptySlice("items", []string{}, businessrules.SeverityError).Check()
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("items"))
	})

	Context("edge cases", func() {
		It("should handle single-element slice", func() {
			Expect(
				businessrules.NotEmptySlice("val", []int{42}, businessrules.SeverityError).Check(),
			).To(Succeed())
		})

		It("should handle single-element map", func() {
			Expect(
				businessrules.NotEmptyMap("val", map[string]int{"key": 42}, businessrules.SeverityError).
					Check(),
			).To(Succeed())
		})

		It("should handle large slice", func() {
			large := make([]int, 1000)
			large[999] = 1
			Expect(
				businessrules.NotEmptySlice("val", large, businessrules.SeverityError).Check(),
			).To(Succeed())
		})
	})
})
