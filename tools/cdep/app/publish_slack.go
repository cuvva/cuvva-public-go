package app

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
)

func (a App) PublishToSlack(ctx context.Context, req *parsers.Params, commitMessage string, updatedFiles []string, repoPath string) error {
	user, err := exec.CommandContext(ctx, "git", "config", "user.name").Output()
	if err != nil {
		fmt.Println(string(user))
		return err
	}

	if !a.DryRun && req.System == "prod" {
		textTemplate := ":wrench: *command*: `%s`\n:technologist: *user*: `%s`"
		text := fmt.Sprintf(textTemplate, req.String("update"), strings.Split(string(user), "\n")[0])
		if req.Message != "" {
			text = text + fmt.Sprintf("\n\n:email: *message*: `%s`", req.Message)
		}

		arn := "arn:aws:sns:eu-west-1:005717268539:cuvva-deployments-prod"
		subject := "A prod deployment is happening"
		_, err := a.SNS.PublishWithContext(ctx, &sns.PublishInput{
			TopicArn: &arn,
			Subject:  &subject,
			Message:  &text,
		})
		if err != nil {
			if err, ok := err.(awserr.Error); !ok || err.Code() != "EndpointDisabled" {
				fmt.Println(err)
			}
		}
	}

	return nil
}
