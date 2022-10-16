package main

import (
	"context"
	"log"

	"github.com/sethvargo/go-githubactions"
	"github.com/yardbirdsax/action-terragrunt/pkg/config"
	"github.com/yardbirdsax/action-terragrunt/pkg/github"
	githubinterface "github.com/yardbirdsax/action-terragrunt/pkg/interfaces/github"
	terragruntinterface "github.com/yardbirdsax/action-terragrunt/pkg/interfaces/terragrunt"
	"github.com/yardbirdsax/action-terragrunt/pkg/terragrunt"
)

func main() {
	action := githubactions.New()
	action.Infof("starting up")
	config, err := config.NewConfig(action)
	if err != nil {
		action.Fatalf("error generating configuration: %v", err)
	}

	terragrunt, err := terragrunt.NewFromConfig(config)
	if err != nil {
		action.Fatalf("error creating Terragrunt configuration: %v", err)
	}

	gitHubClient, err := github.NewClientFromAction(action)
	if err != nil {
		action.Fatalf("error creating GitHub client: %v", err)
	}

	execute(terragrunt, config, gitHubClient)
}

func execute(tg terragruntinterface.Terragrunt, config *config.Config, githubClient githubinterface.Client) {
	ctx := context.TODO()
	log.Println("command is: ", config.Command())
	log.Println("base directory is: ", config.BaseDirectory())
	switch config.Command() {
	case terragrunt.TerragruntCommandPlan:
		output, err := tg.Plan()
		if (err != nil) {
			log.Fatalf("error executing Terragrunt plan: %v", err)
		}
		if (output.HasChanges) {
			if config.GitHubContext().EventName == "pull_request" {
				githubClient.CreateCommentFromPlan(ctx, output.TerragruntOutput.Output)
			}
		}
	case terragrunt.TerragruntCommandApply:
		planOutput, err := tg.Plan()
		if (err != nil) {
			log.Fatalf("error executing Terragrunt plan: %v", err)
		}
		if (planOutput.HasChanges) {
			applyOutput, err := tg.Apply()
			if err != nil {
				log.Fatalf("error executing Terragrunt apply: %v", err)
			}
			log.Printf("apply exit code: %d", applyOutput.ExitCode)
		} else {
			log.Println("no changes found, apply will not be run")
		}
	}
}