package businessrules_test

import (
	businessrules "github.com/LarsArtmann/go-business-rules"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Format Builders", func() {
	Describe("Email", func() {
		DescribeTable(
			"validation",
			func(email string, shouldPass bool) {
				expectRuleResult(
					businessrules.Email("email", email, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("valid email", "test@example.com", true),
			Entry("valid email with subdomain", "user.name+tag@domain.co.uk", true),
			Entry("empty email", "", false),
			Entry("invalid format", "invalid", false),
			Entry("missing local part", "@example.com", false),
		)
	})

	Describe("URL", func() {
		DescribeTable(
			"validation",
			func(url string, shouldPass bool) {
				expectRuleResult(
					businessrules.URL("url", url, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("http URL", "http://example.com", true),
			Entry("https URL with path and query", "https://example.com/path?query=1", true),
			Entry("empty URL", "", false),
			Entry("ftp URL", "ftp://example.com", false),
			Entry("invalid URL format", "not-a-url", false),
		)
	})

	Describe("UUID", func() {
		DescribeTable(
			"validation",
			func(uuid string, shouldPass bool) {
				expectRuleResult(
					businessrules.UUID("id", uuid, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("valid lowercase UUID", "550e8400-e29b-41d4-a716-446655440000", true),
			Entry("valid uppercase UUID", "550E8400-E29B-41D4-A716-446655440000", true),
			Entry("empty UUID", "", false),
			Entry("invalid format", "not-a-uuid", false),
		)
	})
})
