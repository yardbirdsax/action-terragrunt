package main

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v50/github"
	"github.com/sethvargo/go-githubactions"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yardbirdsax/action-terragrunt/internal/config"
	mockgithub "github.com/yardbirdsax/action-terragrunt/internal/mock/github"
	mockterragrunt "github.com/yardbirdsax/action-terragrunt/internal/mock/terragrunt"
	"github.com/yardbirdsax/action-terragrunt/internal/terragrunt"
)

func TestExecute(t *testing.T) {
	Convey("Execute", t, func() {
		ctrl := gomock.NewController(t)
		mockAction := mockgithub.NewMockAction(ctrl)
		mockGitHubClient := mockgithub.NewMockClient(ctrl)

		mockAction.EXPECT().GetInput(config.ActionInputDebug).Return("false")

		Convey("plan", func() {
			expectedBaseDirectory := "base/directory"
			expectedTerraformCommand := "plan"
			expectedToken := "token"

			Convey("pull_request", func() {
				eventName := "pull_request"
				eventData := map[string]interface{}{}
				gitHubContext := &githubactions.GitHubContext{
					EventName: eventName,
					Event:     eventData,
				}
				mockAction.EXPECT().GetInput(config.ActionInputTerraformCommand).Return(expectedTerraformCommand).After(
					mockAction.EXPECT().GetInput(config.ActionInputBaseDirectory).Return(expectedBaseDirectory).Return(expectedBaseDirectory).After(
						mockAction.EXPECT().GetInput(config.ActionInputToken).Times(1).Return(expectedToken),
					),
				)
				mockAction.EXPECT().Context().Return(gitHubContext, nil)
				config, err := config.NewConfig(mockAction)
				So(err, ShouldBeNil)
				mockTerragrunt := mockterragrunt.NewMockTerragrunt(ctrl)

				Convey("when plan shows changes", func() {
					expectedPlanFilePath := "/path/to/plan/file"
					expectedCommandOutput := []string{"plan", "out"}
					expectedPlanOutput := &terragrunt.TerragruntPlanOutput{
						HasChanges:   true,
						PlanFilePath: expectedPlanFilePath,
						TerragruntOutput: terragrunt.TerragruntOutput{
							Output: expectedCommandOutput,
						},
					}
					mockTerragrunt.EXPECT().Plan().Times(1).Return(expectedPlanOutput, nil)
					mockGitHubClient.EXPECT().CreateCommentFromOutput(gomock.Any(), expectedCommandOutput, expectedBaseDirectory).Times(1).Return(nil, &github.Response{Response: &http.Response{Status: "200 OK"}}, nil)
					execute(mockTerragrunt, config, mockGitHubClient)
				})

				Convey("when plan does not show changes", func() {
					expectedPlanOutput := &terragrunt.TerragruntPlanOutput{
						HasChanges: false,
					}
					mockTerragrunt.EXPECT().Plan().Times(1).Return(expectedPlanOutput, nil)
					execute(mockTerragrunt, config, mockGitHubClient)
				})
			})
		})

		Convey("apply", func() {
			expectedBaseDirectory := "base/directory"
			expectedTerraformCommand := "apply"
			expectedToken := "token"

			eventName := "pull_request"
			eventData := map[string]interface{}{}
			gitHubContext := &githubactions.GitHubContext{
				EventName: eventName,
				Event:     eventData,
			}
			mockAction.EXPECT().GetInput(config.ActionInputTerraformCommand).Return(expectedTerraformCommand).After(
				mockAction.EXPECT().GetInput(config.ActionInputBaseDirectory).Return(expectedBaseDirectory).Return(expectedBaseDirectory).After(
					mockAction.EXPECT().GetInput(config.ActionInputToken).Times(1).Return(expectedToken),
				),
			)
			mockAction.EXPECT().Context().Return(gitHubContext, nil)
			config, err := config.NewConfig(mockAction)
			So(err, ShouldBeNil)
			mockTerragrunt := mockterragrunt.NewMockTerragrunt(ctrl)

			Convey("when plan shows changes", func() {
				expectedPlanOutput := &terragrunt.TerragruntPlanOutput{
					HasChanges: true,
				}
				expectedApplyOutput := &terragrunt.TerragruntApplyOutput{
					TerragruntOutput: terragrunt.TerragruntOutput{
						ExitCode: 0,
					},
				}
				mockTerragrunt.EXPECT().Plan().Times(1).Return(expectedPlanOutput, nil)
				mockTerragrunt.EXPECT().Apply().Times(1).Return(expectedApplyOutput, nil)
				execute(mockTerragrunt, config, mockGitHubClient)
			})

			Convey("when plan does not show changes", func() {
				expectedPlanOutput := &terragrunt.TerragruntPlanOutput{
					HasChanges: false,
				}
				mockTerragrunt.EXPECT().Plan().Times(1).Return(expectedPlanOutput, nil)
				execute(mockTerragrunt, config, mockGitHubClient)
			})

		})
	})
}
