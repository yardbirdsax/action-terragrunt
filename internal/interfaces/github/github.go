//go:generate mockgen -source=github.go -destination=../../mock/github/github.go -package=github
/*
Package github contains code that is for interacting with the GitHub API and the Actions interfaces.
*/
package github

import (
	"context"

	"github.com/google/go-github/v50/github"
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

// IssueService is a wrapper around the GitHub IssueService exposed by "github.com/google/go-github/
type IssueService interface {
	CreateComment(ctx context.Context, owner string, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
	EditComment(ctx context.Context, owner string, repo string, commentID int64, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
	ListComments(ctx context.Context, owner string, repo string, number int, opts *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error)
}

// Client is an interface mirroring the internal GitHub wrapper client.
type Client interface {
	CreateCommentFromOutput(ctx context.Context, planOutput []string, path string) (*github.IssueComment, *github.Response, error)
}
