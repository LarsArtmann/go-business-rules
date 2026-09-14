package businessrules_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

var _ = Describe("Time Builders", func() {
	// Entries are evaluated eagerly at tree construction, so relative times
	// are expressed as offsets from a `now` computed inside each spec body.
	var referenceNow = time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

	Describe("NotPast", func() {
		DescribeTable(
			"validation",
			func(offset time.Duration, shouldPass bool) {
				now := referenceNow
				expectRuleResult(
					businessrules.NotPast("val", now.Add(offset), now, businessrules.SeverityError).
						Check(),
					shouldPass,
				)
			},
			Entry("future time", 24*time.Hour, true),
			Entry("exactly now", 0, true),
			Entry("past time", -24*time.Hour, false),
		)
	})

	Describe("NotFuture", func() {
		DescribeTable(
			"validation",
			func(offset time.Duration, shouldPass bool) {
				now := referenceNow
				expectRuleResult(
					businessrules.NotFuture("val", now.Add(offset), now, businessrules.SeverityError).
						Check(),
					shouldPass,
				)
			},
			Entry("past time", -24*time.Hour, true),
			Entry("exactly now", 0, true),
			Entry("future time", 24*time.Hour, false),
		)
	})

	Describe("DateInRange", func() {
		var (
			start time.Time
			end   time.Time
		)

		BeforeEach(func() {
			start = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			end = time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)
		})

		DescribeTable(
			"validation",
			func(value time.Time, shouldPass bool) {
				expectRuleResult(
					businessrules.DateInRange(
						"val",
						value,
						start,
						end,
						businessrules.SeverityError,
					).Check(),
					shouldPass,
				)
			},
			Entry("within range", time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC), true),
			Entry("at start boundary", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), true),
			Entry("at end boundary", time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC), true),
			Entry("before range", time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC), false),
			Entry("after range", time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), false),
		)
	})
})
