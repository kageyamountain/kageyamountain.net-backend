package main

import (
	"context"
	"fmt"
	"log"

	appconfig "github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/dbmodel"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal("Error loading .env file. err:", err)
		return
	}

	appConfig, err := appconfig.Load()
	if err != nil {
		log.Fatal("Error AppConfig Load. err:", err)
		return
	}

	// AWS設定の読み込み
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(appConfig.AWS.DynamoDB.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			appConfig.AWS.AccessKeyID, appConfig.AWS.SecretAccessKey, "",
		)),
	)
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// DynamoDBクライアントの作成
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(appConfig.AWS.DynamoDB.EndpointURL)
	})

	// 取得したい記事のPKを指定
	targetPK := "f7e8c2a1-4b9d-4c3e-8f2a-1d5e6b7c8a9f"

	// GetItem実行
	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(appConfig.AWS.DynamoDB.TableNameArticle),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: targetPK},
		},
	})
	if err != nil {
		log.Fatalf("GetItem error: %v", err)
	}

	// アイテムが存在しない場合のチェック
	if result.Item == nil {
		log.Fatalf("記事が見つかりません: %s", targetPK)
	}

	// 結果をArticle構造体にマッピング
	var article dbmodel.Article
	err = attributevalue.UnmarshalMap(result.Item, &article)
	if err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}

	// statusが"publish"でない場合はエラー
	if article.Status != "publish" {
		log.Fatalf("この記事は公開されていません。pk: %s, Status: %s", targetPK, article.Status)
	}

	// 結果を表示
	fmt.Printf("記事を正常に取得しました\n")
	fmt.Printf("パーティションキー: %s\n", article.PK)
	fmt.Printf("ステータス: %s\n", article.Status)
	fmt.Printf("作成日: %d\n", article.CreatedAt)
	fmt.Printf("公開日: %d\n", article.PublishedAt)
	fmt.Printf("公開年: %s\n", article.PublishedYear)
	fmt.Printf("件名: %s\n", article.Title)
	fmt.Printf("本文: %s\n", article.Contents)
	fmt.Printf("タグ: %v\n", article.Tags)
}
