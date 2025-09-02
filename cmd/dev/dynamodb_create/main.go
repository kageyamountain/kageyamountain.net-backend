package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	appConfig "github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value/enum"
	appDynamoDB "github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/gateway/dynamodb"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/dbmodel"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

	appConfig, err := appConfig.Load()
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

	// DynamoDBクライアントの作成（エンドポイントを直接指定）
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(appConfig.AWS.DynamoDB.EndpointURL)
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
		PK:            "a87ff679a2f3e71d9181a67b7542122c",
		Status:        "publish",
		CreatedAt:     1691575800, // 2024-08-09T15:30:00Z
		PublishedAt:   1691575800, // 2024-08-09T15:30:00Z
		PublishedYear: "2024",
		Title:         "サンプルタイトル",
		Contents:      contents,
		Tags:          []string{enum.TagGo.String(), enum.TagGin.String(), enum.TagAWS.String(), enum.TagDynamoDB.String()},
	}

	// 構造体をDynamoDB AttributeValue形式に変換
	item, err := attributevalue.MarshalMap(article)
	if err != nil {
		log.Fatalf("データの変換に失敗しました: %v", err)
	}

	// PK重複登録防止の条件式
	condition := expression.AttributeNotExists(expression.Name(appDynamoDB.ArticleAttributePK))
	expr, err := expression.NewBuilder().WithCondition(condition).Build()
	if err != nil {
		log.Fatalf("expression build error: %v", err)
	}

	// PutItem実行
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:                aws.String(appConfig.AWS.DynamoDB.TableNameArticle),
		Item:                     item,
		ConditionExpression:      expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
	})

	if err != nil {
		var condCheckErr *types.ConditionalCheckFailedException
		if errors.As(err, &condCheckErr) {
			log.Fatalf("データ登録に失敗しました: %v", err)
		}
		log.Fatalf("データ登録に失敗しました: %v", err)
	}

	fmt.Println("データを正常に登録しました")
}
