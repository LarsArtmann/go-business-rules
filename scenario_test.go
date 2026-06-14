package businessrules_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	businessrules "github.com/LarsArtmann/go-business-rules"
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

		makeRegistrationForm := func(email string, age int) RegistrationForm {
			return RegistrationForm{
				Email:           email,
				Password:        "securepassword123",
				ConfirmPassword: "securepassword123",
				Age:             age,
				TermsAccepted:   true,
			}
		}

		validRegistrationForm := func() RegistrationForm {
			return makeRegistrationForm("user@example.com", 25)
		}

		invalidEmailRegistrationForm := func() RegistrationForm {
			return makeRegistrationForm("not-an-email", 25)
		}

		underAgeRegistrationForm := func() RegistrationForm {
			return makeRegistrationForm("user@example.com", 10)
		}

		validateRegistration := func(form RegistrationForm) businessrules.ValidationResultError {
			validator := businessrules.NewValidator()

			validator.AddRule(businessrules.Email("email", form.Email, businessrules.SeverityError))
			validator.AddRule(
				businessrules.MinLength("password", form.Password, 8, businessrules.SeverityError),
			)
			validator.AddRule(
				businessrules.Equals(
					"confirm",
					form.Password,
					form.ConfirmPassword,
					businessrules.SeverityError,
				),
			)
			validator.AddRule(
				businessrules.GreaterThan(
					"age",
					float64(form.Age),
					13,
					businessrules.SeverityWarning,
				),
			)
			validator.AddRule(businessrules.Custom("terms", func() error {
				if !form.TermsAccepted {
					return errors.New("terms must be accepted")
				}

				return nil
			}, businessrules.SeverityError))

			return validator.Build()
		}

		Context("when all fields are valid", func() {
			It("should pass validation", func() {
				form := validRegistrationForm()
				result := validateRegistration(form)
				expectValid(result)
			})
		})

		Context("when email is invalid", func() {
			It("should fail with error severity", func() {
				form := invalidEmailRegistrationForm()
				result := validateRegistration(form)
				expectInvalid(result)
				expectHasErrors(result)
				Expect(result.Errors()).ToNot(BeEmpty())
			})
		})

		Context("when user is under 13", func() {
			It("should fail with warning severity", func() {
				form := underAgeRegistrationForm()
				result := validateRegistration(form)
				expectInvalid(result)
				expectNoErrors(result)
				expectHasWarnings(result)
			})
		})
	})

	Describe("Severity-Based Decision Making", func() {
		Context("when result has only warnings", func() {
			It("should allow submission with warnings", func() {
				result := businessrules.ValidationResultError{
					Valid: false,
					ViolationErrors: []businessrules.ViolationError{
						newViolationWithContext(
							"age",
							businessrules.SeverityWarning,
							"age is low",
							"user.age",
						),
					},
				}
				expectNoErrors(result)
				expectHasWarnings(result)
				expectNoCritical(result)
			})
		})

		Context("when result has errors", func() {
			It("should block submission on errors", func() {
				result := businessrules.ValidationResultError{
					Valid: false,
					ViolationErrors: []businessrules.ViolationError{
						newViolationWithContext(
							"email",
							businessrules.SeverityError,
							"email required",
							"user.email",
						),
					},
				}
				expectHasErrors(result)
				expectInvalid(result)
			})
		})

		Context("when result has critical severity", func() {
			It("should escalate to alerting system", func() {
				result := businessrules.ValidationResultError{
					Valid: false,
					ViolationErrors: []businessrules.ViolationError{
						newViolationWithContext(
							"security",
							businessrules.SeverityCritical,
							"potential security breach",
							"security.audit",
						),
					},
				}
				expectHasCritical(result)
				expectInvalid(result)
			})
		})

		Context("when result has info messages", func() {
			It("should log but not block", func() {
				result := businessrules.ValidationResultError{
					Valid: false,
					ViolationErrors: []businessrules.ViolationError{
						newViolationWithContext(
							"info",
							businessrules.SeverityInfo,
							"informational message",
							"user.activity",
						),
					},
				}
				Expect(result.HasInfo()).To(BeTrue())
				expectInvalid(result)
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

		makeProduct := func(quantity int, description string) Product {
			return Product{
				Name:        "Widget Pro",
				Price:       29.99,
				Quantity:    quantity,
				Description: description,
			}
		}

		validProduct := func() Product {
			return makeProduct(100, "A great widget")
		}

		zeroQuantityProduct := func() Product {
			return makeProduct(0, "Out of stock")
		}

		validateProduct := func(product Product) businessrules.ValidationResultError {
			validator := businessrules.NewValidator()

			validator.AddRule(
				businessrules.NotBlank("name", product.Name, businessrules.SeverityError),
			)
			validator.AddRule(
				businessrules.MinLength("name", product.Name, 3, businessrules.SeverityError),
			)
			validator.AddRule(
				businessrules.Positive("price", product.Price, businessrules.SeverityError),
			)
			validator.AddRule(
				businessrules.NonNegative(
					"quantity",
					float64(product.Quantity),
					businessrules.SeverityError,
				),
			)
			validator.AddRule(
				businessrules.InRange(
					"quantity",
					float64(product.Quantity),
					1,
					10000,
					businessrules.SeverityWarning,
				),
			)

			return validator.Build()
		}

		Context("valid product", func() {
			It("should pass all validations", func() {
				product := validProduct()
				result := validateProduct(product)
				expectValid(result)
			})
		})

		Context("product with zero quantity", func() {
			It("should pass with warning", func() {
				product := zeroQuantityProduct()
				result := validateProduct(product)
				expectInvalid(result)
				expectNoErrors(result)
				expectHasWarnings(result)
			})
		})
	})
})

func newViolationWithContext(
	name string,
	severity businessrules.Severity,
	msg, ctx string,
) businessrules.ViolationError {
	return businessrules.NewViolation(passingRule(name, severity, msg), ctx)
}
