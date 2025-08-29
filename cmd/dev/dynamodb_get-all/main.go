package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/constant"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/dbmodel"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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
		awsConfig.WithRegion(appConfig.AWS.DynamoDB.Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
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

	// Key Condition Expression: status = "publish"
	keyEx := expression.Key("status").Equal(expression.Value("publish"))

	// Expressionを構築
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Fatalf("expression build error: %v", err)
	}

	// Query実行
	result, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(appConfig.AWS.DynamoDB.TableNameArticle),
		IndexName:                 aws.String(constant.ArticleGSIPublishedArticle), // GSIを指定
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(false), // 降順（新しい順）
	})
	if err != nil {
		log.Fatalf("query error: %v", err)
	}

	// 結果をArticle構造体にマッピング
	var articles []*dbmodel.Article
	err = attributevalue.UnmarshalListOfMaps(result.Items, &articles)
	if err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}

	// 結果を表示
	fmt.Printf("取得した記事数: %d\n", len(articles))
	for i, article := range articles {
		fmt.Printf("%d. タイトル: %s, 公開日: %d\n", i+1, article.Title, article.PublishedAt)
	}
}
