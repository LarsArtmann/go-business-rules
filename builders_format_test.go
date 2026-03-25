package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/artmann/businessrules"
)

var _ = Describe("Format Builders", func() {
	It("should validate Email", func() {
		Expect(
			businessrules.Email("email", "test@example.com", businessrules.SeverityError).
				Check(),
		).To(Succeed())
		Expect(
			businessrules.Email("email", "user.name+tag@domain.co.uk", businessrules.SeverityError).
				Check(),
		).To(Succeed())
		Expect(
			businessrules.Email("email", "", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
		Expect(
			businessrules.Email("email", "invalid", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
		Expect(
			businessrules.Email("email", "@example.com", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should validate URL", func() {
		Expect(
			businessrules.URL("url", "http://example.com", businessrules.SeverityError).Check(),
		).To(Succeed())
		Expect(
			businessrules.URL("url", "https://example.com/path?query=1", businessrules.SeverityError).
				Check(),
		).To(Succeed())
		Expect(
			businessrules.URL("url", "", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
		Expect(
			businessrules.URL("url", "ftp://example.com", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
		Expect(
			businessrules.URL("url", "not-a-url", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})

	It("should validate UUID", func() {
		Expect(
			businessrules.UUID("id", "550e8400-e29b-41d4-a716-446655440000", businessrules.SeverityError).
				Check(),
		).To(Succeed())
		Expect(
			businessrules.UUID("id", "550E8400-E29B-41D4-A716-446655440000", businessrules.SeverityError).
				Check(),
		).To(Succeed())
		Expect(
			businessrules.UUID("id", "", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
		Expect(
			businessrules.UUID("id", "not-a-uuid", businessrules.SeverityError).Check(),
		).ToNot(Succeed())
	})
})
