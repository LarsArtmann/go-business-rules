package businessrules

import (
	"strings"
	"testing"
)

func BenchmarkValidatorBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewValidator().
			AddRule(NonNegative("price", 10.0, SeverityError)).
			AddRule(NotEmpty("name", "test", SeverityError)).
			AddRule(MinLength("code", "ABC", 3, SeverityWarning)).
			Build()
	}
}

func BenchmarkValidationPass(b *testing.B) {
	validator := NewValidator().
		AddRule(NonNegative("price", 10.0, SeverityError)).
		AddRule(NotEmpty("name", "test", SeverityError)).
		AddRule(MinLength("code", "ABC", 3, SeverityWarning)).
		Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.Valid
	}
}

func BenchmarkValidationFail(b *testing.B) {
	validator := NewValidator().
		AddRule(NonNegative("price", -5.0, SeverityError)).
		AddRule(NotEmpty("name", "", SeverityError)).
		AddRule(MinLength("code", "AB", 3, SeverityWarning)).
		Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.Valid
	}
}

func BenchmarkResultFiltering(b *testing.B) {
	violations := make([]Violation, 100)
	for i := range violations {
		severity := SeverityInfo
		if i%4 == 0 {
			severity = SeverityError
		} else if i%4 == 1 {
			severity = SeverityWarning
		} else if i%4 == 2 {
			severity = SeverityCritical
		}
		rule := NewRule("test", func() error { return nil }, severity, "msg")
		violations[i] = NewViolation(rule, "context")
	}
	result := ValidationResult{Violations: violations}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = result.Errors()
	}
}

func BenchmarkResultMerge(b *testing.B) {
	rule1 := NewRule("test1", func() error { return nil }, SeverityError, "msg1")
	rule2 := NewRule("test2", func() error { return nil }, SeverityWarning, "msg2")
	result1 := ValidationResult{Valid: true, Violations: []Violation{NewViolation(rule1, "")}}
	result2 := ValidationResult{Valid: true, Violations: []Violation{NewViolation(rule2, "")}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = result1.Merge(result2)
	}
}

func BenchmarkNotBlank(b *testing.B) {
	cases := []string{
		"valid content",
		"",
		"   ",
		"\t\t",
		"\n\n",
		"  x  ",
		strings.Repeat(" ", 100),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range cases {
			rule := NotBlank("field", c, SeverityError)
			_ = rule.Check()
		}
	}
}

func BenchmarkEquals(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule := Equals("status", "active", "active", SeverityError)
		_ = rule.Check()
	}
}

func BenchmarkEqualsFail(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule := Equals("status", "inactive", "active", SeverityError)
		_ = rule.Check()
	}
}
