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

**Commit:** `feat: initialize project with Go module, dependencies, and config system`

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

**Commit:** `feat: implement core models and filesystem operations for projects and worktrees`

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

**Commit:** `feat: implement plist and wlist commands`

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

**Commit:** `feat: implement pcreate and wcreate commands`

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

**Commit:** `feat: implement pcd and wcd commands`

---

## Phase 6: Remove Command

### Step 6.1: `wremove` Command (`cmd/wremove.go`)
- Accept `[project-name]` (optional) and `[worktree-name]` (optional)
- Use pickers as needed
- **Skip** main worktree (detect via `git worktree list`)
- Prompt for confirmation before removal
- Remove via `git worktree remove`

**Commit:** `feat: implement wremove command with confirmation`

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

**Commit:** `feat: add interactive fuzzy-search picker for project/worktree selection`

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

**Commit:** `feat: add shell completion and setup command`

---

## Phase 9: Testing
All tests should be ran in a containerized environment if they require filesystem ($HOME or projects dir) or git operations to ensure consistency and isolation.

### Step 9.1: Containerized Testing Environment
- Create `Dockerfile` for test environment
- Build container with Go, Git, and daycare binary

### Step 9.2: Unit Tests
- Test config loading
- Test project listing/scanning
- Test worktree operations (mock git commands)
- Test picker logic

### Step 9.3: Integration Tests
- Create temp directories for test projects
- Test full command execution
- Test worktree lifecycle (create, list, remove)
- Test picker integration with mock data

### Step 9.4: E2E Tests (Containerized)
- Test scenarios:
  1. Create project, verify `.daycare/` structure
  2. Create worktree, verify git worktree list output
  3. Remove worktree with confirmation
  4. Test `plist` and `wlist` output
  5. Test `pcd`/`wcd` shell script output
- Use `docker compose` for test orchestration
- Clean up containers after tests

**Commit:** `test: add unit, integration, and E2E tests`

---

## Phase 10: CI/CD

### Step 10.1: GitHub Actions Workflow
- **Lint/Format Check**: Run `gofmt`, `golint`, or `staticcheck`
- **Tests**: Run unit and integration tests on multiple Go versions
- **E2E Tests**: Build Docker image, run containerized E2E tests
- **Cross-compile**: Build binaries for Linux, macOS, Windows

**Commit:** `ci: add GitHub Actions workflow for lint, test, and E2E`

---

## Phase 11: Documentation

### Step 11.1: README.md
- Project description
- Quick start / installation instructions
- Basic usage examples for each command
- Shell setup instructions
- Configuration reference

### Step 11.2: AGENTS.md & CLAUDE.md
- Create `AGENTS.md` with:
  - Project overview
  - Key files and their purposes
  - Common tasks and how to accomplish them
  - Testing strategy
- Create `CLAUDE.md` as symlink to `AGENTS.md`

**Commit:** `docs: add README, AGENTS.md, and CLAUDE.md`

---

## Summary

| Phase | Description | Files Changed |
|-------|-------------|---------------|
| 1 | Project scaffolding | go.mod, main.go, config |
| 2 | Core models & operations | internal/* |
| 3 | List commands | cmd/plist.go, cmd/wlist.go |
| 4 | Create commands | cmd/pcreate.go, cmd/wcreate.go |
| 5 | CD commands | cmd/pcd.go, cmd/wcd.go |
| 6 | Remove command | cmd/wremove.go |
| 7 | TUI picker | internal/picker/, integrate |
| 8 | Shell integration | cmd/completion.go, cmd/setup.go |
| 9 | Testing | *_test.go, Dockerfile, docker-compose.yml |
| 10 | CI/CD | .github/workflows/*.yml |
| 11 | Documentation | README.md, AGENTS.md, CLAUDE.md |

Total: **11 commits** (can be merged or split further if needed)
