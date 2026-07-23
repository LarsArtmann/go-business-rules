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

	Describe("Running branching-flow all", func() {
		It("should execute without errors", func() {
			cmd := exec.Command(bfBin, "all", modulePath)
			cmd.Env = os.Environ()
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred(),
				"branching-flow all should complete successfully.\nOutput: %s", string(output))
		})

		It("should analyze all Go source files", func() {
			expectBFOutputContains("all", "Files Analyzed")
		})
	})

	Describe("PHANTOM violations (documented false positives)", func() {
		It("should report exactly 12 PHANTOM violations", func() {
			result := runPhantomCommand()
			Expect(result.Summary.Total).To(Equal(12),
				"Expected 12 PHANTOM violations (documented false positives)")
		})

		expectSeverityCount := func(severity string, expected int) {
			result := runPhantomCommand()
			count := countSeverity(result.Findings, severity)
			Expect(count).To(Equal(expected))
		}

		It("should have 5 critical severity violations", func() {
			expectSeverityCount("critical", 5)
		})

		It("should have 1 error severity violation", func() {
			expectSeverityCount("error", 1)
		})

		It("should have 6 info severity violations", func() {
			expectSeverityCount("info", 6)
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

		It("should report the current total of 36 issues", func() {
			expectBFOutputContains("stats", "36")
		})
	})
})

// findingResult represents the JSON output from branching-flow phantom --format finding.
type findingResult struct {
	Findings []findingEntry `json:"findings"`
	Summary  findingSummary  `json:"summary"`
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
