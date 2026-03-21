package models

import "time"

// Project represents a software project managed by daycare.
// A project is a directory containing a Git repository under the configured projects directory.
type Project struct {
	// Name is the directory name of the project.
	Name string
	// Path is the absolute filesystem path to the project directory.
	Path string
	// LastCommit is the time of the most recent commit in the repository.
	// It is zero if the repository has no commits yet.
	LastCommit time.Time
}

// Worktree represents a Git worktree associated with a project.
type Worktree struct {
	// Name is the name of the worktree (also used as the branch name).
	Name string
	// Branch is the Git branch checked out in this worktree.
	Branch string
	// Path is the absolute filesystem path to the worktree directory.
	Path string
	// ProjectName is the name of the parent project.
	ProjectName string
	// IsMain indicates whether this is the main worktree (i.e., the project directory itself).
	IsMain bool
}
