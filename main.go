package main

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-githubactions"
	"github.com/yardbirdsax/action-terragrunt/internal/config"
	"github.com/yardbirdsax/action-terragrunt/internal/github"
	githubinterface "github.com/yardbirdsax/action-terragrunt/internal/interfaces/github"
	terragruntinterface "github.com/yardbirdsax/action-terragrunt/internal/interfaces/terragrunt"
	"github.com/yardbirdsax/action-terragrunt/internal/terragrunt"
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
	action := githubactions.New()
	action.Infof("command is: %s", config.Command())
	action.Infof("base directory is: %s", config.BaseDirectory())
	switch config.Command() {
	case terragrunt.TerragruntCommandPlan:
		output, err := tg.Plan()
		if err != nil {
			action.Fatalf("error executing Terragrunt plan: %v", err)
		}
		if output.HasChanges {
			action.Debugf("event name is: %s", config.GitHubContext().EventName)
			if config.GitHubContext().EventName == "pull_request" {
				_, _, err = githubClient.CreateCommentFromOutput(ctx, output.TerragruntOutput.Output, config.BaseDirectory())
				if err != nil {
					action.Warningf(fmt.Errorf("error creating GitHub comment: %w", err).Error())
				}
			}
		}
	case terragrunt.TerragruntCommandApply:
		planOutput, err := tg.Plan()
		if err != nil {
			action.Fatalf("error executing Terragrunt plan: %v", err)
		}
		if planOutput.HasChanges {
			applyOutput, err := tg.Apply()
			if err != nil {
				action.Fatalf("error executing Terragrunt apply: %v", err)
			}
			action.Debugf("apply exit code: %d", applyOutput.ExitCode)
		} else {
			action.Infof("no changes found, apply will not be run")
		}
	}
}
