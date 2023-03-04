package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sethvargo/go-githubactions"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yardbirdsax/action-terragrunt/pkg/config"
	mockgithub "github.com/yardbirdsax/action-terragrunt/pkg/mock/github"
	mockterragrunt "github.com/yardbirdsax/action-terragrunt/pkg/mock/terragrunt"
	"github.com/yardbirdsax/action-terragrunt/pkg/terragrunt"
)

func TestExecute(t *testing.T) {
	Convey("Execute", t, func() {
		ctrl := gomock.NewController(t)
		mockAction := mockgithub.NewMockAction(ctrl)
		mockGitHubClient := mockgithub.NewMockClient(ctrl)

		Convey("plan", func() {
			expectedBaseDirectory := "base/directory"
			expectedTerraformCommand := "plan"

			Convey("pull_request", func() {
				eventName := "pull_request"
				eventData := map[string]interface{}{}
				gitHubContext := &githubactions.GitHubContext{
					EventName: eventName,
					Event:     eventData,
				}
				mockAction.EXPECT().GetInput(config.ActionTerraformCommand).Return(expectedTerraformCommand).After(
					mockAction.EXPECT().GetInput(config.ActionInputBaseDirectory).Return(expectedBaseDirectory).Return(expectedBaseDirectory),
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
					mockGitHubClient.EXPECT().CreateCommentFromOutput(gomock.Any(), expectedCommandOutput, expectedBaseDirectory).Times(1)
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

			eventName := "pull_request"
			eventData := map[string]interface{}{}
			gitHubContext := &githubactions.GitHubContext{
				EventName: eventName,
				Event:     eventData,
			}
			mockAction.EXPECT().GetInput(config.ActionTerraformCommand).Return(expectedTerraformCommand).After(
				mockAction.EXPECT().GetInput(config.ActionInputBaseDirectory).Return(expectedBaseDirectory).Return(expectedBaseDirectory),
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
