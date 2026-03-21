# Daycare

A CLI tool for managing software projects and Git worktrees. Daycare helps you
organize projects in a central directory and quickly switch between different
worktrees using an interactive fuzzy picker.

## Features

- **Project Management**: Create and list projects in a centralized projects directory
- **Worktree Management**: Create, list, and remove Git worktrees with ease
- **Interactive Picker**: Fuzzy search for selecting projects and worktrees
- **Shell Integration**: Seamless `cd` integration with zsh
- **Shell Completion**: Tab completion for project and worktree names

## Installation

### From Source

```bash
git clone https://github.com/utagai/daycare.git
cd daycare
go install
```

### Build Manually

```bash
go build -o daycare .
```

## Configuration

Daycare reads configuration from `~/.config/daycare/config.yaml`.

### Create Config Directory

```bash
mkdir -p ~/.config/daycare
```

### Example Config

```yaml
# ~/.config/daycare/config.yaml

# Directory where projects are stored
projects_dir: ~/projects
```

If no configuration file exists, Daycare defaults to `~/projects` (or `~/code` as seen in the codebase).

## Quick Start

### 1. Set Up Your Projects Directory

```bash
mkdir -p ~/projects
echo 'projects_dir: ~/projects' > ~/.config/daycare/config.yaml
```

### 2. Create a Project

```bash
daycare pcreate my-api
```

### 3. Navigate to Your Project

```bash
eval "$(daycare pcd my-api)"
```

### 4. Create a Worktree

```bash
daycare wcreate my-api feature-login
```

### 5. Switch to Worktree

```bash
eval "$(daycare wcd my-api feature-login)"
```

## Commands

### Project Commands

#### `plist`

List all projects in your projects directory.

```bash
daycare plist
```

Output:

```
┌──────────┬──────────────────────┬──────────────────┐
│ Name     │ Path                 │ Last Commit      │
├──────────┼──────────────────────┼──────────────────┤
│ my-api   │ ~/projects/my-api    │ 2 hours ago      │
│ web-app  │ ~/projects/web-app   │ 1 day ago        │
└──────────┴──────────────────────┴──────────────────┘
```

#### `pcreate`

Create a new project.

```bash
# With name argument
daycare pcreate my-project

# Without argument (prompts for name)
daycare pcreate
```

When creating a project, Daycare:
1. Creates the project directory
2. Initializes a Git repository
3. Creates a `.daycare/` subdirectory
4. Adds an initial commit ignoring the `.daycare/` directory

#### `pcd`

Output shell script to cd into a project.

```bash
# Use with eval
eval "$(daycare pcd my-project)"

# Interactive (prompts for project)
eval "$(daycare pcd)"
```

### Worktree Commands

#### `wlist`

List all worktrees for a project.

```bash
# List worktrees for specific project
daycare wlist my-project

# Interactive (prompts for project)
daycare wlist
```

Output:

```
┌───────────────┬────────────────────────────────┬──────────────────┐
│ Name          │ Path                           │ Last Commit      │
├───────────────┼────────────────────────────────┼──────────────────┤
│ main          │ ~/projects/my-project          │ 2 hours ago      │
│ feature-login │ ~/projects/my-project/feature  │ 5 minutes ago    │
└───────────────┴────────────────────────────────┴──────────────────┘
```

#### `wcreate`

Create a new worktree.

```bash
# Provide both project and worktree name
daycare wcreate my-project feature-dashboard

# Just worktree name (prompts for project)
daycare wcreate feature-dashboard
```

After creating a worktree, Daycare executes the `.daycare/wcreate` script if it exists.

#### `wcd`

Output shell script to cd into a worktree.

```bash
# Use with eval
eval "$(daycare wcd my-project feature-login)"

# Interactive (prompts for project then worktree)
eval "$(daycare wcd)"

# Just project (prompts for worktree)
eval "$(daycare wcd my-project)"
```

#### `wremove`

Remove a worktree from a project.

```bash
# Remove specific worktree
daycare wremove my-project feature-login

# Interactive (prompts for project then worktree)
daycare wremove

# Skip confirmation (for scripts)
daycare wremove my-project feature-login
```

You cannot remove the main worktree.

### Shell Integration

#### `setup`

Print shell setup instructions for zsh.

```bash
# Add to your ~/.zshrc
echo "$(daycare setup)" >> ~/.zshrc

# Or source directly
source <(daycare setup)
```

This creates a wrapper function that automatically evaluates the output of `pcd` and `wcd` commands.

#### `completion`

Generate shell completion scripts.

```bash
# For zsh
daycare completion zsh >> ~/.zshrc
```

### Global Options

```bash
# Show version
daycare --version

# Show help
daycare --help
daycare [command] --help
```

## Shell Integration Setup

Daycare's `pcd` and `wcd` commands output shell script that changes directories. To use them seamlessly, add the setup function to your shell:

```bash
echo "$(daycare setup)" >> ~/.zshrc
source ~/.zshrc
```

Then you can use `daycare pcd` and `daycare wcd` without `eval`:

```bash
daycare pcd my-project
daycare wcd my-project feature-login
```

## Custom Scripts

Projects can include custom scripts in the `.daycare/` directory:

- `.daycare/pwcd`: Executed before changing to the project/worktree directory. Use this to set up environment variables, activate virtual environments, etc.
- `.daycare/wcreate`: Executed after creating a worktree. The environment variable `PROJECT_NAME` is set when this script runs.

Example `.daycare/pwcd`:

```bash
#!/bin/bash
# Activate virtual environment
source .venv/bin/activate

# Set environment variables
export DEBUG=true
```

Example `.daycare/wcreate`:

```bash
#!/bin/bash
# Install dependencies for new worktree
npm install
```
