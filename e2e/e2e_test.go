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

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return stdout.String(), nil
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func containsString(s, substr string) bool {
	return strings.Contains(stripANSI(s), substr)
}

func TestProjectCreate(t *testing.T) {
	out, err := runCmd(daycareBin, "pcreate", projectName)
	if err != nil {
		t.Fatalf("pcreate failed: %v, stderr: %s", err, string(out))
	}

	if !containsString(out, projectName) {
		t.Errorf("pcreate output should contain project name %q, got: %s", projectName, out)
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

	out, err := runCmd(daycareBin, "plist")
	if err != nil {
		t.Fatalf("plist failed: %v", err)
	}

	if !containsString(out, projectName) {
		t.Errorf("plist output should contain project name %q, got: %s", projectName, out)
	}
}

func TestProjectCD(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)

	out, err := runCmd(daycareBin, "pcd", projectName)
	if err != nil {
		t.Fatalf("pcd failed: %v", err)
	}

	projectPath := filepath.Join(projectsDir, projectName)
	if !containsString(out, projectPath) {
		t.Errorf("pcd output should contain project path %q, got: %s", projectPath, out)
	}

	if !containsString(out, "cd ") {
		t.Errorf("pcd output should contain a cd command, got: %s", out)
	}
}

func TestWorktreeCreate(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)

	out, err := runCmd(daycareBin, "wcreate", projectName, worktreeName)
	if err != nil {
		t.Fatalf("wcreate failed: %v", err)
	}

	if !containsString(out, worktreeName) {
		t.Errorf("wcreate output should contain worktree name %q, got: %s", worktreeName, out)
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

	out, err := runCmd(daycareBin, "wlist", projectName)
	if err != nil {
		t.Fatalf("wlist failed: %v", err)
	}

	if !containsString(out, projectName) {
		t.Errorf("wlist output should contain project name %q, got: %s", projectName, out)
	}

	if !containsString(out, worktreeName) {
		t.Errorf("wlist output should contain worktree name %q, got: %s", worktreeName, out)
	}
}

func TestWorktreeCD(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)
	runCmd(daycareBin, "wcreate", projectName, worktreeName)

	out, err := runCmd(daycareBin, "wcd", projectName, worktreeName)
	if err != nil {
		t.Fatalf("wcd failed: %v", err)
	}

	worktreePath := filepath.Join(projectsDir, "worktrees", projectName, worktreeName)
	if !containsString(out, worktreePath) {
		t.Errorf("wcd output should contain worktree path %q, got: %s", worktreePath, out)
	}

	if !containsString(out, "cd ") {
		t.Errorf("wcd output should contain a cd command, got: %s", out)
	}
}

func TestCompletionZsh(t *testing.T) {
	out, err := runCmd(daycareBin, "completion", "zsh")
	if err != nil {
		t.Fatalf("completion zsh failed: %v", err)
	}

	if !containsString(out, "#compdef daycare") {
		t.Errorf("zsh completion should contain compdef directive, got: %s", out)
	}
}

func TestSetup(t *testing.T) {
	out, err := runCmd(daycareBin, "setup")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if !containsString(out, "zsh") && !containsString(out, "bash") {
		t.Errorf("setup output should mention shell setup instructions, got: %s", out)
	}
}

func TestProjectCreateDuplicate(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)

	_, err := runCmd(daycareBin, "pcreate", projectName)
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

	out, err := runCmd(daycareBin, "wcd", scriptTestProject, "test-branch")
	if err != nil {
		t.Fatalf("wcd failed: %v", err)
	}

	expectedSource := fmt.Sprintf(". %q", pwcdScript)
	if !containsString(out, expectedSource) {
		t.Errorf("wcd output should contain the pwcd sourcing command %q, got: %s", expectedSource, out)
	}

	os.RemoveAll(filepath.Join(projectsDir, scriptTestProject))
}

func TestWorktreeCreateDuplicate(t *testing.T) {
	runCmd(daycareBin, "pcreate", projectName)
	runCmd(daycareBin, "wcreate", projectName, worktreeName)

	_, err := runCmd(daycareBin, "wcreate", projectName, worktreeName)
	if err == nil {
		t.Errorf("wcreate should fail when creating duplicate worktree")
	}
}
