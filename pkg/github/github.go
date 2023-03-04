/*
Package github contains code that is for interacting with the GitHub API and the Actions interfaces.
*/
package github

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/google/go-github/v47/github"
	githubinterface "github.com/yardbirdsax/action-terragrunt/pkg/interfaces/github"
)

var (
	//go:embed plan_comment.go.tmpl
	commentTemplate string
)

func NewClientFromAction(githubinterface.Action) (githubinterface.Client, error) {
	ghClient := github.NewClient(nil)
	client := &Client{
		pullRequestService: ghClient.PullRequests,
	}
	return client, nil
}

type Client struct {
	pullRequestService githubinterface.PullRequestService
}

func (c *Client) CreateCommentFromOutput(ctx context.Context, planOutput []string, path string) (*github.PullRequestComment, *github.Response, error) {
	buf := bytes.Buffer{}
	commentData := struct {
		Path    string
		Summary string
		Output  string
	}{
		Path:    path,
		Output:  strings.Join(cleanOutput(planOutput), "\n"),
		Summary: getSummaryFromPlanOutput(strings.Join(planOutput, "\n")),
	}
	commentTemplate := template.Must(template.New("comment").
		Funcs(template.FuncMap{"indent": indent}).
		Parse(commentTemplate))

	err := commentTemplate.Execute(&buf, commentData)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing comment template: %w", err)
	}
	commentBody := buf.String()
	comment := &github.PullRequestComment{
		Body: &commentBody,
	}
	return c.pullRequestService.CreateComment(context.TODO(), "something", "something", 1, comment)
}

func cleanOutput(output []string) []string {
	cleanedOutput := []string{}

	for _, s := range output {
		cleaned := strings.TrimSpace(s)
		cleanedOutput = append(cleanedOutput, cleaned)
	}
	return cleanedOutput
}

func getSummaryFromPlanOutput(output string) string {
	re := regexp.MustCompile(`\|\s+Plan:\s+(\d+ to add, \d+ to change, \d+ to destroy)`)
	submatches := re.FindStringSubmatch(output)
	if len(submatches) > 1 {
		return submatches[1]
	} else {
		return ""
	}
}

func indent(length int, v string) string {
	return strings.ReplaceAll(v, "\n", "\n"+strings.Repeat(" ", length))
}
