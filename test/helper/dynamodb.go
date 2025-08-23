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
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/gateway"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/constant"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/dbmodel"
	"github.com/stretchr/testify/require"
)

func InsertTestArticles(
	t testing.TB,
	ctx context.Context,
	appConfig *config.AppConfig,
	dynamoDB *gateway.DynamoDB,
	articles []entity.Article,
) {
	t.Helper()

	// 並列テスト実行のためにユニークなテーブル名に変更
	appConfig.AWS.DynamoDB.TableNameArticle = fmt.Sprintf("%s_test_%s", appConfig.AWS.DynamoDB.TableNameArticle, uuid.New().String())
	fmt.Println(appConfig.AWS.DynamoDB.TableNameArticle)

	// テスト用のテーブル作成（テスト終了時のテーブル削除処理含む）
	createTestTableArticle(t, ctx, dynamoDB, appConfig.AWS.DynamoDB.TableNameArticle)

	// テストデータ登録
	for _, article := range articles {
		dbModel := dbmodel.Article{
			PK:            article.ID.Value(),
			Status:        article.Status.String(),
			CreatedAt:     article.CreatedAt.Unix(),
			PublishedAt:   article.PublishedAt.Unix(),
			PublishedYear: article.PublishedYear,
			Title:         article.Title,
			Contents:      article.Contents,
			Tags:          make([]string, len(article.Tags)),
		}
		for i, tag := range article.Tags {
			dbModel.Tags[i] = tag.String()
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

func createTestTableArticle(t testing.TB, ctx context.Context, dynamoDB *gateway.DynamoDB, tableName string) {
	t.Helper()

	input := &dynamodb.CreateTableInput{
		BillingMode: types.BillingModePayPerRequest,
		TableName:   aws.String(tableName),

		// KeySchema
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(constant.ArticleAttributePK),
				KeyType:       types.KeyTypeHash,
			},
		},

		// AttributeDefinitions
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(constant.ArticleAttributePK),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(constant.ArticleAttributeStatus),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(constant.ArticleAttributePublishedAt),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String(constant.ArticleAttributePublishedYear),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(constant.ArticleAttributeCreatedAt),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},

		// GlobalSecondaryIndexes
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			// publishedArticleIndex
			{
				IndexName: aws.String(constant.ArticleGSIPublishedArticle),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String(constant.ArticleAttributeStatus),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String(constant.ArticleAttributePublishedAt),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
			// draftArticleIndex
			{
				IndexName: aws.String(constant.ArticleGSIDraftArticle),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String(constant.ArticleAttributeStatus),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String(constant.ArticleAttributeCreatedAt),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
			// publishedYearIndex
			{
				IndexName: aws.String(constant.ArticleGSIPublishedYear),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String(constant.ArticleAttributePublishedYear),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String(constant.ArticleAttributePublishedAt),
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
		_, err := dynamoDB.Client().DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		require.NoError(t, err)
	})
}
