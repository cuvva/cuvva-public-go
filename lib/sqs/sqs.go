package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/client"
	awssqs "github.com/aws/aws-sdk-go/service/sqs"
)

type SQS struct {
	*awssqs.SQS
	QueueURL string
}

func New(cfg client.ConfigProvider, queueURL string) *SQS {
	return &SQS{
		SQS:      awssqs.New(cfg),
		QueueURL: queueURL,
	}
}

func (sqs *SQS) SendMessage(input *awssqs.SendMessageInput) (*awssqs.SendMessageOutput, error) {
	input.QueueUrl = &sqs.QueueURL
	return sqs.SQS.SendMessage(input)
}

func (sqs *SQS) SendMessageWithContext(ctx context.Context, input *awssqs.SendMessageInput) (*awssqs.SendMessageOutput, error) {
	input.QueueUrl = &sqs.QueueURL
	return sqs.SQS.SendMessageWithContext(ctx, input)
}
