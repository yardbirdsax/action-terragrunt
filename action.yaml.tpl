name: Action-Terragrunt
description: Apply Terragrunt configurations with GitHub Actions.
author: Josh Feierman <josh@sqljosh.com>
runs:
  using: 'docker'
  image: 'docker://ghcr.io/{{ getenv "GITHUB_OWNER" }}/action-terragrunt:{{ getenv "GIT_TAG" }}'
inputs:
  base-directory:
    description: |
      The base directory for all Terragrunt configurations. Defaults to the root of the
      repository.
    default: "."
    required: false
  terraform-command:
    description: |
      The Terraform command to run.
    required: true
  token:
    description: |
      The token to use when interacting with the GitHub API. This can generally be set to the
      default token.
    required: true
    default: $\{\{ github.token \}\}
  terragrunt-version:
    description: >
      The version of Terragrunt to use. This sets the `TG_VERSION` environment variable, causing the
      `tgswitch` tool to select the specified version of Terragrunt. If this is unset, `tgswitch`
      will attempt to determine the correct version of Terragrunt to install using the logic
      detailed in the [`tgswitch` repository](https://github.com/warrensbox/tgswitch).
    required: false
  terraform-version:
    description: >
      The version of Terraform to use. This sets the `TFENV_TERRAFORM_VERSION` environment variable,
      causing the `tfenv` tool to select the specified version of Terraform. If this is unset,
      `tfenv` will attempt to determine the correct version of Terraform to install using the logic
      detailed in the [`tfenv` repository](https://github.com/tfutils/tfenv).
  enable-debug-logging:
    required: false
    default: "false"
    description: >
      If set to 'true', then debug logging for Terragrunt will be enabled. It is recommended that
      you set this input to the 'ACTIONS_RUNNER_DEBUG' secret, so that if debug logging is turned on
      within GitHub Action's UI, debug logging for Terragrunt will also be turned on.
  ref:
    description: The ref of the Action. Don't set this!
    default: $\{\{ github.action_ref \}\}