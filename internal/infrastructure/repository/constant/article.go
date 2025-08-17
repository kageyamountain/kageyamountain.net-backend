package constant

type article struct {
	TableName    string
	PartitionKey string
	Attribute    articleAttributes
	GSI          gsi
}

type articleAttributes struct {
	PK            string
	Status        string
	CreatedAt     string
	PublishedAt   string
	PublishedYear string
	Title         string
	Contents      string
	Tags          string
}

type gsi struct {
	PublishedArticleIndex string
	DraftArticleIndex     string
	PublishedYearIndex    string
}

var Article = article{
	TableName:    "article",
	PartitionKey: "pk",
	Attribute: articleAttributes{
		PK:            "pk",
		Status:        "status",
		CreatedAt:     "createdAt",
		PublishedAt:   "publishedAt",
		PublishedYear: "publishedYear",
		Title:         "title",
		Contents:      "contents",
		Tags:          "tags",
	},
	GSI: gsi{
		PublishedArticleIndex: "publishedArticleIndex",
		DraftArticleIndex:     "draftArticleIndex",
		PublishedYearIndex:    "publishedYearIndex",
	},
}
