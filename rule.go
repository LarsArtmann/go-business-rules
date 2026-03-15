package businessrules

type Rule interface {
	Name() string
	Check() error
	Severity() Severity
	Message() string
}

type baseRule struct {
	name     string
	check    func() error
	severity Severity
	message  string
}

func (r baseRule) Name() string {
	return r.name
}

func (r baseRule) Check() error {
	return r.check()
}

func (r baseRule) Severity() Severity {
	return r.severity
}

func (r baseRule) Message() string {
	return r.message
}

func NewRule(name string, check func() error, severity Severity, message string) Rule {
	return baseRule{
		name:     name,
		check:    check,
		severity: severity,
		message:  message,
	}
}
