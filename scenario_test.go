package businessrules_test

import (
	"fmt"

	"github.com/artmann/businessrules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Scenarios", func() {
	Describe("Registration Form Validation", func() {
		type RegistrationForm struct {
			Email           string
			Password        string
			ConfirmPassword string
			Age             int
			TermsAccepted   bool
		}

		validRegistrationForm := func() RegistrationForm {
			return RegistrationForm{
				Email:           "user@example.com",
				Password:        "securepassword123",
				ConfirmPassword: "securepassword123",
				Age:             25,
				TermsAccepted:   true,
			}
		}

		invalidEmailRegistrationForm := func() RegistrationForm {
			return RegistrationForm{
				Email:           "not-an-email",
				Password:        "securepassword123",
				ConfirmPassword: "securepassword123",
				Age:             25,
				TermsAccepted:   true,
			}
		}

		underAgeRegistrationForm := func() RegistrationForm {
			return RegistrationForm{
				Email:           "user@example.com",
				Password:        "securepassword123",
				ConfirmPassword: "securepassword123",
				Age:             10,
				TermsAccepted:   true,
			}
		}

		validateRegistration := func(form RegistrationForm) businessrules.ValidationResultError {
			v := businessrules.NewValidator()

			v.AddRule(businessrules.Email("email", form.Email, businessrules.SeverityError))
			v.AddRule(
				businessrules.MinLength("password", form.Password, 8, businessrules.SeverityError),
			)
			v.AddRule(
				businessrules.Equals(
					"confirm",
					form.Password,
					form.ConfirmPassword,
					businessrules.SeverityError,
				),
			)
			v.AddRule(
				businessrules.GreaterThan(
					"age",
					float64(form.Age),
					13,
					businessrules.SeverityWarning,
				),
			)
			v.AddRule(businessrules.Custom("terms", func() error {
				if !form.TermsAccepted {
					return fmt.Errorf("terms must be accepted")
				}
				return nil
			}, businessrules.SeverityError))

			return v.Build()
		}

		Context("when all fields are valid", func() {
			It("should pass validation", func() {
				form := validRegistrationForm()
				result := validateRegistration(form)
				Expect(result.Valid).To(BeTrue())
			})
		})

		Context("when email is invalid", func() {
			It("should fail with error severity", func() {
				form := invalidEmailRegistrationForm()
				result := validateRegistration(form)
				Expect(result.Valid).To(BeFalse())
				Expect(result.HasErrors()).To(BeTrue())
				Expect(result.Errors()).ToNot(BeEmpty())
			})
		})

		Context("when user is under 13", func() {
			It("should fail with warning severity", func() {
				form := underAgeRegistrationForm()
				result := validateRegistration(form)
				Expect(result.Valid).To(BeFalse())
				Expect(result.HasErrors()).To(BeFalse())
				Expect(result.HasWarnings()).To(BeTrue())
			})
		})
	})

	Describe("Severity-Based Decision Making", func() {
		Context("when result has only warnings", func() {
			It("should allow submission with warnings", func() {
				result := businessrules.ValidationResultError{
					Valid: false,
					ViolationErrors: []businessrules.ViolationError{
						newViolationWithContext("age", businessrules.SeverityWarning, "age is low", "user.age"),
					},
				}
				Expect(result.HasErrors()).To(BeFalse())
				Expect(result.HasWarnings()).To(BeTrue())
				Expect(result.HasCritical()).To(BeFalse())
			})
		})

		Context("when result has errors", func() {
			It("should block submission on errors", func() {
				result := businessrules.ValidationResultError{
					Valid: false,
					ViolationErrors: []businessrules.ViolationError{
						newViolationWithContext("email", businessrules.SeverityError, "email required", "user.email"),
					},
				}
				Expect(result.HasErrors()).To(BeTrue())
				Expect(result.Valid).To(BeFalse())
			})
		})

		Context("when result has critical severity", func() {
			It("should escalate to alerting system", func() {
				result := businessrules.ValidationResultError{
					Valid: false,
					ViolationErrors: []businessrules.ViolationError{
						newViolationWithContext("security", businessrules.SeverityCritical, "potential security breach", "security.audit"),
					},
				}
				Expect(result.HasCritical()).To(BeTrue())
				Expect(result.Valid).To(BeFalse())
			})
		})

		Context("when result has info messages", func() {
			It("should log but not block", func() {
				result := businessrules.ValidationResultError{
					Valid: false,
					ViolationErrors: []businessrules.ViolationError{
						newViolationWithContext("info", businessrules.SeverityInfo, "informational message", "user.activity"),
					},
				}
				Expect(result.HasInfo()).To(BeTrue())
				Expect(result.Valid).To(BeFalse())
			})
		})
	})

	Describe("Product Validation", func() {
		type Product struct {
			Name        string
			Price       float64
			Quantity    int
			Description string
		}

		validProduct := func() Product {
			return Product{
				Name:        "Widget Pro",
				Price:       29.99,
				Quantity:    100,
				Description: "A great widget",
			}
		}

		zeroQuantityProduct := func() Product {
			return Product{
				Name:        "Widget Pro",
				Price:       29.99,
				Quantity:    0,
				Description: "Out of stock",
			}
		}

		validateProduct := func(p Product) businessrules.ValidationResultError {
			v := businessrules.NewValidator()

			v.AddRule(businessrules.NotBlank("name", p.Name, businessrules.SeverityError))
			v.AddRule(businessrules.MinLength("name", p.Name, 3, businessrules.SeverityError))
			v.AddRule(businessrules.Positive("price", p.Price, businessrules.SeverityError))
			v.AddRule(
				businessrules.NonNegative(
					"quantity",
					float64(p.Quantity),
					businessrules.SeverityError,
				),
			)
			v.AddRule(
				businessrules.InRange(
					"quantity",
					float64(p.Quantity),
					1,
					10000,
					businessrules.SeverityWarning,
				),
			)

			return v.Build()
		}

		Context("valid product", func() {
			It("should pass all validations", func() {
				product := validProduct()
				result := validateProduct(product)
				Expect(result.Valid).To(BeTrue())
			})
		})

		Context("product with zero quantity", func() {
			It("should pass with warning", func() {
				product := zeroQuantityProduct()
				result := validateProduct(product)
				Expect(result.Valid).To(BeFalse())
				Expect(result.HasErrors()).To(BeFalse())
				Expect(result.HasWarnings()).To(BeTrue())
			})
		})
	})
})

func newViolationWithContext(name string, severity businessrules.Severity, msg, ctx string) businessrules.ViolationError {
	return businessrules.NewViolation(
		businessrules.NewRule(name, func() error { return nil }, severity, msg),
		ctx,
	)
}
