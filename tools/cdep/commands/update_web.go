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
	UpdateWebCmd.Flags().StringP("branch", "b", cdep.DefaultBranch, "Branch to deploy")
	UpdateWebCmd.Flags().BoolP("prod", "", false, "Work on prod")
	UpdateWebCmd.Flags().BoolP("dry-run", "", false, "Dry run only?")
	UpdateWebCmd.Flags().StringSliceP("overrule-checks", "", []string{}, "Overrule checks the tool does")

	UpdateWebCmd.Flags().MarkHidden("overrule-checks")

	var envs string
	types := strings.Join(cdep.ListWebTypes(), ", ")

	for _, env := range cdep.ListEnvironments("*") {
		envs = fmt.Sprintf("%s - %s\n", envs, env)
	}

	UpdateWebCmd.SetHelpTemplate(fmt.Sprintf(helpTemplateUpdateWeb, types, envs))
}

// UpdateWebCmd is the initiator for the update command
var UpdateWebCmd = &cobra.Command{
	Use:   "update-web [type] [env] [services ...]",
	Short: "Update the deployment identifier for web applications ",
	Long:  "Please read the README.md file",
	Example: strings.Join([]string{
		"update-web cloudfront all website",
		"uw cf prod website --prod",
		"uw cf prod website -b new-landing-page --prod",
	}, "\n"),
	Aliases: []string{"uw"},
	Args:    updateWebArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		branch, err := cmd.Flags().GetString("branch")
		if err != nil {
			return err
		}

		useProd, err := cmd.Flags().GetBool("prod")
		if err != nil {
			return err
		}

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return nil
		}

		params, err := parsers.ParseWeb(args, branch, useProd)
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

		return a.UpdateWeb(ctx, params, overruleChecks)
	},
}

func updateWebArgs(cmd *cobra.Command, args []string) error {
	system := "nonprod"
	if useProd, err := cmd.Flags().GetBool("prod"); err == nil && useProd {
		system = "prod"
	}

	switch len(args) {
	case 1:
		if _, err := cdep.ParseWebTypeArg(args[0]); err != nil {
			return err
		}
	case 2:
		if err := cdep.ValidateSystemEnvironment(system, args[1]); err != nil {
			return err
		}
	}

	switch true {
	case len(args) == 0:
		return cher.New("missing_type", cher.M{"allowed": cdep.ListWebTypes()})
	case len(args) == 1:
		return cher.New("missing_environment", cher.M{"allowed": cdep.ListEnvironments(system)})
	case len(args) == 2:
		return cher.New("missing_app", nil)
	case len(args) >= 3:
		return nil
	default:
		return cher.New("unknown_arguments", nil)
	}
}

const helpTemplateUpdateWeb = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
Allowed types: %s

Allowed environments:
%s`
