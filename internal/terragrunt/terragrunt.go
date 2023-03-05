/*
Package terragrunt contains all logic around the execution of the Terragrunt binary.
*/
package terragrunt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yardbirdsax/action-terragrunt/internal/config"
	"github.com/yardbirdsax/action-terragrunt/internal/exec"
	interfaces "github.com/yardbirdsax/action-terragrunt/internal/interfaces/exec"
)

const (
	terragruntDefaultBinary            string = "terragrunt"
	terragruntWorkingDirectoryArgument string = "--terragrunt-working-dir"
	terragruntExitCodeNoChanges        int    = 0
	terragruntExitCodeError            int    = 1
	terragruntExitCodeWithChanges      int    = 2
	TerragruntCommandApply             string = "apply"
	TerragruntCommandPlan              string = "plan"
	terraformArgumentAutoApprove       string = "-auto-approve"
	terraformArgumentDetailedExitCode  string = "-detailed-exitcode"
	terraformArgumentInputFalse        string = "-input=false"
	terraformArgumentOut               string = "-out"
)

type Terragrunt struct {
	// exec is an instance of the Exec interface used to run the Terragrunt commands
	exec interfaces.Exec

	// workingDirectory is the working directory passed to Terragrunt via the --terragrunt-working-dir
	// argument
	workingDirectory string
}

// terragruntOptFns defines optional functions for the Terragrunt struct
type terragruntOptFns func(*Terragrunt) error

// WithExec specifies an Exec object for use with the struct. It's internal because it should
// only be used to mock out executions for tests.
func WithExec(exec interfaces.Exec) terragruntOptFns {
	return func(t *Terragrunt) error {
		t.exec = exec
		return nil
	}
}

// WithCommand sets the command and arguments passed to Terragrunt
// func WithCommand(command string, arguments ...string) terragruntOptFns {
// 	return func(t *Terragrunt) error {
// 		t.command = command
// 		t.arguments = arguments
// 		return nil
// 	}
// }

// WithWorkingDirectory sets the working directory where Terragrunt will be executed
func WithWorkingDirectory(path string) terragruntOptFns {
	return func(t *Terragrunt) error {
		t.workingDirectory = path
		return nil
	}
}

// NewTerragrunt is used to create a new Terragrunt struct
func NewTerragrunt(opts ...terragruntOptFns) (*Terragrunt, error) {
	terragrunt := &Terragrunt{
		exec: exec.NewExecutor(),
	}
	for _, f := range opts {
		err := f(terragrunt)
		if err != nil {
			return nil, err
		}
	}
	return terragrunt, nil
}

func (t *Terragrunt) run(command string, arguments ...string) (*TerragruntOutput, error) {
	output := &TerragruntOutput{}

	gitCommand := "git"
	gitArguments := []string{
		"config",
		"--global",
		"--add",
		"safe.directory",
		"/github/workspace",
	}
	execOutput, exitCode, err := t.exec.ExecCommand(gitCommand, true, gitArguments...)
	if err != nil || exitCode != terragruntExitCodeNoChanges {
		output.Output = strings.Split(execOutput, "\n")
		output.ExitCode = exitCode
		err = fmt.Errorf("Error configuring Git safe.directory setting (exit code: %d): %w", exitCode, err)
		return output, err
	}

	combinedArguments := []string{command, terragruntWorkingDirectoryArgument, t.workingDirectory}
	combinedArguments = append(combinedArguments, arguments...)
	execOutput, exitCode, err = t.exec.ExecCommand(terragruntDefaultBinary, true, combinedArguments...)
	output.Output = strings.Split(execOutput, "\n")
	output.ExitCode = exitCode
	return output, err
}

// NewFromConfig is used to create a new Terragrunt struct from a Config object
func NewFromConfig(config *config.Config, opts ...terragruntOptFns) (*Terragrunt, error) {
	terragrunt, err := NewTerragrunt(opts...)
	if err != nil {
		return nil, err
	}
	terragrunt.workingDirectory = config.BaseDirectory()
	return terragrunt, nil
}

func (t *Terragrunt) Plan() (*TerragruntPlanOutput, error) {
	planOutput := &TerragruntPlanOutput{}
	planFilePath := t.GetPlanFilePath()
	output, err := t.run(TerragruntCommandPlan, terraformArgumentDetailedExitCode, terraformArgumentOut, planFilePath, terraformArgumentInputFalse)
	planOutput.TerragruntOutput = *output
	if output.ExitCode == terragruntExitCodeWithChanges {
		planOutput.HasChanges = true
		// deliberately don't return an error if the exit code is 2, since that's
		// expected behavior
		return planOutput, nil
	}
	return planOutput, err
}

func (t *Terragrunt) Apply() (*TerragruntApplyOutput, error) {
	applyOutput := &TerragruntApplyOutput{}
	output, err := t.run(TerragruntCommandApply, terraformArgumentAutoApprove, terraformArgumentInputFalse)
	applyOutput.TerragruntOutput = *output
	return applyOutput, err
}

// GetPlanFilePath generates the path at which a plan file should be placed in a standard way.
func (t *Terragrunt) GetPlanFilePath() string {
	regexSlashSpace := regexp.MustCompile(`[/ ]`)
	regexLeadingDash := regexp.MustCompile(`^[- /]`)
	sanitizedWorkingDirectory := regexSlashSpace.ReplaceAllString(regexLeadingDash.ReplaceAllString(t.workingDirectory, ""), "-")

	return fmt.Sprintf("/tmp/%s.tfplan", sanitizedWorkingDirectory)
}
