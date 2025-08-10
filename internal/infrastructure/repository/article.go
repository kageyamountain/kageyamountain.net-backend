package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/repository"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/gateway"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/dbmodel"
)

type articleRepository struct {
	dynamoDB *gateway.DynamoDB
}

func NewArticleRepository(dynamoDB *gateway.DynamoDB) repository.ArticleRepository {
	return &articleRepository{
		dynamoDB: dynamoDB,
	}
}

func (a articleRepository) FindAllForList(ctx context.Context) ([]*entity.Article, error) {
	// データ取得仕様の定義
	keyCondition := expression.Key("status").Equal(expression.Value("publish"))
	projection := expression.NamesList(
		expression.Name("pk"),
		expression.Name("status"),
		expression.Name("publishedAt"),
		expression.Name("title"),
		expression.Name("tags"),
	)
	exp, err := expression.NewBuilder().WithKeyCondition(keyCondition).WithProjection(projection).Build()
	if err != nil {
		return nil, err
	}

	// データ取得
	result, err := a.dynamoDB.Client().Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String("article"),
		IndexName:                 aws.String("publishedArticleIndex"), // GSIを指定
		KeyConditionExpression:    exp.KeyCondition(),
		ProjectionExpression:      exp.Projection(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		ScanIndexForward:          aws.Bool(false), // 公開日の降順
	})
	if err != nil {
		return nil, err
	}

	// DBModelにマッピング
	var dbModels []*dbmodel.Article
	err = attributevalue.UnmarshalListOfMaps(result.Items, &dbModels)
	if err != nil {
		return nil, err
	}

	// DomainModelに変換
	var domainModels []*entity.Article
	for _, dbModel := range dbModels {
		domainModel, err := entity.NewArticle(&entity.NewArticleInput{
			PK:          dbModel.PK,
			Status:      dbModel.Status,
			PublishedAt: dbModel.PublishedAt,
			Title:       dbModel.Title,
			Tags:        dbModel.Tags,
		})
		if err != nil {
			return nil, err
		}

		domainModels = append(domainModels, domainModel)
	}

	return domainModels, nil
}
