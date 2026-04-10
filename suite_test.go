package businessrules_test

import (
	"testing"

	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBusinessRules(t *testing.T) {
	t.Parallel()
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

func expectRuleResult(result error, shouldPass bool) {
	if shouldPass {
		Expect(result).To(Succeed())
	} else {
		Expect(result).ToNot(Succeed())
	}
}

func passingRule(name string, severity businessrules.Severity, msg string) businessrules.Rule {
	return businessrules.NewRule(name, func() error { return nil }, severity, msg)
}

func failingRule(
	name string,
	fn func() error,
	severity businessrules.Severity,
	msg string,
) businessrules.Rule {
	return businessrules.NewRule(name, fn, severity, msg)
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
			rule := passingRule("test_rule", businessrules.SeverityError, "test message")
			Expect(rule.Name()).To(Equal("test_rule"))
			Expect(rule.Severity()).To(Equal(businessrules.SeverityError))
			Expect(rule.Message()).To(Equal("test message"))
			Expect(rule.Check()).To(Succeed())
		})
	})

	Describe("ViolationError", func() {
		var rule businessrules.Rule

		BeforeEach(func() {
			rule = passingRule("test_rule", businessrules.SeverityError, "test message")
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
			Expect(errStr).To(ContainSubstring("context: ctx"))
		})

		It("should format error without context", func() {
			violation := businessrules.NewViolation(rule, "")
			errStr := violation.Error()
			Expect(errStr).To(ContainSubstring("[ERROR]"))
			Expect(errStr).To(ContainSubstring("test_rule"))
			Expect(errStr).ToNot(ContainSubstring("context:"))
		})

		It("should create violation with updated context", func() {
			violation := businessrules.NewViolation(rule, "original")
			updated := violation.WithContext("updated")
			Expect(updated.Context).To(Equal("updated"))
			name := updated.Rule.Name()
			Expect(name).To(Equal("test_rule"))
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
			Expect(businessrules.Version).To(Equal("1.1.0"))
		})
	})
})
