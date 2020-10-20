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

	bytes, err := json.Marshal(tag)
	if err != nil {
		return err
	}

	// // Save tag by content ID
	// if err := store.Write(&store.Record{
	// 	Key:   fmt.Sprintf("%v:%v", idPrefix, tag.ResourceID),
	// 	Value: bytes,
	// }); err != nil {
	// 	return err
	// }

	// Save tag by slug
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", slugPrefix, tag.Slug),
		Value: bytes,
	}); err != nil {
		return err
	}

	err = store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v:%v", tagCountPrefix, tag.Type, tag.Slug),
		Value: bytes,
	})
	if err != nil {
		return err
	}

	return nil
}

//IncresePostTagCount increase :1 ,tags count for post ,2,total tags count
func (p *DB) IncresePostTagCount(resourceID string, tag *Tag) error {

	//tagCountPrefix:tagslug:resourceID
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
	return nil
}
