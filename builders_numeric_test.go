package businessrules_test

import (
	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Builders", func() {
	Describe("Numeric Builders", func() {
		It("should validate NonNegative", func() {
			Expect(
				businessrules.NonNegative("val", 0, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.NonNegative("val", 10.5, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.NonNegative("val", -1, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
		})

		It("should validate Positive", func() {
			Expect(
				businessrules.Positive("val", 1, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.Positive("val", 0.1, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.Positive("val", 0, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
			Expect(
				businessrules.Positive("val", -1, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
		})

		It("should validate InRange", func() {
			Expect(
				businessrules.InRange("val", 5, 0, 10, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.InRange("val", 0, 0, 10, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.InRange("val", 10, 0, 10, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.InRange("val", -1, 0, 10, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
			Expect(
				businessrules.InRange("val", 11, 0, 10, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
		})

		It("should validate MinInt", func() {
			Expect(
				businessrules.MinInt("val", 5, 3, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.MinInt("val", 3, 3, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.MinInt("val", 2, 3, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
		})

		It("should validate MaxInt", func() {
			Expect(
				businessrules.MaxInt("val", 5, 10, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.MaxInt("val", 10, 10, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.MaxInt("val", 15, 10, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
		})
	})

	Describe("Additional Numeric Builders", func() {
		It("should validate GreaterThan", func() {
			Expect(
				businessrules.GreaterThan("val", 10, 5, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.GreaterThan("val", 5, 5, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
			Expect(
				businessrules.GreaterThan("val", 3, 5, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
		})

		It("should validate LessThan", func() {
			Expect(
				businessrules.LessThan("val", 3, 5, businessrules.SeverityError).Check(),
			).To(Succeed())
			Expect(
				businessrules.LessThan("val", 5, 5, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
			Expect(
				businessrules.LessThan("val", 10, 5, businessrules.SeverityError).Check(),
			).ToNot(Succeed())
		})
	})
})
