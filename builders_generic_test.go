package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/artmann/businessrules"
)

var _ = Describe("Generic Builders", func() {
	It("should validate OneOf", func() {
		Expect(
			businessrules.OneOf("val", "a", []string{"a", "b"}, businessrules.SeverityError).
				Check(),
		).To(Succeed())
		Expect(
			businessrules.OneOf("val", "c", []string{"a", "b"}, businessrules.SeverityError).
				Check(),
		).ToNot(Succeed())
	})

	It("should validate OneOf with integers", func() {
		Expect(
			businessrules.OneOf("val", 1, []int{1, 2, 3}, businessrules.SeverityError).Check(),
		).To(Succeed())
		Expect(
			businessrules.OneOf("val", 4, []int{1, 2, 3}, businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should validate Custom", func() {
		Expect(
			businessrules.Custom("val", func() error { return nil }, businessrules.SeverityError).
				Check(),
		).To(Succeed())
		Expect(
			businessrules.Custom("val", func() error { return assertError("failed") }, businessrules.SeverityError).
				Check(),
		).ToNot(Succeed())
	})

	It("should validate Equals with various types", func() {
		// Test strings
		Expect(
			businessrules.Equals("val", "active", "active", businessrules.SeverityError).Check(),
		).To(Succeed())
		Expect(
			businessrules.Equals("val", "inactive", "active", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
		// Test integers
		Expect(
			businessrules.Equals("val", 42, 42, businessrules.SeverityError).Check(),
		).To(Succeed())
		Expect(
			businessrules.Equals("val", 43, 42, businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})
})
