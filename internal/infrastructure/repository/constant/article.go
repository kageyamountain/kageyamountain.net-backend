package constant

// テーブル名はAppConfigから取得する
const (
	ArticleAttributePK            = "pk"
	ArticleAttributeStatus        = "status"
	ArticleAttributeCreatedAt     = "createdAt"
	ArticleAttributePublishedAt   = "publishedAt"
	ArticleAttributePublishedYear = "publishedYear"
	ArticleAttributeTitle         = "title"
	ArticleAttributeContents      = "contents"
	ArticleAttributeTags          = "tags"

	ArticleGSIPublishedArticle = "publishedArticleIndex"
	ArticleGSIDraftArticle     = "draftArticleIndex"
	ArticleGSIPublishedYear    = "publishedYearIndex"
)
