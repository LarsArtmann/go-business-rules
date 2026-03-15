package businessrules_test

import (
	"fmt"
	"regexp"

	"github.com/artmann/businessrules"
)

func ExampleNonNegative() {
	rule := businessrules.NonNegative("age", 25, businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExamplePositive() {
	rule := businessrules.Positive("count", 5, businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleInRange() {
	rule := businessrules.InRange("percentage", 75, 0, 100, businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleNotEmpty() {
	rule := businessrules.NotEmpty("name", "John Doe", businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleMinLength() {
	rule := businessrules.MinLength("password", "secret123", 8, businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleMaxLength() {
	rule := businessrules.MaxLength("username", "johndoe", 20, businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleEmail() {
	rule := businessrules.Email("email", "user@example.com", businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleURL() {
	rule := businessrules.URL("website", "https://example.com", businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleUUID() {
	rule := businessrules.UUID("id", "550e8400-e29b-41d4-a716-446655440000", businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleOneOf() {
	rule := businessrules.OneOf("status", "active", []string{"active", "inactive", "pending"}, businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleCustom() {
	rule := businessrules.Custom("custom", func() error {
		return nil // Custom validation logic
	}, businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
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
	rule := businessrules.When("admin-check", isAdmin,
		businessrules.NotEmpty("admin-key", "admin-123", businessrules.SeverityError),
	)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

func ExampleMatches() {
	pattern := regexp.MustCompile(`^[A-Z]{2}\d{4}$`)
	rule := businessrules.Matches("code", "AB1234", pattern, businessrules.SeverityError)
	if err := rule.Check(); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}
