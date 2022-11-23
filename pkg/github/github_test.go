package github

import (
	"context"
	"strings"
	"testing"

	gogithub "github.com/google/go-github/v47/github"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yardbirdsax/action-terragrunt/pkg/mock/github"
	"github.com/yardbirdsax/action-terragrunt/pkg/terragrunt"
)

func TestCreateCommentFromPlan(t *testing.T) {
	Convey("CreateCommentFromPlan", t, func() {
		rawPlanOutput := `
Initializing the backend...
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
		planOutput := &terragrunt.TerragruntPlanOutput {
			TerragruntOutput: terragrunt.TerragruntOutput{
				Output: strings.Split(rawPlanOutput, "\n"),
			},
		}
		expectedCommentText := "## Title \n" +
			"<details>\n" +
			"  <summary>something</summary>\n" +
			"  <p>\n" +
			"  ```diff\n" +
			rawPlanOutput + "\n" +
			"  ```\n" +
			"	 </p>\n" +
			"</details"
		ctrl := gomock.NewController(t)
		mockPullRequestService := github.NewMockPullRequestService(ctrl)
		client := &Client{
			pullRequestService: mockPullRequestService,
		}
		mockPullRequestService.EXPECT().CreateComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, owner string, repo string, number int, comment *gogithub.PullRequestComment) (interface{}, interface{}, interface{}) {
				Convey("should use the right comment text", func() {
					So(*comment.Body, ShouldEqual, expectedCommentText)
				})
				return comment, nil, nil
			},
		)

		client.CreateCommentFromPlan(context.TODO(), planOutput.Output)
	})
}
