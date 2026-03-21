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
- Interactive fuzzy-search picker (using Charmbracelet Lip Gloss) for project/worktree selection when arguments are omitted

### Data Handling
This project does not persist any data beyond what the filesystem provides. It relies on the directory structure to determine projects and worktrees.

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
```

### Global Options
| Flag | Description | Default |
|------|-------------|---------|
| -h, --help | Show help | false |
| -v, --version | Show version | false |

### Commands

#### `plist`
Lists all projects found under the common directory.

**Examples:**
```
daycare plist
```

#### `wlist`
Lists worktrees for a project. If no project name is provided, opens an interactive fuzzy-search picker to select the project.

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
Creates a new project. Prompts for the project name if not provided.

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
Creates a new worktree within an existing project. The worktree name is also used as the branch name. If no project name is provided, opens an interactive fuzzy-search picker to select the project.

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
Removes a worktree from a project. If no project name is provided, opens an interactive fuzzy-search picker to select the project. Prompts for confirmation before removal.

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
Prints the path to a project. If no project name is provided, opens an interactive fuzzy-search picker to select the project. Intended to be used with command substitution: `cd $(daycare pcd <args>)`.

**Arguments:**
| Argument | Description | Required |
|----------|-------------|----------|
| `[project-name]` | Project name | No (opens picker if omitted) |

**Examples:**
```
cd $(daycare pcd)                           # Opens picker to select project
cd $(daycare pcd my-project)
```

#### `wcd`
Prints the path to a worktree. If no project name is provided, opens an interactive fuzzy-search picker to select the project. Intended to be used with command substitution: `cd $(daycare wcd <args>)`.

**Arguments:**
| Argument | Description | Required |
|----------|-------------|----------|
| `[project-name]` | Project name | No (opens picker if omitted) |
| `[worktree-name]` | Worktree name | No (opens picker if omitted) |

**Examples:**
```
cd $(daycare wcd)                           # Opens pickers for project then worktree
cd $(daycare wcd my-project)                # Opens picker for worktree
cd $(daycare wcd my-project feature-branch)
```

## Configuration
A file can be placed at ~/.config/daycare/config.yaml to configure the common
directory where projects are stored, and any other future settings.

Projects can contain a hidden shell script named `.daycare` that `daycare` will
execute when the user `cd`s into that project or worktree. This allows users to
set up project-specific environments, such as activating a virtual environment
or setting environment variables.

## Shell Integration
Daycare uses Cobra's built-in shell completion support to provide completions
for all subcommands, flags, and dynamic project/worktree names from the
filesystem.

Supported shells: **Bash**, **Zsh**, **Fish**, **PowerShell**.

**Setup:**
```
daycare completion bash > /etc/bash_completion.d/daycare      # Linux
daycare completion zsh > "${fpath[1]}/_daycare"              # Zsh
daycare completion fish > ~/.config/fish/completions/daycare.fish
```

Daycare also ships with a setup script that configures your shell (currently zsh)
to wrap the `pcd` and `wcd` commands so they actually change the shell's
working directory.

**Setup:**
```
daycare setup
```

This appends a wrapper function to `~/.zshrc`:
```zsh
daycare() {
  if [[ "$1" == "pcd" || "$1" == "wcd" ]]; then
    cd "$(command daycare "$@")"
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
