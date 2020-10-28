package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

//QueryPostByID from db
func (r *Repository) QueryPostByID(ctx context.Context, id string) ([]*Post, error) {

	logger.Infof("Reading post by id: %v", id)

	//post := &Post{}
	//return post, r.posts.Read(r.idIndex.ToQuery(id), post)

	key := fmt.Sprintf("%v:%v", IDPrefix, id)
	records, err := store.Read("", store.Prefix(key))

	if err != nil {
		return nil, errors.BadRequest("posts.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	posts := buildPostsFromRecords(records)
	return posts, nil
}

//QueryPostBySlug from db
func (r *Repository) QueryPostBySlug(ctx context.Context, slug string) ([]*Post, error) {

	logger.Infof("Reading post by slug: %v", slug)

	key := fmt.Sprintf("%v:%v", SlugPrefix, slug)
	records, err := store.Read("", store.Prefix(key))

	if err != nil {
		return nil, errors.BadRequest("posts.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	posts := buildPostsFromRecords(records)
	return posts, nil
}

//QueryPostByTimeStamp from db
func (r *Repository) QueryPostByTimeStamp(ctx context.Context, qLimit, qOffset int64) ([]*Post, error) {

	key := fmt.Sprintf("%v:", TimeStampPrefix)
	var limit uint
	limit = 20 //default if without limition in req
	if qLimit > 0 {
		limit = uint(qLimit)
	}
	logger.Infof("Listing posts, offset: %v, limit: %v", qOffset, limit)
	records, err := store.Read("", store.Prefix(key), store.Offset(uint(qOffset)), store.Limit(limit))
	if err != nil {
		return nil, errors.BadRequest("posts.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	posts := buildPostsFromRecords(records)
	return posts, nil
}

func buildPostsFromRecords(records []*store.Record) []*Post {

	if records == nil {
		return nil
	}
	posts := make([]*Post, len(records))
	for i, record := range records {
		//dto proc to handle po to bo
		postRecord := &Post{}
		if err := json.Unmarshal(record.Value, postRecord); err != nil {
			return nil
		}
		posts[i] = postRecord
	}
	return posts

}
