package businessrules_test

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sort"

	. "github.com/onsi/ginkgo/v2"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	businessrules "github.com/LarsArtmann/go-business-rules/v2"
)

// rulesFromValues maps generated values to deterministic rules: value < 0
// fails, value >= 0 passes. Names are unique per index so violation sets
// never deduplicate.
func rulesFromValues(values []int) []businessrules.Rule {
	rules := make([]businessrules.Rule, 0, len(values))

	for index, value := range values {
		captured := value
		name := fmt.Sprintf("rule_%d", index)
		rules = append(rules, businessrules.NewRule(
			name,
			func() error {
				if captured < 0 {
					return fmt.Errorf("%s must be non-negative, got %d", name, captured)
				}

				return nil
			},
			businessrules.SeverityError,
			name+" must be non-negative",
		))
	}

	return rules
}

func violationIdentities(result businessrules.ValidationResultError) []string {
	names := make([]string, 0, len(result.ViolationErrors))

	for _, violation := range result.ViolationErrors {
		names = append(names, violation.Rule.Name())
	}

	return names
}

func drainStreamEvents(events <-chan businessrules.Event) businessrules.ValidationResultError {
	var result businessrules.ValidationResultError

	for event := range events {
		if completed, ok := event.(businessrules.ValidationCompleted); ok {
			result = completed.Result
		}
	}

	return result
}

var _ = Describe("Build and Stream equivalence", Label("property"), func() {
	BeforeEach(func() {
		if os.Getenv("GBR_PROPERTY") == "" {
			Skip("opt-in: run with GBR_PROPERTY=1 to execute property-based tests")
		}
	})

	It("agree on the violation set for arbitrary rule values", func() {
		property := prop.ForAll(
			func(values []int) bool {
				rules := rulesFromValues(values)

				buildResult := businessrules.NewValidator().AddRules(rules...).Build()
				streamResult := drainStreamEvents(
					businessrules.NewValidator().AddRules(rules...).Stream(context.Background()),
				)

				buildNames := violationIdentities(buildResult)
				streamNames := violationIdentities(streamResult)

				return len(buildNames) == len(streamNames) &&
					sortedEqual(buildNames, streamNames)
			},
			gen.SliceOf(gen.IntRange(-100, 100)),
		)

		result := property.Check(gopter.NewDefaultRunner())
		Expect(result.Passed()).To(BeTrue(), "property failed: %v", result.Error)
	})
})

func sortedEqual(first, second []string) bool {
	sortedFirst := append([]string(nil), first...)
	sortedSecond := append([]string(nil), second...)

	sort.Strings(sortedFirst)
	sort.Strings(sortedSecond)

	return slices.Equal(sortedFirst, sortedSecond)
}
