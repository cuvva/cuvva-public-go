package app

import (
	"github.com/aws/aws-sdk-go/service/sns"
)

type App struct {
	DryRun bool
	SNS    *sns.SNS
}
