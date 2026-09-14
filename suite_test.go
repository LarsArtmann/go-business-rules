package businessrules_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
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

func checkValidity(r businessrules.ValidationResultError) {
	if r.Valid {
		Expect(r.Valid).To(BeTrue())
	} else {
		Expect(r.Valid).To(BeFalse())
	}
}

func checkErrorsExist(r businessrules.ValidationResultError) {
	found := r.HasErrors()
	if found {
		Expect(found).To(BeTrue())
	} else {
		Expect(found).To(BeFalse())
	}
}

func checkWarningsExist(r businessrules.ValidationResultError) {
	present := r.HasWarnings()
	if present {
		Expect(present).To(BeTrue())
	} else {
		Expect(present).To(BeFalse())
	}
}

func checkCriticalExists(r businessrules.ValidationResultError) {
	exists := r.HasCritical()
	if exists {
		Expect(exists).To(BeTrue())
	} else {
		Expect(exists).To(BeFalse())
	}
}

func expectValid(r businessrules.ValidationResultError) {
	checkValidity(r)
}

func expectInvalid(r businessrules.ValidationResultError) {
	checkValidity(r)
}

func expectHasErrors(r businessrules.ValidationResultError) {
	checkErrorsExist(r)
}

func expectNoErrors(r businessrules.ValidationResultError) {
	checkErrorsExist(r)
}

func expectHasWarnings(r businessrules.ValidationResultError) {
	checkWarningsExist(r)
}

func expectNoWarnings(r businessrules.ValidationResultError) {
	checkWarningsExist(r)
}

func expectHasCritical(r businessrules.ValidationResultError) {
	checkCriticalExists(r)
}

func expectNoCritical(r businessrules.ValidationResultError) {
	checkCriticalExists(r)
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
			Expect(businessrules.SeverityInfo).To(Equal(businessrules.Severity("info")))
			Expect(businessrules.SeverityWarning).To(Equal(businessrules.Severity("warning")))
			Expect(businessrules.SeverityError).To(Equal(businessrules.Severity("error")))
			Expect(businessrules.SeverityCritical).To(Equal(businessrules.Severity("critical")))
		})

		It("should return correct string representations", func() {
			Expect(businessrules.SeverityInfo.String()).To(Equal("info"))
			Expect(businessrules.SeverityWarning.String()).To(Equal("warning"))
			Expect(businessrules.SeverityError.String()).To(Equal("error"))
			Expect(businessrules.SeverityCritical.String()).To(Equal("critical"))
		})

		It("should handle unknown severity", func() {
			invalidSeverity := businessrules.Severity("invalid")
			Expect(invalidSeverity.String()).To(Equal("invalid"))
			Expect(invalidSeverity.IsValid()).To(BeFalse())
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
			Expect(errStr).To(ContainSubstring("[error]"))
			Expect(errStr).To(ContainSubstring("test_rule"))
			Expect(errStr).To(ContainSubstring("context: ctx"))
		})

		It("should format error without context", func() {
			violation := businessrules.NewViolation(rule, "")
			errStr := violation.Error()
			Expect(errStr).To(ContainSubstring("[error]"))
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
			Expect(string(data)).To(ContainSubstring(`"ruleName":"test_rule"`))
			Expect(string(data)).To(ContainSubstring(`"context":"ctx"`))
		})
	})

	Describe("Version", func() {
		It("should have a version constant", func() {
			Expect(businessrules.Version).To(Equal("2.1.0"))
		})
	})
})
