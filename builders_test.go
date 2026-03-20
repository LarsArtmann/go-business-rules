package businessrules_test

import (
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/artmann/businessrules"
)

var _ = Describe("Builders", func() {
	Describe("Numeric Builders", func() {
		It("should validate NonNegative", func() {
			Expect(businessrules.NonNegative("val", 0, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.NonNegative("val", 10.5, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.NonNegative("val", -1, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate Positive", func() {
			Expect(businessrules.Positive("val", 1, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.Positive("val", 0.1, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.Positive("val", 0, businessrules.SeverityError).Check()).ToNot(BeNil())
			Expect(businessrules.Positive("val", -1, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate InRange", func() {
			Expect(businessrules.InRange("val", 5, 0, 10, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.InRange("val", 0, 0, 10, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.InRange("val", 10, 0, 10, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.InRange("val", -1, 0, 10, businessrules.SeverityError).Check()).ToNot(BeNil())
			Expect(businessrules.InRange("val", 11, 0, 10, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate MinInt", func() {
			Expect(businessrules.MinInt("val", 5, 5, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.MinInt("val", 10, 5, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.MinInt("val", 3, 5, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate MaxInt", func() {
			Expect(businessrules.MaxInt("val", 5, 10, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.MaxInt("val", 10, 10, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.MaxInt("val", 15, 10, businessrules.SeverityError).Check()).ToNot(BeNil())
		})
	})

	Describe("String Builders", func() {
		It("should validate NotEmpty", func() {
			Expect(businessrules.NotEmpty("val", "hello", businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.NotEmpty("val", "", businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate NotBlank", func() {
			Expect(businessrules.NotBlank("val", "hello", businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.NotBlank("val", "  \t\n  ", businessrules.SeverityError).Check()).ToNot(BeNil())
			Expect(businessrules.NotBlank("val", "", businessrules.SeverityError).Check()).ToNot(BeNil())
			Expect(businessrules.NotBlank("val", "x", businessrules.SeverityError).Check()).To(BeNil())
		})

		It("should validate MinLength", func() {
			Expect(businessrules.MinLength("val", "hello", 3, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.MinLength("val", "hi", 3, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate MaxLength", func() {
			Expect(businessrules.MaxLength("val", "hi", 5, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.MaxLength("val", "hello world", 5, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate Matches", func() {
			pattern := regexp.MustCompile(`^[a-z]+$`)
			Expect(businessrules.Matches("val", "hello", pattern, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.Matches("val", "Hello123", pattern, businessrules.SeverityError).Check()).ToNot(BeNil())
		})
	})

	Describe("Format Builders", func() {
		It("should validate Email", func() {
			Expect(businessrules.Email("email", "test@example.com", businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.Email("email", "user.name+tag@domain.co.uk", businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.Email("email", "", businessrules.SeverityError).Check()).ToNot(BeNil())
			Expect(businessrules.Email("email", "invalid", businessrules.SeverityError).Check()).ToNot(BeNil())
			Expect(businessrules.Email("email", "@example.com", businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate URL", func() {
			Expect(businessrules.URL("url", "http://example.com", businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.URL("url", "https://example.com/path?query=1", businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.URL("url", "", businessrules.SeverityError).Check()).ToNot(BeNil())
			Expect(businessrules.URL("url", "ftp://example.com", businessrules.SeverityError).Check()).ToNot(BeNil())
			Expect(businessrules.URL("url", "not-a-url", businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate UUID", func() {
			Expect(businessrules.UUID("id", "550e8400-e29b-41d4-a716-446655440000", businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.UUID("id", "550E8400-E29B-41D4-A716-446655440000", businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.UUID("id", "", businessrules.SeverityError).Check()).ToNot(BeNil())
			Expect(businessrules.UUID("id", "not-a-uuid", businessrules.SeverityError).Check()).ToNot(BeNil())
		})
	})

	Describe("Generic Builders", func() {
		It("should validate OneOf", func() {
			Expect(businessrules.OneOf("val", "a", []string{"a", "b"}, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.OneOf("val", "c", []string{"a", "b"}, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate OneOf with integers", func() {
			Expect(businessrules.OneOf("val", 1, []int{1, 2, 3}, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.OneOf("val", 4, []int{1, 2, 3}, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate Custom", func() {
			Expect(businessrules.Custom("val", func() error { return nil }, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.Custom("val", func() error { return assertError("failed") }, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate Equals", func() {
			Expect(businessrules.Equals("val", "active", "active", businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.Equals("val", "inactive", "active", businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate Equals with integers", func() {
			Expect(businessrules.Equals("val", 42, 42, businessrules.SeverityError).Check()).To(BeNil())
			Expect(businessrules.Equals("val", 43, 42, businessrules.SeverityError).Check()).ToNot(BeNil())
		})
	})

	Describe("Composite Builders", func() {
		It("should validate All - all pass", func() {
			rules := []businessrules.Rule{
				businessrules.NonNegative("a", 1, businessrules.SeverityError),
				businessrules.Positive("b", 1, businessrules.SeverityError),
			}
			Expect(businessrules.All("all", rules, businessrules.SeverityError).Check()).To(BeNil())
		})

		It("should validate All - one fails", func() {
			rules := []businessrules.Rule{
				businessrules.NonNegative("a", 1, businessrules.SeverityError),
				businessrules.Positive("b", -1, businessrules.SeverityError),
			}
			Expect(businessrules.All("all", rules, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate Any - one passes", func() {
			rules := []businessrules.Rule{
				businessrules.Positive("a", -1, businessrules.SeverityError),
				businessrules.Positive("b", 1, businessrules.SeverityError),
			}
			Expect(businessrules.Any("any", rules, businessrules.SeverityError).Check()).To(BeNil())
		})

		It("should validate Any - all fail", func() {
			rules := []businessrules.Rule{
				businessrules.Positive("a", -1, businessrules.SeverityError),
				businessrules.Positive("b", 0, businessrules.SeverityError),
			}
			Expect(businessrules.Any("any", rules, businessrules.SeverityError).Check()).ToNot(BeNil())
		})

		It("should validate When - condition true", func() {
			rule := businessrules.NotEmpty("val", "", businessrules.SeverityError)
			conditional := businessrules.When("conditional", true, rule)
			Expect(conditional.Check()).ToNot(BeNil())
		})

		It("should validate When - condition false", func() {
			rule := businessrules.NotEmpty("val", "", businessrules.SeverityError)
			conditional := businessrules.When("conditional", false, rule)
			Expect(conditional.Check()).To(BeNil())
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
})
