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

	Describe("Running branching-flow all", func() {
		It("should execute without errors", func() {
			cmd := exec.Command(bfBin, "all", modulePath)
			cmd.Env = os.Environ()
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred(),
				"branching-flow all should complete successfully.\nOutput: %s", string(output))
		})

		It("should analyze all Go source files", func() {
			cmd := exec.Command(bfBin, "all", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("Files Analyzed"),
				"Should report files analyzed")
		})
	})

	Describe("PHANTOM violations (documented false positives)", func() {
		// These violations are documented in AGENTS.md as false positives.
		// This library validates raw primitives - forcing branded types would
		// defeat the library's purpose of accepting any comparable value.

		It("should report exactly 15 PHANTOM violations", func() {
			cmd := exec.Command(bfBin, "phantom", "--format", "json", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())

			var result phantomResult
			err = json.Unmarshal(output, &result)
			Expect(err).ToNot(HaveOccurred())

			Expect(result.Count).To(Equal(15),
				"Expected 15 PHANTOM violations (documented false positives)")
		})

		It("should have 5 critical severity violations", func() {
			cmd := exec.Command(bfBin, "phantom", "--format", "json", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())

			var result phantomResult
			err = json.Unmarshal(output, &result)
			Expect(err).ToNot(HaveOccurred())

			count := countSeverity(result.Violations, "critical")
			Expect(count).To(Equal(5),
				"Expected 5 critical PHANTOM violations")
		})

		It("should have 7 low severity violations", func() {
			cmd := exec.Command(bfBin, "phantom", "--format", "json", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())

			var result phantomResult
			err = json.Unmarshal(output, &result)
			Expect(err).ToNot(HaveOccurred())

			count := countSeverity(result.Violations, "low")
			Expect(count).To(Equal(7),
				"Expected 7 low PHANTOM violations")
		})

		It("should flag string parameters in builders.go", func() {
			cmd := exec.Command(bfBin, "phantom", "--format", "json", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())

			var result phantomResult
			err = json.Unmarshal(output, &result)
			Expect(err).ToNot(HaveOccurred())

			found := false
			for _, v := range result.Violations {
				if strings.Contains(v.Location, "builders.go") && v.Name == "errMsg" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue(),
				"Should flag errMsg parameter in builders.go")
		})

		It("should flag context parameter in errors.go", func() {
			cmd := exec.Command(bfBin, "phantom", "--format", "json", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())

			var result phantomResult
			err = json.Unmarshal(output, &result)
			Expect(err).ToNot(HaveOccurred())

			found := false
			for _, v := range result.Violations {
				if strings.Contains(v.Location, "errors.go") && v.Name == "context" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue(),
				"Should flag context parameter in errors.go")
		})

		It("should flag name/message fields in rule.go", func() {
			cmd := exec.Command(bfBin, "phantom", "--format", "json", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())

			var result phantomResult
			err = json.Unmarshal(output, &result)
			Expect(err).ToNot(HaveOccurred())

			hasName := false
			hasMessage := false
			for _, v := range result.Violations {
				if strings.Contains(v.Location, "rule.go") && v.Name == "name" {
					hasName = true
				}
				if strings.Contains(v.Location, "rule.go") && v.Name == "message" {
					hasMessage = true
				}
			}
			Expect(hasName).To(BeTrue(), "Should flag name field in rule.go")
			Expect(hasMessage).To(BeTrue(), "Should flag message field in rule.go")
		})

		It("should flag bool condition parameter in builders.go and builders_composite.go", func() {
			cmd := exec.Command(bfBin, "phantom", "--format", "json", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())

			var result phantomResult
			err = json.Unmarshal(output, &result)
			Expect(err).ToNot(HaveOccurred())

			hasBoolCondition := false
			for _, v := range result.Violations {
				if v.Name == "condition" &&
					(strings.Contains(v.Location, "builders.go") ||
						strings.Contains(v.Location, "builders_composite.go")) {
					hasBoolCondition = true
					break
				}
			}
			Expect(hasBoolCondition).To(BeTrue(),
				"Should flag condition parameter as bool that should be enum")
		})
	})

	Describe("Other linters (no violations expected)", func() {
		It("should report no context violations", func() {
			cmd := exec.Command(bfBin, "context", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("No semantic context issues"),
				"Should have no semantic context issues")
		})

		It("should report no duplicate type violations", func() {
			cmd := exec.Command(bfBin, "dupe", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("No duplicates found"),
				"Should have no duplicate types")
		})

		It("should report no panic conditions", func() {
			cmd := exec.Command(bfBin, "panic", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("No panic conditions"),
				"Should have no panic conditions")
		})

		It("should report no strong-id violations", func() {
			cmd := exec.Command(bfBin, "strong-id", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("No string ID parameters"),
				"Should have no strong ID violations")
		})

		It("should report no boolblind violations", func() {
			cmd := exec.Command(bfBin, "boolblind", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("No boolean blindness"),
				"Should have no boolean blindness violations")
		})
	})

	Describe("Stats command", func() {
		It("should run stats successfully", func() {
			cmd := exec.Command(bfBin, "stats", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("Total Issues"),
				"Should show total issues in stats")
		})

		It("should report 15 total issues", func() {
			cmd := exec.Command(bfBin, "stats", modulePath)
			output, err := cmd.CombinedOutput()

			Expect(err).ToNot(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("Total Issues: 15"),
				"Should report 15 total issues (all PHANTOM)")
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
		if _, err := os.Stat(path); err == nil {
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
		if _, err := os.Stat(path); err == nil {
			return path
		}
		if _, err := exec.LookPath(path); err == nil {
			return path
		}
	}

	Fail("Could not find branching-flow binary")
	return ""
}
