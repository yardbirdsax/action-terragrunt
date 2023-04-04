//go:generate mockgen -source=config.go -destination=../../mock/config/config.go -package=config
package config

import (
	"github.com/sethvargo/go-githubactions"
)

type Config interface {
	BaseDirectory() string
	Command() string
	GitHubContext() githubactions.GitHubContext
	DebugEnabled() bool
	Debugf(message string, args ...any)
	Infof(message string, args ...any)
	Warningf(message string, args ...any)
	Errorf(message string, args ...any)
	Fatalf(message string, args ...any)
}
