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

//UpdatePost to db
func (r *Repository) UpdatePost(ctx context.Context, oldPost, post *Post) error {

	bytes, err := json.Marshal(post)
	if err != nil {
		return err
	}

	// Save post by content ID
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", IDPrefix, post.ID),
		Value: bytes,
	}); err != nil {
		return err
	}

	// Delete old by slug index if the slug has changed
	if oldPost != nil && oldPost.Slug != post.Slug {
		if err := store.Delete(fmt.Sprintf("%v:%v", SlugPrefix, oldPost.Slug)); err != nil {
			return err
		}
	}

	// Save post by slug
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", SlugPrefix, post.Slug),
		Value: bytes,
	}); err != nil {
		return err
	}

	// Save post by timeStamp
	if err := store.Write(&store.Record{
		// We revert the timestamp so the order is chronologically reversed
		Key:   fmt.Sprintf("%v:%v", TimeStampPrefix, math.MaxInt64-post.Created),
		Value: bytes,
	}); err != nil {
		return err
	}

	//update tags
	return r.diffTags(ctx, post.ID, oldPost.Tags, post.Tags)
}

//diffTags to update tags
func (r *Repository) diffTags(ctx context.Context, resourceID string, oldTagNames, newTagNames []string) error {

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
		_, err := r.Tags.Remove(ctx, &tags.RemoveRequest{
			ResourceID: resourceID,
			Type:       TagType,
			Titles:     tags2remove,
		})
		if err != nil {
			logger.Errorf("Error decreasing count for tag '%v' with type '%v' for Post '%v'", tags2remove, TagType, resourceID)
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

		_, err := r.Tags.Add(ctx, &tags.AddRequest{
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
