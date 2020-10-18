package model

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/micro-community/micro-blog/common/protos/tags"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
)

//CheckByTagID from store(db,cache etc.)
func (p *DB) CheckByTagID(postID string) (*Tag, error) {
	records, err := store.Read(fmt.Sprintf("%v:%v", IDPrefix, postID))
	if err != nil && err != store.ErrNotFound {
		return nil, errors.InternalServerError("posts.Save.store-id-read", "Failed to check tag by id: %v", err.Error())
	}

	if len(records) == 0 {
		return nil, nil
	}

	// there is some posts with this id, so we update current tag
	record := records[0]
	oldTag := &Tag{}
	err = json.Unmarshal(record.Value, oldTag)
	if err != nil {
		return nil, errors.InternalServerError("posts.save.unmarshal", "Failed to unmarshal old tag: %v", err.Error())
	}

	return oldTag, nil
}

//CheckBySlug from store(db,cache etc.)
func (p *DB) CheckBySlug(postSlug, oldTagID string) error {
	recordsBySlug, err := store.Read(fmt.Sprintf("%v:%v", SlugPrefix, postSlug))
	if err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("posts.Save.store-read", "Failed to read tag by slug: %v", err.Error())
	}

	if len(recordsBySlug) > 0 {
		otherSlugTag := &Tag{}
		err := json.Unmarshal(recordsBySlug[0].Value, otherSlugTag)
		if oldTagID != otherSlugTag.ID {
			if err != nil {
				return errors.InternalServerError("posts.Save.slug-unmarshal", "Error un-marshalling other tag with same slug: %v", err.Error())
			}
		}
		return errors.BadRequest("posts.Save.slug-check", "An other tag with this slug already exists")
	}
	return nil
}

//CreateTag to db
func (p *DB) CreateTag(ctx context.Context, tag *Tag) error {

	bytes, err := json.Marshal(tag)
	if err != nil {
		return err
	}

	// Save tag by content ID
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", IDPrefix, tag.ID),
		Value: bytes,
	}); err != nil {
		return err
	}

	// Save tag by slug
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", SlugPrefix, tag.Slug),
		Value: bytes,
	}); err != nil {
		return err
	}

	// Save tag by timeStamp
	if err := store.Write(&store.Record{
		// We revert the timestamp so the order is chronologically reversed
		Key:   fmt.Sprintf("%v:%v", TimeStampPrefix, math.MaxInt64-tag.CreateTimestamp),
		Value: bytes,
	}); err != nil {
		return err
	}

	//Add New Tags
	var tagNames []string
	for _, tagName := range tag.Tags {
		tagNames = append(tagNames, tagName)
	}
	if _, err := p.Tags.Add(ctx, &tags.AddRequest{
		ResourceID: tag.ID,
		Type:       TagType,
		Titles:     tagNames,
	}); err != nil {
		return err
	}

	return nil
}
