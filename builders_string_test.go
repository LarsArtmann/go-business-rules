package businessrules_test

import (
	"regexp"

	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("String Builders", func() {
	It("should validate NotEmpty", func() {
		Expect(
			businessrules.NotEmpty("val", "hello", businessrules.SeverityError).Check(),
		).To(Succeed())
		Expect(
			businessrules.NotEmpty("val", "", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should validate NotBlank", func() {
		Expect(
			businessrules.NotBlank("val", "hello", businessrules.SeverityError).Check(),
		).To(Succeed())
		Expect(
			businessrules.NotBlank("val", "  \t\n  ", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
		Expect(
			businessrules.NotBlank("val", "", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
		Expect(
			businessrules.NotBlank("val", "x", businessrules.SeverityError).Check(),
		).To(Succeed())
	})

	It("should validate MinLength", func() {
		Expect(
			businessrules.MinLength("val", "hello", 3, businessrules.SeverityError).Check(),
		).To(Succeed())
		Expect(
			businessrules.MinLength("val", "hi", 3, businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should validate MaxLength", func() {
		Expect(
			businessrules.MaxLength("val", "hi", 5, businessrules.SeverityError).Check(),
		).To(Succeed())
		Expect(
			businessrules.MaxLength("val", "hello world", 5, businessrules.SeverityError).
				Check(),
		).ToNot(Succeed())
	})

	It("should validate Matches", func() {
		pattern := regexp.MustCompile(`^[a-z]+$`)
		Expect(
			businessrules.Matches("val", "hello", pattern, businessrules.SeverityError).Check(),
		).To(Succeed())
		Expect(
			businessrules.Matches("val", "Hello123", pattern, businessrules.SeverityError).
				Check(),
		).ToNot(Succeed())
	})
})
