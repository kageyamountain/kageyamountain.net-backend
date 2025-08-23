package gateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
)

type DynamoDB struct {
	client *dynamodb.Client
}

func NewDynamoDB(ctx context.Context, appConfig *config.AppConfig) (*DynamoDB, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(appConfig.AWS.DynamoDB.Region),
		awsConfig.WithCredentialsProvider(
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

	return &DynamoDB{
		client: client,
	}, nil
}

func (d *DynamoDB) Client() *dynamodb.Client {
	return d.client
}
