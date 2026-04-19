package businessrules_test

import (
	"encoding/json"
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

	findViolation := func(result phantomResult, file, name string) bool {
		for _, v := range result.Violations {
			if strings.Contains(v.Location, file) && v.Name == name {
				return true
			}
		}

		return false
	}

	runPhantomCommand := func() phantomResult {
		output, _ := runBFCommand("phantom", "--format", "json", modulePath)

		var result phantomResult

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
		It("should report exactly 15 PHANTOM violations", func() {
			result := runPhantomCommand()
			Expect(result.Count).To(Equal(15),
				"Expected 15 PHANTOM violations (documented false positives)")
		})

		expectSeverityCount := func(severity string, expected int) {
			result := runPhantomCommand()
			count := countSeverity(result.Violations, severity)
			Expect(count).To(Equal(expected))
		}

		It("should have 5 critical severity violations", func() {
			expectSeverityCount("critical", 5)
		})

		It("should have 8 low severity violations", func() {
			expectSeverityCount("low", 8)
		})

		expectViolation := func(file, name string) {
			result := runPhantomCommand()
			Expect(findViolation(result, file, name)).To(BeTrue())
		}

		DescribeTable("should flag violations",
			func(file, name string) {
				expectViolation(file, name)
			},
			Entry("errMsg in builders.go", "builders.go", "errMsg"),
			Entry("context in errors.go", "errors.go", "context"),
		)

		It("should flag bool condition parameter in builders.go and builders_composite.go", func() {
			result := runPhantomCommand()
			found := false

			for _, v := range result.Violations {
				if v.Name == "condition" &&
					(strings.Contains(v.Location, "builders.go") ||
						strings.Contains(v.Location, "builders_composite.go")) {
					found = true

					break
				}
			}

			Expect(found).To(BeTrue())
		})
	})

	Describe("Other linters (no violations expected)", func() {
		It("should report no context violations", func() {
			checkNoViolations("context", "No semantic context issues")
		})

		It("should report no duplicate type violations", func() {
			checkNoViolations("dupe", "No duplicates found")
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
			expectBFOutputContains("stats", "Total Issues")
		})

		It("should report 15 total issues", func() {
			expectBFOutputContains("stats", "Total Issues: 15")
		})
	})
})

// phantomResult represents the JSON output from branching-flow phantom.
type phantomResult struct {
	Target     string             `json:"target"`
	Count      int                `json:"count"`
	Violations []phantomViolation `json:"violations"`
}

type phantomViolation struct {
	Location      string `json:"location"`
	Kind          string `json:"kind"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	SuggestedType string `json:"suggestedType,omitempty"`
	Severity      string `json:"severity"`
	Message       string `json:"message"`
}

func countSeverity(violations []phantomViolation, severity string) int {
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
