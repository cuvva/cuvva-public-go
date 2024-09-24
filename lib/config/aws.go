package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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
