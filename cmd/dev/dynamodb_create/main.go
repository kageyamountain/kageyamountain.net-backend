package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kageyamountain/kageyamountain.net-backend/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/dbmodel"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal("Error loading .env file. err:", err)
		return
	}

	appConfig, err := config.Load()
	if err != nil {
		log.Fatal("Error AppConfig Load. err:", err)
		return
	}

	// AWS設定の読み込み
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(appConfig.AWS.DynamoDBRegion),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			appConfig.AWS.AccessKeyID, appConfig.AWS.SecretAccessKey, "",
		)),
	)
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// DynamoDBクライアントの作成（エンドポイントを直接指定）
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(appConfig.AWS.DynamoDBEndpointURL)
	})

	contents := `
# 見出し1
## 見出し1-1
- aaa
- bbb
- ccc
`

	// 構造体でデータを準備
	article := &dbmodel.Article{
		PK:            "article#sample-003",
		Status:        "publish",
		CreatedAt:     "2024-08-09T15:30:00Z",
		PublishedAt:   "2024-08-10T10:00:00Z",
		PublishedYear: "2024",
		Title:         "サンプルタイトル",
		Contents:      contents,
		Tags:          []string{value.TagGo, value.TagGin, value.TagAWS, value.TagDynamoDB},
	}

	// 構造体をDynamoDB AttributeValue形式に変換
	item, err := attributevalue.MarshalMap(article)
	if err != nil {
		log.Fatalf("データの変換に失敗しました: %v", err)
	}

	// PutItem実行
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("article"),
		Item:      item,
	})

	if err != nil {
		log.Fatalf("データ登録に失敗しました: %v", err)
	}

	fmt.Println("データを正常に登録しました")
}
