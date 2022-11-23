//go:generate mockgen -source=github.go -destination=../../mock/github/github.go -package=github
/*
Package github contains code that is for interacting with the GitHub API and the Actions interfaces.
*/
package github

import (
	"context"

	"github.com/google/go-github/v47/github"
	"github.com/sethvargo/go-githubactions"
)

// Action is an interface that mirrors the functionality exposed by the
// github.com/sethvargo/go-githubactions library's root client.
type Action interface {
	Context() (*githubactions.GitHubContext, error)
	GetInput(string) string
}

// PullRequestService is a wrapper around the GitHub PullRequestService exposed by "github.com/google/go-github/
type PullRequestService interface {
	CreateComment(ctx context.Context, owner string, repo string, number int, comment *github.PullRequestComment) (*github.PullRequestComment, *github.Response, error)
}

// Client is an interface mirroring the internal GitHub wrapper client.
type Client interface {
	CreateCommentFromPlan(ctx context.Context, planText []string) (*github.PullRequestComment, *github.Response, error)
}