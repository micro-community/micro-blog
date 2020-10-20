package model

import (
	"context"
	"encoding/json"
	"fmt"

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

	//this is a new tag
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

	return p.saveTag(tag)
}

//IncreseTagCount increase  a tag count basing on all post with same tag
func (p *DB) IncreseTagCount(resourceID string, tag *Tag) error {

	//tagCountPrefix:tagslug:resourceID ,add resource in some slug
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v%v", tagCountPrefix, tag.Slug, resourceID),
		Value: nil,
	}); err != nil {
		return err
	}
	oldTagCount := tag.Count
	// get tag count
	recs, err := store.List(store.Prefix(fmt.Sprintf("%v:%v", tagCountPrefix, tag.Slug)), store.Limit(1000))
	if err != nil {
		return err
	}

	tag.Count = int64(len(recs))

	if tag.Count == oldTagCount {
		return fmt.Errorf("Tag count for tag %v is unchanged, was: %v, now: %v", tag.Slug, oldTagCount, tag.Count)
	}
	tagJSON, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	err = store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v:%v", resourcePrefix, resourceID, tag.Slug),
		Value: tagJSON,
	})
	if err != nil {
		return err
	}

	return nil
}
