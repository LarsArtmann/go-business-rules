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

	Describe("Result", func() {
		createViolation := func(severity businessrules.Severity) businessrules.Violation {
			rule := businessrules.NewRule("test", func() error { return nil }, severity, "msg")
			return businessrules.NewViolation(rule, "")
		}

		It("should filter violations by severity", func() {
			result := businessrules.Result{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityInfo),
					createViolation(businessrules.SeverityWarning),
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityCritical),
				},
			}

			Expect(result.Errors()).To(HaveLen(2))
			Expect(result.Warnings()).To(HaveLen(1))
			Expect(result.Info()).To(HaveLen(1))
			Expect(result.Critical()).To(HaveLen(1))
		})

		It("should check for presence of severities", func() {
			result := businessrules.Result{
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityError)},
			}
			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.HasWarnings()).To(BeFalse())
		})

		It("should count violations", func() {
			result := businessrules.Result{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityWarning),
				},
			}
			Expect(result.Count()).To(Equal(2))
		})

		It("should return first error", func() {
			result := businessrules.Result{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityWarning),
					createViolation(businessrules.SeverityError),
				},
			}
			first := result.FirstError()
			Expect(first.Rule.Name()).To(Equal("test"))
		})

		It("should return empty violation when no errors", func() {
			result := businessrules.Result{
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityWarning)},
			}
			first := result.FirstError()
			Expect(first.Rule).To(BeNil())
		})

		It("should return first critical", func() {
			result := businessrules.Result{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityCritical),
				},
			}
			first := result.FirstCritical()
			Expect(first.Rule.Severity()).To(Equal(businessrules.SeverityCritical))
		})

		It("should iterate with ForEach", func() {
			result := businessrules.Result{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityWarning),
				},
			}
			count := 0
			result.ForEach(func(v businessrules.Violation) {
				count++
			})
			Expect(count).To(Equal(2))
		})

		It("should merge results", func() {
			result1 := businessrules.Result{
				Valid:      true,
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityError)},
			}
			result2 := businessrules.Result{
				Valid:      true,
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityWarning)},
			}
			merged := result1.Merge(result2)
			Expect(merged.Valid).To(BeTrue())
			Expect(merged.Count()).To(Equal(2))
		})

		It("should merge with invalid result", func() {
			result1 := businessrules.Result{Valid: true}
			result2 := businessrules.Result{Valid: false}
			merged := result1.Merge(result2)
			Expect(merged.Valid).To(BeFalse())
		})

		It("should marshal result to JSON", func() {
			result := businessrules.Result{
				Valid:      false,
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityError)},
			}
			data, err := result.MarshalJSON()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(ContainSubstring(`"valid":false`))
		})
	})

	Describe("ValidatorBuilder", func() {
		It("should build valid result when all rules pass", func() {
			rule := businessrules.NewRule("pass", func() error { return nil }, businessrules.SeverityInfo, "pass")
			result := businessrules.NewValidator().AddRule(rule).Build()
			Expect(result.Valid).To(BeTrue())
		})

		It("should build invalid result when rule fails", func() {
			rule := businessrules.NewRule("fail", func() error { return assertError("failed") }, businessrules.SeverityError, "fail")
			result := businessrules.NewValidator().AddRule(rule).Build()
			Expect(result.Valid).To(BeFalse())
			Expect(result.HasErrors()).To(BeTrue())
		})

		It("should support AddRules for multiple rules", func() {
			rule1 := businessrules.NewRule("r1", func() error { return nil }, businessrules.SeverityInfo, "r1")
			rule2 := businessrules.NewRule("r2", func() error { return nil }, businessrules.SeverityInfo, "r2")
			result := businessrules.NewValidator().AddRules(rule1, rule2).Build()
			Expect(result.Valid).To(BeTrue())
		})
	})

	Describe("Version", func() {
		It("should have a version constant", func() {
			Expect(businessrules.Version).To(Equal("1.0.0"))
		})
	})
})
