//go:generate mockgen -source=exec.go -destination=../../mock/exec/exec.go -package=exec
package interfaces

import (
	"os/exec"
)

// Cmd is a wrapper interface for the `exec.Cmd` type.
type Cmd interface {
	Run() error
}

// CommandFunc is a wrapper interface for the `exec.Command` function.
type CommandFunc func(string, ...string) *exec.Cmd

// Exec is a wrapper interface for the `exec` package.
type Exec interface {
	ExecCommand(string, bool, ...string) (string, int, error)
}