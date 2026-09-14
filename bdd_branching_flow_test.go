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
	})

	runBFCommand := func(args ...string) (string, error) {
		cmd := exec.Command(bfBin, args...)
		cmd.Env = os.Environ()
		output, err := cmd.CombinedOutput()
		Expect(err).ToNot(HaveOccurred())

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

			Expect(err).ToNot(HaveOccurred(),
				"branching-flow all should complete successfully.\nOutput: %s", string(output))
		})

		It("should analyze all Go source files", func() {
			expectBFOutputContains("all", "MULTI-LINTER ANALYSIS")
		})
	})

	Describe("PHANTOM violations (documented false positives)", func() {
		// Counts re-pinned 2026-09-14 (18 total): the false-positive class is
		// "primitives used on purpose" — the library validates raw primitives,
		// the SSE example validates raw primitives, and the builders take raw
		// primitive parameters by design. The analyzer has no path-exclude flag
		// and scans the whole directory tree, so ANY new module, example, or
		// builder changes these counts; that is the known pin fragility
		// documented in AGENTS.md. A changed count means: diff the findings,
		// confirm they are still the same false-positive class, then re-pin.
		It("should report exactly 18 PHANTOM violations", func() {
			result := runPhantomCommand()
			Expect(result.Summary.Total).To(Equal(18),
				"Expected 18 PHANTOM violations (library + example + builder primitives, all deliberate)")
		})

		expectSeverityCount := func(severity string, expected int) {
			result := runPhantomCommand()
			count := countSeverity(result.Findings, severity)
			Expect(count).To(Equal(expected))
		}

		It("should have 6 critical severity violations", func() {
			expectSeverityCount("critical", 6)
		})

		It("should have 4 error severity violations", func() {
			expectSeverityCount("error", 4)
		})

		It("should have 7 info severity violations", func() {
			expectSeverityCount("info", 7)
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
			Entry("errMsg in builders.go", "builders.go", "errMsg"),
			Entry("context in errors.go", "errors.go", "context"),
		)

		It("should flag bool condition parameter in builders.go and builders_composite.go", func() {
			result := runPhantomCommand()
			found := false

			for _, v := range result.Findings {
				if v.Metadata.Name == "condition" &&
					(strings.Contains(v.Position.File, "builders.go") ||
						strings.Contains(v.Position.File, "builders_composite.go")) {
					found = true

					break
				}
			}

			Expect(found).To(BeTrue())
		})
	})

	Describe("Other linters (no violations expected)", func() {
		It("should report no duplicate type violations", func() {
			checkNoViolations("dupe", "No duplicate types detected")
		})

		It("should report no panic conditions", func() {
			checkNoViolations("panic", "No panic conditions")
		})

		It("should report no strong-id violations", func() {
			checkNoViolations("strong-id", "No string ID parameters")
		})

		It("should report no boolblind violations", func() {
			checkNoViolations("boolblind", "No boolean blindness")
		})
	})

	Describe("Stats command", func() {
		It("should run stats successfully", func() {
			expectBFOutputContains("stats", "Total")
		})

		It("should report the current total of 22 issues", func() {
			result := runStatsCommand()
			Expect(result.TotalIssues).To(Equal(22),
				"Expected 22 total issues across all linters (parsed from stats --format json); "+
					"a mismatch almost always means analyzer binary drift, not new violations")
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

	Fail("Could not find branching-flow binary")

	return ""
}
