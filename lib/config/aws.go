package config

import (
	"context"

	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	credentials2 "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

// AWS configures credentials for access to Amazon Web Services.
// It is intended to be used in composition rather than a key.
type AWS struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`

	Region string `json:"region,omitempty"`
}

// Credentials returns a configured set of AWS credentials.
func (a AWS) Credentials() *credentials.Credentials {
	if a.AccessKeyID != "" && a.AccessKeySecret != "" {
		return credentials.NewStaticCredentials(a.AccessKeyID, a.AccessKeySecret, "")
	}

	return nil
}

// Session returns an AWS Session configured with region and credentials.
func (a AWS) Session() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(a.Region),
		Credentials: a.Credentials(),
	})
}

func (a AWS) SessionV2(ctx context.Context) (aws2.Config, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(a.Region),
	}

	if a.AccessKeyID != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials2.NewStaticCredentialsProvider(
				a.AccessKeyID,
				a.AccessKeySecret,
				"",
			)),
		)
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		opts...,
	)
	if err != nil {
		return aws2.Config{}, errors.Wrap(err, "aws v2 config")
	}

	return cfg, nil
}
