package config

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sethvargo/go-githubactions"
	. "github.com/smartystreets/goconvey/convey"
	mock "github.com/yardbirdsax/action-terragrunt/internal/mock/github"
)

func TestNewConfig(t *testing.T) {
	Convey("TestNewConfig", t, func() {
		ctrl := gomock.NewController(t)
		mockAction := mock.NewMockAction(ctrl)
		tests := []struct {
			name            string
			action          *mock.MockAction
			optFns          []configOptsFn
			expectedContext *githubactions.GitHubContext
			expectedConfig  *Config
			actionStepDebug bool
		}{
			{
				name:   "WithDefaultPath",
				action: mockAction,
				optFns: []configOptsFn{},
				expectedConfig: &Config{
					gitHubContext: &githubactions.GitHubContext{
						EventName: "pull_request",
						Event:     map[string]any{},
					},
					baseDirectory: "path",
					command:       "plan",
					token:         "token",
					action:        mockAction,
				},
			},
			{
				name:   "WithDebugLogging",
				action: mockAction,
				optFns: []configOptsFn{},
				expectedConfig: &Config{
					gitHubContext: &githubactions.GitHubContext{
						EventName: "pull_request",
						Event:     map[string]any{},
					},
					baseDirectory: "path",
					command:       "plan",
					token:         "token",
					action:        mockAction,
					debug:         true,
				},
				actionStepDebug: true,
			},
		}

		for _, test := range tests {
			Convey(test.name, func() {
				mockAction.EXPECT().Context().Times(1).Return(test.expectedConfig.gitHubContext, nil)
				mockAction.EXPECT().GetInput(ActionInputDebug).Return(strconv.FormatBool(test.actionStepDebug)).Times(1).After(
					mockAction.EXPECT().GetInput(ActionInputTerraformCommand).Times(1).Return(test.expectedConfig.command).After(
						mockAction.EXPECT().GetInput(ActionInputBaseDirectory).Times(1).Return(test.expectedConfig.baseDirectory).After(
							mockAction.EXPECT().GetInput(ActionInputToken).Times(1).Return(test.expectedConfig.token),
						),
					),
				)
				Convey(test.name, func() {
					config, err := NewConfig(test.action, test.optFns...)
					Convey("should not return an error", func() {
						So(err, ShouldBeNil)
					})
					Convey("should return the expected config", func() {
						So(config, ShouldResemble, test.expectedConfig)
					})
				})
			})
		}
	})
}

func TestBaseDirectory(t *testing.T) {
	Convey("BaseDirectory", t, func() {
		expectedBaseDirectory := "base"
		config := &Config{
			baseDirectory: expectedBaseDirectory,
		}

		actualBaseDirectory := config.BaseDirectory()

		Convey("should return the correct value", func() {
			So(actualBaseDirectory, ShouldEqual, expectedBaseDirectory)
		})
	})
}

func TestCommand(t *testing.T) {
	Convey("Command", t, func() {
		expectedCommand := "plan"
		config := &Config{
			command: expectedCommand,
		}

		actualCommand := config.Command()

		Convey("should return the correct value", func() {
			So(actualCommand, ShouldEqual, expectedCommand)
		})
	})
}

func TestGitHubContext(t *testing.T) {
	Convey("GitHubContext", t, func() {
		expectedGitHubContext := &githubactions.GitHubContext{
			EventName: "pull_request",
			RunNumber: 1,
		}
		config := &Config{
			gitHubContext: expectedGitHubContext,
		}

		actualGitHubContext := config.GitHubContext()

		Convey("should return the correct value", func() {
			So(actualGitHubContext, ShouldResemble, *expectedGitHubContext)
		})
	})
}

func TestDebugEnabled(t *testing.T) {
	Convey("DebugEnabled", t, func() {
		config := &Config{
			debug: true,
		}

		Convey("should return the correct value", func() {
			So(config.DebugEnabled(), ShouldBeTrue)
		})
	})
}
