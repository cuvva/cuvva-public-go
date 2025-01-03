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
}

var OverruleChecks = map[string]string{
	"working_copy_dirty": "Does not check if the working copy is dirty first",
}
