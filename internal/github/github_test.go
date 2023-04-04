package github

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	gogithub "github.com/google/go-github/v50/github"
	"github.com/sethvargo/go-githubactions"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yardbirdsax/action-terragrunt/internal/mock/github"
	"github.com/yardbirdsax/action-terragrunt/internal/terragrunt"
)

const (
	rawPlanOutput = `Initializing the backend...
|
| Initializing provider plugins...
| - Finding latest version of hashicorp/local...
| - Installing hashicorp/local v2.2.3...
| - Installed hashicorp/local v2.2.3 (signed by HashiCorp)
|
| Terraform has created a lock file .terraform.lock.hcl to record the provider
| selections it made above. Include this file in your version control repository
| so that Terraform can guarantee to make the same selections by default when
| you run "terraform init" in the future.
|
| Terraform has been successfully initialized!
|
| You may now begin working with Terraform. Try running "terraform plan" to see
| any changes that are required for your infrastructure. All Terraform commands
| should now work.
|
| If you ever set or change modules or backend configuration for Terraform,
| rerun this command to reinitialize your working directory. If you forget, other
| commands will detect it and remind you to do so if necessary.
|
| Terraform used the selected providers to generate the following execution
| plan. Resource actions are indicated with the following symbols:
|   + create
|
| Terraform will perform the following actions:
|
|   # local_file.name will be created
|   + resource "local_file" "name" {
|       + content              = "hello world"
|       + directory_permission = "0777"
|       + file_permission      = "0777"
|       + filename             = "out.txt"
|       + id                   = (known after apply)
|     }
|
| Plan: 1 to add, 0 to change, 0 to destroy.
|
| Changes to Outputs:
|   + output = "hello"
`
)

func TestNewClientFromAction(t *testing.T) {
	Convey("NewClientFromAction", t, func() {
		t.Setenv("INPUT_TOKEN", "token")
		ctrl := gomock.NewController(t)
		mockAction := github.NewMockAction(ctrl)
		mockAction.EXPECT().Context().Times(1).Return(&githubactions.GitHubContext{
			Repository: "owner/repo",
			EventName:  "pull_request",
			Event: map[string]any{
				"number": float64(2),
			},
		}, nil)
		mockAction.EXPECT().GetInput("token").Times(1).Return("token")

		clientInterface, err := NewClientFromAction(mockAction)

		Convey("should not return an error", func() {
			So(err, ShouldBeNil)
		})
		client := clientInterface.(*Client)
		Convey("should have the right properties", func() {
			So(client.owner, ShouldEqual, "owner")
			So(client.repository, ShouldEqual, "repo")
			So(client.pullRequestNumber, ShouldEqual, 2)
			So(client.token, ShouldEqual, "token")
		})
	})
}

func TestCreateCommentFromPlan(t *testing.T) {
	Convey("CreateCommentFromPlan", t, func() {
		planOutput := &terragrunt.TerragruntPlanOutput{
			TerragruntOutput: terragrunt.TerragruntOutput{
				Output: strings.Split(rawPlanOutput, "\n"),
				Path:   "some/path",
			},
		}
		expectedCommentText := "## Terragrunt Execution for `some/path`\n\n" +
			"<details>\n" +
			"  <summary>1 to add, 0 to change, 0 to destroy</summary>\n\n" +
			"  ```diff\n" +
			"  " + indent(2, rawPlanOutput) + "\n" +
			"  ```\n\n" +
			"</details>"
		ctrl := gomock.NewController(t)
		mockIssueService := github.NewMockIssueService(ctrl)
		client := &Client{
			issueService:      mockIssueService,
			owner:             "owner",
			repository:        "repo",
			pullRequestNumber: 2,
		}
		Convey("when there is no existing comment", func() {
			mockIssueService.EXPECT().ListComments(gomock.Any(), "owner", "repo", 2, gomock.Any()).DoAndReturn(
				func(ctx context.Context, owner string, repo string, number int, opts *gogithub.IssueListCommentsOptions) (interface{}, interface{}, interface{}) {
					return []*gogithub.IssueComment{}, nil, nil
				},
			)
			mockIssueService.EXPECT().CreateComment(gomock.Any(), "owner", "repo", 2, gomock.Any()).DoAndReturn(
				func(ctx context.Context, owner string, repo string, number int, comment *gogithub.IssueComment) (interface{}, interface{}, interface{}) {
					Convey("should use the right comment text", func() {
						So(*comment.Body, ShouldEqual, expectedCommentText)
					})
					return comment, nil, nil
				},
			)

			_, _, err := client.CreateCommentFromOutput(context.TODO(), planOutput.Output, planOutput.Path)

			So(err, ShouldBeNil)
		})
		Convey("when there is an existing comment", func() {
			var expectedCommentID int64 = 1
			originalCommentText := strings.ReplaceAll(expectedCommentText, "0 to change", "1 to change")
			mockIssueService.EXPECT().ListComments(gomock.Any(), "owner", "repo", 2, gomock.Any()).DoAndReturn(
				func(ctx context.Context, owner string, repo string, number int, opts *gogithub.IssueListCommentsOptions) (interface{}, interface{}, interface{}) {
					comments := []*gogithub.IssueComment{
						{
							ID:   &expectedCommentID,
							Body: &originalCommentText,
						},
					}
					return comments, nil, nil
				},
			)
			mockIssueService.EXPECT().DeleteComment(gomock.Any(), "owner", "repo", expectedCommentID).DoAndReturn(
				func(ctx context.Context, owner string, repo string, number int64) (interface{}, interface{}) {
					return &gogithub.Response{Response: &http.Response{Status: "200 OK"}}, nil
				},
			)
			mockIssueService.EXPECT().CreateComment(gomock.Any(), "owner", "repo", 2, gomock.Any()).DoAndReturn(
				func(ctx context.Context, owner string, repo string, number int, comment *gogithub.IssueComment) (interface{}, interface{}, interface{}) {
					Convey("should use the right comment text", func() {
						So(*comment.Body, ShouldEqual, expectedCommentText)
					})
					return comment, nil, nil
				},
			)

			_, _, err := client.CreateCommentFromOutput(context.TODO(), planOutput.Output, planOutput.Path)

			So(err, ShouldBeNil)
		})
	})
}

func TestGetSummaryFromPlanOutput(t *testing.T) {
	Convey("GetSummaryFromPlanOutput", t, func() {
		expectedOutput := "1 to add, 0 to change, 0 to destroy"

		actualOutput := getSummaryFromPlanOutput(rawPlanOutput)

		Convey("should give the right output", func() {
			So(actualOutput, ShouldEqual, expectedOutput)
		})
	})
}
