# Daycare

## Overview
A simple program that tracks software projects found under a common directory,
and their worktrees as the core elements. It provides a CLI interface to list,
create, remove, and `cd` into these elements. The guiding principles of this
program are simplicity, efficiency and aesthetics.

## Tech Stack
- **Language/Runtime:** Go
- **Dependencies:** Charmbracelet, Cobra
- **Target Platform:** Unix-like systems (Linux, macOS, WSL)

## Functionality

### Core Features
- List projects and worktrees
- Create new projects & worktrees
- Remove worktrees
    - Note, this project does NOT remove projects.
- `cd` into projects and worktrees
- Display project and worktree information (e.g., name, path, last modified)
- Interactive fuzzy-search picker (using Charmbracelet Bubble Tea) for project/worktree selection when arguments are omitted

### Data Handling
This project does not maintain its own database. Instead, it relies on the filesystem
and Git to store all state. Projects are directories within the configured common
directory, and worktrees are Git worktrees created within a `worktrees/` subdirectory
at the same level as projects (`worktrees/<project-name>/<worktree-name>`). This means
if data is deleted or modified externally, daycare will reflect those changes on next execution.

```
<common-dir>/
├── my-project/           # Main project (the main worktree)
└── worktrees/
    └── my-project/
        └── feature/      # Worktree for feature branch
```

## CLI Interface
```
# List projects.
daycare plist
# List worktrees (opens picker if project omitted).
daycare wlist [project-name]
# Create a new project.
daycare pcreate
# Create a new worktree (opens picker if project omitted).
daycare wcreate [project-name] <name>
# Remove a worktree (opens picker(s) if args omitted).
daycare wremove [project-name] [worktree-name]  # Prompts for confirmation
# Change directory into a project.
daycare pcd [project-name]
# Change directory into a worktree (opens picker(s) if args omitted).
daycare wcd [project-name] [worktree-name]
# Setup shell integration (currently zsh only).
daycare setup
```

### Global Options
| Flag | Description | Default |
|------|-------------|---------|
| -h, --help | Show help | false |
| -v, --version | Show version | false |

### Commands

#### `plist`
Lists all projects found under the common directory. Output is displayed as a
table with columns for name, path, and time since last commit.

**Examples:**
```
daycare plist
```

#### `wlist`
Lists worktrees for a project. Output is displayed as a table with columns for
name, branch, and path. If no project name is provided, opens an interactive
fuzzy-search picker to select the project.

**Arguments:**
| Argument | Description | Required |
|----------|-------------|----------|
| `[project-name]` | Project to list worktrees for | No (opens picker if omitted) |

**Examples:**
```
daycare wlist                  # Opens picker to select project
daycare wlist my-project       # List worktrees for a specific project
```

#### `pcreate`
Creates a new project. Creates the project directory, initializes a Git repository,
creates a `.daycare/` directory with `pwcd` and `wcreate` script files, and adds
`.daycare/` to `.git/info/exclude` for local ignore. Prompts for the project name
if not provided.

**Arguments:**
| Argument | Description | Required |
|----------|-------------|----------|
| `[project-name]` | Name of the project to create | No |

**Examples:**
```
daycare pcreate
daycare pcreate my-new-project
```

#### `wcreate`
Creates a new worktree within an existing project. The worktree is created at
`worktrees/<project-name>/<name>` (at the same level as the project, not inside it)
and uses `<name>` as both the worktree directory and branch name. If no project name
is provided, opens an interactive fuzzy-search picker to select the project. Errors
if a worktree with the given name already exists. After creating the worktree,
executes the project's `.daycare/wcreate` script (if present).

**Arguments:**
| Argument | Description | Required |
|----------|-------------|----------|
| `[project-name]` | Project to add worktree to | No (opens picker if omitted) |
| `<name>` | Name for both the worktree and branch | Yes |

**Examples:**
```
daycare wcreate feature/new-feature             # Opens picker to select project
daycare wcreate my-project feature/new-feature
```

#### `wremove`
Removes a worktree from a project. If no project name is provided, opens an
interactive fuzzy-search picker to select the project. Prompts for confirmation
before removal. The main project directory (the main worktree) cannot be removed.

**Arguments:**
| Argument | Description | Required |
|----------|-------------|----------|
| `[project-name]` | Project containing the worktree | No (opens picker if omitted) |
| `[worktree-name]` | Worktree to remove | No (opens picker if omitted) |

**Examples:**
```
daycare wremove                             # Opens pickers for project then worktree
daycare wremove my-project                  # Opens picker for worktree
daycare wremove my-project feature-branch
```

#### `pcd`
Outputs an eval-able script that executes the project's `.daycare/pwcd` script
(if present) and changes to the project directory. If no project name is provided,
opens an interactive fuzzy-search picker to select the project.

**Arguments:**
| Argument | Description | Required |
|----------|-------------|----------|
| `[project-name]` | Project name | No (opens picker if omitted) |

**Examples:**
```
eval "$(daycare pcd)"                       # Opens picker to select project
eval "$(daycare pcd my-project)"
```

#### `wcd`
Outputs an eval-able script that executes the project's `.daycare/pwcd` script
(if present) and changes to the worktree directory. If no project name is provided,
opens an interactive fuzzy-search picker to select the project.

**Arguments:**
| Argument | Description | Required |
|----------|-------------|----------|
| `[project-name]` | Project name | No (opens picker if omitted) |
| `[worktree-name]` | Worktree name | No (opens picker if omitted) |

**Examples:**
```
eval "$(daycare wcd)"                       # Opens pickers for project then worktree
eval "$(daycare wcd my-project)"            # Opens picker for worktree
eval "$(daycare wcd my-project feature-branch)"
```

#### `setup`
Configures the shell to handle `pcd` and `wcd` commands by evaluating their output.
Appends a wrapper function to the shell's rc file (currently zsh only).

**Examples:**
```
daycare setup
```

## Configuration
A file can be placed at ~/.config/daycare/config.yaml to configure the common
directory where projects are stored, and any other future settings.

```yaml
# ~/.config/daycare/config.yaml
common_dir: /path/to/projects  # Required. The directory containing projects.
```

Projects can contain a `.daycare/` directory with script files for project-specific
customization:

- `.daycare/pwcd` - Executed when running `daycare pcd` or `daycare wcd`. Useful
  for setting up project-specific environments, such as activating a virtual
  environment or setting environment variables.
- `.daycare/wcreate` - Executed when running `daycare wcreate` to provision a new
  worktree. Useful for copying non-git-tracked files (e.g., `.cache/`, `.claude/`)
  from the main project into the new worktree.

Both scripts are optional. If present, they are executed by the shell (sourced)
rather than run as a subprocess.

## Shell Integration
Daycare uses Cobra's built-in shell completion support to provide completions
for all subcommands, flags, and dynamic project/worktree names from the
filesystem.

**Setup:**
Daycare also ships with a setup script that configures your shell (currently zsh)
to handle the `pcd` and `wcd` commands by evaluating their output.

**Setup:**
```
daycare setup
```

This appends a wrapper function to `~/.zshrc`:
```zsh
daycare() {
  if [[ "$1" == "pcd" || "$1" == "wcd" ]]; then
    eval "$(command daycare "$@")"
  else
    command daycare "$@"
  fi
}
```

After running `daycare setup`, you can use:
```
daycare pcd my-project
daycare wcd my-project feature-branch
```
