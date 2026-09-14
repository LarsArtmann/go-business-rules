package businessrules_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

func newTestViolation(ruleName, msg, ctx string) businessrules.ViolationError {
	return businessrules.NewViolation(passingRule(ruleName, businessrules.SeverityError, msg), ctx)
}

func newViolationWithSeverity(
	ruleName string,
	severity businessrules.Severity,
	msg, ctx string,
) businessrules.ViolationError {
	return businessrules.NewViolation(passingRule(ruleName, severity, msg), ctx)
}

var _ = Describe("Context in Validation", func() {
	Describe("WithContext", func() {
		It("should update context to field path", func() {
			updated := newTestViolation(
				"email",
				"email validation",
				"original",
			).WithContext("user.profile.email")
			Expect(updated.Context).To(Equal("user.profile.email"))
		})

		It("should preserve rule name when updating context", func() {
			rule := passingRule("age_validator", businessrules.SeverityWarning, "age must be valid")
			violation := businessrules.NewViolation(rule, "input")
			updated := violation.WithContext("registration.age")
			Expect(updated.Rule.Name()).To(Equal("age_validator"))
		})

		It("should preserve timestamp when updating context", func() {
			violation := newTestViolation("test", "test", "original")
			originalTime := violation.Timestamp
			updated := violation.WithContext("new_context")
			Expect(updated.Timestamp).To(Equal(originalTime))
		})

		It("should support request ID context", func() {
			violation := newTestViolation("validator", "validation failed", "req-123")
			Expect(violation.Context).To(Equal("req-123"))
		})

		It("should support hierarchical context", func() {
			violation := newTestViolation("required", "field required", "form")
			l1 := violation.WithContext("form.order")
			l2 := l1.WithContext("form.order.shipping_address")
			l3 := l2.WithContext("form.order.shipping_address.zipcode")

			Expect(l1.Context).To(Equal("form.order"))
			Expect(l2.Context).To(Equal("form.order.shipping_address"))
			Expect(l3.Context).To(Equal("form.order.shipping_address.zipcode"))
		})
	})

	Describe("Context in error messages", func() {
		It("should include context in JSON marshaling", func() {
			violation := newViolationWithSeverity(
				"price",
				businessrules.SeverityError,
				"price must be positive",
				"checkout.total",
			)
			data, err := violation.MarshalJSON()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(ContainSubstring(`"context":"checkout.total"`))
		})

		It("should handle empty context in JSON", func() {
			violation := newTestViolation("name", "name required", "")
			data, err := violation.MarshalJSON()
			Expect(err).ToNot(HaveOccurred())

			jsonStr := string(data)
			Expect(jsonStr).ToNot(ContainSubstring(`"context":""`))
		})

		It("should handle special characters in context", func() {
			violation := newTestViolation("test", "test", `path/with/special.chars["bracket"]`)
			data, err := violation.MarshalJSON()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(ContainSubstring("special"))
		})

		It("should handle long context paths", func() {
			longPath := strings.Repeat("nested.", 20) + "field"
			violation := newViolationWithSeverity(
				"deep",
				businessrules.SeverityError,
				"deep validation",
				longPath,
			)
			updated := violation.WithContext(longPath)
			Expect(len(updated.Context)).To(BeNumerically(">", 100))
		})
	})

	Describe("Context propagation through validators", func() {
		It("should allow context enrichment at each validation layer", func() {
			baseViolation := newTestViolation("format", "invalid format", "validation")
			formContext := baseViolation.WithContext("registration_form")
			fieldContext := formContext.WithContext("registration_form.email")
			batchContext := fieldContext.WithContext("registration_form.email.batch_upload")

			Expect(batchContext.Rule.Name()).To(Equal("format"))
			Expect(batchContext.Context).To(Equal("registration_form.email.batch_upload"))
		})
	})
})
