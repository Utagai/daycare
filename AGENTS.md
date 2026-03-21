# Daycare - AI Agent Documentation

## Project Overview

**Daycare** is a CLI tool for managing software projects and Git worktrees. It helps developers:
- Organize projects in a central directory
- Create, list, and switch between Git worktrees using an interactive fuzzy picker
- Seamlessly `cd` into projects and worktrees with shell integration
- Get tab completion for project and worktree names

## Project Structure

```
agents-md/          # Root directory
├── cmd/            # CLI commands (Cobra)
│   ├── root.go     # Root command, version handling
│   ├── plist.go    # List projects
│   ├── pcreate.go  # Create project
│   ├── pcd.go      # cd into project
│   ├── wlist.go    # List worktrees
│   ├── wcreate.go  # Create worktree
│   ├── wcd.go      # cd into worktree
│   ├── wremove.go  # Remove worktree
│   ├── setup.go    # Shell integration
│   ├── picker_helpers.go  # Picker helpers
│   └── styles.go   # Terminal styling
├── internal/       # Business logic
│   ├── config/     # Configuration (YAML)
│   ├── models/      # Data models
│   ├── project/     # Project operations
│   ├── worktree/    # Worktree operations
│   └── picker/      # Interactive TUI picker
├── e2e/            # End-to-end tests
├── main.go         # Entry point
└── Makefile        # Build automation
```

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25+ |
| CLI Framework | Cobra |
| TUI/Picker | Bubble Tea |
| Styling | Lipgloss |
| Config | YAML |

## Commands

```bash
# Build and run
make build          # Build binary
make run ARGS="<command>"  # Run with args

# Testing
make test-unit      # Unit tests
make test-e2e       # E2E tests in Docker
make test           # All tests
```

## Code Conventions

### Dependency Injection
- Interfaces for `GitExecuter` and `FilesystemOps` enable testing
- Use `SetGitExecuter()`, `SetFilesystemOps()` with `ResetDeps()` for tests

### Data Models

```go
type Project struct {
    Name      string
    Path      string
    LastCommit time.Time
}

type Worktree struct {
    Name        string
    Branch      string
    Path        string
    ProjectName string
    IsMain      bool
    LastCommit  time.Time
}
```

### Terminal UI
- Uses Lipgloss for styling
- Color constants in `cmd/styles.go`

## Testing

### Unit Tests
- Mock interfaces for filesystem and Git operations
- Located in `*_test.go` files in each package

### E2E Tests
- Run in Docker container
- Environment variables: `DAYCARE_BIN`, `DAYCARE_PROJECTS_DIR`

## Configuration

- Config file: `~/.config/daycare/config.yaml`
- Example: `config.example.yaml`

## Key Files

| File | Purpose |
|------|---------|
| `cmd/root.go` | Version, config initialization |
| `internal/config/config.go` | Config loading |
| `internal/models/models.go` | Project/Worktree structs |
| `internal/project/project.go` | Project CRUD |
| `internal/worktree/worktree.go` | Worktree operations |
| `internal/picker/picker.go` | Fuzzy search TUI |
