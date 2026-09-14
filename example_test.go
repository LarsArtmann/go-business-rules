package businessrules_test

import (
	"fmt"
	"regexp"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

func checkAndPrint(rule businessrules.Rule) {
	err := rule.Check()
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
}

func buildValidator(rules ...businessrules.Rule) businessrules.ValidationResultError {
	v := businessrules.NewValidator()
	for _, r := range rules {
		v.AddRule(r)
	}

	return v.Build()
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
	result := buildValidator(
		businessrules.NotEmpty("name", "John", businessrules.SeverityError),
		businessrules.Email("email", "john@example.com", businessrules.SeverityError),
		businessrules.MinLength("password", "secret123", 8, businessrules.SeverityError),
	)

	if result.Valid {
		fmt.Println("All validations passed")
	} else {
		fmt.Println("Validation errors:", len(result.Errors()))
	}
	// Output: All validations passed
}

func ExampleValidationResultError_HasErrors() {
	result := buildValidator(
		businessrules.NotEmpty("name", "", businessrules.SeverityError),
	)

	if result.HasErrors() {
		fmt.Println("Has errors:", len(result.Errors()))
	}
	// Output: Has errors: 1
}

func ExampleWhen() {
	isAdmin := true
	checkAndPrint(businessrules.When(
		"admin-check", isAdmin,
		businessrules.NotEmpty("admin-key", "admin-123", businessrules.SeverityError),
	))
	// Output: Validation passed
}

func ExampleMatches() {
	pattern := regexp.MustCompile(`^[A-Z]{2}\d{4}$`)
	checkAndPrint(businessrules.Matches("code", "AB1234", pattern, businessrules.SeverityError))
	// Output: Validation passed
}

func ExampleValidatorBuilder_WithListener() {
	validator := businessrules.NewValidator().
		WithListener(func(e businessrules.Event) {
			if evaluated, ok := e.(businessrules.RuleEvaluated); ok && !evaluated.Passed() {
				fmt.Printf("rule %s failed: %s\n", evaluated.RuleName, evaluated.Err)
			}
		}).
		AddRule(businessrules.NonNegative("price", -5, businessrules.SeverityError))

	result := validator.Build()

	fmt.Println("valid:", result.Valid)
	// Output: rule price failed: price must be non-negative, got -5.000000
	// valid: false
}
