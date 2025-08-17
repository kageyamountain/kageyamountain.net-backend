package repository

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/repository"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/gateway"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository/constant"
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
	keyCondition := expression.Key(constant.AttributeStatus).Equal(expression.Value("publish"))
	projection := expression.NamesList(
		expression.Name(constant.AttributePK),
		expression.Name(constant.AttributeStatus),
		expression.Name(constant.AttributePublishedAt),
		expression.Name(constant.AttributeTitle),
		expression.Name(constant.AttributeTags),
	)
	exp, err := expression.NewBuilder().WithKeyCondition(keyCondition).WithProjection(projection).Build()
	if err != nil {
		return nil, err
	}

	// データ取得
	result, err := a.dynamoDB.Client().Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(constant.TableArticle),
		IndexName:                 aws.String(constant.GSIPublishedArticleIndex), // GSIを指定
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
			ID:          dbModel.PK,
			Status:      dbModel.Status,
			PublishedAt: time.Unix(dbModel.PublishedAt, 0).UTC(),
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

func (a articleRepository) FindByID(ctx context.Context, articleID string) (*entity.Article, error) {
	// データ取得
	result, err := a.dynamoDB.Client().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(constant.TableArticle),
		Key: map[string]types.AttributeValue{
			constant.AttributePK: &types.AttributeValueMemberS{
				Value: articleID,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// 指定ID記事の存在チェック
	if result.Item == nil {
		return nil, nil
	}

	// DBModelにマッピング
	var dbModel dbmodel.Article
	err = attributevalue.UnmarshalMap(result.Item, &dbModel)
	if err != nil {
		return nil, err
	}

	// DomainModelに変換
	domainModel, err := entity.NewArticle(&entity.NewArticleInput{
		ID:          dbModel.PK,
		Status:      dbModel.Status,
		PublishedAt: time.Unix(dbModel.PublishedAt, 0).UTC(),
		Title:       dbModel.Title,
		Contents:    dbModel.Contents,
		Tags:        dbModel.Tags,
	})
	if err != nil {
		return nil, err
	}

	return domainModel, nil
}
