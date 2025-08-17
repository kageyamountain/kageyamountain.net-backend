package constant

const TableArticle = "article"

const (
	AttributePK            = "pk"
	AttributeStatus        = "status"
	AttributeCreatedAt     = "createdAt"
	AttributePublishedAt   = "publishedAt"
	AttributePublishedYear = "publishedYear"
	AttributeTitle         = "title"
	AttributeContents      = "contents"
	AttributeTags          = "tags"
)

const (
	GSIPublishedArticleIndex = "publishedArticleIndex"
	GSIDraftArticleIndex     = "draftArticleIndex"
	GSIPublishedYearIndex    = "publishedYearIndex"
)
