package helper

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	appDynamoDB "github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/gateway/dynamodb"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/dbmodel"
	"github.com/stretchr/testify/require"
)

func InsertTestArticles(
	t testing.TB,
	ctx context.Context,
	appConfig *config.AppConfig,
	dynamoDB *appDynamoDB.Client,
	articles []entity.Article,
) {
	t.Helper()

	// 並列テスト実行のためにユニークなテーブル名に変更
	appConfig.AWS.DynamoDB.TableNameArticle = fmt.Sprintf("%s_test_%s", appConfig.AWS.DynamoDB.TableNameArticle, uuid.New().String())

	// テスト用のテーブル作成（テスト終了時のテーブル削除処理含む）
	createTestTableArticle(t, ctx, dynamoDB, appConfig.AWS.DynamoDB.TableNameArticle)

	// テストデータ登録
	for i := range articles {
		dbModel := dbmodel.Article{
			PK:            articles[i].ID.Value(),
			Status:        articles[i].Status.String(),
			CreatedAt:     articles[i].CreatedAt.Unix(),
			UpdatedAt:     articles[i].UpdatedAt.Unix(),
			PublishedAt:   articles[i].PublishedAt.Unix(),
			PublishedYear: articles[i].PublishedYear,
			Title:         articles[i].Title,
			Contents:      articles[i].Contents,
			Tags:          make([]string, len(articles[i].Tags)),
		}
		for j, tag := range articles[i].Tags {
			dbModel.Tags[j] = tag.String()
		}

		item, err := attributevalue.MarshalMap(dbModel)
		require.NoError(t, err)

		_, err = dynamoDB.Client().PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(appConfig.AWS.DynamoDB.TableNameArticle),
			Item:      item,
		})
		require.NoError(t, err)
	}
}

func createTestTableArticle(t testing.TB, ctx context.Context, dynamoDB *appDynamoDB.Client, tableName string) {
	t.Helper()

	input := &dynamodb.CreateTableInput{
		BillingMode: types.BillingModePayPerRequest,
		TableName:   aws.String(tableName),

		// KeySchema
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(appDynamoDB.ArticleAttributePK),
				KeyType:       types.KeyTypeHash,
			},
		},

		// AttributeDefinitions
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(appDynamoDB.ArticleAttributePK),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(appDynamoDB.ArticleAttributeStatus),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(appDynamoDB.ArticleAttributePublishedAt),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String(appDynamoDB.ArticleAttributePublishedYear),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(appDynamoDB.ArticleAttributeCreatedAt),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},

		// GlobalSecondaryIndexes
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			// publishedArticleIndex
			{
				IndexName: aws.String(appDynamoDB.ArticleGSIPublishedArticle),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String(appDynamoDB.ArticleAttributeStatus),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String(appDynamoDB.ArticleAttributePublishedAt),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
			// draftArticleIndex
			{
				IndexName: aws.String(appDynamoDB.ArticleGSIDraftArticle),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String(appDynamoDB.ArticleAttributeStatus),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String(appDynamoDB.ArticleAttributeCreatedAt),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
			// publishedYearIndex
			{
				IndexName: aws.String(appDynamoDB.ArticleGSIPublishedYear),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String(appDynamoDB.ArticleAttributePublishedYear),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String(appDynamoDB.ArticleAttributePublishedAt),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
	}

	// テーブル作成
	_, err := dynamoDB.Client().CreateTable(ctx, input)
	require.NoError(t, err)

	// テーブル作成完了まで待機
	waiter := dynamodb.NewTableExistsWaiter(dynamoDB.Client())
	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 30*time.Second)
	require.NoError(t, err)

	// テスト終了時にテーブル削除
	t.Cleanup(func() {
		_, err2 := dynamoDB.Client().DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		require.NoError(t, err2)
	})
}
