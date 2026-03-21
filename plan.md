# Daycare Implementation Plan

## Overview

This plan outlines the implementation of **Daycare**, a Go CLI tool for managing software projects and Git worktrees. The implementation follows the specification in `spec.md`.

## Project Structure

```
daycare/
├── cmd/
│   └── root.go
│   └── plist.go
│   └── pcreate.go
│   └── pcd.go
│   └── wlist.go
│   └── wcreate.go
│   └── wremove.go
│   └── wcd.go
│   └── setup.go
│   └── completion.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   └── models.go
│   ├── picker/
│   │   └── picker.go
│   ├── project/
│   │   └── project.go
│   └── worktree/
│       └── worktree.go
├── main.go
├── go.mod
└── ...
```

---

## Phase 1: Project Scaffolding & Configuration

### Step 1.1: Initialize Go Module
- Create `go.mod` with module name, Go version, and dependencies:
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/charmbracelet/bubbletea` - TUI
  - `gopkg.in/yaml.v3` - Config parsing
- Create `main.go` entry point

### Step 1.2: Configuration System
- Implement `internal/config/config.go`:
  - Load from `~/.config/daycare/config.yaml`
  - Default `ProjectsDir` to `~/projects`
  - Support custom projects directory path
- Create example config file in repo

### Step 1.3: Cobra Root Setup
- Implement `cmd/root.go`:
  - Global flags (`-h`, `-v`, `--version`)
  - Version command
  - Register all subcommands

---

## Phase 2: Core Models & Filesystem Operations

### Step 2.1: Data Models
- Define `Project` struct: Name, Path, LastCommit (time)
- Define `Worktree` struct: Name, Branch, Path, ProjectName

### Step 2.2: Project Operations (`internal/project/`)
- `ListProjects(config)` - Scan projects directory for Git repos
- `GetProject(name)` - Get single project by name
- `CreateProject(name)` - Create directory, init Git, create `.daycare/` subdir
- `GetProjectInfo(project)` - Get last commit time via `git log -1`

### Step 2.3: Worktree Operations (`internal/worktree/`)
- `ListWorktrees(project)` - Parse `git worktree list` output
- `CreateWorktree(project, name)` - Run `git worktree add`
- `RemoveWorktree(project, name)` - Run `git worktree remove` (with confirmation logic in command)
- Validate worktree doesn't already exist

---

## Phase 3: List Commands

### Step 3.1: `plist` Command (`cmd/plist.go`)
- Execute project listing
- Output as ASCII table (name, path, last commit time)
- Use `tablewriter` or similar for formatting

### Step 3.2: `wlist` Command (`cmd/wlist.go`)
- Accept optional `[project-name]` argument
- If omitted, use picker to select project
- Output worktrees as table (name, branch, path)

---

## Phase 4: Create Commands

### Step 4.1: `pcreate` Command (`cmd/pcreate.go`)
- Accept optional `[project-name]` argument
- If omitted, prompt for name
- Create project via `internal/project`
- Initialize `.daycare/pwcd` and `.daycare/wcreate` as empty scripts
- Add `.daycare/` to `.git/info/exclude`
- Create initial commit

### Step 4.2: `wcreate` Command (`cmd/wcreate.go`)
- Accept `[project-name]` (optional) and `<name>` (required)
- Use picker for project if not provided
- Create worktree via `internal/worktree`
- Execute `.daycare/wcreate` script after creation
- Error if worktree name already exists

---

## Phase 5: CD Commands

### Step 5.1: `pcd` Command (`cmd/pcd.go`)
- Accept optional `[project-name]`
- Use picker if omitted
- Output shell script to stdout:
  ```sh
  # Source pwcd if exists
  # cd to project path
  ```
- Designed for `eval "$(daycare pcd)"`

### Step 5.2: `wcd` Command (`cmd/wcd.go`)
- Accept `[project-name]` (optional) and `[worktree-name]` (optional)
- Use pickers as needed
- Output similar shell script for worktree path
- Execute `.daycare/pwcd` before cd

---

## Phase 6: Remove Command

### Step 6.1: `wremove` Command (`cmd/wremove.go`)
- Accept `[project-name]` (optional) and `[worktree-name]` (optional)
- Use pickers as needed
- **Skip** main worktree (detect via `git worktree list`)
- Prompt for confirmation before removal
- Remove via `git worktree remove`

---

## Phase 7: TUI Picker Integration

### Step 7.1: Fuzzy Picker (`internal/picker/picker.go`)
- Create reusable Bubble Tea model for fuzzy search
- Accept list of items with names
- Return selected item
- Handle ESC/cancel gracefully
- Support searching by project name

### Step 7.2: Integrate Pickers
- Replace picker calls in all commands that need them
- Ensure picker works in terminal with proper key handling

---

## Phase 8: Shell Integration

### Step 8.1: Completion Command (`cmd/completion.go`)
- Use Cobra's built-in shell completion generation
- Register completion command for zsh
- Add dynamic completions for project/worktree names

### Step 8.2: Setup Command (`cmd/setup.go`)
- Print out wrapper function to `~/.zshrc` that user can append to their ~/.zshrc.
- Wrapper intercepts `pcd` and `wcd` to auto-eval their output
- Preserve existing `daycare` command for other invocations

---

## Phase 9: Unit Testing

Introduce/use mocks/interfaces when needed to isolate filesystem and git operations for testing.

NOTE: The program has interactive elements. We may need to implement flags
to bypass the interactive picker for testing purposes (e.g. -y for wremove to
skip confirmation).

### Step 9.1: Unit Tests
- Test config loading
- Test project listing/scanning
- Test worktree operations (mock git commands -- may need to introduce interfaces for easier testing)
- Test picker logic

## Phase 10: Integration/E2E Testing

All tests should be ran in a containerized environment if they require
filesystem ($HOME or projects dir) or git operations to ensure consistency and
isolation.

NOTE: The program has interactive elements. We may need to implement flags
to bypass the interactive picker for testing purposes (e.g. -y for wremove to
skip confirmation).

### Step 10.1: Containerized Testing Environment
- Create `Dockerfile` for test environment
- Build container with Go, Git, and daycare binary

### Step 10.2: Integration/E2E Tests
- Test scenarios:
  1. Create project, verify `.daycare/` structure
  2. Create worktree, verify git worktree list output and .daycare/wcreate execution.
  3. Remove worktree with confirmation
  4. Test `plist` and `wlist` output
  5. Test `pcd`/`wcd` shell script output
  6. Test shell completion generation
  7. Test setup command output
- Clean up containers after tests

---

## Phase 11: CI/CD

### Step 11.1: GitHub Actions Workflow
- **Lint/Format Check**: Run `gofmt`, `golint`, or `staticcheck`
- **Tests**: Run unit and integration tests on multiple Go versions
- **E2E Tests**: Build Docker image, run containerized E2E tests
- **Cross-compile**: Build binaries for Linux, macOS, Windows

---

## Phase 12: Documentation

### Step 12.1: README.md
- Project description
- Quick start / installation instructions
- Basic usage examples for each command
- Shell setup instructions
- Configuration reference

### Step 12.2: AGENTS.md & CLAUDE.md
- Create `AGENTS.md` with:
  - Project overview
  - Key files and their purposes
  - Common tasks and how to accomplish them
  - Testing strategy
- Create `CLAUDE.md` as symlink to `AGENTS.md`

