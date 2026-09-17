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
//   - Event / Listener: Observability hooks for every rule evaluation
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
//   - Required: Validates string has visible content (Required = NotEmpty + NotBlank)
//   - MinLength: Validates minimum string length
//   - MaxLength: Validates maximum string length
//   - LengthRange: Validates minimum <= length <= maximum
//   - Matches: Validates string matches a regex pattern
//   - Contains: Validates string contains a substring
//   - MatchesFunc: Validates string with a custom predicate
//
// Time rules (injectable "now" for deterministic tests):
//   - NotPast: Validates time is not before the given reference time
//   - NotFuture: Validates time is not after the given reference time
//   - DateInRange: Validates start <= time <= end
//
// Network / identifier rules:
//   - IPAddress: Validates an IPv4 or IPv6 address
//   - CreditCard: Validates 13-19 digits passing the Luhn checksum
//   - PhoneNumber: Validates a generic phone-number shape
//   - PostalCode: Validates a generic postal-code shape
//
// Precision rules:
//   - MaxDecimalPlaces: Validates a float carries at most N decimal places
//   - DivisibleBy: Validates an integer divides without remainder
//
// Generic rules:
//   - Equals: Validates value equals expected
//   - OneOf[T]: Validates value is in allowed set
//   - Custom: Custom validation function
//
// Composite rules:
//   - All: Passes when every sub-rule passes
//   - Any: Passes when at least one sub-rule passes
//   - When: Evaluates the sub-rule only when a condition holds
//   - Not: Passes when the inner rule fails
//   - Or: Any over variadic sub-rules
//   - Xor: Passes when exactly one of two sub-rules passes
//
// Rules can carry optional metadata — WithDescription and WithTags — which is
// surfaced on RuleEvaluated events for observability without affecting checks.
//
// # Validation Events
//
// Register listeners to observe a validation run while it happens. Every
// rule check emits a RuleEvaluated event (passes included), and the run ends
// with a ValidationCompleted event carrying the returned result:
//
//	result := businessrules.NewValidator().
//	    WithListener(func(e businessrules.Event) {
//	        if re, ok := e.(businessrules.RuleEvaluated); ok && !re.Passed() {
//	            log.Printf("rule %s failed after %s", re.RuleName, re.Duration)
//	        }
//	    }).
//	    AddRule(businessrules.NonNegative("age", user.Age, businessrules.SeverityError)).
//	    Build()
//
// Listeners are called synchronously in registration order before Build
// returns; they must not panic. Without listeners, no events are constructed
// and the build path has no added cost.
//
// For slow (e.g. I/O-bound) rules, Stream executes all rules concurrently and
// delivers the same events on a channel in completion order:
//
//	events := businessrules.NewValidator().
//	    AddRules(rules...).
//	    Stream(ctx)
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
// Stream may be called once per builder; the returned channel must be drained
// or its context canceled to release the internal goroutine.
package businessrules

// Version is the current library version.
const Version = "2.2.0"
