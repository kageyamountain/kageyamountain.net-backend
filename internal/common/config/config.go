package config

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AWS      AWSConfig      `envconfig:"AWS"`
	Frontend FrontendConfig `envconfig:"FRONTEND"`
}

type AWSConfig struct {
	AccessKeyID     string         `envconfig:"ACCESS_KEY_ID" required:"true"`
	SecretAccessKey string         `envconfig:"SECRET_ACCESS_KEY" required:"true"`
	DynamoDB        DynamoDBConfig `envconfig:"DYNAMODB"`
}

type DynamoDBConfig struct {
	Region           string `envconfig:"REGION" required:"true"`
	EndpointURL      string `envconfig:"ENDPOINT_URL" required:"true"`
	TableNameArticle string `envconfig:"TABLE_NAME_ARTICLE" required:"true"`
}

type FrontendConfig struct {
	HostURL string `envconfig:"HOST_URL" required:"true"`
}

// APPプレフィックスを持つ環境変数を構造体に読み込む
func Load() (*AppConfig, error) {
	var appConfig AppConfig
	err := envconfig.Process("APP", &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
