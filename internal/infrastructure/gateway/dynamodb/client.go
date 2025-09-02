package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	appconfig "github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
)

type Client struct {
	client *dynamodb.Client
}

func NewClient(ctx context.Context, appConfig *appconfig.AppConfig) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(appConfig.AWS.DynamoDB.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				appConfig.AWS.AccessKeyID, appConfig.AWS.SecretAccessKey, "",
			),
		),
	)
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(appConfig.AWS.DynamoDB.EndpointURL)
	})

	return &Client{
		client: client,
	}, nil
}

func (d *Client) Client() *dynamodb.Client {
	return d.client
}
