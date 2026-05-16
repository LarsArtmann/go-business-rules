package businessrules

import (
	"fmt"
	"net/url"
	"regexp"
)

var (
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	uuidPattern  = regexp.MustCompile(
		`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`,
	)
)

// checkNonEmpty validates that the value is not empty.
func checkNonEmpty(name, value string) error {
	if value == "" {
		return fmt.Errorf("%s must not be empty, got %q", name, value)
	}

	return nil
}

// Email creates a rule that validates the string is a valid email address.
// Uses a basic RFC 5322-compatible pattern for validation.
func Email(name, value string, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			err := checkNonEmpty(name, value)
			if err != nil {
				return err
			}

			if !emailPattern.MatchString(value) {
				return fmt.Errorf("%s must be a valid email address, got %q", name, value)
			}

			return nil
		},
		severity,
		name+" must be a valid email address",
	)
}

// URL creates a rule that validates the string is a valid HTTP/HTTPS URL.
// Use for validating web addresses and API endpoints.
func URL(name, value string, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			err := checkNonEmpty(name, value)
			if err != nil {
				return err
			}

			parsed, err := url.Parse(value)
			if err != nil {
				return fmt.Errorf("%s must be a valid URL: %w", name, err)
			}

			if parsed.Scheme != "http" && parsed.Scheme != "https" {
				return fmt.Errorf(
					"%s must be an HTTP or HTTPS URL, got scheme %q",
					name,
					parsed.Scheme,
				)
			}

			if parsed.Host == "" {
				return fmt.Errorf("%s must have a host, got %q", name, value)
			}

			return nil
		},
		severity,
		name+" must be a valid URL",
	)
}

// UUID creates a rule that validates the string is a valid UUID.
// Supports both uppercase and lowercase formats.
func UUID(name, value string, severity Severity) RuleImpl {
	return NewRule(
		name,
		func() error {
			err := checkNonEmpty(name, value)
			if err != nil {
				return err
			}

			if !uuidPattern.MatchString(value) {
				return fmt.Errorf("%s must be a valid UUID, got %q", name, value)
			}

			return nil
		},
		severity,
		name+" must be a valid UUID",
	)
}
