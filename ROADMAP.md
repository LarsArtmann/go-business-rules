# Roadmap — businessrules

> Long-term direction and raw ideas **not yet refined into actionable tasks**.
>
> When an item here becomes bounded and short-term, it graduates into [`TODO_LIST.md`](TODO_LIST.md). When an item ships, it is recorded in [`CHANGELOG.md`](CHANGELOG.md) and reflected in [`FEATURES.md`](FEATURES.md).

**Last reviewed:** 2026-07-26 against `master`.

---

## Already shipped (not roadmap — recorded for context)

The following were once roadmap candidates and are now `FULLY_FUNCTIONAL` (see FEATURES.md):

- Collection rules: `NotEmptySlice`, `NotEmptyMap`
- Extended numeric: `GreaterThan`, `LessThan`
- String: `NotBlank`
- Generic: `Equals`
- Composite: `All`, `Any`, `When`

## Additional Rule Builders (candidates)

### Time / Date

- `NotPast(name, value, severity)` — date/time not in the past
- `NotFuture(name, value, severity)` — date/time not in the future
- `DateInRange(name, value, min, max, severity)` — date/time within range

### Numeric Precision

- `MaxDecimalPlaces(name, value, max, severity)` — limit decimal places
- `DivisibleBy(name, value, divisor, severity)` — divisible by a number

### String

- `NoWhitespace(name, value, severity)` — no whitespace characters
- `Alphanumeric(name, value, severity)` — only alphanumeric characters
- `Numeric(name, value, severity)` — only numeric characters
- `Alpha(name, value, severity)` — only alphabetic characters

### Network / Identifier

- `IPAddress(name, value, severity)` — IPv4 or IPv6
- `CreditCard(name, value, severity)` — Luhn algorithm validation
- `PhoneNumber(name, value, severity)` — phone number format
- `PostalCode(name, value, country, severity)` — country-specific postal codes

## Advanced Features (candidates)

### Rule Composition

- `Not(rule)` — negate a rule
- `Or(rules...)` — logical OR composition (distinct from `Any` by signature shape)
- `Xor(rule1, rule2)` — exactly one rule passes

### Context-Aware Validation

- `WithContext(key, value)` — pass context to rules
- Rules that can access shared context during a validation session

### Asynchronous Rules

- `Async(name, asyncCheck, severity)` — rules that check asynchronously
- `Await()` method on `ValidationResultError`

### Rule Metadata

- `Description()` — human-readable description
- `Tags()` — categorize rules (e.g. "security", "performance")
- `Priority()` — execution-order hints

### Validation Groups

- Groups of rules that can be enabled/disabled together
- Group inheritance for complex validation scenarios

## Ecosystem & Developer Experience

- Tighter integration examples with `sivchari/govalid`, `go-playground/validator`, `asaskevich/validator`
- Context propagation for tracing / OpenTelemetry
- Generate rule documentation from code
- CLI tool for running validations
- Interactive playground / more real-world examples

## Performance

- Benchmark validation performance baseline and publish numbers
- Parallel rule execution option
- Caching for expensive regex compilations

## Versioning & Stability

- Decide whether to adopt `/v2` module suffix if breaking changes continue
- API stability review once a real-world consumer (Polish-Customs) is integrated
- Re-evaluate the `encoding/json/v2` / `GOEXPERIMENT=jsonv2` downstream constraint when the package graduates from experimental

---

_Last Updated: 2026-07-26_
