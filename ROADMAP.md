# Roadmap — businessrules

> Aspirational items with no timeline. Ideas for future enhancement.

---

## Additional Rule Builders

### Collection Rules

- `NotEmptySlice(name, value, severity)` — slice has elements
- `MinSliceLength(name, value, minimum, severity)` — minimum slice length
- `MaxSliceLength(name, value, maximum, severity)` — maximum slice length
- `Unique(name, value, severity)` — all slice elements unique

### Time/Date Rules

- `NotPast(name, value, severity)` — date/time not in the past
- `NotFuture(name, value, severity)` — date/time not in the future
- `InRange(name, value, min, max, severity)` — date/time within range

### Numeric Precision Rules

- `MaxDecimalPlaces(name, value, max, severity)` — limit decimal places
- `DivisibleBy(name, value, divisor, severity)` — divisible by number

### String Rules

- `NoWhitespace(name, value, severity)` — no whitespace characters
- `Alphanumeric(name, value, severity)` — only alphanumeric characters
- `Numeric(name, value, severity)` — only numeric characters
- `Alpha(name, value, severity)` — only alphabetic characters

### Network/Identifier Rules

- `IPAddress(name, value, severity)` — IPv4 or IPv6
- `CreditCard(name, value, severity)` — Luhn algorithm validation
- `PhoneNumber(name, value, severity)` — phone number format
- `PostalCode(name, value, country, severity)` — country-specific postal codes

## Advanced Features

### Rule Composition

- `Not(rule)` — negate a rule
- `Or(rules...)` — logical OR composition
- `Xor(rule1, rule2)` — exactly one rule passes

### Context-Aware Validation

- `WithContext(key, value)` — pass context to rules
- Rules that can access shared context during validation

### Asynchronous Rules

- `Async(name, asyncCheck, severity)` — rules that check asynchronously
- `Await()` method on ValidationResult

### Rule Metadata

- `Description()` — human-readable description
- `Tags()` — categorize rules (e.g., "security", "performance")
- `Priority()` — execution order hints

### Validation Groups

- Groups of rules that can be enabled/disabled together
- Group inheritance for complex validation scenarios

### Integration Suggestions

- `sivchari/govalid` — structural validation (compile-time safe)
- `asaskevich/validator` — comprehensive format validation
- `go-playground/validator` — struct tag-based validation
- Context propagation for tracing

### Performance

- Benchmark validation performance
- Parallel rule execution option
- Caching for expensive regex compilations

### Developer Experience

- Generate rule documentation from code
- CLI tool for running validations
- Interactive playground/examples

---

_Last Updated: 2026-04-05_
