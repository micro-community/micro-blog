package model

import (
	"github.com/micro-community/micro-blog/common/protos/tags"
	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/store"
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

	posts model.Table

	slugIndex model.Index
	idIndex   model.Index //resource Index
}

//NewService return a model context
func NewService(cli client.Client) *Repository {

	idIndex := model.ByEquality("id")
	idIndex.Unique = true
	idIndex.Order.Type = model.OrderTypeUnordered

	slugIndex := model.ByEquality("slug")
	slugIndex.Unique = true
	slugIndex.Order.Type = model.OrderTypeUnordered

	return &Repository{
		Tags:      tags.NewTagsService("tags", cli),
		posts:     model.NewTable(store.DefaultStore, "posts", model.Indexes(idIndex, slugIndex), nil),
		slugIndex: slugIndex,
		idIndex:   idIndex,
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
