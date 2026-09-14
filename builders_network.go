package businessrules

import (
	"fmt"
	"net/netip"
	"regexp"
	"strings"
)

var (
	phonePattern  = regexp.MustCompile(`^\+?[0-9]{7,15}$`)
	postalPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 -]{2,9}$`)
)

const (
	minCardDigits        = 13
	maxCardDigits        = 19
	luhnDoubledThreshold = 9
)

// IPAddress creates a rule that validates the string parses as an IPv4 or
// IPv6 address. Use for validating network configuration or audit-log fields.
func IPAddress(name, value string, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			err := checkNonEmpty(name, value)
			if err != nil {
				return err
			}

			if _, err := netip.ParseAddr(value); err != nil {
				return fmt.Errorf("%s must be a valid IP address, got %q", name, value)
			}

			return nil
		},
		severity,
		name+" must be a valid IP address",
	)
}

// CreditCard creates a rule that validates the string is a syntactically
// plausible payment card number: 13-19 digits passing the Luhn checksum.
// It does NOT verify that a card exists or is active.
func CreditCard(name, value string, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			err := checkNonEmpty(name, value)
			if err != nil {
				return err
			}

			digits := strings.ReplaceAll(value, " ", "")
			if len(digits) < minCardDigits || len(digits) > maxCardDigits {
				return fmt.Errorf(
					"%s must be %d-%d digits, got %d characters",
					name,
					minCardDigits,
					maxCardDigits,
					len(digits),
				)
			}

			for _, r := range digits {
				if r < '0' || r > '9' {
					return fmt.Errorf("%s must contain only digits, got %q", name, value)
				}
			}

			if !passesLuhn(digits) {
				return fmt.Errorf("%s must pass the Luhn checksum, got %q", name, value)
			}

			return nil
		},
		severity,
		name+" must be a valid card number",
	)
}

// passesLuhn reports whether the digit string satisfies the Luhn checksum.
func passesLuhn(digits string) bool {
	sum := 0
	parity := len(digits) % 2

	for index, r := range digits {
		digit := int(r - '0')
		if index%2 == parity {
			digit *= 2
			if digit > luhnDoubledThreshold {
				digit -= luhnDoubledThreshold
			}
		}

		sum += digit
	}

	return sum%10 == 0
}

// PhoneNumber creates a rule that validates the string is a plausible phone
// number: an optional leading '+' followed by 7-15 digits, with spaces,
// dashes, dots, and parentheses allowed as separators. Country-specific
// numbering plans are out of scope.
func PhoneNumber(name, value string, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			err := checkNonEmpty(name, value)
			if err != nil {
				return err
			}

			normalized := strings.Map(func(r rune) rune {
				switch r {
				case ' ', '-', '.', '(', ')':
					return -1
				default:
					return r
				}
			}, value)

			if !phonePattern.MatchString(normalized) {
				return fmt.Errorf(
					"%s must be a phone number with 7-15 digits, got %q",
					name,
					value,
				)
			}

			return nil
		},
		severity,
		name+" must be a valid phone number",
	)
}

// PostalCode creates a rule that validates the string is a plausible postal
// code: 3-10 characters from letters, digits, spaces, and hyphens. It cannot
// verify country-specific formats; pair it with MatchesFunc when a strict
// national pattern is required.
func PostalCode(name, value string, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			err := checkNonEmpty(name, value)
			if err != nil {
				return err
			}

			if !postalPattern.MatchString(value) {
				return fmt.Errorf("%s must be a valid postal code, got %q", name, value)
			}

			return nil
		},
		severity,
		name+" must be a valid postal code",
	)
}
