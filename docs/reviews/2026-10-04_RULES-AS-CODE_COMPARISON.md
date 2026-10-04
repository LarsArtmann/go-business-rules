# Rules-as-Code Landscape Comparison: go-business-rules vs OpenFisca, Catala, smucclaw/L4, Blawx

**Date:** 2026-10-04
**Type:** Point-in-time landscape review
**Scope:** How go-business-rules compares to four prominent rules-as-code / computational-law systems, and where we fit.

---

## TL;DR

go-business-rules is in a **different category** from all four systems. They are
_legislation-as-code_ engines: they compute what the law entitles or obliges,
with period-aware time handling, default logic (exceptions), and provenance to
legal text. We are a _severity-aware validation library_: we check whether an
input violates a rule and how badly (Info → Critical). Complementary, not
competing: we sit at the edge (validating inputs before they reach such
engines) or encode simple eligibility gates. Becoming comparable would require
temporal rules, default logic, and provenance — essentially building Catala in
Go.

---

## The Four Systems

### OpenFisca (openfisca/openfisca-core)

- **What:** Microsimulation framework that models tax and social-benefit
  legislation as executable code ("legislation as code", AGPL-3.0). Core is
  country-agnostic; actual rules live in country packages (OpenFisca-France et
  al.).
- **Signature features:** first-class **time periods** (retroactive simulation
  over arbitrary dates), legislation **parameters** (thresholds, rates) served
  separately from formulas via an API, built-in REST API (`openfisca serve`),
  YAML scenario tests.
- **Stack:** Python, tightly coupled to NumPy (vectorized computation).
- **Users:** governments publishing rules as APIs, economists/social-policy
  researchers, civic tech.

### Catala (CatalaLang/catala)

- **What:** A domain-specific language for deriving "faithful-by-construction"
  algorithms from legislative texts (Inria research project, Apache-2.0,
  ~2.4k stars). Literate programming of law: code interleaved verbatim with the
  legal article it implements.
- **Signature features:** **default logic as a language primitive** (formalized
  from Sarah Lawsky's _A Logic for Statutes_) — exceptions and
  exceptions-to-exceptions are native; every computation traces back to the
  legal article justifying it (**provenance**); compiler emits a
  **lawyer-readable PDF** for review by domain experts.
- **Stack:** OCaml compiler (formally verified core), JS runtime backend, LSP,
  tree-sitter grammar. Flagship example: `CatalaLang/french-law`.
- **Users:** lawyers reviewing law-as-code for correctness + developers of
  socio-fiscal computation systems.

### smucclaw (SMU Centre for Computational Law) — the L4 language family

- **What:** A research organization (73 repos), not a single tool. Main
  artifact is the **L4** DSL family: Natural L4 (spreadsheet "LegalSS" surface
  syntax with English keywords: EVERY/WHO/MUST/MEANS/HENCE/LEST), Core/Baby L4
  for rigorous decision logic, plus `nlg` (natural-language generation of legal
  text from formalisms).
- **Signature features:** **transpilation to many targets** — PureScript,
  Haskell, Prolog (SWIPL/Clingo/s(CASP)), ASP, TypeScript, Petri nets, with
  DMN, Alloy, Uppaal, Catala and more on the roadmap. Explicitly distinguishes
  _legal engineers_ (author rules) from _toolchain developers_.
- **Stack:** primarily Haskell (compilers), Clojure (in-browser IDE `l4-lp`),
  Grammatical Framework (baby-l4, nlg).
- **Users:** legal engineers and computational-law researchers.

### Blawx (Lexpedite/blawx)

- **What:** Web-based, visual rules-as-code tool by Jason Morris (MIT).
  Drag-and-drop rule authoring on Google Blockly; reasoning via SWI-Prolog and
  s(CASP).
- **Signature features:** **explanations for answers**, **hypothetical
  reasoning** (what-if queries), Scenario Explorer, REST API. Claims to be the
  only open-source RaC environment combining these. Note: representation is
  Blockly-based; the README does **not** mention LegalRuleML.
- **Stack:** Python/Django server, Blockly frontend, Docker deployment.
- **Users:** explicitly educational and experimental — **not production-ready**
  per its own README.

---

## Comparison

| Dimension                     | go-business-rules                                          | OpenFisca                                | Catala                                               | smucclaw / L4                         | Blawx                                                  |
| ----------------------------- | ---------------------------------------------------------- | ---------------------------------------- | ---------------------------------------------------- | ------------------------------------- | ------------------------------------------------------ |
| Category                      | Validation library                                         | Microsimulation engine                   | Legal DSL (statutes → code)                          | DSL family + transpilers              | Visual RaC web tool                                    |
| Question answered             | "Does input X violate rule Y, how badly?"                  | "What tax/benefit is owed for period P?" | "What does this statute compute?"                    | "How do norms formalize and execute?" | "Is this scenario consistent with the rules, and why?" |
| Reasoning model               | Boolean checks + severity axis                             | Formulas over parameters, vectorized     | Default logic (exceptions)                           | Decision logic, abduction, workflows  | ASP / s(CASP), Prolog                                  |
| Severity/outcome nuance       | **First-class** (Info/Warning/Error/Critical)              | Pass/fail per test scenario              | Exception resolution, not severities                 | Norm-level (obligation/permission)    | Answer sets with justification                         |
| Time modeling                 | Point-in-time                                              | First-class periods, retroactive         | Legal-change aware, temporal logic in use cases      | Varies by transpile target            | Scenario-based                                         |
| Provenance                    | Violation carries rule name/message/timestamp              | Inspectable formulas + parameter API     | **Literate: code ↔ legal article**                   | NLG back to natural language          | Answer explanations                                    |
| Explanations                  | Violation messages + event stream                          | API inspectability                       | Lawyer-readable PDF from code                        | nlg grammars                          | **Core feature**                                       |
| Streaming/async               | Event stream, cancellation, bounded concurrency            | Batch simulation                         | Compiled programs                                    | Target-dependent                      | Query/response API                                     |
| Language/ecosystem            | Go, ~zero runtime deps                                     | Python + NumPy                           | OCaml (JS backend)                                   | Haskell, Clojure, GF                  | Python/Django, Blockly                                 |
| Embeddability in a Go service | **Native library**                                         | External service (REST)                  | Compiled artifacts via API                           | Transpiled outputs                    | External service (REST)                                |
| Maturity claim                | v2.2.0, 251 root specs, published & consumed in production | Production (multi-country deployments)   | Research-grade, compiler "still unstable" per README | Research programme                    | Educational, not production (README)                   |

---

## Gap Analysis: What They Have That We Don't

1. **Period/time-aware rules.** OpenFisca computes over arbitrary time spans;
   Catala handles law changes over time. Our rules are point-in-time. Any
   retroactivity or "as of date D" question is out of scope today.
2. **Default logic.** Law is exceptions-to-exceptions. Catala makes this a
   language primitive; Blawx and s(CASP)-based tools get it from ASP.
   Our composites (`All`, `Any`, `When`, `Not`, `Or`, `Xor`) cover flat
   boolean logic, not defeasible reasoning with priority.
3. **Normative computation.** They compute amounts, tiers, entitlements.
   We return violations, never values. Nothing in our API accumulates a
   calculated entitlement.
4. **Legal provenance.** Catala's code interleaves with the statute; every
   result cites its article. Our `ViolationError` carries rule name/message —
   not a citation to an authoritative text.
5. **Non-programmer audiences.** Blawx (Blockly) and L4 (spreadsheets) target
   lawyers directly. We are a Go API for developers.

## What We Offer That They Don't

- A **lightweight, embeddable Go primitive**: import and call in-process.
  All four are engines/toolchains/services — none is a library you'd put in a
  Go service's request path.
- **Severity as a first axis** (Info → Critical) with filtering — none of them
  model graded violation severities; their outcomes are entitlements or answers,
  not graded violations.
- **Streaming validation with cancellation and bounded concurrency** —
  `Stream(ctx)` with events, semaphore-bounded checks, zero-cost default path.
  Batch simulators have nothing comparable for live request validation.
- Near-zero dependency footprint and Go-native events (OTel/CQRS adapters
  already exist as sibling modules).

---

## Positioning

**Complementary, not competing.** Realistic placements:

- **Edge validation:** validate entities _before_ feeding them to an OpenFisca/
  Catala/Blawx simulation; fail fast with severity instead of computing on
  garbage.
- **Simple eligibility gates:** encode binary eligibility rules (age bounds,
  residency flags) with severities in-process, deferring amount computation to
  a proper engine.
- **Operational instrumentation:** our OTel adapter gives validation the
  observability these engines lack.

**Not our lane (and would require building Catala in Go to claim):** period
math, defeasible reasoning, entitlement computation, legal-text provenance.

---

## Sources

- https://github.com/openfisca/openfisca-core (fetched 2026-10-04)
- https://openfisca.org/doc/
- https://github.com/CatalaLang/catala (fetched 2026-10-04)
- https://book.catala-lang.org/
- https://github.com/smucclaw, https://github.com/smucclaw/dsl, https://github.com/smucclaw/nlg (fetched 2026-10-04)
- https://github.com/Lexpedite/blawx (fetched 2026-10-04)

_Point-in-time document: reflects the state of the external projects as of
2026-10-04. Do not treat as a backlog._
