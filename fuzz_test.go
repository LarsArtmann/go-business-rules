package businessrules

import (
	"regexp"
	"testing"
)

func addSeedCases(f *testing.F, cases []string) {
	f.Helper()

	for _, c := range cases {
		f.Add(c)
	}
}

func FuzzEmail(f *testing.F) {
	testCases := []string{
		"test@example.com",
		"user.name+tag@domain.co.uk",
		"invalid",
		"",
		"@example.com",
	}
	addSeedCases(f, testCases)

	f.Fuzz(func(t *testing.T, email string) {
		rule := Email("email", email, SeverityError)
		_ = rule.Check()
	})
}

func FuzzURL(f *testing.F) {
	testCases := []string{
		"http://example.com",
		"https://example.com/path?query=1",
		"ftp://example.com",
		"",
		"not-a-url",
	}
	addSeedCases(f, testCases)

	f.Fuzz(func(t *testing.T, url string) {
		rule := URL("url", url, SeverityError)
		_ = rule.Check()
	})
}

func FuzzUUID(f *testing.F) {
	testCases := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"550E8400-E29B-41D4-A716-446655440000",
		"not-a-uuid",
		"",
	}
	addSeedCases(f, testCases)

	f.Fuzz(func(t *testing.T, uuid string) {
		rule := UUID("id", uuid, SeverityError)
		_ = rule.Check()
	})
}

func FuzzNotBlank(f *testing.F) {
	testCases := []string{
		"valid content",
		"",
		"   ",
		"\t\t",
		"\n\n",
		"  x  ",
	}
	addSeedCases(f, testCases)

	f.Fuzz(func(t *testing.T, val string) {
		rule := NotBlank("field", val, SeverityError)
		_ = rule.Check()
	})
}

func FuzzMatches(f *testing.F) {
	pattern := regexp.MustCompile(`^[a-z]+$`)
	f.Add("hello", pattern.String())
	f.Add("Hello123", pattern.String())
	f.Add("", pattern.String())
	f.Add("test", pattern.String())
	f.Fuzz(func(t *testing.T, val, patternStr string) {
		pat, err := regexp.Compile(patternStr)
		if err != nil {
			return
		}

		rule := Matches("field", val, pat, SeverityError)
		_ = rule.Check()
	})
}

func FuzzEquals(f *testing.F) {
	f.Add("active", "active")
	f.Add("inactive", "active")
	f.Add("42", "42")
	f.Add("43", "42")
	f.Fuzz(func(t *testing.T, val, expected string) {
		rule := Equals("field", val, expected, SeverityError)
		_ = rule.Check()
	})
}

func FuzzOneOf(f *testing.F) {
	f.Add("active", "active", "inactive", "pending")
	f.Add("unknown", "active", "inactive", "pending")
	f.Fuzz(func(t *testing.T, val, a0, a1, a2 string) {
		allowed := []string{a0, a1, a2}
		rule := OneOf("field", val, allowed, SeverityError)
		_ = rule.Check()
	})
}
