package businessrules_test

import (
	"fmt"
	"regexp"

	"github.com/artmann/businessrules"
)

func checkAndPrint(rule businessrules.Rule) {
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
}

func ExampleNonNegative() {
	checkAndPrint(businessrules.NonNegative("age", 25, businessrules.SeverityError))
	// Output: Validation passed
}

func ExamplePositive() {
	checkAndPrint(businessrules.Positive("count", 5, businessrules.SeverityError))
	// Output: Validation passed
}

func ExampleInRange() {
	checkAndPrint(businessrules.InRange("percentage", 75, 0, 100, businessrules.SeverityError))
	// Output: Validation passed
}

func ExampleNotEmpty() {
	checkAndPrint(businessrules.NotEmpty("name", "John Doe", businessrules.SeverityError))
	// Output: Validation passed
}

func ExampleMinLength() {
	checkAndPrint(businessrules.MinLength("password", "secret123", 8, businessrules.SeverityError))
	// Output: Validation passed
}

func ExampleMaxLength() {
	checkAndPrint(businessrules.MaxLength("username", "johndoe", 20, businessrules.SeverityError))
	// Output: Validation passed
}

func ExampleEmail() {
	checkAndPrint(businessrules.Email("email", "user@example.com", businessrules.SeverityError))
	// Output: Validation passed
}

func ExampleURL() {
	checkAndPrint(businessrules.URL("website", "https://example.com", businessrules.SeverityError))
	// Output: Validation passed
}

func ExampleUUID() {
	checkAndPrint(
		businessrules.UUID(
			"id",
			"550e8400-e29b-41d4-a716-446655440000",
			businessrules.SeverityError,
		),
	)
	// Output: Validation passed
}

func ExampleOneOf() {
	checkAndPrint(
		businessrules.OneOf(
			"status",
			"active",
			[]string{"active", "inactive", "pending"},
			businessrules.SeverityError,
		),
	)
	// Output: Validation passed
}

func ExampleCustom() {
	checkAndPrint(businessrules.Custom("custom", func() error {
		return nil
	}, businessrules.SeverityError))
	// Output: Validation passed
}

func ExampleValidatorBuilder() {
	result := businessrules.NewValidator().
		AddRule(businessrules.NotEmpty("name", "John", businessrules.SeverityError)).
		AddRule(businessrules.Email("email", "john@example.com", businessrules.SeverityError)).
		AddRule(businessrules.MinLength("password", "secret123", 8, businessrules.SeverityError)).
		Build()

	if result.Valid {
		fmt.Println("All validations passed")
	} else {
		fmt.Println("Validation errors:", len(result.Errors()))
	}
	// Output: All validations passed
}

func ExampleValidationResult_HasErrors() {
	result := businessrules.NewValidator().
		AddRule(businessrules.NotEmpty("name", "", businessrules.SeverityError)).
		Build()

	if result.HasErrors() {
		fmt.Println("Has errors:", len(result.Errors()))
	}
	// Output: Has errors: 1
}

func ExampleWhen() {
	isAdmin := true
	checkAndPrint(businessrules.When("admin-check", isAdmin,
		businessrules.NotEmpty("admin-key", "admin-123", businessrules.SeverityError),
	))
	// Output: Validation passed
}

func ExampleMatches() {
	pattern := regexp.MustCompile(`^[A-Z]{2}\d{4}$`)
	checkAndPrint(businessrules.Matches("code", "AB1234", pattern, businessrules.SeverityError))
	// Output: Validation passed
}
