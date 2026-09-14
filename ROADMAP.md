# Roadmap — businessrules

> Long-term direction and raw ideas **not yet refined into actionable tasks**.
>
> When an item here becomes bounded and short-term, it graduates into [`TODO_LIST.md`](TODO_LIST.md). When an item ships, it is recorded in [`CHANGELOG.md`](CHANGELOG.md) and reflected in [`FEATURES.md`](FEATURES.md).

**Last reviewed:** 2026-09-14 against `master`.

---

## Additional Rule Builders (candidates)

### Time / Date

- `NotPast(name, value, severity)` — date/time not in the past
- `NotFuture(name, value, severity)` — date/time not in the future
- `DateInRange(name, value, min, max, severity)` — date/time within range

### Numeric Precision

- `MaxDecimalPlaces(name, value, max, severity)` — limit decimal places
- `DivisibleBy(name, value, divisor, severity)` — divisible by a number

### String

- `Contains(name, value, substring, severity)` — substring containment
- `LengthRange(name, value, minimum, maximum, severity)` — bounded length (from the 2026-03-21 lint-session report)
- `Required(name, value, severity)` — `NotEmpty` + `NotBlank` combined
- `MatchesFunc(name, value, fn func(string) bool, severity)` — function-based matching

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

## Event-Driven Validation (candidates, from 2026-09-14)

- Rule scheduler with priorities / short-circuit on `Critical` during `Stream`
- Listener backpressure / async queue semantics (listeners are synchronous today, by design)
- Signing validation events (go-cqrs-lite recipes) for tamper-evident audit trails
- Run ID (unique per validation run) on `ValidationCompleted` for distributed provenance
- Rule metadata (tags) surfaced on `RuleEvaluated`
- Property-based `Build` ≡ `Stream` equivalence invariant (e.g. gopter)
- Event JSON marshaling (deliberate YAGNI today — add only when a consumer asks)
- `RuleID` branded type via `go-composable-business-types/id` (analysis in `docs/planning/archived/`; breaking change, deferred since 2026-03 — revisit at the next major version)

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

## Open questions (user decisions needed)

1. **Module versioning endgame** — rename to `.../v2` + re-tag (every consumer import changes), or fresh compatible tag on the current path? (TODO_LIST carries the actionable item; the direction is a release-engineering call.)
2. **Branching-flow stale pins** — re-pin to actual, keep red as a signal, or move to nightly? (TODO_LIST.)
3. **CI on a private repo** — re-enable the disabled workflow (Actions minutes cost) or rely on local flake gates?
4. **Adapter home** — keep `adapters/cqrslite` and `examples/sse` as nested modules here, or promote to sibling repos (collector-extraction pattern)?

---

_Last Updated: 2026-09-14_
