//go:generate mockgen -source=terragrunt.go -destination=../../mock/terragrunt/terragrunt.go -package=terragrunt
package terragrunt

import (
	"github.com/yardbirdsax/action-terragrunt/pkg/terragrunt"
)

type Terragrunt interface {
	Plan() (*terragrunt.TerragruntPlanOutput, error)
	Apply() (*terragrunt.TerragruntApplyOutput, error)
}
