package terragrunt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	mockconfig "github.com/yardbirdsax/action-terragrunt/internal/mock/config"
	mockexec "github.com/yardbirdsax/action-terragrunt/internal/mock/exec"
)

var (
	expectedGitCommand   string   = "git"
	expectedGitArguments []string = []string{
		"config",
		"--global",
		"--add",
		"safe.directory",
		"/github/workspace",
	}
	expectedTGSwitchCommand string = "tgswitch"
)

func TestWithExecutor(t *testing.T) {
	Convey("WithExecutor", t, func() {
		ctrl := gomock.NewController(t)
		mockExecutor := mockexec.NewMockExec(ctrl)
		terragrunt := &Terragrunt{}
		f := WithExec(mockExecutor)

		err := f(terragrunt)
		So(err, ShouldBeNil)

		Convey("should set the executor correctory", func() {
			So(terragrunt.exec, ShouldEqual, mockExecutor)
		})
	})
}

func TestWithWorkingDirectory(t *testing.T) {
	Convey("WithCommand", t, func() {
		expectedPath := "path"
		terragrunt := &Terragrunt{}
		f := WithWorkingDirectory(expectedPath)

		err := f(terragrunt)

		So(err, ShouldBeNil)
		Convey("sets the working directory field", func() {
			So(terragrunt.workingDirectory, ShouldEqual, expectedPath)
		})
	})
}

func TestNewTerragrunt(t *testing.T) {
	Convey("NewTerragrunt", t, func() {
		ctrl := gomock.NewController(t)
		mockExecutor := mockexec.NewMockExec(ctrl)
		expectedPath := "path"

		terragrunt, err := NewTerragrunt(
			WithExec(mockExecutor),
			WithWorkingDirectory(expectedPath),
		)

		So(err, ShouldBeNil)
		Convey("should execute all specified functional opts", func() {
			So(terragrunt.exec, ShouldEqual, mockExecutor)
			So(terragrunt.workingDirectory, ShouldEqual, expectedPath)
		})

		Convey("should return an error if a functional opt fails", func() {
			expectedError := fmt.Errorf("this is an error")
			var errorFunc terragruntOptFns = func(t *Terragrunt) error {
				return expectedError
			}

			_, err := NewTerragrunt(errorFunc)

			So(err, ShouldEqual, expectedError)
		})
	})
}

func TestNewFromConfig(t *testing.T) {
	Convey("NewFromConfig", t, func() {
		ctrl := gomock.NewController(t)
		mockExec := mockexec.NewMockExec(ctrl)
		mockConfig := mockconfig.NewMockConfig(ctrl)
		expectedBaseDirectory := "base/directory"

		mockConfig.EXPECT().BaseDirectory().Times(1).Return(expectedBaseDirectory)
		Convey("debug logging enabled", func() {
			mockConfig.EXPECT().DebugEnabled().Times(1).Return(true)

			terragrunt, err := NewFromConfig(mockConfig, WithExec(mockExec))
			So(err, ShouldBeNil)

			Convey("should set the base directory", func() {
				So(terragrunt.workingDirectory, ShouldEqual, expectedBaseDirectory)
			})
			Convey("should have debug logging enabled", func() {
				So(terragrunt.enableDebugLogging, ShouldBeTrue)
			})
		})

	})
}

func TestRun(t *testing.T) {
	Convey("Run", t, func() {
		ctrl := gomock.NewController(t)
		mockExecutor := mockexec.NewMockExec(ctrl)
		mockConfig := mockconfig.NewMockConfig(ctrl)
		expectedCommand := TerragruntCommandPlan
		expectedWorkingDirectory := "/path/to/terragrunt"
		expectedArguments := []string{
			expectedCommand,
			terragruntWorkingDirectoryArgument,
			expectedWorkingDirectory,
		}
		expectedOutput := []string{
			"hello",
			"world",
		}
		expectedJoinedOutput := strings.Join(expectedOutput, "\n")

		Convey("when Terragrunt is executed", func() {
			expectedError := fmt.Errorf("this is an error")
			mockExecutor.EXPECT().ExecCommand(terragruntDefaultBinary, true, "", expectedArguments).Return(expectedJoinedOutput, terragruntExitCodeWithChanges, expectedError).Times(1)
			mockExecutor.EXPECT().ExecCommand(expectedGitCommand, true, "", expectedGitArguments).Return("", terragruntExitCodeNoChanges, nil).Times(1)
			mockExecutor.EXPECT().ExecCommand(expectedTGSwitchCommand, true, expectedWorkingDirectory).Times(1)
			terragrunt, err := NewTerragrunt(WithExec(mockExecutor), WithWorkingDirectory(expectedWorkingDirectory))
			So(err, ShouldBeNil)

			output, err := terragrunt.run(expectedCommand)
			Convey("should return the correct exit code", func() {
				So(output.ExitCode, ShouldEqual, terragruntExitCodeWithChanges)
			})
			Convey("should return the expected output", func() {
				So(output.Output, ShouldResemble, expectedOutput)
			})
			Convey("should return the expected error", func() {
				So(err, ShouldResemble, expectedError)
			})
		})

		Convey("when Terragrunt is executed and the debug input is set", func() {
			expectedError := fmt.Errorf("this is an error")
			expectedArgumentsPlusDebug := append(expectedArguments, "--terragrunt-log-level=trace")
			mockExecutor.EXPECT().ExecCommand(terragruntDefaultBinary, true, "", expectedArgumentsPlusDebug).Return(expectedJoinedOutput, terragruntExitCodeWithChanges, expectedError).Times(1)
			mockExecutor.EXPECT().ExecCommand(expectedGitCommand, true, "", expectedGitArguments).Return("", terragruntExitCodeNoChanges, nil).Times(1)
			mockExecutor.EXPECT().ExecCommand(expectedTGSwitchCommand, true, expectedWorkingDirectory).Times(1)
			mockConfig.EXPECT().BaseDirectory().Return(expectedWorkingDirectory)
			mockConfig.EXPECT().Infof(gomock.Any(), gomock.Any())
			mockConfig.EXPECT().DebugEnabled().Return(true)
			terragrunt, err := NewFromConfig(mockConfig, WithExec(mockExecutor), WithWorkingDirectory(expectedWorkingDirectory))
			So(err, ShouldBeNil)

			output, err := terragrunt.run(expectedCommand)
			Convey("should return the correct exit code", func() {
				So(output.ExitCode, ShouldEqual, terragruntExitCodeWithChanges)
			})
			Convey("should return the expected output", func() {
				So(output.Output, ShouldResemble, expectedOutput)
			})
			Convey("should return the expected error", func() {
				So(err, ShouldResemble, expectedError)
			})
		})

		Convey("when git config throws an error", func() {
			expectedBaseError := fmt.Errorf("this is an error")
			expectedError := fmt.Errorf("error configuring Git safe.directory setting (exit code: %d): %w", terragruntExitCodeError, expectedBaseError)
			mockExecutor.EXPECT().ExecCommand(expectedGitCommand, true, "", expectedGitArguments).Return(expectedJoinedOutput, terragruntExitCodeError, expectedBaseError).Times(1)
			terragrunt, err := NewTerragrunt(WithExec(mockExecutor), WithWorkingDirectory(expectedWorkingDirectory))
			So(err, ShouldBeNil)

			output, err := terragrunt.run(expectedCommand)

			Convey("should return the correct exit code", func() {
				So(output.ExitCode, ShouldEqual, terragruntExitCodeError)
			})
			Convey("should return the expected output", func() {
				So(output.Output, ShouldResemble, expectedOutput)
			})
			Convey("should return the expected error", func() {
				So(err, ShouldResemble, expectedError)
			})
		})

		Convey("when tgswitch throws an error", func() {
			expectedBaseError := fmt.Errorf("this is an error")
			expectedError := fmt.Errorf("error executing 'tgswitch' (exit code: %d): %w", terragruntExitCodeError, expectedBaseError)
			mockExecutor.EXPECT().ExecCommand(expectedGitCommand, true, "", expectedGitArguments).Return(expectedJoinedOutput, terragruntExitCodeNoChanges, nil).Times(1)
			mockExecutor.EXPECT().ExecCommand(expectedTGSwitchCommand, true, expectedWorkingDirectory).Return(expectedJoinedOutput, terragruntExitCodeError, expectedBaseError).Times(1)
			terragrunt, err := NewTerragrunt(WithExec(mockExecutor), WithWorkingDirectory(expectedWorkingDirectory))
			So(err, ShouldBeNil)

			output, err := terragrunt.run(expectedCommand)

			Convey("should return the correct exit code", func() {
				So(output.ExitCode, ShouldEqual, terragruntExitCodeError)
			})
			Convey("should return the expected output", func() {
				So(output.Output, ShouldResemble, expectedOutput)
			})
			Convey("should return the expected error", func() {
				So(err, ShouldResemble, expectedError)
			})
		})
	})
}

func TestPlan(t *testing.T) {
	Convey("Plan", t, func() {
		ctrl := gomock.NewController(t)
		mockExecutor := mockexec.NewMockExec(ctrl)
		expectedWorkingDirectory := "/path/to/terragrunt"
		expectedArguments := []string{
			TerragruntCommandPlan,
			terragruntWorkingDirectoryArgument,
			expectedWorkingDirectory,
			terraformArgumentDetailedExitCode,
			terraformArgumentOut,
			"/tmp/path-to-terragrunt.tfplan",
			terraformArgumentInputFalse,
		}
		expectedOutput := []string{
			"hello",
			"world",
		}
		Convey("when changes are detected", func() {
			unexpectedError := fmt.Errorf("i'm an error that you shouldn't see")
			mockExecutor.EXPECT().ExecCommand(expectedGitCommand, true, "", expectedGitArguments).Return("", terragruntExitCodeNoChanges, nil).Times(1)
			mockExecutor.EXPECT().ExecCommand(terragruntDefaultBinary, true, "", expectedArguments).Return(strings.Join(expectedOutput, "\n"), terragruntExitCodeWithChanges, unexpectedError).Times(1)
			mockExecutor.EXPECT().ExecCommand(expectedTGSwitchCommand, true, expectedWorkingDirectory).Times(1)
			terragrunt, err := NewTerragrunt(WithExec(mockExecutor), WithWorkingDirectory(expectedWorkingDirectory))
			So(err, ShouldBeNil)

			output, err := terragrunt.Plan()
			Convey("should not return an error when exit code is 2", func() {
				So(err, ShouldBeNil)
			})
			Convey("should return the correct exit code", func() {
				So(output.ExitCode, ShouldEqual, terragruntExitCodeWithChanges)
			})
			Convey("should return the expected output", func() {
				So(output.Output, ShouldResemble, expectedOutput)
			})
			Convey("should show HasChanges true when exit code is 2", func() {
				So(output.HasChanges, ShouldEqual, true)
			})
		})

		Convey("should not return error when the exit code is 0", func() {
			mockExecutor.EXPECT().ExecCommand(expectedGitCommand, true, "", expectedGitArguments).Return("", terragruntExitCodeNoChanges, nil).Times(1)
			mockExecutor.EXPECT().ExecCommand(terragruntDefaultBinary, true, "", expectedArguments).Return(strings.Join(expectedOutput, "\n"), terragruntExitCodeNoChanges, nil).Times(1)
			mockExecutor.EXPECT().ExecCommand(expectedTGSwitchCommand, true, expectedWorkingDirectory).Times(1)

			terragrunt, _ := NewTerragrunt(WithExec(mockExecutor), WithWorkingDirectory(expectedWorkingDirectory))
			output, err := terragrunt.Plan()

			So(err, ShouldBeNil)
			So(output.ExitCode, ShouldEqual, terragruntExitCodeNoChanges)
			So(output.Output, ShouldResemble, expectedOutput)
		})
	})
}

func TestApply(t *testing.T) {
	Convey("Apply", t, func() {
		ctrl := gomock.NewController(t)
		mockExecutor := mockexec.NewMockExec(ctrl)
		expectedWorkingDirectory := "/path/to/terragrunt"
		expectedArguments := []string{
			TerragruntCommandApply,
			terragruntWorkingDirectoryArgument,
			expectedWorkingDirectory,
			terraformArgumentAutoApprove,
			terraformArgumentInputFalse,
		}
		expectedOutput := []string{
			"hello",
			"world",
		}
		mockExecutor.EXPECT().ExecCommand(expectedGitCommand, true, "", expectedGitArguments).Return("", terragruntExitCodeNoChanges, nil).Times(1)
		mockExecutor.EXPECT().ExecCommand(terragruntDefaultBinary, true, "", expectedArguments).Return(strings.Join(expectedOutput, "\n"), terragruntExitCodeNoChanges, nil)
		mockExecutor.EXPECT().ExecCommand(expectedTGSwitchCommand, true, expectedWorkingDirectory).Times(1)
		terragrunt, err := NewTerragrunt(WithExec(mockExecutor), WithWorkingDirectory(expectedWorkingDirectory))
		So(err, ShouldBeNil)

		output, err := terragrunt.Apply()

		Convey("should return the expected output", func() {
			So(output.Output, ShouldResemble, expectedOutput)
		})
		Convey("should return the expected exit code", func() {
			So(output.ExitCode, ShouldEqual, terragruntExitCodeNoChanges)
		})
		Convey("should not return an error", func() {
			So(err, ShouldBeNil)
		})
	})
}

func TestGetPlanFilePath(t *testing.T) {
	Convey("GetPlanFilePath", t, func() {
		expectedWorkingDirectory := "/path/to/terragrunt"
		expectedPlanFilePath := "/tmp/path-to-terragrunt.tfplan"

		terragrunt, err := NewTerragrunt(WithWorkingDirectory(expectedWorkingDirectory))
		So(err, ShouldBeNil)

		actualPlanFilePath := terragrunt.GetPlanFilePath()
		Convey("should provide the correct path", func() {
			So(actualPlanFilePath, ShouldEqual, expectedPlanFilePath)
		})

	})
}
