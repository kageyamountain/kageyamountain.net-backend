package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	t.Parallel()
	helper.InitializeIntegrationTest(t)

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		t.Run("公開ステータスの記事一覧情報が公開日の降順ソートで取得されること", func(t *testing.T) {
			t.Parallel()
			// Arrange
			ctx := context.Background()
			appConfig, err := config.Load()
			require.NoError(t, err)

			// テストデータ作成
			// テストデータ1（ステータス: publish）
			wantArticleID1 := value.GenerateArticleID().Value()
			wantStatus1 := enum.StatusPublish
			wantCreatedAt1 := time.Date(2024, 12, 14, 10, 14, 14, 0, time.UTC)
			wantPublishedAt1 := time.Date(2024, 12, 14, 10, 14, 14, 0, time.UTC)
			wantPublishedYear1 := "2024"
			wantTitle1 := "テストデータ1"
			wantContents1 := `# テストデータ1
## 見出し1
- aaa
- bbb
- ccc 
`
			wantTags1 := []string{enum.TagGo.String(), enum.TagGin.String(), enum.TagAWS.String(), enum.TagDynamoDB.String()}

			// テストデータ2（ステータス: publish）
			wantArticleID2 := value.GenerateArticleID().Value()
			wantStatus2 := enum.StatusPublish
			wantCreatedAt2 := time.Date(2024, 11, 14, 10, 14, 14, 0, time.UTC)
			wantPublishedAt2 := time.Date(2024, 11, 14, 10, 14, 14, 0, time.UTC)
			wantPublishedYear2 := "2024"
			wantTitle2 := "テストデータ2"
			wantContents2 := `# テストデータ2
## 見出し1
- aaa
- bbb
- ccc 
`
			wantTags2 := []string{enum.TagGo.String(), enum.TagGin.String()}

			// テストデータ3（ステータス: draft）
			wantArticleID3 := value.GenerateArticleID().Value()
			wantStatus3 := enum.StatusDraft
			wantPublishedAt3 := time.Time{}
			wantPublishedYear3 := "0000"
			wantCreatedAt3 := time.Date(2024, 10, 14, 10, 14, 14, 0, time.UTC)
			wantTitle3 := "テストデータ3"
			wantContents3 := `# テストデータ3
## 見出し1
- aaa
- bbb
- ccc 
`
			wantTags3 := []string{enum.TagAWS.String()}

			article1 := helper.NewTestArticle(t, wantArticleID1, wantStatus1.String(), wantCreatedAt1, wantPublishedAt1, wantPublishedYear1, wantTitle1, wantContents1, wantTags1)
			article2 := helper.NewTestArticle(t, wantArticleID2, wantStatus2.String(), wantCreatedAt2, wantPublishedAt2, wantPublishedYear2, wantTitle2, wantContents2, wantTags2)
			article3 := helper.NewTestArticle(t, wantArticleID3, wantStatus3.String(), wantCreatedAt3, wantPublishedAt3, wantPublishedYear3, wantTitle3, wantContents3, wantTags3)

			var articles []entity.Article
			articles = append(articles, *article1)
			articles = append(articles, *article2)
			articles = append(articles, *article3)

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
			req, err := http.NewRequest(http.MethodGet, testServer.URL+"/articles", nil)
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
			// Notice: 取得するArticle配列はステータス="publish"で、ソートは公開日の降順
			var wantArticles []openapi.Article
			wantArticles = append(wantArticles, openapi.Article{
				Id:          wantArticleID1,
				PublishedAt: wantPublishedAt1,
				Title:       wantTitle1,
				Tags:        wantTags1,
			})
			wantArticles = append(wantArticles, openapi.Article{
				Id:          wantArticleID2,
				PublishedAt: wantPublishedAt2,
				Title:       wantTitle2,
				Tags:        wantTags2,
			})
			wantResponseBody := openapi.ArticlesGetResponseBody{
				Articles: wantArticles,
			}

			var decodedResponseBody openapi.ArticlesGetResponseBody
			err = json.NewDecoder(gotResponse.Body).Decode(&decodedResponseBody)
			require.NoError(t, err)

			a := assert.New(t)
			a.Equal(http.StatusOK, gotResponse.StatusCode)
			a.Equal(wantResponseBody, decodedResponseBody)
		})
	})
}
