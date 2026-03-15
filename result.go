package businessrules

type Result struct {
	Valid      bool
	Violations []Violation
}

func (r Result) Errors() []Violation {
	return r.BySeverity(SeverityError, SeverityCritical)
}

func (r Result) Warnings() []Violation {
	return r.BySeverity(SeverityWarning)
}

func (r Result) Info() []Violation {
	return r.BySeverity(SeverityInfo)
}

func (r Result) Critical() []Violation {
	return r.BySeverity(SeverityCritical)
}

func (r Result) BySeverity(severities ...Severity) []Violation {
	severitySet := make(map[Severity]bool, len(severities))
	for _, s := range severities {
		severitySet[s] = true
	}

	var result []Violation
	for _, v := range r.Violations {
		if severitySet[v.Rule.Severity()] {
			result = append(result, v)
		}
	}
	return result
}

func (r Result) HasErrors() bool {
	return len(r.Errors()) > 0
}

func (r Result) HasWarnings() bool {
	return len(r.Warnings()) > 0
}

func (r Result) HasCritical() bool {
	return len(r.Critical()) > 0
}

func (r Result) HasInfo() bool {
	return len(r.Info()) > 0
}
