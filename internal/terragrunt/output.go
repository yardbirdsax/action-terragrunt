package terragrunt

// TerragruntOutput is used to pass back the results of executing Terragrunt
type TerragruntOutput struct {
	// The exit code of Terragrunt
	ExitCode int

	// The combined standard output and input of Terragrunt
	Output []string

	// The path at which Terragrunt was executed.
	Path string
}

type TerragruntPlanOutput struct {
	TerragruntOutput

	// Indicates if there are changes if the operation was a `plan`
	HasChanges bool

	// The path to the plan file
	PlanFilePath string
}

type TerragruntApplyOutput struct {
	TerragruntOutput
}