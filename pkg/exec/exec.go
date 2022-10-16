package exec

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (*Executor) ExecCommand(command string, writeToConsole bool, args ...string) (output string, exitCode int, err error) {
	var buffer bytes.Buffer
	exitCode = 0
	stdOutWriters := []io.Writer{&buffer}
	stdErrWriters := []io.Writer{&buffer}
	if writeToConsole {
		stdOutWriters = append(stdOutWriters, os.Stdout)
	}
	stdErrWriters = append(stdErrWriters, os.Stderr)
	stdOutW := io.MultiWriter(stdOutWriters...)
	stdErrW := io.MultiWriter(stdErrWriters...)

	cmd := exec.Command(command, args...)
	cmd.Stdout = stdOutW
	cmd.Stdin = os.Stdin
	cmd.Stderr = stdErrW

	err = cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}

	output = buffer.String()
	return
}