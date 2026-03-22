//go:build e2e

package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	daycareBin  string
	projectsDir string
)

func init() {
	daycareBin = os.Getenv("DAYCARE_BIN")
	if daycareBin == "" {
		daycareBin = "/usr/local/bin/daycare"
	}

	projectsDir = os.Getenv("DAYCARE_PROJECTS_DIR")
	if projectsDir == "" {
		projectsDir = "/root/projects"
	}
}

const (
	projectName  = "e2e-test-project"
	worktreeName = "feature-branch"
)

func TestMain(m *testing.M) {
	if _, err := os.Stat(daycareBin); os.IsNotExist(err) {
		panic("daycare binary not found at " + daycareBin)
	}

	cleanupTestEnv()
	os.Exit(m.Run())
}

func cleanupTestEnv() {
	os.RemoveAll(projectsDir)
	os.MkdirAll(projectsDir, 0755)
	configContent := "projects_dir: " + projectsDir + "\n"
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "daycare", "config.yaml")
	os.MkdirAll(filepath.Dir(configPath), 0755)
	os.WriteFile(configPath, []byte(configContent), 0644)
}

func runCmd(name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func containsString(s, substr string) bool {
	return strings.Contains(stripANSI(s), substr)
}

func TestProjectCreate(t *testing.T) {
	stdout, stderr, err := runCmd(daycareBin, "pcreate", projectName)
	if err != nil {
		t.Fatalf("pcreate failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	if !containsString(stdout, projectName) {
		t.Errorf("pcreate output should contain project name %q, got: %s", projectName, stdout)
	}

	projectPath := filepath.Join(projectsDir, projectName)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Errorf("project directory should exist at %s", projectPath)
	}

	gitDir := filepath.Join(projectPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf("git repository should exist at %s", gitDir)
	}

	daycareDir := filepath.Join(projectPath, ".daycare")
	if _, err := os.Stat(daycareDir); os.IsNotExist(err) {
		t.Errorf(".daycare directory should exist at %s", daycareDir)
	}

	pwcdPath := filepath.Join(daycareDir, "pwcd")
	if _, err := os.Stat(pwcdPath); os.IsNotExist(err) {
		t.Errorf("pwcd script should exist at %s", pwcdPath)
	}

	wcreatePath := filepath.Join(daycareDir, "wcreate")
	if _, err := os.Stat(wcreatePath); os.IsNotExist(err) {
		t.Errorf("wcreate script should exist at %s", wcreatePath)
	}
}

func TestProjectList(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)

	stdout, stderr, err := runCmd(daycareBin, "plist")
	if err != nil {
		t.Fatalf("plist failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	if !containsString(stdout, projectName) {
		t.Errorf("plist output should contain project name %q, got: %s", projectName, stdout)
	}
}

func TestProjectCD(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)

	stdout, stderr, err := runCmd(daycareBin, "pcd", projectName)
	if err != nil {
		t.Fatalf("pcd failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	projectPath := filepath.Join(projectsDir, projectName)
	if !containsString(stdout, projectPath) {
		t.Errorf("pcd output should contain project path %q, got: %s", projectPath, stdout)
	}

	if !containsString(stdout, "cd ") {
		t.Errorf("pcd output should contain a cd command, got: %s", stdout)
	}
}

func TestWorktreeCreate(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)

	stdout, stderr, err := runCmd(daycareBin, "wcreate", projectName, worktreeName)
	if err != nil {
		t.Fatalf("wcreate failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	if !containsString(stdout, worktreeName) {
		t.Errorf("wcreate output should contain worktree name %q, got: %s", worktreeName, stdout)
	}

	worktreePath := filepath.Join(projectsDir, "worktrees", projectName, worktreeName)
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Errorf("worktree directory should exist at %s", worktreePath)
	}

	gitDir := filepath.Join(worktreePath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf("worktree should be a git repository at %s", gitDir)
	}
}

func TestWorktreeList(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)
	runCmd(daycareBin, "wcreate", projectName, worktreeName)

	stdout, stderr, err := runCmd(daycareBin, "wlist", projectName)
	if err != nil {
		t.Fatalf("wlist failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	if !containsString(stdout, projectName) {
		t.Errorf("wlist output should contain project name %q, got: %s", projectName, stdout)
	}

	if !containsString(stdout, worktreeName) {
		t.Errorf("wlist output should contain worktree name %q, got: %s", worktreeName, stdout)
	}
}

func TestWorktreeCD(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)
	runCmd(daycareBin, "wcreate", projectName, worktreeName)

	stdout, stderr, err := runCmd(daycareBin, "wcd", projectName, worktreeName)
	if err != nil {
		t.Fatalf("wcd failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	worktreePath := filepath.Join(projectsDir, "worktrees", projectName, worktreeName)
	if !containsString(stdout, worktreePath) {
		t.Errorf("wcd output should contain worktree path %q, got: %s", worktreePath, stdout)
	}

	if !containsString(stdout, "cd ") {
		t.Errorf("wcd output should contain a cd command, got: %s", stdout)
	}
}

func TestCompletionZsh(t *testing.T) {
	stdout, stderr, err := runCmd(daycareBin, "completion", "zsh")
	if err != nil {
		t.Fatalf("completion zsh failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	if !containsString(stdout, "#compdef daycare") {
		t.Errorf("zsh completion should contain compdef directive, got: %s", stdout)
	}
}

func TestSetup(t *testing.T) {
	stdout, stderr, err := runCmd(daycareBin, "setup")
	if err != nil {
		t.Fatalf("setup failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	expectedPatterns := []string{
		"Daycare shell integration",
		"daycare() {",
		"pcd|wcd",
		"command daycare",
		`eval "$output"`,
	}

	for _, pattern := range expectedPatterns {
		if !containsString(stdout, pattern) {
			t.Errorf("setup output should contain %q, got: %s", pattern, stdout)
		}
	}
}

func TestProjectCreateDuplicate(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)

	_, _, err := runCmd(daycareBin, "pcreate", projectName)
	if err == nil {
		t.Errorf("pcreate should fail when creating duplicate project")
	}
}

func TestWcreateScriptExecution(t *testing.T) {
	scriptTestProject := projectName + "-wcreate-script"
	markerPath := filepath.Join(projectsDir, scriptTestProject, "marker.txt")

	runCmd(daycareBin, "pcreate", scriptTestProject)

	wcreateScript := filepath.Join(projectsDir, scriptTestProject, ".daycare", "wcreate")
	scriptContent := fmt.Sprintf("#!/bin/bash\necho \"wcreate ran\" > %q\n", markerPath)
	if err := os.WriteFile(wcreateScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to write wcreate script: %v", err)
	}

	runCmd(daycareBin, "wcreate", scriptTestProject, "test-branch")

	if _, err := os.Stat(markerPath); os.IsNotExist(err) {
		t.Errorf(".daycare/wcreate script should have been executed, marker file not found at %s", markerPath)
	}

	os.RemoveAll(filepath.Join(projectsDir, scriptTestProject))
}

func TestWcreateScriptEnv(t *testing.T) {
	scriptTestProject := projectName + "-wcreate-env"
	runCmd(daycareBin, "pcreate", scriptTestProject)

	wcreateScript := filepath.Join(projectsDir, scriptTestProject, ".daycare", "wcreate")
	scriptContent := `#!/bin/bash
echo "PROJECT_NAME=$PROJECT_NAME"
echo "PROJECT_DIR=$PROJECT_DIR"
echo "WORKTREE_DIR=$WORKTREE_DIR"
`
	if err := os.WriteFile(wcreateScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to write wcreate script: %v", err)
	}

	stdout, stderr, err := runCmd(daycareBin, "wcreate", scriptTestProject, "test-branch")
	if err != nil {
		t.Fatalf("wcreate failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	expectedVars := []string{
		fmt.Sprintf("PROJECT_NAME=%s", scriptTestProject),
		fmt.Sprintf("PROJECT_DIR=%s", filepath.Join(projectsDir, scriptTestProject)),
		fmt.Sprintf("WORKTREE_DIR=%s", filepath.Join(projectsDir, "worktrees", scriptTestProject, "test-branch")),
	}

	for _, expected := range expectedVars {
		if !containsString(stdout, expected) {
			t.Errorf("wcreate script output should contain %q, got: %s", expected, stdout)
		}
	}

	os.RemoveAll(filepath.Join(projectsDir, scriptTestProject))
}

func TestWcdScriptOutput(t *testing.T) {
	scriptTestProject := projectName + "-wcd-script"

	runCmd(daycareBin, "pcreate", scriptTestProject)

	pwcdScript := filepath.Join(projectsDir, scriptTestProject, ".daycare", "pwcd")
	pwcdContent := `#!/bin/bash
echo "custom pwcd setup"
`
	if err := os.WriteFile(pwcdScript, []byte(pwcdContent), 0755); err != nil {
		t.Fatalf("failed to write pwcd script: %v", err)
	}

	runCmd(daycareBin, "wcreate", scriptTestProject, "test-branch")

	stdout, stderr, err := runCmd(daycareBin, "wcd", scriptTestProject, "test-branch")
	if err != nil {
		t.Fatalf("wcd failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	expectedSource := fmt.Sprintf(". %q", pwcdScript)
	if !containsString(stdout, expectedSource) {
		t.Errorf("wcd output should contain the pwcd sourcing command %q, got: %s", expectedSource, stdout)
	}

	os.RemoveAll(filepath.Join(projectsDir, scriptTestProject))
}

func TestWorktreeCreateDuplicate(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)
	runCmd(daycareBin, "wcreate", projectName, worktreeName)

	_, _, err := runCmd(daycareBin, "wcreate", projectName, worktreeName)
	if err == nil {
		t.Errorf("wcreate should fail when creating duplicate worktree")
	}
}

func TestWorktreeRemoveDeletesBranch(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)
	runCmd(daycareBin, "wcreate", projectName, worktreeName)

	out, stderr, err := runCmd(daycareBin, "wremove", projectName, worktreeName, "-y")
	if err != nil {
		t.Fatalf("wremove failed: %v, stdout: %s, stderr: %s", err, out, stderr)
	}

	if !containsString(out, "Removed worktree") {
		t.Errorf("wremove output should contain 'Removed worktree', got: %s", out)
	}

	worktreePath := filepath.Join(projectsDir, "worktrees", projectName, worktreeName)
	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		t.Errorf("worktree directory should be removed at %s", worktreePath)
	}

	projectPath := filepath.Join(projectsDir, projectName)
	branchCheck := exec.Command("git", "-C", projectPath, "branch", "--list", worktreeName)
	var branchOut bytes.Buffer
	branchCheck.Stdout = &branchOut
	if err := branchCheck.Run(); err != nil {
		t.Fatalf("git branch check failed: %v", err)
	}
	if branchOut.String() != "" {
		t.Errorf("branch %q should be deleted, but still exists", worktreeName)
	}
}
