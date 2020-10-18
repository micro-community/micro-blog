package model

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/micro-community/micro-blog/common/protos/tags"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

//DB to handle DB
type DB struct {
	Tags tags.TagsService
}

//NewService return a model context
func NewService(tagsService tags.TagsService) *DB {
	return &DB{
		Tags: tagsService,
	}
}

//CheckByPostID from store(db,cache etc.)
func (p *DB) CheckByPostID(postID string) (*Post, error) {
	records, err := store.Read(fmt.Sprintf("%v:%v", IDPrefix, postID))
	if err != nil && err != store.ErrNotFound {
		return nil, errors.InternalServerError("posts.Save.store-id-read", "Failed to read post by id: %v", err.Error())
	}

	if len(records) == 0 {
		return nil, nil
	}

	// there is some posts with this id, so we update current post
	record := records[0]
	oldPost := &Post{}
	err = json.Unmarshal(record.Value, oldPost)
	if err != nil {
		return nil, errors.InternalServerError("posts.save.unmarshal", "Failed to unmarshal old post: %v", err.Error())
	}

	return oldPost, nil
}

//CheckBySlug from store(db,cache etc.)
func (p *DB) CheckBySlug(postSlug, oldPostID string) error {
	recordsBySlug, err := store.Read(fmt.Sprintf("%v:%v", SlugPrefix, postSlug))
	if err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("posts.Save.store-read", "Failed to read post by slug: %v", err.Error())
	}

	if len(recordsBySlug) > 0 {
		otherSlugPost := &Post{}
		err := json.Unmarshal(recordsBySlug[0].Value, otherSlugPost)
		if oldPostID != otherSlugPost.ID {
			if err != nil {
				return errors.InternalServerError("posts.Save.slug-unmarshal", "Error un-marshalling other post with same slug: %v", err.Error())
			}
		}
		return errors.BadRequest("posts.Save.slug-check", "An other post with this slug already exists")
	}
	return nil
}

//SavePost to db
func (p *DB) SavePost(ctx context.Context, oldPost, post *Post) error {

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
		Key:   fmt.Sprintf("%v:%v", TimeStampPrefix, math.MaxInt64-post.CreateTimestamp),
		Value: bytes,
	}); err != nil {
		return err
	}

	//this is a new post
	if oldPost == nil {

		var tagNames []string
		for _, tagName := range post.Tags {
			tagNames = append(tagNames, tagName)
		}
		if _, err := p.Tags.Add(ctx, &tags.AddRequest{
			ResourceID: post.ID,
			Type:       TagType,
			Titles:     tagNames,
		}); err != nil {
			return err
		}

		// this is all
		return nil
	}

	//update tags
	return p.diffTags(ctx, post.ID, oldPost.Tags, post.Tags)

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
