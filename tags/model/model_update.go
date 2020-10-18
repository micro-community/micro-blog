package model

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/micro-community/micro-blog/common/protos/tags"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

//UpdateTag to db
func (p *DB) UpdateTag(ctx context.Context, oldTag, tag *Tag) error {

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

	// Delete old by slug index if the slug has changed
	if oldTag != nil && oldTag.Slug != tag.Slug {
		if err := store.Delete(fmt.Sprintf("%v:%v", SlugPrefix, oldTag.Slug)); err != nil {
			return err
		}
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

	//update tags
	return p.diffTags(ctx, tag.ID, oldTag.Tags, tag.Tags)
}

//diffTags to update tags
func (p *DB) diffTags(ctx context.Context, resourceID string, oldTagNames, newTagNames []string) error {

	oldTags := map[string]struct{}{}
	for _, v := range oldTagNames {
		oldTags[v] = struct{}{}
	}

	newTags := map[string]struct{}{}
	for _, v := range newTagNames {
		newTags[v] = struct{}{}
	}

	//find removed tags
	var tags2remove []string
	for i := range oldTags {
		_, stillThere := newTags[i]
		if !stillThere {
			tags2remove = append(tags2remove, i)
		}
	}

	if len(tags2remove) > 0 {
		_, err := p.Tags.Remove(ctx, &tags.RemoveRequest{
			ResourceID: resourceID,
			Type:       TagType,
			Titles:     tags2remove,
		})
		if err != nil {
			logger.Errorf("Error decreasing count for tag '%v' with type '%v' for Tag '%v'", tags2remove, TagType, resourceID)
		}

	}

	//find added tags
	var tags2add []string
	for i := range newTags {
		_, exist := oldTags[i]
		if !exist {
			tags2add = append(tags2add, i)
		}
	}

	if len(tags2add) > 0 {

		_, err := p.Tags.Add(ctx, &tags.AddRequest{
			ResourceID: resourceID,
			Type:       TagType,
			Titles:     tags2add,
		})

		if err != nil {
			logger.Errorf("Error increasing count for tag '%v' with type '%v' for parent '%v': %v", tags2add, TagType, resourceID, err)
		}

	}

	return nil

}
