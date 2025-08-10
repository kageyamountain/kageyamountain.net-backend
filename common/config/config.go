package config

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AWS AWSConfig `envconfig:"AWS"`
}

type AWSConfig struct {
	AccessKeyID         string `envconfig:"ACCESS_KEY_ID" required:"true"`
	SecretAccessKey     string `envconfig:"SECRET_ACCESS_KEY" required:"true"`
	DynamoDBRegion      string `envconfig:"DYNAMODB_REGION" required:"true"`
	DynamoDBEndpointURL string `envconfig:"DYNAMODB_ENDPOINT_URL" required:"true"`
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
