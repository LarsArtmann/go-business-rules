package businessrules_test

import (
	. "github.com/onsi/ginkgo/v2"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

var _ = Describe("Network Builders", func() {
	Describe("IPAddress", func() {
		DescribeTable(
			"validation",
			func(value string, shouldPass bool) {
				expectRuleResult(
					businessrules.IPAddress("ip", value, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("IPv4", "192.168.1.1", true),
			Entry("IPv6", "2001:db8::1", true),
			Entry("IPv4 loopback", "127.0.0.1", true),
			Entry("empty", "", false),
			Entry("not an address", "not-an-ip", false),
			Entry("out of range octet", "999.1.1.1", false),
		)
	})

	Describe("CreditCard", func() {
		DescribeTable(
			"validation",
			func(value string, shouldPass bool) {
				expectRuleResult(
					businessrules.CreditCard("card", value, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("valid Visa test number", "4111111111111111", true),
			Entry("valid with spaces", "4111 1111 1111 1111", true),
			Entry("valid Mastercard test number", "5500005555555559", true),
			Entry("empty", "", false),
			Entry("too short", "411111111111", false),
			Entry("fails Luhn checksum", "4111111111111112", false),
			Entry("contains letters", "411111111111111a", false),
		)
	})

	Describe("PhoneNumber", func() {
		DescribeTable(
			"validation",
			func(value string, shouldPass bool) {
				expectRuleResult(
					businessrules.PhoneNumber("phone", value, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("plain digits", "49301234567", true),
			Entry("with country code", "+49 30 1234567", true),
			Entry("with dashes and parens", "+1 (555) 123-4567", true),
			Entry("empty", "", false),
			Entry("too few digits", "+49 30 12", false),
			Entry("contains letters", "+49 30 abcdefg", false),
		)
	})

	Describe("PostalCode", func() {
		DescribeTable(
			"validation",
			func(value string, shouldPass bool) {
				expectRuleResult(
					businessrules.PostalCode("postal", value, businessrules.SeverityError).Check(),
					shouldPass,
				)
			},
			Entry("German style", "10115", true),
			Entry("UK style with space", "CB1 1TN", true),
			Entry("Canadian style with hyphen letters", "K1A-0B1", true),
			Entry("empty", "", false),
			Entry("too short", "12", false),
			Entry("contains symbols", "10115!", false),
		)
	})
})
