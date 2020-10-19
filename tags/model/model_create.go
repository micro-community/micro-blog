package model

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
)

//CheckByResourceID from store(db,cache etc.)
func (p *DB) CheckByResourceID(resourceID string) (*Tag, error) {
	records, err := store.Read(fmt.Sprintf("%v:%v", idPrefix, resourceID))
	if err != nil && err != store.ErrNotFound {
		return nil, errors.InternalServerError("tags.Save.store-id-read", "Failed to check tag by id: %v", err.Error())
	}

	if len(records) == 0 {
		return nil, nil
	}

	// there is some tags with this id, so we update current tag
	oldTag := &Tag{}
	err = json.Unmarshal(records[0].Value, oldTag)
	if err != nil {
		return nil, errors.InternalServerError("tags.save.unmarshal", "Failed to unmarshal old tag: %v", err.Error())
	}

	return oldTag, nil
}

//CheckBySlug from store(db,cache etc.)
func (p *DB) CheckBySlug(tagSlug string) (*Tag, error) {

	recordsBySlug, err := store.Read(fmt.Sprintf("%v:%v", slugPrefix, tagSlug))
	if err != nil && err != store.ErrNotFound {
		return nil, errors.InternalServerError("tags.Save.store-read", "Failed to read tag by slug: %v", err.Error())
	}

	if len(recordsBySlug) == 0 {
		return nil, nil
	}
	retrivedTag := &Tag{}
	if len(recordsBySlug) > 0 {
		if err := json.Unmarshal(recordsBySlug[0].Value, retrivedTag); err != nil {
			return nil, errors.InternalServerError("tags.Save.slug-unmarshal", "Error un-marshalling other tag with same slug: %v", err.Error())
		}
		return nil, errors.BadRequest("tags.Save.slug-check", "An other tag with this slug already exists")
	}
	return retrivedTag, nil
}

//CreateTag to db
func (p *DB) CreateTag(ctx context.Context, tag *Tag) error {

	bytes, err := json.Marshal(tag)
	if err != nil {
		return err
	}

	// Save tag by content ID
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", idPrefix, tag.ID),
		Value: bytes,
	}); err != nil {
		return err
	}

	// Save tag by slug
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", slugPrefix, tag.Slug),
		Value: bytes,
	}); err != nil {
		return err
	}

	// Save tag by timeStamp
	if err := store.Write(&store.Record{
		// We revert the timestamp so the order is chronologically reversed
		Key:   fmt.Sprintf("%v:%v", timeStampPrefix, math.MaxInt64-tag.CreateTimestamp),
		Value: bytes,
	}); err != nil {
		return err
	}

	return nil
}

func (p *DB) saveTag(tag *Tag) error {

	key := fmt.Sprintf("%v:%v", slugPrefix, tag.Slug)
	typeKey := fmt.Sprintf("%v:%v:%v", typePrefix, tag.Type, tag.Slug)

	bytes, err := json.Marshal(tag)
	if err != nil {
		return err
	}

	// write resourceId:slug to enable prefix listing based on type
	err = store.Write(&store.Record{
		Key:   key,
		Value: bytes,
	})
	if err != nil {
		return err
	}
	return store.Write(&store.Record{
		Key:   typeKey,
		Value: bytes,
	})
}
