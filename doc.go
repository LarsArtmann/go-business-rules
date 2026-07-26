// Package businessrules provides a severity-aware validation framework for Go.
//
// This library implements business rule validation with severity levels, allowing
// applications to distinguish between critical errors, warnings, and informational
// issues. It supports both built-in rule constructors and custom validation functions.
//
// # Overview
//
// The core concepts are:
//   - Rule: A validation check with a name, severity, and message
//   - ViolationError: A failed rule check with context and timestamp
//   - ValidationResultError: The outcome of validating multiple rules
//   - ValidatorBuilder: A fluent API for building validators
//
// # Quick Start
//
//	result := businessrules.NewValidator().
//	    AddRule(businessrules.NonNegative("age", user.Age, businessrules.SeverityError)).
//	    AddRule(businessrules.NotEmpty("name", user.Name, businessrules.SeverityError)).
//	    AddRule(businessrules.MaxLength("bio", user.Bio, 500, businessrules.SeverityWarning)).
//	    Build()
//
//	if result.HasErrors() {
//	    for _, v := range result.Errors() {
//	        log.Error(v.Error())
//	    }
//	}
//
// # Severity Levels
//
// The library defines four severity levels:
//   - SeverityInfo: Informational issues that don't block processing
//   - SeverityWarning: Warnings that should be reviewed but don't block
//   - SeverityError: Errors that should block processing
//   - SeverityCritical: Critical errors requiring immediate attention
//
// # Built-in Rules
//
// Numeric rules:
//   - NonNegative: Validates value >= 0
//   - Positive: Validates value > 0
//   - InRange: Validates min <= value <= max
//   - MinInt: Validates value >= min (for integers)
//   - MaxInt: Validates value <= max (for integers)
//
// String rules:
//   - NotEmpty: Validates string is not empty
//   - NotBlank: Validates string is not empty or whitespace-only
//   - MinLength: Validates minimum string length
//   - MaxLength: Validates maximum string length
//   - Matches: Validates string matches a regex pattern
//
// Generic rules:
//   - Equals: Validates value equals expected
//   - OneOf[T]: Validates value is in allowed set
//   - Custom: Custom validation function
//
// # Custom Rules
//
// Create custom rules using the Custom constructor or by implementing the Rule interface:
//
//	rule := businessrules.Custom("custom-check", func() error {
//	    if someCondition {
//	        return errors.New("validation failed")
//	    }
//	    return nil
//	}, businessrules.SeverityError)
//
// # Thread Safety
//
// Rule instances are immutable and safe for concurrent use. ValidatorBuilder
// should not be shared across goroutines; create a new builder per validation session.
package businessrules

// Version is the current library version.
const Version = "1.1.0"
