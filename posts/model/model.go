package model

// Some Const data for post
const (
	TagType         = "post-tag"
	SlugPrefix      = "slug"
	IDPrefix        = "id"
	TimeStampPrefix = "timestamp"
)

// QueryType type
type QueryType int

// QueryType type
const (
	QueryByID QueryType = iota
	QueryBySlug
	QueryByTimestamp
)

//Post for article
type Post struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Content         string   `json:"content"`
	CreateTimestamp int64    `json:"create_timestamp"`
	UpdateTimestamp int64    `json:"update_timestamp"`
	Tags            []string `json:"tags"`
}
