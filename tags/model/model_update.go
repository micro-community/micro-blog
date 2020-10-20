package model

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/micro/micro/v3/service/store"
)

//UpdateTag to db
func (p *DB) UpdateTag(ctx context.Context, oldTag, tag *Tag) error {

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

	// Delete old by slug index if the slug has changed
	if oldTag != nil && oldTag.Slug != tag.Slug {
		if err := store.Delete(fmt.Sprintf("%v:%v", slugPrefix, oldTag.Slug)); err != nil {
			return err
		}
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

	//update tags
	return nil
}
