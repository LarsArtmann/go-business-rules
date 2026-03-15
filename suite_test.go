package businessrules_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/artmann/businessrules"
)

func TestBusinessRules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BusinessRules Suite")
}

func assertError(msg string) error {
	return &testError{msg: msg}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

var _ = Describe("Core Types", func() {
	Describe("Severity", func() {
		It("should define correct constants", func() {
			Expect(businessrules.SeverityInfo).To(BeNumerically("==", 0))
			Expect(businessrules.SeverityWarning).To(BeNumerically("==", 1))
			Expect(businessrules.SeverityError).To(BeNumerically("==", 2))
			Expect(businessrules.SeverityCritical).To(BeNumerically("==", 3))
		})

		It("should return correct string representations", func() {
			Expect(businessrules.SeverityInfo.String()).To(Equal("INFO"))
			Expect(businessrules.SeverityWarning.String()).To(Equal("WARNING"))
			Expect(businessrules.SeverityError.String()).To(Equal("ERROR"))
			Expect(businessrules.SeverityCritical.String()).To(Equal("CRITICAL"))
		})

		It("should handle unknown severity", func() {
			invalidSeverity := businessrules.Severity(999)
			Expect(invalidSeverity.String()).To(ContainSubstring("UNKNOWN"))
		})
	})

	Describe("Rule", func() {
		It("should create rule with all fields", func() {
			rule := businessrules.NewRule(
				"test_rule",
				func() error { return nil },
				businessrules.SeverityError,
				"test message",
			)

			Expect(rule.Name()).To(Equal("test_rule"))
			Expect(rule.Severity()).To(Equal(businessrules.SeverityError))
			Expect(rule.Message()).To(Equal("test message"))
			Expect(rule.Check()).To(BeNil())
		})
	})

	Describe("Violation", func() {
		var rule businessrules.Rule

		BeforeEach(func() {
			rule = businessrules.NewRule(
				"test_rule",
				func() error { return nil },
				businessrules.SeverityError,
				"test message",
			)
		})

		It("should create violation with context", func() {
			violation := businessrules.NewViolation(rule, "ctx")
			Expect(violation.Rule.Name()).To(Equal("test_rule"))
			Expect(violation.Context).To(Equal("ctx"))
			Expect(violation.Timestamp).ToNot(BeZero())
		})

		It("should format error correctly", func() {
			violation := businessrules.NewViolation(rule, "ctx")
			errStr := violation.Error()
			Expect(errStr).To(ContainSubstring("[ERROR]"))
			Expect(errStr).To(ContainSubstring("test_rule"))
		})

		It("should create violation with updated context", func() {
			violation := businessrules.NewViolation(rule, "original")
			updated := violation.WithContext("updated")
			Expect(updated.Context).To(Equal("updated"))
			Expect(updated.Rule.Name()).To(Equal("test_rule"))
		})

		It("should marshal violation to JSON", func() {
			violation := businessrules.NewViolation(rule, "ctx")
			data, err := violation.MarshalJSON()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(ContainSubstring(`"rule_name":"test_rule"`))
			Expect(string(data)).To(ContainSubstring(`"context":"ctx"`))
		})
	})

	Describe("Version", func() {
		It("should have a version constant", func() {
			Expect(businessrules.Version).To(Equal("1.0.0"))
		})
	})
})
