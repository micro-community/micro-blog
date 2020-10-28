package model

import (
	"github.com/micro-community/micro-blog/common/protos/tags"
	"github.com/micro/micro/v3/service/client"
)

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

//Repository to handle DB
type Repository struct {
	Tags tags.TagsService
}

//NewService return a model context
func NewService(cli client.Client) *Repository {
	return &Repository{
		Tags: tags.NewTagsService("tags", cli),
	}
}

//Post for article
type Post struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Slug    string   `json:"slug"`
	Content string   `json:"content"`
	Created int64    `json:"created"`
	Updated int64    `json:"updated"`
	Tags    []string `json:"tags"`
}
