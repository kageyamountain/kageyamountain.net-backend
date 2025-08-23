package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value/enum"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/gateway"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/openapi"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/router"
	"github.com/kageyamountain/kageyamountain.net-backend/test/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArticlesGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()
	helper.InitializeIntegrationTest(t)

	t.Run("正常系: 公開ステータスの記事一覧情報が公開日の降順ソートで取得される", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name              string
			testArticleInputs []entity.NewArticleInput
		}{
			{
				name: "ステータス=draftは取得されないこと",
				testArticleInputs: []entity.NewArticleInput{
					{
						ID:            value.GenerateArticleID().Value(),
						Status:        enum.StatusPublish.String(),
						CreatedAt:     time.Date(2024, 12, 14, 10, 14, 14, 0, time.UTC),
						PublishedAt:   time.Date(2024, 12, 14, 10, 14, 14, 0, time.UTC),
						PublishedYear: "2024",
						Title:         "テストデータ1",
						Contents: `# テストデータ1
## 見出し1
- aaa
- bbb
- ccc 
`,
						Tags: []string{enum.TagGo.String(), enum.TagGin.String(), enum.TagAWS.String(), enum.TagDynamoDB.String()},
					},
					{
						ID:            value.GenerateArticleID().Value(),
						Status:        enum.StatusPublish.String(),
						CreatedAt:     time.Date(2024, 11, 14, 10, 14, 14, 0, time.UTC),
						PublishedAt:   time.Date(2024, 11, 14, 10, 14, 14, 0, time.UTC),
						PublishedYear: "2024",
						Title:         "テストデータ2",
						Contents: `# テストデータ2
## 見出し1
- aaa
- bbb
- ccc 
`,
						Tags: []string{enum.TagGo.String(), enum.TagGin.String(), enum.TagAWS.String(), enum.TagDynamoDB.String()},
					},
					{
						ID:            value.GenerateArticleID().Value(),
						Status:        enum.StatusDraft.String(),
						CreatedAt:     time.Date(2024, 10, 14, 10, 14, 14, 0, time.UTC),
						PublishedAt:   time.Date(2024, 10, 14, 10, 14, 14, 0, time.UTC),
						PublishedYear: "2024",
						Title:         "テストデータ3",
						Contents: `# テストデータ3
## 見出し1
- aaa
- bbb
- ccc 
`,
						Tags: []string{enum.TagAWS.String()},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// Arrange
				ctx := context.Background()
				appConfig, err := config.Load()
				require.NoError(t, err)

				// テストデータ作成
				var articles []entity.Article
				var wantArticles []openapi.Article
				for _, testArticleInput := range tt.testArticleInputs {
					// DynamoDB登録用のEntity作成
					article, err := entity.NewArticle(&testArticleInput)
					require.NoError(t, err)
					articles = append(articles, *article)

					// テスト期待値の作成
					// Notice: レスポンスボディのArticle配列はステータス="publish"で、ソートは公開日の降順）
					if testArticleInput.Status != enum.StatusPublish.String() {
						continue
					}
					wantArticles = append(wantArticles, openapi.Article{
						Id:          testArticleInput.ID,
						PublishedAt: testArticleInput.PublishedAt,
						Title:       testArticleInput.Title,
						Tags:        testArticleInput.Tags,
					})
				}

				// テストデータをDynamoDBへ登録
				dynamoDB, err := gateway.NewDynamoDB(ctx, appConfig)
				require.NoError(t, err)
				helper.InsertTestArticles(t, ctx, appConfig, dynamoDB, articles)

				// ルーターのセットアップ
				r, err := router.Setup(ctx, appConfig)
				require.NoError(t, err)

				testServer := httptest.NewServer(r)
				t.Cleanup(func() { testServer.Close() })

				// リクエストの作成
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/articles", testServer.URL), http.NoBody)
				require.NoError(t, err)

				// Act
				client := &http.Client{}
				gotResponse, err := client.Do(req)
				require.NoError(t, err)
				t.Cleanup(func() {
					err := gotResponse.Body.Close()
					require.NoError(t, err)
				})

				// Assert
				wantResponseBody := openapi.ArticlesGetResponseBody{
					Articles: wantArticles,
				}

				var decodedGotResponseBody openapi.ArticlesGetResponseBody
				err = json.NewDecoder(gotResponse.Body).Decode(&decodedGotResponseBody)
				require.NoError(t, err)

				a := assert.New(t)
				a.Equal(http.StatusOK, gotResponse.StatusCode)
				a.Equal(wantResponseBody, decodedGotResponseBody)
			})
		}
	})
}
