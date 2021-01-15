package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Client struct {
	*awss3.S3
	Bucket *string
}

// New creates a new S3 client
func New(cfg client.ConfigProvider, bucket string) *Client {
	return &Client{
		S3:     awss3.New(cfg),
		Bucket: aws.String(bucket),
	}
}

func (c *Client) ListObjectsV2WithContext(ctx context.Context, input *awss3.ListObjectsV2Input) (*awss3.ListObjectsV2Output, error) {
	input.Bucket = c.Bucket
	return c.S3.ListObjectsV2WithContext(ctx, input)
}

func (c *Client) GetObjectWithContext(ctx context.Context, input *awss3.GetObjectInput) (*awss3.GetObjectOutput, error) {
	input.Bucket = c.Bucket
	return c.S3.GetObjectWithContext(ctx, input)
}

func (c *Client) GetObjectRequest(input *awss3.GetObjectInput) (*request.Request, *awss3.GetObjectOutput) {
	input.Bucket = c.Bucket
	return c.S3.GetObjectRequest(input)
}

func (c *Client) CreateMultipartUploadWithContext(ctx context.Context, input *awss3.CreateMultipartUploadInput) (output *awss3.CreateMultipartUploadOutput, err error) {
	input.Bucket = c.Bucket
	return c.S3.CreateMultipartUploadWithContext(ctx, input)
}

func (c *Client) UploadPartWithContext(ctx context.Context, input *awss3.UploadPartInput) (output *awss3.UploadPartOutput, err error) {
	input.Bucket = c.Bucket
	return c.S3.UploadPartWithContext(ctx, input)
}

func (c *Client) CompleteMultipartUploadWithContext(ctx context.Context, input *awss3.CompleteMultipartUploadInput) (output *awss3.CompleteMultipartUploadOutput, err error) {
	input.Bucket = c.Bucket
	return c.S3.CompleteMultipartUploadWithContext(ctx, input)
}

func (c *Client) PutObjectWithContext(ctx context.Context, input *awss3.PutObjectInput) (*awss3.PutObjectOutput, error) {
	input.Bucket = c.Bucket
	return c.S3.PutObjectWithContext(ctx, input)
}

func (c *Client) UploadWithContext(ctx context.Context, options func(*s3manager.Uploader), input *s3manager.UploadInput) (*s3manager.UploadOutput, error) {
	input.Bucket = c.Bucket
	uploader := s3manager.NewUploaderWithClient(c.S3, options)

	return uploader.UploadWithContext(ctx, input)
}

func (c *Client) PutObjectRequest(input *awss3.PutObjectInput) (*request.Request, *awss3.PutObjectOutput) {
	input.Bucket = c.Bucket
	return c.S3.PutObjectRequest(input)
}

func (c *Client) HeadObjectWithContext(ctx context.Context, input *awss3.HeadObjectInput) (*awss3.HeadObjectOutput, error) {
	input.Bucket = c.Bucket
	return c.S3.HeadObjectWithContext(ctx, input)
}

func (c *Client) DeleteObjectWithContext(ctx context.Context, input *awss3.DeleteObjectInput) (*awss3.DeleteObjectOutput, error) {
	input.Bucket = c.Bucket
	return c.S3.DeleteObjectWithContext(ctx, input)
}
