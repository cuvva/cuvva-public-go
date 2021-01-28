package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/app"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
	"github.com/spf13/cobra"
)

func init() {
	UpdateDefaultCmd.Flags().BoolP("prod", "", false, "Work on prod")
	UpdateDefaultCmd.Flags().BoolP("dry-run", "", false, "Dry run only?")
	UpdateDefaultCmd.Flags().StringSliceP("overrule-checks", "", []string{}, "Overrule checks the tool does")
	UpdateDefaultCmd.Flags().StringP("message", "m", "", "More details about the deployment")

	UpdateDefaultCmd.Flags().MarkHidden("overrule-checks")

	var envs string
	types := strings.Join(cdep.ListTypes(), ", ")

	for _, env := range cdep.ListEnvironments("*") {
		envs = fmt.Sprintf("%s - %s\n", envs, env)
	}
	UpdateDefaultCmd.SetHelpTemplate(fmt.Sprintf(helpTemplateUpdateDefault, types, envs))
}

// UpdateDefaultCmd is the update-default command
var UpdateDefaultCmd = &cobra.Command{
	Use:   "update-default [type] [env]",
	Short: "Update the default commit for all services or lambdas",
	Long:  "Please read the README.md file",
	Example: strings.Join([]string{
		"update-default services avocado",
		"update-default lambda all",
	}, "\n"),
	Args: updateDefaultArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		useProd, err := cmd.Flags().GetBool("prod")
		if err != nil {
			return err
		}

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return nil
		}

		message, err := cmd.Flags().GetString("message")
		if err != nil {
			return err
		}

		params, err := parsers.Parse(args, cdep.DefaultBranch, useProd, message)
		if err != nil {
			return err
		}

		a := &app.App{
			DryRun: dryRun,
		}

		overruleChecks, err := cmd.Flags().GetStringSlice("overrule-checks")
		if err != nil {
			return err
		}

		return a.UpdateDefault(ctx, params, overruleChecks)
	},
}

func updateDefaultArgs(cmd *cobra.Command, args []string) error {
	system := "nonprod"
	if useProd, err := cmd.Flags().GetBool("prod"); err == nil && useProd {
		system = "prod"
	}

	switch len(args) {
	case 1:
		if _, err := cdep.ParseTypeArg(args[0]); err != nil {
			return err
		}
	case 2:
		if err := cdep.ValidateSystemEnvironment(system, args[1]); err != nil {
			return err
		}
	}

	switch true {
	case len(args) == 0:
		return cher.New("missing_type", cher.M{"allowed": cdep.ListTypes()})
	case len(args) == 1:
		return cher.New("missing_environment", cher.M{"allowed": cdep.ListEnvironments(system)})
	case len(args) == 2:
		return nil
	default:
		return cher.New("unknown_arguments", nil)
	}
}

const helpTemplateUpdateDefault = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
Allowed types: %s

Allowed environments:
%s`
