package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/artmann/businessrules"
)

var createViolation = func(severity businessrules.Severity) businessrules.Violation {
	rule := businessrules.NewRule("test", func() error { return nil }, severity, "msg")
	return businessrules.NewViolation(rule, "")
}

var _ = Describe("ValidationResult", func() {
	Describe("Filtering", func() {
		It("should filter violations by severity", func() {
			result := businessrules.ValidationResult{
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
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityError)},
			}
			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.HasWarnings()).To(BeFalse())
		})

		It("should check for critical severity", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityCritical),
				},
			}
			Expect(result.HasCritical()).To(BeTrue())
			Expect(result.HasErrors()).To(BeTrue())
		})

		It("should check for info severity", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityInfo)},
			}
			Expect(result.HasInfo()).To(BeTrue())
			Expect(result.HasWarnings()).To(BeFalse())
		})
	})

	Describe("Accessors", func() {
		It("should count violations", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityWarning),
				},
			}
			Expect(result.Count()).To(Equal(2))
		})

		It("should return first error", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityWarning),
					createViolation(businessrules.SeverityError),
				},
			}
			first := result.FirstError()
			Expect(first.Rule.Name()).To(Equal("test"))
		})

		It("should return empty violation when no errors", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityWarning),
				},
			}
			first := result.FirstError()
			Expect(first.Rule).To(BeNil())
		})

		It("should return first critical", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityCritical),
				},
			}
			first := result.FirstCritical()
			Expect(first.Rule.Severity()).To(Equal(businessrules.SeverityCritical))
		})

		It("should return first warning", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityInfo),
					createViolation(businessrules.SeverityWarning),
				},
			}
			first := result.FirstWarning()
			Expect(first.Rule.Severity()).To(Equal(businessrules.SeverityWarning))
		})

		It("should return first info", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityInfo),
				},
			}
			first := result.FirstInfo()
			Expect(first.Rule.Severity()).To(Equal(businessrules.SeverityInfo))
		})
	})

	Describe("Iteration", func() {
		It("should iterate with ForEach", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityWarning),
				},
			}
			count := 0
			result.ForEach(func(_ businessrules.Violation) {
				count++
			})
			Expect(count).To(Equal(2))
		})

		It("should filter violations with predicate", func() {
			result := businessrules.ValidationResult{
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityWarning),
					createViolation(businessrules.SeverityInfo),
				},
			}
			filtered := result.Filter(func(v businessrules.Violation) bool {
				return v.Rule.Severity() >= businessrules.SeverityWarning
			})
			Expect(filtered).To(HaveLen(2))
		})
	})

	Describe("Merge", func() {
		It("should merge results", func() {
			result1 := businessrules.ValidationResult{
				Valid:      true,
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityError)},
			}
			result2 := businessrules.ValidationResult{
				Valid: true,
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityWarning),
				},
			}
			merged := result1.Merge(result2)
			Expect(merged.Valid).To(BeTrue())
			Expect(merged.Count()).To(Equal(2))
		})

		It("should merge with invalid result", func() {
			result1 := businessrules.ValidationResult{Valid: true}
			result2 := businessrules.ValidationResult{Valid: false}
			merged := result1.Merge(result2)
			Expect(merged.Valid).To(BeFalse())
		})
	})

	Describe("Error", func() {
		It("should return empty string for valid result", func() {
			result := businessrules.ValidationResult{Valid: true}
			Expect(result.Error()).To(Equal(""))
		})

		It("should return generic message when invalid with no violations", func() {
			result := businessrules.ValidationResult{Valid: false}
			Expect(result.Error()).To(Equal("validation failed"))
		})

		It("should return single violation error", func() {
			result := businessrules.ValidationResult{
				Valid:      false,
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityError)},
			}
			Expect(result.Error()).To(ContainSubstring("[ERROR]"))
		})

		It("should return formatted multi-violation error", func() {
			result := businessrules.ValidationResult{
				Valid: false,
				Violations: []businessrules.Violation{
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityWarning),
				},
			}
			Expect(result.Error()).To(ContainSubstring("validation failed with 2 violations"))
		})
	})

	Describe("JSON", func() {
		It("should marshal result to JSON", func() {
			result := businessrules.ValidationResult{
				Valid:      false,
				Violations: []businessrules.Violation{createViolation(businessrules.SeverityError)},
			}
			data, err := result.MarshalJSON()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(ContainSubstring(`"valid":false`))
		})
	})
})
