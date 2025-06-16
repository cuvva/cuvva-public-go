package cdep

const Version = "0.6"

// DefaultBranch is the default main branch
const DefaultBranch = "master"

var ErrorCodeMapping = map[string]string{
	"config_not_on_master":      "Your config repo is not on master, please swap your HEAD back to master.",
	"frozen_without_commit":     "We found a resource where the config is locked, but no commit is specified.",
	"frozen":                    "The resource you're trying to update is currently frozen.",
	"nothing_changed":           "Running this tool has resulted in no change.",
	"unknown_environment":       "You've provided an environment that does not exist in the provided system.",
	"unknown_system":            "You've provided a system that does not exist.",
	"unknown_type":              "You're trying to update something this tool cannot handle.",
	"working_copy_dirty":        "Your config repo working copy is dirty, please clean it up and try again.",
	"too_many_apps":             "You can only specify one application for web updates",
	"web_deployment_not_found":  "The commit hash discovered has not been pushed to s3 yet",
	"terraform_token_not_found": "No Terraform token found. Please generate one on Terraform (https://app.terraform.io/app/settings/tokens) and put it into the environment variable \"CUVVA_TERRAFORM_TOKEN\".",
	"conflicting_flags":         "Cannot specify both --go-only and --js-only flags at the same time.",
}

var OverruleChecks = map[string]string{
	"working_copy_dirty": "Does not check if the working copy is dirty first",
}
