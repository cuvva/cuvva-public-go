package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/config"
	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/app"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	UpdateCmd.Flags().StringP("branch", "b", cdep.DefaultBranch, "Branch to deploy")
	UpdateCmd.Flags().StringP("commit", "c", "", "Commit to deploy instead of the latest")
	UpdateCmd.Flags().BoolP("prod", "", false, "Work on prod")
	UpdateCmd.Flags().BoolP("dry-run", "", false, "Dry run only?")
	UpdateCmd.Flags().StringSliceP("overrule-checks", "", []string{}, "Overrule checks the tool does")
	UpdateCmd.Flags().StringP("message", "m", "", "More details about the deployment")

	UpdateCmd.Flags().MarkHidden("overrule-checks")

	var envs string
	types := strings.Join(cdep.ListTypes(), ", ")

	for _, env := range cdep.ListEnvironments("*") {
		envs = fmt.Sprintf("%s - %s\n", envs, env)
	}

	UpdateCmd.SetHelpTemplate(fmt.Sprintf(helpTemplateUpdate, types, envs))
}

// UpdateCmd is the initiator for the update command
var UpdateCmd = &cobra.Command{
	Use:   "update [type] [env] [services ...]",
	Short: "Update the branch and commit for a selection of services, lambdas, or cloudfront distributions",
	Long:  "Please read the README.md file",
	Example: strings.Join([]string{
		"update service avocado sms email -b extra-logging",
		"update service avocado sms email -b extra-logging -c f1ec178befe6ed26ce9cec0aa419c763c203bc92",
		"update service all sms email -c 1ed6fd7450031a5240584f8bbe8ec527f9020b5b",
		"update service prod email --prod",
		"update lambda basil ltm-proxy",
		"update cloudfront prod website --prod",
		"update terra avocado aws-env",
	}, "\n"),
	Aliases: []string{"u"},
	Args:    updateArgs,
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

		message, err := cmd.Flags().GetString("message")
		if err != nil {
			return err
		}

		commit, err := cmd.Flags().GetString("commit")
		if err != nil {
			return err
		}

		params, err := parsers.Parse(args, useProd)
		if err != nil {
			return err
		}

		params.Branch = branch
		params.Message = message
		params.Commit = commit

		awsSession, err := session.NewSessionWithOptions(session.Options{
			Profile: "root",
			Config: aws.Config{
				Region:      ptr.String("eu-west-1"),
				Credentials: config.AWS{}.Credentials(),
			},
		})
		if err != nil {
			return errors.Wrap(err, "aws:")
		}

		sns := sns.New(awsSession)

		a := &app.App{
			DryRun: dryRun,
			SNS:    sns,
		}

		overruleChecks, err := cmd.Flags().GetStringSlice("overrule-checks")
		if err != nil {
			return err
		}

		return a.Update(ctx, params, overruleChecks)
	},
}

func updateArgs(cmd *cobra.Command, args []string) error {
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
		return cher.New("missing_services", nil)
	case len(args) >= 3:
		return nil
	default:
		return cher.New("unknown_arguments", nil)
	}
}

const helpTemplateUpdate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
Allowed types: %s

Allowed environments:
%s`
