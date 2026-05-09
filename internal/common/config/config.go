package config

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AWS      AWSConfig
	Frontend FrontendConfig
}

type AWSConfig struct {
	AccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID" required:"true"`
	SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY" required:"true"`
	DynamoDB        DynamoDBConfig
}

type DynamoDBConfig struct {
	Region           string `envconfig:"AWS_DYNAMODB_REGION" required:"true"`
	EndpointURL      string `envconfig:"AWS_DYNAMODB_ENDPOINT_URL" required:"true"`
	TableNameArticle string `envconfig:"AWS_DYNAMODB_TABLE_NAME_ARTICLE" required:"true"`
}

type FrontendConfig struct {
	HostURL string `envconfig:"FRONTEND_HOST_URL" required:"true"`
}

// APPプレフィックスを持つ環境変数を構造体に読み込む
func Load() (*AppConfig, error) {
	var appConfig AppConfig
	err := envconfig.Process("", &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
