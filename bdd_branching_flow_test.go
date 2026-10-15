package businessrules_test

import (
	"encoding/json/v2"
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// BDD User Experience: Analyzing Code Quality with Branching-Flow
// User Goal: Run static analysis to validate code quality and understand expected violations.

var _ = Describe("Branching-Flow Integration", func() {
	var (
		modulePath string
		bfBin      string
	)

	BeforeEach(func() {
		modulePath = getModuleRoot()
		bfBin = findBranchingFlowBinary()

		if bfBin == "" {
			Skip("branching-flow binary not found; skipping analyzer integration specs")
		}
	})

	runBFCommand := func(args ...string) (string, error) {
		var stderr strings.Builder
		cmd := exec.Command(bfBin, args...)
		cmd.Env = os.Environ()
		cmd.Stderr = &stderr
		output, err := cmd.Output()
		// branching-flow (>= 0.6.7, 2026-10-15) exits 1 when linters report
		// findings, the golangci-lint convention, instead of only on
		// invocation errors, and it logs that error to stderr, which would
		// corrupt the JSON output the specs parse if merged in. A findings
		// exit is therefore expected here; missing binaries, usage errors,
		// and crashes still fail the Expect, with stderr as context.
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 1 {
			err = nil
		}
		Expect(err).ToNot(HaveOccurred(),
			"branching-flow %v failed: %s", args, stderr.String())

		return string(output), err
	}

	checkNoViolations := func(linter, expectedMsg string) {
		output, _ := runBFCommand(linter, modulePath)
		Expect(output).To(ContainSubstring(expectedMsg))
	}

	expectBFOutputContains := func(command, substr string) {
		output, _ := runBFCommand(command, modulePath)
		Expect(output).To(ContainSubstring(substr))
	}

	findViolation := func(result findingResult, file, name string) bool {
		for _, v := range result.Findings {
			if strings.Contains(v.Position.File, file) && v.Metadata.Name == name {
				return true
			}
		}

		return false
	}

	runPhantomCommand := func() findingResult {
		output, _ := runBFCommand("phantom", "--format", "finding", modulePath)

		var result findingResult

		err := json.Unmarshal([]byte(output), &result)
		Expect(err).ToNot(HaveOccurred())

		return result
	}

	runStatsCommand := func() statsResult {
		output, _ := runBFCommand("stats", "--format", "json", modulePath)

		var result statsResult

		err := json.Unmarshal([]byte(output), &result)
		Expect(err).ToNot(HaveOccurred())

		return result
	}

	Describe("Running branching-flow all", func() {
		It("should execute without errors", func() {
			cmd := exec.Command(bfBin, "all", modulePath)
			cmd.Env = os.Environ()
			output, err := cmd.CombinedOutput()
			if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 1 {
				err = nil
			}

			Expect(err).ToNot(HaveOccurred(),
				"branching-flow all should complete successfully.\nOutput: %s", string(output))
		})

		It("should analyze all Go source files", func() {
			expectBFOutputContains("all", "MULTI-LINTER ANALYSIS")
		})
	})

	Describe("PHANTOM violations (documented false positives)", func() {
		// Counts re-pinned 2026-10-11 against branching-flow 0.6.4 (192031d)
		// and re-verified unchanged 2026-10-15 under 0.6.7 (aa1611d, flake
		// bump): 20 findings, all PHANTOM-TRANSPOSE same-type parameter
		// pairs, all warning severity, across the four builder files. The
		// false-positive class is unchanged — "primitives used on purpose": the library
		// validates raw primitives and the builders take raw primitive
		// parameters by design. (The previous pin was 25 findings in mixed
		// severities under 0.2.0; 0.6.4 rewrote the rule set and metadata
		// shape, so both counts and matched names moved. Binary drift, not new
		// violations.) The analyzer has no path-exclude flag and scans the
		// whole directory tree, so ANY new module, example, or builder changes
		// these counts; that is the known pin fragility documented in
		// AGENTS.md. A changed count means: diff the findings, confirm they are
		// still the same false-positive class, then re-pin.
		It("should report exactly 20 PHANTOM violations", func() {
			result := runPhantomCommand()
			Expect(result.Summary.Total).To(Equal(20),
				"Expected 20 PHANTOM violations (builder swap-hazard pairs on deliberate primitives)")
		})

		expectSeverityCount := func(severity string, expected int) {
			result := runPhantomCommand()
			count := countSeverity(result.Findings, severity)
			Expect(count).To(Equal(expected))
		}

		It("should have 20 warning severity violations", func() {
			expectSeverityCount("warning", 20)
		})

		It("should have no critical, error, or info violations", func() {
			result := runPhantomCommand()
			Expect(countSeverity(result.Findings, "critical")).To(Equal(0))
			Expect(countSeverity(result.Findings, "error")).To(Equal(0))
			Expect(countSeverity(result.Findings, "info")).To(Equal(0))
		})

		expectViolation := func(file, name string) {
			result := runPhantomCommand()
			Expect(findViolation(result, file, name)).To(BeTrue())
		}

		DescribeTable(
			"should flag violations",
			func(file, name string) {
				expectViolation(file, name)
			},
			Entry("value/minimum pair in builders.go", "builders.go", "value, minimum"),
			Entry("value/divisor pair in builders_collection.go", "builders_collection.go", "value, divisor"),
			Entry("name/value pair in builders_network.go", "builders_network.go", "name, value"),
		)

		It("should flag same-type parameter pairs in every builder file", func() {
			result := runPhantomCommand()
			flaggedFiles := map[string]bool{}

			for _, v := range result.Findings {
				flaggedFiles[v.Position.File] = true
			}

			for _, file := range []string{"builders.go", "builders_collection.go", "builders_format.go", "builders_network.go"} {
				Expect(flaggedFiles).To(HaveKey(ContainSubstring(file)),
					"expected PHANTOM-TRANSPOSE findings in %s", file)
			}
		})
	})

	Describe("Other linters (no violations expected)", func() {
		It("should report no duplicate type violations", func() {
			checkNoViolations("dupe", "No duplicate types detected")
		})

		It("should report no panic conditions", func() {
			checkNoViolations("panic", "No panic conditions")
		})

		It("should report no primitive-ID violations (brandid-lint BD003)", func() {
			// branching-flow's strong-id rule is deprecated (moved to
			// brandid-lint BD003, linter/v0.4.0, parity-gated). The gate
			// shells out to brandid-lint directly: exit 0 = no findings.
			brandidBin := findBrandidLintBinary()
			if brandidBin == "" {
				Skip("brandid-lint binary not found; skipping BD003 integration spec")
			}

			cmd := exec.Command(brandidBin, "-rules", "BD003", modulePath)
			output, err := cmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred(),
				"brandid-lint BD003 findings (brand the IDs or suppress with "+
					"//brandid-lint:ignore(BD003) <reason>): %s", output)
		})

		It("should report no boolblind violations", func() {
			checkNoViolations("boolblind", "No boolean blindness")
		})
	})

	Describe("Stats command", func() {
		It("should run stats successfully", func() {
			expectBFOutputContains("stats", "Total")
		})

		It("should report the current total of 24 issues", func() {
			result := runStatsCommand()
			Expect(result.TotalIssues).To(Equal(24),
				"Expected 24 total issues across all linters (parsed from stats --format json): "+
					"20 phantom + 2 flagparam + 1 ifacecomplete + 1 mixins (pinned 2026-10-11 "+
					"against branching-flow 0.6.4, re-verified unchanged 2026-10-15 under 0.6.7: "+
					"phantom shrank 25→20 under the rewritten rule "+
					"set and the panic analyzer now reports zero, so the two former nolint-suppressed "+
					"panic findings dropped out of the total; flagparam/ifacecomplete/mixins unchanged), "+
					"all documented deliberate-structure findings; a mismatch almost always means "+
					"analyzer binary drift, not new violations")
		})
	})
})

// findingResult represents the JSON output from branching-flow phantom --format finding.
type findingResult struct {
	Findings []findingEntry `json:"findings"`
	Summary  findingSummary `json:"summary"`
}

// statsResult represents the relevant fields from branching-flow stats --format json.
type statsResult struct {
	TotalIssues int `json:"totalIssues"`
	ErrorCount  int `json:"errorCount"`
}

type findingEntry struct {
	Severity string          `json:"severity"`
	Position findingPosition `json:"position"`
	Metadata findingMetadata `json:"metadata"`
}

type findingPosition struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

type findingMetadata struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type findingSummary struct {
	Total      int            `json:"total"`
	BySeverity map[string]int `json:"bySeverity"`
}

func countSeverity(violations []findingEntry, severity string) int {
	count := 0

	for _, v := range violations {
		if v.Severity == severity {
			count++
		}
	}

	return count
}

func getModuleRoot() string {
	return findFileInParents("go.mod")
}

func findFileInParents(filename string) string {
	dir, err := os.Getwd()
	if err != nil {
		Fail("Failed to get current directory: " + err.Error())

		return ""
	}

	for {
		path := dir + "/" + filename

		_, err := os.Stat(path)
		if err == nil {
			return dir
		}

		parent := parentDir(dir)
		if parent == dir {
			Fail("Could not find " + filename + " in parent directories")

			return ""
		}

		dir = parent
	}
}

func parentDir(dir string) string {
	idx := strings.LastIndex(dir, "/")
	if idx <= 0 {
		return dir
	}

	return dir[:idx]
}

func findBrandidLintBinary() string {
	paths := []string{
		"brandid-lint",
		"/usr/local/bin/brandid-lint",
		os.Getenv("HOME") + "/go/bin/brandid-lint",
		os.Getenv("GOBIN") + "/brandid-lint",
	}

	for _, path := range paths {
		_, err := os.Stat(path)
		if err == nil {
			return path
		}

		_, err = exec.LookPath(path)
		if err == nil {
			return path
		}
	}

	return ""
}

func findBranchingFlowBinary() string {
	paths := []string{
		"branching-flow",
		"/usr/local/bin/branching-flow",
		os.Getenv("HOME") + "/go/bin/branching-flow",
		os.Getenv("GOBIN") + "/branching-flow",
	}

	for _, path := range paths {
		_, err := os.Stat(path)
		if err == nil {
			return path
		}

		_, err = exec.LookPath(path)
		if err == nil {
			return path
		}
	}

	return ""
}
