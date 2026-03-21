package worktree

import (
	"os"
	"testing"

	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
)

type mockGitExecuter struct {
	runError  error
	output    []byte
	outputErr error
}

func (g *mockGitExecuter) Run(dir string, name string, args ...string) error {
	return g.runError
}

func (g *mockGitExecuter) Output(dir string, name string, args ...string) ([]byte, error) {
	return g.output, g.outputErr
}

type mockFilesystemOps struct {
	mkdirAllError error
}

func (f *mockFilesystemOps) MkdirAll(path string, perm os.FileMode) error {
	return f.mkdirAllError
}

func TestParseWorktreePorcelain(t *testing.T) {
	input := `worktree /projects/myproject
HEAD abc123def456
branch refs/heads/main

worktree /projects/myproject/worktrees/feature1
HEAD def456ghi789
branch refs/heads/feature1

worktree /projects/myproject/worktrees/feature2
HEAD ghi789jkl012
detached
`

	worktrees := parseWorktreePorcelain("myproject", []byte(input))

	if len(worktrees) != 3 {
		t.Fatalf("parseWorktreePorcelain() returned %d worktrees, want 3", len(worktrees))
	}

	// First worktree should be main
	if !worktrees[0].IsMain {
		t.Error("First worktree should be main")
	}
	if worktrees[0].Branch != "main" {
		t.Errorf("First worktree branch = %v, want main", worktrees[0].Branch)
	}

	// Second worktree
	if worktrees[1].Name != "feature1" {
		t.Errorf("Second worktree name = %v, want feature1", worktrees[1].Name)
	}
	if worktrees[1].Branch != "feature1" {
		t.Errorf("Second worktree branch = %v, want feature1", worktrees[1].Branch)
	}

	// Third worktree is detached
	if worktrees[2].Branch != "(detached)" {
		t.Errorf("Third worktree branch = %v, want (detached)", worktrees[2].Branch)
	}
}

func TestParseWorktreePorcelainEmpty(t *testing.T) {
	worktrees := parseWorktreePorcelain("myproject", []byte(""))

	if len(worktrees) != 0 {
		t.Errorf("parseWorktreePorcelain() returned %d worktrees, want 0", len(worktrees))
	}
}

func TestParseWorktreePorcelainSingleMain(t *testing.T) {
	input := `worktree /projects/myproject
HEAD abc123def456
branch refs/heads/main
`

	worktrees := parseWorktreePorcelain("myproject", []byte(input))

	if len(worktrees) != 1 {
		t.Fatalf("parseWorktreePorcelain() returned %d worktrees, want 1", len(worktrees))
	}

	if !worktrees[0].IsMain {
		t.Error("Worktree should be main")
	}
}

func TestListWorktrees(t *testing.T) {
	mockGit := &mockGitExecuter{
		output: []byte(`worktree /projects/myproject
HEAD abc123def456
branch refs/heads/main
`),
		outputErr: nil,
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/projects",
	}
	p := &models.Project{
		Name: "myproject",
		Path: "/projects/myproject",
	}

	worktrees, err := ListWorktrees(cfg, p)
	if err != nil {
		t.Fatalf("ListWorktrees() error = %v", err)
	}

	if len(worktrees) != 1 {
		t.Fatalf("ListWorktrees() returned %d worktrees, want 1", len(worktrees))
	}
}

func TestListWorktreesGitError(t *testing.T) {
	mockGit := &mockGitExecuter{
		outputErr: &mockError{},
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/projects",
	}
	p := &models.Project{
		Name: "myproject",
		Path: "/projects/myproject",
	}

	_, err := ListWorktrees(cfg, p)
	if err == nil {
		t.Error("ListWorktrees() should return error when git command fails")
	}
}

func TestGetWorktreeInfo(t *testing.T) {
	mockGit := &mockGitExecuter{
		output:    []byte("1700000000\n"),
		outputErr: nil,
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	wt := &models.Worktree{
		Name: "feature1",
		Path: "/projects/myproject/worktrees/feature1",
	}

	err := GetWorktreeInfo(wt)
	if err != nil {
		t.Fatalf("GetWorktreeInfo() error = %v", err)
	}

	if wt.LastCommit.IsZero() {
		t.Error("GetWorktreeInfo() should populate LastCommit")
	}
}

func TestGetWorktreeInfoEmpty(t *testing.T) {
	mockGit := &mockGitExecuter{
		output:    []byte("\n"),
		outputErr: nil,
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	wt := &models.Worktree{
		Name: "feature1",
		Path: "/projects/myproject/worktrees/feature1",
	}

	err := GetWorktreeInfo(wt)
	if err != nil {
		t.Fatalf("GetWorktreeInfo() error = %v", err)
	}

	// Should not panic, LastCommit should remain zero
}

func TestCreateWorktreeSuccess(t *testing.T) {
	callCount := 0
	mockGit := &mockGitExecuter{
		output: []byte(`worktree /projects/myproject
HEAD abc123def456
branch refs/heads/main

worktree /projects/myproject/worktrees/feature1
HEAD def456ghi789
branch refs/heads/feature1
`),
		outputErr: nil,
		runError:  nil,
	}

	// Wrap to track calls
	origGit := mockGit
	mockGitWrapper := &mockGitExecuterWrapper{
		mockGit: origGit,
		onRun: func(dir string, name string, args ...string) error {
			callCount++
			return origGit.Run(dir, name, args...)
		},
		onOutput: func(dir string, name string, args ...string) ([]byte, error) {
			return origGit.Output(dir, name, args...)
		},
	}

	SetGitExecuter(mockGitWrapper)
	defer ResetDeps()

	mockFs := &mockFilesystemOps{}
	SetFilesystemOps(mockFs)

	cfg := &config.Config{
		ProjectsDir: "/projects",
	}
	p := &models.Project{
		Name: "myproject",
		Path: "/projects/myproject",
	}

	wt, err := CreateWorktree(cfg, p, "newbranch")
	if err != nil {
		t.Fatalf("CreateWorktree() error = %v", err)
	}

	if wt.Name != "newbranch" {
		t.Errorf("CreateWorktree().Name = %v, want newbranch", wt.Name)
	}

	if callCount == 0 {
		t.Error("CreateWorktree should have called git worktree add")
	}
}

func TestCreateWorktreeAlreadyExists(t *testing.T) {
	mockGit := &mockGitExecuter{
		output: []byte(`worktree /projects/myproject
HEAD abc123def456
branch refs/heads/main

worktree /projects/myproject/worktrees/feature1
HEAD def456ghi789
branch refs/heads/feature1
`),
		outputErr: nil,
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	mockFs := &mockFilesystemOps{}
	SetFilesystemOps(mockFs)

	cfg := &config.Config{
		ProjectsDir: "/projects",
	}
	p := &models.Project{
		Name: "myproject",
		Path: "/projects/myproject",
	}

	_, err := CreateWorktree(cfg, p, "feature1")
	if err == nil {
		t.Error("CreateWorktree() should return error when worktree already exists")
	}
}

func TestRemoveWorktreeNotFound(t *testing.T) {
	mockGit := &mockGitExecuter{
		output: []byte(`worktree /projects/myproject
HEAD abc123def456
branch refs/heads/main
`),
		outputErr: nil,
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/projects",
	}
	p := &models.Project{
		Name: "myproject",
		Path: "/projects/myproject",
	}

	err := RemoveWorktree(cfg, p, "nonexistent")
	if err == nil {
		t.Error("RemoveWorktree() should return error when worktree not found")
	}
}

func TestRemoveWorktreeMain(t *testing.T) {
	mockGit := &mockGitExecuter{
		output: []byte(`worktree /projects/myproject
HEAD abc123def456
branch refs/heads/main
`),
		outputErr: nil,
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/projects",
	}
	p := &models.Project{
		Name: "myproject",
		Path: "/projects/myproject",
	}

	err := RemoveWorktree(cfg, p, "myproject")
	if err == nil {
		t.Error("RemoveWorktree() should return error when trying to remove main worktree")
	}
}

func TestRemoveWorktreeSuccess(t *testing.T) {
	mockGit := &mockGitExecuter{
		output: []byte(`worktree /projects/myproject
HEAD abc123def456
branch refs/heads/main

worktree /projects/myproject/worktrees/feature1
HEAD def456ghi789
branch refs/heads/feature1
`),
		outputErr: nil,
		runError:  nil,
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/projects",
	}
	p := &models.Project{
		Name: "myproject",
		Path: "/projects/myproject",
	}

	err := RemoveWorktree(cfg, p, "feature1")
	if err != nil {
		t.Fatalf("RemoveWorktree() error = %v", err)
	}
}

type mockError struct{}

func (e *mockError) Error() string { return "mock error" }

type mockGitExecuterWrapper struct {
	mockGit  *mockGitExecuter
	onRun    func(dir string, name string, args ...string) error
	onOutput func(dir string, name string, args ...string) ([]byte, error)
}

func (g *mockGitExecuterWrapper) Run(dir string, name string, args ...string) error {
	if g.onRun != nil {
		return g.onRun(dir, name, args...)
	}
	return g.mockGit.Run(dir, name, args...)
}

func (g *mockGitExecuterWrapper) Output(dir string, name string, args ...string) ([]byte, error) {
	if g.onOutput != nil {
		return g.onOutput(dir, name, args...)
	}
	return g.mockGit.Output(dir, name, args...)
}
