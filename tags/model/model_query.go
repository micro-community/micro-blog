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

	key := fmt.Sprintf("%v:%v", resourcePrefix, id)
	logger.Infof("Reading tag by id: %v", id)
	records, err := store.Read("", store.Prefix(key))

	if err != nil {
		return nil, errors.BadRequest("tags.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	tags := buildTagsFromRecords(records)
	return tags, nil
}

//QueryTagBySlug from db
func (p *DB) QueryTagBySlug(ctx context.Context, slug string) ([]*Tag, error) {

	key := fmt.Sprintf("%v:%v", slugPrefix, slug)
	logger.Infof("Reading tag by slug: %v", slug)
	records, err := store.Read("", store.Prefix(key))

	if err != nil {
		return nil, errors.BadRequest("tags.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	tags := buildTagsFromRecords(records)
	return tags, nil
}

// QueryTagsByType query type
func (p *DB) QueryTagsByType(ctx context.Context, types string) ([]*Tag, error) {

	key := fmt.Sprintf("%v:%v", typePrefix, types)
	logger.Infof("Reading tag by slug: %v", types)
	records, err := store.Read("", store.Prefix(key))

	if err != nil {
		return nil, errors.BadRequest("tags.Query.QueryTagByType", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	tags := buildTagsFromRecords(records)
	return tags, nil
}

//QueryTagByTimeStamp from db
func (p *DB) QueryTagByTimeStamp(ctx context.Context, qLimit, qOffset int64) ([]*Tag, error) {

	key := fmt.Sprintf("%v:", timeStampPrefix)
	var limit uint
	limit = 20 //default if without limition in req
	if qLimit > 0 {
		limit = uint(qLimit)
	}
	logger.Infof("Listing tags, offset: %v, limit: %v", qOffset, limit)
	records, err := store.Read("", store.Prefix(key), store.Offset(uint(qOffset)), store.Limit(limit))
	if err != nil {
		return nil, errors.BadRequest("tags.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	tags := buildTagsFromRecords(records)
	return tags, nil
}

func buildTagsFromRecords(records []*store.Record) []*Tag {

	if records == nil {
		return nil
	}

	tags := make([]*Tag, len(records))

	for i, record := range records {
		//dto proc to handle po to bo
		tagRecord := &Tag{}
		if err := json.Unmarshal(record.Value, tagRecord); err != nil {
			return nil
		}
		tags[i] = tagRecord
	}
	return tags

}
