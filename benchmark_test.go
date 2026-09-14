package businessrules

import (
	"context"
	"strings"
	"testing"
)

func BenchmarkValidatorBuilder(b *testing.B) {
	for b.Loop() {
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

	for b.Loop() {
		_ = validator.Valid
	}
}

func BenchmarkValidationFail(b *testing.B) {
	validator := NewValidator().
		AddRule(NonNegative("price", -5.0, SeverityError)).
		AddRule(NotEmpty("name", "", SeverityError)).
		AddRule(MinLength("code", "AB", 3, SeverityWarning)).
		Build()

	for b.Loop() {
		_ = validator.Valid
	}
}

func BenchmarkResultFiltering(b *testing.B) {
	violations := make([]ViolationError, 100)
	for index := range violations {
		var severity Severity

		switch index % 4 {
		case 0:
			severity = SeverityError
		case 1:
			severity = SeverityWarning
		case 2:
			severity = SeverityCritical
		default:
			severity = SeverityInfo
		}

		rule := NewRule("test", func() error { return nil }, severity, "msg")
		violations[index] = NewViolation(rule, "context")
	}

	result := ValidationResultError{Valid: false, ViolationErrors: violations}

	for b.Loop() {
		_ = result.Errors()
	}
}

func BenchmarkResultMerge(b *testing.B) {
	rule1 := NewRule("test1", func() error { return nil }, SeverityError, "msg1")
	rule2 := NewRule("test2", func() error { return nil }, SeverityWarning, "msg2")
	result1 := ValidationResultError{
		Valid:           true,
		ViolationErrors: []ViolationError{NewViolation(rule1, "")},
	}
	result2 := ValidationResultError{
		Valid:           true,
		ViolationErrors: []ViolationError{NewViolation(rule2, "")},
	}

	for b.Loop() {
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

	for b.Loop() {
		for _, c := range cases {
			rule := NotBlank("field", c, SeverityError)
			_ = rule.Check()
		}
	}
}

func BenchmarkEquals(b *testing.B) {
	cases := []struct {
		name  string
		value string
	}{
		{"pass", "active"},
		{"fail", "inactive"},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for b.Loop() {
				rule := Equals("status", c.value, "active", SeverityError)
				_ = rule.Check()
			}
		})
	}
}

func BenchmarkValidatorNoListener(b *testing.B) {
	for b.Loop() {
		_ = NewValidator().
			AddRule(NonNegative("price", 10.0, SeverityError)).
			AddRule(NotEmpty("name", "test", SeverityError)).
			Build()
	}
}

func BenchmarkValidatorOneListener(b *testing.B) {
	noop := func(Event) {}

	for b.Loop() {
		_ = NewValidator().
			WithListener(noop).
			AddRule(NonNegative("price", 10.0, SeverityError)).
			AddRule(NotEmpty("name", "test", SeverityError)).
			Build()
	}
}

func BenchmarkValidatorThreeListeners(b *testing.B) {
	noop := func(Event) {}

	for b.Loop() {
		_ = NewValidator().
			WithListener(noop, noop, noop).
			AddRule(NonNegative("price", 10.0, SeverityError)).
			AddRule(NotEmpty("name", "test", SeverityError)).
			Build()
	}
}

func benchmarkStream(b *testing.B, ruleCount, concurrencyLimit int) {
	builder := NewValidator()

	for range ruleCount {
		builder.AddRule(NonNegative("price", 10.0, SeverityError))
	}

	if concurrencyLimit > 0 {
		builder.WithConcurrency(concurrencyLimit)
	}

	b.ReportAllocs()

	for b.Loop() {
		for range builder.Stream(context.Background()) {
		}
	}
}

func BenchmarkStream2Rules(b *testing.B) {
	benchmarkStream(b, 2, 0)
}

func BenchmarkStream10Rules(b *testing.B) {
	benchmarkStream(b, 10, 0)
}

func BenchmarkStream10RulesConcurrency4(b *testing.B) {
	benchmarkStream(b, 10, 4)
}
