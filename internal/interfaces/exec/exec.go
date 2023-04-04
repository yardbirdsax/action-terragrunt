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
type CommandFunc func(command string, args ...string) *exec.Cmd

// Exec is a wrapper interface for the `exec` package.
type Exec interface {
	ExecCommand(command string, outputToStdOut bool, workingDirectory string, args ...string) (output string, exitCode int, err error)
}