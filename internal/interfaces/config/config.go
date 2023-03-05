//go:generate mockgen -source=config.go -destination=../../mock/config/config.go -package=config
package config

import (
	"github.com/sethvargo/go-githubactions"
)

type Config interface {
	BaseDirectory() string
	Command() string
	GitHubContext() githubactions.GitHubContext
}
