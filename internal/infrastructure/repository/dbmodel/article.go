package dbmodel

type Article struct {
	PK            string   `dynamodbav:"pk"`
	Status        string   `dynamodbav:"status"`
	CreatedAt     int64    `dynamodbav:"createdAt"`
	PublishedAt   int64    `dynamodbav:"publishedAt"`
	PublishedYear string   `dynamodbav:"publishedYear"`
	Title         string   `dynamodbav:"title"`
	Contents      string   `dynamodbav:"contents"`
	Tags          []string `dynamodbav:"tags"`
}
