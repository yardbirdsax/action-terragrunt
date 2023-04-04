/*
Package config is used to store configuration data for the Action.
*/
package config

import (
	"fmt"
	"strconv"

	"github.com/sethvargo/go-githubactions"
	"github.com/yardbirdsax/action-terragrunt/internal/interfaces/github"
)

const (
	ActionInputBaseDirectory    string = "base-directory"
	ActionInputTerraformCommand string = "terraform-command"
	ActionInputToken            string = "token"
	ActionInputDebug            string = "enable-debug-logging"
)

// Config is a struct that contains the required elements for configuring the Action.
type Config struct {
	// The GitHub Action Context associated with the run
	gitHubContext *githubactions.GitHubContext
	// The path where the base configuration resides
	baseDirectory string
	// The Terragrunt command to run
	command string
	// The GitHub token to use when interacting with the API
	token string
	// The underlying Action client
	action github.Action
	// Whether debug logging is turned on
	debug bool
}

// configOptsFn is used for functional options operating on the Config struct.
type configOptsFn func(*Config)

// NewConfig initializes a new Config object from an Action struct.
func NewConfig(action github.Action, optFns ...configOptsFn) (*Config, error) {
	config := &Config{}
	context, _ := action.Context()
	config.gitHubContext = context
	config.token = action.GetInput(ActionInputToken)
	config.baseDirectory = action.GetInput(ActionInputBaseDirectory)
	config.command = action.GetInput(ActionInputTerraformCommand)
	enableDebugInput := action.GetInput(ActionInputDebug)
	enableDebug, err := strconv.ParseBool(enableDebugInput)
	if err != nil {
		return nil, fmt.Errorf("unable to parse enable debug input (%q), value must be 'true' or 'false': %w", enableDebugInput, err)
	}
	config.debug = enableDebug
	config.action = action
	for _, f := range optFns {
		f(config)
	}
	return config, nil
}

// BaseDirectory gets the base directory for the struct.
func (c *Config) BaseDirectory() string {
	return c.baseDirectory
}

func (c *Config) Command() string {
	return c.command
}

func (c *Config) GitHubContext() githubactions.GitHubContext {
	return *c.gitHubContext
}

func (c *Config) DebugEnabled() bool {
	return c.debug
}

func (c *Config) Debugf(msg string, args ...any)   { c.action.Debugf(msg, args) }
func (c *Config) Infof(msg string, args ...any)    { c.action.Infof(msg, args) }
func (c *Config) Warningf(msg string, args ...any) { c.action.Warningf(msg, args) }
func (c *Config) Errorf(msg string, args ...any)   { c.action.Errorf(msg, args) }
func (c *Config) Fatalf(msg string, args ...any)   { c.action.Fatalf(msg, args) }
