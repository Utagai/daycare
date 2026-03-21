# Daycare - AI Agent Documentation

## Project Overview

**Daycare** is a CLI tool for managing software projects and Git worktrees. It helps developers:
- Organize projects in a central directory
- Create, list, and switch between Git worktrees using an interactive fuzzy picker
- Seamlessly `cd` into projects and worktrees with shell integration
- Get tab completion for project and worktree names

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
# Build
make build          # Build binary

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
NOTE: DO NOT RUN E2E TESTS ON YOUR LOCAL MACHINE AS THEY WILL MODIFY YOUR FILESYSTEM AND GIT REPOSITORIES. USE THE PROVIDED DOCKER ENVIRONMENT.
- Run in Docker container
- Environment variables: `DAYCARE_BIN`, `DAYCARE_PROJECTS_DIR`

## Configuration

- Config file: `~/.config/daycare/config.yaml`
- Example: `config.example.yaml`
