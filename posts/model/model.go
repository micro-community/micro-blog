package model

import "github.com/micro-community/micro-blog/common/protos/tags"

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

//DB to handle DB
type DB struct {
	Tags tags.TagsService
}

//NewService return a model context
func NewService(tagsService tags.TagsService) *DB {
	return &DB{
		Tags: tagsService,
	}
}

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
