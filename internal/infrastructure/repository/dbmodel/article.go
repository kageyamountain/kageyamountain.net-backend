package dbmodel

type Article struct {
	PK            string   `dynamodbav:"pk"`
	Status        string   `dynamodbav:"status"`
	CreatedAt     string   `dynamodbav:"createdAt"`
	PublishedAt   string   `dynamodbav:"publishedAt"`
	PublishedYear string   `dynamodbav:"publishedYear"`
	Title         string   `dynamodbav:"title"`
	Contents      string   `dynamodbav:"contents"`
	Tags          []string `dynamodbav:"tags"`
}
