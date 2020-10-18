package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

//QueryTagByID from db
func (p *DB) QueryTagByID(ctx context.Context, id string) ([]*Tag, error) {

	key := fmt.Sprintf("%v:%v", SlugPrefix, id)
	logger.Infof("Reading tag by id: %v", id)
	records, err := store.Read("", store.Prefix(key))

	if err != nil {
		return nil, errors.BadRequest("posts.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	posts := buildTagsFromRecords(records)
	return posts, nil
}

//QueryTagBySlug from db
func (p *DB) QueryTagBySlug(ctx context.Context, slug string) ([]*Tag, error) {

	key := fmt.Sprintf("%v:%v", SlugPrefix, slug)
	logger.Infof("Reading tag by slug: %v", slug)
	records, err := store.Read("", store.Prefix(key))

	if err != nil {
		return nil, errors.BadRequest("posts.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	posts := buildTagsFromRecords(records)
	return posts, nil
}

//QueryTagByTimeStamp from db
func (p *DB) QueryTagByTimeStamp(ctx context.Context, qLimit, qOffset int64) ([]*Tag, error) {

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
	posts := buildTagsFromRecords(records)
	return posts, nil
}

func buildTagsFromRecords(records []*store.Record) []*Tag {

	if records == nil {
		return nil
	}

	posts := make([]*Tag, len(records))

	for i, record := range records {
		//dto proc to handle po to bo
		postRecord := &Tag{}
		if err := json.Unmarshal(record.Value, postRecord); err != nil {
			return nil
		}
		posts[i] = postRecord
	}
	return posts

}
