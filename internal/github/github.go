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

	"github.com/google/go-github/v50/github"
	githubinterface "github.com/yardbirdsax/action-terragrunt/internal/interfaces/github"

	"golang.org/x/oauth2"
)

var (
	//go:embed plan_comment.go.tmpl
	commentTemplate string
)

type Client struct {
	issueService      githubinterface.IssueService
	owner             string
	repository        string
	pullRequestNumber int
	token             string
}

func NewClientFromAction(action githubinterface.Action) (githubinterface.Client, error) {

	actionContext, err := action.Context()
	if err != nil {
		return nil, fmt.Errorf("error getting GitHub context: %w", err)
	}
	actionRepository := strings.Split(actionContext.Repository, "/")
	if len(actionRepository) < 2 {
		return nil, fmt.Errorf("action repository string (%s) is not of the correct pattern", actionContext.Repository)
	}
	client := &Client{
		owner:      actionRepository[0],
		repository: actionRepository[1],
		token:      action.GetInput("token"),
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: client.token,
		},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)
	client.issueService = ghClient.Issues

	switch actionContext.EventName {
	case "pull_request":
		client.pullRequestNumber = int(actionContext.Event["number"].(float64))
	}
	return client, nil
}

func (c *Client) CreateCommentFromOutput(ctx context.Context, planOutput []string, path string) (*github.IssueComment, *github.Response, error) {
	cleanedPlanOutput := strings.Join(cleanOutput(planOutput), "\n")
	title := fmt.Sprintf("Terragrunt Execution for `%s`", path)
	buf := bytes.Buffer{}
	commentData := struct {
		Path    string
		Summary string
		Output  string
		Title   string
	}{
		Path:    path,
		Output:  cleanedPlanOutput,
		Summary: getSummaryFromPlanOutput(cleanedPlanOutput),
		Title:   title,
	}
	commentTemplate := template.Must(template.New("comment").
		Funcs(template.FuncMap{"indent": indent}).
		Parse(commentTemplate))

	err := commentTemplate.Execute(&buf, commentData)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing comment template: %w", err)
	}
	commentBody := buf.String()
	comment := &github.IssueComment{
		Body: &commentBody,
	}

	existingComments, _, err := c.issueService.ListComments(context.TODO(), c.owner, c.repository, c.pullRequestNumber, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving existing issue comments: %w", err)
	}
	var existingCommentID int64 = 0
	for _, c := range existingComments {
		if strings.Contains(*c.Body, title) {
			existingCommentID = *c.ID
		}
	}

	if existingCommentID != 0 {
		resp, err := c.issueService.DeleteComment(context.TODO(), c.owner, c.repository, existingCommentID)
		if err != nil {
			return nil, resp, fmt.Errorf("error deleting comment with ID %q: %w", existingCommentID, err)
		}
	}
	return c.issueService.CreateComment(context.TODO(), c.owner, c.repository, c.pullRequestNumber, comment)
}

func cleanOutput(output []string) []string {
	cleanedOutput := []string{}

	for _, s := range output {
		cleaned := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(strings.TrimSpace(s), "")
		cleanedOutput = append(cleanedOutput, cleaned)
	}
	return cleanedOutput
}

func getSummaryFromPlanOutput(output string) string {
	re := regexp.MustCompile(`Plan:\s+(\d+ to add, \d+ to change, \d+ to destroy)`)
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
