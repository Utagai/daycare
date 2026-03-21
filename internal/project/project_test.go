package project

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/utagai/daycare/internal/config"
	"github.com/utagai/daycare/internal/models"
)

type mockDirEntry struct {
	name    string
	isDir   bool
	isValid bool
}

type mockFileInfo struct {
	name  string
	isDir bool
}

func (m mockDirEntry) Name() string      { return m.name }
func (m mockDirEntry) IsDir() bool       { return m.isDir }
func (m mockDirEntry) Type() os.FileMode { return 0 }
func (m mockDirEntry) Info() (os.FileInfo, error) {
	if !m.isValid {
		return nil, os.ErrInvalid
	}
	return mockFileInfo{name: m.name, isDir: m.isDir}, nil
}

type mockFilesystemOps struct {
	readDirResult []os.DirEntry
	readDirError  error
	statResult    os.FileInfo
	statError     error
	mkdirError    error
	readFileData  []byte
	readFileError error
}

func (f *mockFilesystemOps) ReadDir(path string) ([]os.DirEntry, error) {
	return f.readDirResult, f.readDirError
}

func (f *mockFilesystemOps) Stat(path string) (os.FileInfo, error) {
	return f.statResult, f.statError
}

func (f *mockFilesystemOps) MkdirAll(path string, perm os.FileMode) error {
	return f.mkdirError
}

func (f *mockFilesystemOps) Mkdir(path string, perm os.FileMode) error {
	return f.mkdirError
}

func (f *mockFilesystemOps) ReadFile(path string) ([]byte, error) {
	return f.readFileData, f.readFileError
}

func (f *mockFilesystemOps) WriteFile(path string, data []byte, perm os.FileMode) error {
	return nil
}

func (f *mockFilesystemOps) OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return nil, nil
}

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

func (m mockFileInfo) Name() string { return m.name }
func (m mockFileInfo) Size() int64  { return 0 }
func (m mockFileInfo) Mode() os.FileMode {
	if m.isDir {
		return os.ModeDir
	}
	return 0644
}
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }

func TestListProjectsEmpty(t *testing.T) {
	mockFs := &mockFilesystemOps{
		readDirResult: []os.DirEntry{},
		readDirError:  nil,
	}

	SetFilesystemOps(mockFs)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/test/projects",
	}

	projects, err := ListProjects(cfg)
	if err != nil {
		t.Fatalf("ListProjects() error = %v", err)
	}

	if len(projects) != 0 {
		t.Errorf("ListProjects() = %v projects, want 0", len(projects))
	}
}

func TestListProjectsWithGitRepos(t *testing.T) {
	mockFs := &mockFilesystemOps{
		readDirResult: []os.DirEntry{
			mockDirEntry{name: "project1", isDir: true, isValid: true},
			mockDirEntry{name: "project2", isDir: true, isValid: true},
			mockDirEntry{name: "not-a-git", isDir: true, isValid: true},
			mockDirEntry{name: "worktrees", isDir: true, isValid: true},
		},
		readDirError: nil,
		statResult:   mockFileInfo{name: ".git", isDir: true},
	}

	// Wrap the mock to handle stat calls for .git directory
	wrappedFs := &wrappedFilesystemOps{
		mock: mockFs,
	}

	SetFilesystemOps(wrappedFs)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/test/projects",
	}

	projects, err := ListProjects(cfg)
	if err != nil {
		t.Fatalf("ListProjects() error = %v", err)
	}

	// project1, project2 are git repos (have .git), not-a-git doesn't have .git stat result, worktrees is skipped
	if len(projects) != 2 {
		t.Errorf("ListProjects() = %v projects, want 2", len(projects))
	}

	// Check that worktrees was excluded
	for _, p := range projects {
		if p.Name == "worktrees" {
			t.Error("worktrees directory should be excluded")
		}
	}
}

type wrappedFilesystemOps struct {
	mock *mockFilesystemOps
}

func (f *wrappedFilesystemOps) ReadDir(path string) ([]os.DirEntry, error) {
	return f.mock.ReadDir(path)
}

func (f *wrappedFilesystemOps) Stat(path string) (os.FileInfo, error) {
	// Check if this is a .git path - return valid git dir for git repos
	// return error for not-a-git so it's not treated as git repo
	if path == "/test/projects/project1/.git" || path == "/test/projects/project2/.git" {
		return mockFileInfo{name: ".git", isDir: true}, nil
	}
	if path == "/test/projects/not-a-git/.git" {
		return nil, os.ErrNotExist
	}
	return f.mock.Stat(path)
}

func (f *wrappedFilesystemOps) MkdirAll(path string, perm os.FileMode) error {
	return f.mock.MkdirAll(path, perm)
}

func (f *wrappedFilesystemOps) Mkdir(path string, perm os.FileMode) error {
	return f.mock.Mkdir(path, perm)
}

func (f *wrappedFilesystemOps) ReadFile(path string) ([]byte, error) {
	return f.mock.ReadFile(path)
}

func (f *wrappedFilesystemOps) WriteFile(path string, data []byte, perm os.FileMode) error {
	return f.mock.WriteFile(path, data, perm)
}

func (f *wrappedFilesystemOps) OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return f.mock.OpenFile(path, flag, perm)
}

func TestListProjectsDirNotExist(t *testing.T) {
	mockFs := &mockFilesystemOps{
		readDirError: os.ErrNotExist,
	}

	SetFilesystemOps(mockFs)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/nonexistent/path",
	}

	_, err := ListProjects(cfg)
	if err == nil {
		t.Error("ListProjects() should return error when directory doesn't exist")
	}
}

func TestGetProject(t *testing.T) {
	mockFs := &mockFilesystemOps{
		statResult: mockFileInfo{name: "myproject", isDir: true},
		statError:  nil,
	}

	wrappedFs := &wrappedFilesystemOps{mock: mockFs}
	SetFilesystemOps(wrappedFs)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/test/projects",
	}

	project, err := GetProject(cfg, "myproject")
	if err != nil {
		t.Fatalf("GetProject() error = %v", err)
	}

	if project.Name != "myproject" {
		t.Errorf("GetProject().Name = %v, want myproject", project.Name)
	}

	expectedPath := filepath.Join("/test/projects", "myproject")
	if project.Path != expectedPath {
		t.Errorf("GetProject().Path = %v, want %v", project.Path, expectedPath)
	}
}

func TestGetProjectNotFound(t *testing.T) {
	mockFs := &mockFilesystemOps{
		statError: os.ErrNotExist,
	}

	SetFilesystemOps(mockFs)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/test/projects",
	}

	_, err := GetProject(cfg, "nonexistent")
	if err == nil {
		t.Error("GetProject() should return error when project doesn't exist")
	}
}

func TestFindProjectByPath(t *testing.T) {
	mockFs := &mockFilesystemOps{
		statResult: mockFileInfo{name: "myproject", isDir: true},
	}

	wrappedFs := &wrappedFilesystemOps{mock: mockFs}
	SetFilesystemOps(wrappedFs)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/test/projects",
	}

	project, err := FindProjectByPath(cfg, "/test/projects/myproject")
	if err != nil {
		t.Fatalf("FindProjectByPath() error = %v", err)
	}

	if project == nil {
		t.Fatal("FindProjectByPath() returned nil")
	}

	if project.Name != "myproject" {
		t.Errorf("FindProjectByPath().Name = %v, want myproject", project.Name)
	}
}

func TestFindProjectByPathOutsideProjectsDir(t *testing.T) {
	mockFs := &mockFilesystemOps{}
	SetFilesystemOps(mockFs)
	defer ResetDeps()

	cfg := &config.Config{
		ProjectsDir: "/test/projects",
	}

	project, err := FindProjectByPath(cfg, "/home/user")
	if err != nil {
		t.Fatalf("FindProjectByPath() error = %v", err)
	}

	if project != nil {
		t.Error("FindProjectByPath() should return nil for path outside projects dir")
	}
}

func TestGetProjectInfo(t *testing.T) {
	mockGit := &mockGitExecuter{
		output:    []byte("1700000000\n"),
		outputErr: nil,
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	p := &models.Project{
		Name: "testproject",
		Path: "/test/projects/testproject",
	}

	err := GetProjectInfo(p)
	if err != nil {
		t.Fatalf("GetProjectInfo() error = %v", err)
	}

	if p.LastCommit.IsZero() {
		t.Error("GetProjectInfo() should populate LastCommit")
	}
}

func TestGetProjectInfoEmptyRepo(t *testing.T) {
	mockGit := &mockGitExecuter{
		output:    []byte("\n"),
		outputErr: nil,
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	p := &models.Project{
		Name: "testproject",
		Path: "/test/projects/testproject",
	}

	err := GetProjectInfo(p)
	if err != nil {
		t.Fatalf("GetProjectInfo() error = %v", err)
	}

	// Should not panic, LastCommit should remain zero
}

func TestGetProjectInfoNoCommits(t *testing.T) {
	mockGit := &mockGitExecuter{
		outputErr: &mockError{},
	}

	SetGitExecuter(mockGit)
	defer ResetDeps()

	p := &models.Project{
		Name: "testproject",
		Path: "/test/projects/testproject",
	}

	err := GetProjectInfo(p)
	if err != nil {
		t.Fatalf("GetProjectInfo() error = %v", err)
	}

	// Should not panic, LastCommit should remain zero
}

type mockError struct{}

func (e *mockError) Error() string { return "mock error" }
