package appsflyer

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/cuvva/cuvva-public-go/lib/s3"
)

type ReportRequester interface {
	ResourceName() string
	Get(context.Context, *string) (*Scanner, *string, error)
}

type ReportRequest struct {
	c            *Client
	t            interface{}
	resourceName string
}

func (r ReportRequest) Get(ctx context.Context, lastS3Key *string) (*Scanner, *string, error) {
	return r.c.get(ctx, r.resourceName, lastS3Key, r.t)
}

func (r ReportRequest) ResourceName() string {
	return r.resourceName
}

func (c *Client) NewClicksReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            Clicks{},
		resourceName: "clicks",
	}
}

func (c *Client) NewClicksRetargetingReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            ClicksRetargeting{},
		resourceName: "clicks_retargeting",
	}
}

func (c *Client) NewConversionsRetargetingReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            ConversionsRetargeting{},
		resourceName: "conversions_retargeting",
	}
}

func (c *Client) NewImpressionsReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            Impressions{},
		resourceName: "impressions",
	}
}

func (c *Client) NewImpressionsRetargetingReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            ImpressionsRetargeting{},
		resourceName: "impressions_retargeting",
	}
}

func (c *Client) NewInAppsReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            InApps{},
		resourceName: "inapps",
	}
}

func (c *Client) NewInAppsRetargetingReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            InAppsRetargeting{},
		resourceName: "inapps_retargeting",
	}
}

func (c *Client) NewInstallsReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            Installs{},
		resourceName: "installs",
	}
}

func (c *Client) NewOrganicUninstallsReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            OrganicUninstalls{},
		resourceName: "organic_uninstalls",
	}
}

func (c *Client) NewSessionsReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            Sessions{},
		resourceName: "sessions",
	}
}

func (c *Client) NewUninstallsReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            Uninstalls{},
		resourceName: "uninstalls",
	}
}

func (c *Client) NewWebEventsReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            WebEvents{},
		resourceName: "web_events",
	}
}

func (c *Client) NewWebToAppReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            WebToApp{},
		resourceName: "web_to_app",
	}
}

func (c *Client) NewWebTouchPointsReportRequest() ReportRequest {
	return ReportRequest{
		c:            c,
		t:            WebTouchPoints{},
		resourceName: "web_touch_points",
	}
}

type Client struct {
	s3Client   *s3.Client
	homeFolder string
}

// Parameters are of the same name you would find on AppsFlyer website.
func NewClient(awsAccessKey, homeFolder, bucketName, bucketSecret string) (*Client, error) {
	awsCredentials := credentials.NewStaticCredentials(awsAccessKey, bucketSecret, "")

	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: awsCredentials,
	})
	if err != nil {
		return nil, fmt.Errorf("aws: %w", err)
	}

	s3Client := s3.New(awsSession, bucketName)

	client := Client{
		s3Client:   s3Client,
		homeFolder: homeFolder,
	}

	return &client, nil
}

func GetReportFilePathPrefix(homeFolder, reportName string) string {
	return fmt.Sprintf("%s/data-locker-hourly/t=%s/", homeFolder, reportName)
}

func (c Client) get(ctx context.Context, reportName string, lastS3Key *string, objectType interface{}) (*Scanner, *string, error) {
	s3Prefix := GetReportFilePathPrefix(c.homeFolder, reportName)

	s3Key, err := c.getNextKey(ctx, lastS3Key, s3Prefix)
	if err != nil {
		if err == NoData {
			return nil, nil, NoData
		}
		if err == NoNewData {
			return nil, nil, NoNewData
		}
		return nil, nil, fmt.Errorf("failed to get next s3 key: %w", err)
	}

	getObjectOutput, err := c.s3Client.GetObjectWithContext(ctx,
		&awss3.GetObjectInput{Key: s3Key},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download s3 object: %w", err)
	}

	scanner, err := NewScanner(getObjectOutput.Body, objectType)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create scanner: %w", err)
	}

	paginationToken := makeS3KeyOrderable(*s3Key)

	return scanner, &paginationToken, nil
}
