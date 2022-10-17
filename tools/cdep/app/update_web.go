package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/cuvva/cuvva-public-go/tools/cdep"
	"github.com/cuvva/cuvva-public-go/tools/cdep/git"
	"github.com/cuvva/cuvva-public-go/tools/cdep/parsers"
	"github.com/hashicorp/go-tfe"
	log "github.com/sirupsen/logrus"
)

const (
	TERRAFORM_TOKEN_ENV = "CUVVA_TERRAFORM_TOKEN"
)

var (
	distributions = map[string]map[string]string{
		"prod": {
			"website": "E2YSXQ43EF8JNW",
		},
		"nonprod": {
			"website": "E1XTHIBWUL0BQ",
		},
	}

	workspaces = map[string]string{
		"prod": "ws-aknKczrco9TGXn6o",
	}
)

func (a App) UpdateWeb(ctx context.Context, req *parsers.Params, overruleChecks []string) error {
	if req.Environment == "prod" && req.Branch != cdep.DefaultBranch {
		return cher.New("invalid_operation", nil)
	}

	if len(req.Items) != 1 {
		return cher.New("too_many_apps", nil)
	}

	log.Info("creating aws sessions")

	profileName := "nonprod"

	if req.System == "prod" {
		profileName = "root"
	}

	awsSession, err := session.NewSessionWithOptions(session.Options{
		Profile: profileName,
	})
	if err != nil {
		return err
	}

	var tf *tfe.Client
	webDeploymentBucket := fmt.Sprintf("cuvva-web-deployments-%s", req.System)
	app := req.Items[0]
	cfClient := cloudfront.New(awsSession)
	s3Client := s3.New(awsSession)

	distributionID, err := getDistributionID(app, req.System)
	if err != nil {
		return err
	}

	log.Info("getting latest commit hash")

	latestHash, err := git.GetLatestCommitHash(ctx, req.Branch)
	if err != nil {
		return err
	}

	// Due to the "cuvva-cloudfront-sites" bucket existing in root, we'll only do this
	// check on production deploys for now
	if req.System == "prod" {
		log.Info("ensuring latest commit hash has been deployed")

		cfSitesClient := s3.New(awsSession)
		ensuredFileKey := fmt.Sprintf("%s/%s/index", app, latestHash)
		_, err = cfSitesClient.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
			Bucket: aws.String("cuvva-cloudfront-sites"),
			Key:    &ensuredFileKey,
		})
		if err != nil {
			// Can't use s3.ErrCodeNoSuchKey it is not returned by the API, relies on
			// https://github.com/aws/aws-sdk-go/issues/1208
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
				return cher.New("web_deployment_not_found", nil)
			}

			return err
		}
	}

	log.Info("fetching s3 deployment file")

	deployFileResp, err := s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &webDeploymentBucket,
		Key:    aws.String(app),
	})
	if err != nil {
		return err
	}

	deployFileBytes, err := io.ReadAll(deployFileResp.Body)
	if err != nil {
		return err
	}

	currentIdentifier := strings.TrimSpace(string(deployFileBytes))

	log.Info("getting cloudfront distribution config")

	distOutput, err := cfClient.GetDistributionConfigWithContext(ctx, &cloudfront.GetDistributionConfigInput{
		Id: &distributionID,
	})
	if err != nil {
		return err
	}
	if len(distOutput.DistributionConfig.Origins.Items) != 1 {
		return cher.New("incorrect_origin_count", nil)
	}

	existingOriginPath := distOutput.DistributionConfig.Origins.Items[0].OriginPath
	if strings.HasSuffix(*existingOriginPath, latestHash) && latestHash == currentIdentifier {
		return cher.New("nothing_changed", nil)
	}

	if a.DryRun {
		log.Info("Dry run only, stopping now")
		log.Infof("current identifier: (%s), proposed identifier: (%s)\n", currentIdentifier, latestHash)
		return nil
	}

	if req.System == "prod" {
		log.Info("locking terraform workspace")

		token := os.Getenv(TERRAFORM_TOKEN_ENV)
		if token == "" {
			return cher.New("terraform_token_not_found", nil)
		}

		tf, err = tfe.NewClient(&tfe.Config{
			Token: token,
		})
		if err != nil {
			return err
		}

		_, err = tf.Workspaces.Lock(ctx, workspaces[req.System], tfe.WorkspaceLockOptions{
			Reason: ptr.String("cdep web deployment"),
		})
		if err != nil {
			return err
		}
	}

	log.Info("updating s3 deployment file")

	_, err = s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      &webDeploymentBucket,
		Key:         aws.String(app),
		ContentType: aws.String("text/plain"),
		Body:        strings.NewReader(latestHash),
	})
	if err != nil {
		return err
	}

	log.Info("updating cloudfront distribution config")

	originPath := fmt.Sprintf("/%s/%s", app, latestHash)
	distOutput.DistributionConfig.Origins.Items[0].OriginPath = &originPath

	_, err = cfClient.UpdateDistributionWithContext(ctx, &cloudfront.UpdateDistributionInput{
		Id:                 &distributionID,
		IfMatch:            distOutput.ETag,
		DistributionConfig: distOutput.DistributionConfig,
	})
	if err != nil {
		return err
	}

	if req.System == "prod" {
		log.Info("unlocking terraform workspace")

		if _, err := tf.Workspaces.Unlock(ctx, workspaces[req.System]); err != nil {
			return err
		}
	}

	log.Info("creating cloudfront invalidation")

	callerRef := fmt.Sprintf("cdep:update-web:%d", time.Now().Unix())
	_, err = cfClient.CreateInvalidationWithContext(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: &distributionID,
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: &callerRef,
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(1),
				Items:    []*string{aws.String("/*")},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func getDistributionID(app, system string) (string, error) {
	apps, ok := distributions[system]
	if !ok {
		return "", cher.New("unknown_distribution_system", nil)
	}

	distID, ok := apps[app]
	if !ok {
		return "", cher.New("unknown_distribution_app", nil)
	}

	return distID, nil
}
