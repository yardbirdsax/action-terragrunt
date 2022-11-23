/*
Package github contains code that is for interacting with the GitHub API and the Actions interfaces.
*/
package github

import (
	"context"
	"log"

	"github.com/google/go-github/v47/github"
	githubinterface "github.com/yardbirdsax/action-terragrunt/pkg/interfaces/github"
)

func NewClientFromAction(githubinterface.Action) (githubinterface.Client, error) {
	client := &Client{}
	return client, nil
}

type Client struct {
	pullRequestService githubinterface.PullRequestService
}

func (c *Client) CreateCommentFromPlan(ctx context.Context, planOutput []string) (*github.PullRequestComment, *github.Response, error) {
	commentBody := "something"
	comment := &github.PullRequestComment{
		Body: &commentBody,
	}
	response := &github.Response{}
	c.pullRequestService.CreateComment(context.TODO(), "something", "something", 1, comment)
	log.Println("WARN: CreateCommentFromPlan method not yet implemented.")
	return comment, response, nil
}