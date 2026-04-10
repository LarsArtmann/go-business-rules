package businessrules_test

import (
	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidationResult", func() {
	createViolation := func(severity businessrules.Severity) businessrules.ViolationError {
		rule := businessrules.NewRule("test", func() error { return nil }, severity, "msg")
		return businessrules.NewViolation(rule, "")
	}

	newResultWithViolations := func(valid bool, severities ...businessrules.Severity) businessrules.ValidationResultError {
		violations := make([]businessrules.ViolationError, len(severities))
		for i, s := range severities {
			violations[i] = createViolation(s)
		}
		return businessrules.ValidationResultError{Valid: valid, ViolationErrors: violations}
	}

	Describe("Filtering", func() {
		It("should filter violations by severity", func() {
			result := businessrules.ValidationResultError{
				Valid: false,
				ViolationErrors: []businessrules.ViolationError{
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
			result := newResultWithViolations(false, businessrules.SeverityError)
			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.HasWarnings()).To(BeFalse())
		})

		It("should check for critical severity", func() {
			result := businessrules.ValidationResultError{
				Valid: false,
				ViolationErrors: []businessrules.ViolationError{
					createViolation(businessrules.SeverityCritical),
				},
			}
			Expect(result.HasCritical()).To(BeTrue())
			Expect(result.HasErrors()).To(BeTrue())
		})

		It("should check for info severity", func() {
			result := businessrules.ValidationResultError{
				Valid: false,
				ViolationErrors: []businessrules.ViolationError{
					createViolation(businessrules.SeverityInfo),
				},
			}
			Expect(result.HasInfo()).To(BeTrue())
			Expect(result.HasWarnings()).To(BeFalse())
		})
	})

	Describe("Accessors", func() {
		It("should count violations", func() {
			result := newResultWithViolations(
				false,
				businessrules.SeverityError,
				businessrules.SeverityWarning,
			)
			Expect(result.Count()).To(Equal(2))
		})

		It("should return first error", func() {
			result := newResultWithViolations(
				false,
				businessrules.SeverityWarning,
				businessrules.SeverityError,
			)
			first := result.FirstError()
			Expect(first.Rule.Name()).To(Equal("test"))
		})

		It("should return empty violation when no errors", func() {
			result := newResultWithViolations(true, businessrules.SeverityWarning)
			first := result.FirstError()
			Expect(first.Rule).To(BeNil())
		})

		It("should return first critical", func() {
			result := newResultWithViolations(
				false,
				businessrules.SeverityError,
				businessrules.SeverityCritical,
			)
			first := result.FirstCritical()
			Expect(first.Rule.Severity()).To(Equal(businessrules.SeverityCritical))
		})

		It("should return first warning", func() {
			result := newResultWithViolations(
				false,
				businessrules.SeverityInfo,
				businessrules.SeverityWarning,
			)
			first := result.FirstWarning()
			Expect(first.Rule.Severity()).To(Equal(businessrules.SeverityWarning))
		})

		It("should return first info", func() {
			result := newResultWithViolations(false, businessrules.SeverityInfo)
			first := result.FirstInfo()
			Expect(first.Rule.Severity()).To(Equal(businessrules.SeverityInfo))
		})
	})

	Describe("Iteration", func() {
		It("should iterate with ForEach", func() {
			result := newResultWithViolations(
				false,
				businessrules.SeverityError,
				businessrules.SeverityWarning,
			)
			count := 0
			result.ForEach(func(_ businessrules.ViolationError) {
				count++
			})
			Expect(count).To(Equal(2))
		})

		It("should filter violations with predicate", func() {
			result := businessrules.ValidationResultError{
				Valid: false,
				ViolationErrors: []businessrules.ViolationError{
					createViolation(businessrules.SeverityError),
					createViolation(businessrules.SeverityWarning),
					createViolation(businessrules.SeverityInfo),
				},
			}
			filtered := result.Filter(func(v businessrules.ViolationError) bool {
				return v.Rule.Severity() >= businessrules.SeverityWarning
			})
			Expect(filtered).To(HaveLen(2))
		})
	})

	Describe("Merge", func() {
		It("should merge results", func() {
			result1 := businessrules.ValidationResultError{
				Valid: true,
				ViolationErrors: []businessrules.ViolationError{
					createViolation(businessrules.SeverityError),
				},
			}
			result2 := businessrules.ValidationResultError{
				Valid: true,
				ViolationErrors: []businessrules.ViolationError{
					createViolation(businessrules.SeverityWarning),
				},
			}
			merged := result1.Merge(result2)
			Expect(merged.Valid).To(BeTrue())
			Expect(merged.Count()).To(Equal(2))
		})

		It("should merge with invalid result", func() {
			result1 := businessrules.ValidationResultError{Valid: true, ViolationErrors: nil}
			result2 := businessrules.ValidationResultError{Valid: false, ViolationErrors: nil}
			merged := result1.Merge(result2)
			Expect(merged.Valid).To(BeFalse())
		})
	})

	Describe("Error", func() {
		It("should return empty string for valid result", func() {
			result := businessrules.ValidationResultError{Valid: true, ViolationErrors: nil}
			Expect(result.Error()).To(Equal(""))
		})

		It("should return generic message when invalid with no violations", func() {
			result := businessrules.ValidationResultError{Valid: false, ViolationErrors: nil}
			Expect(result.Error()).To(Equal("validation failed"))
		})

		It("should return single violation error", func() {
			result := newResultWithViolations(false, businessrules.SeverityError)
			Expect(result.Error()).To(ContainSubstring("[ERROR]"))
		})

		It("should return formatted multi-violation error", func() {
			result := newResultWithViolations(
				false,
				businessrules.SeverityError,
				businessrules.SeverityWarning,
			)
			Expect(result.Error()).To(ContainSubstring("validation failed with 2 violations"))
		})
	})

	Describe("JSON", func() {
		It("should marshal result to JSON", func() {
			result := newResultWithViolations(false, businessrules.SeverityError)
			data, err := result.MarshalJSON()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(ContainSubstring(`"valid":false`))
		})
	})
})
