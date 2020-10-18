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

//CheckByPostID from store(db,cache etc.)
func (p *DB) CheckByPostID(postID string) (*Post, error) {
	records, err := store.Read(fmt.Sprintf("%v:%v", IDPrefix, postID))
	if err != nil && err != store.ErrNotFound {
		return nil, errors.InternalServerError("posts.Save.store-id-read", "Failed to check post by id: %v", err.Error())
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

//CreatePost to db
func (p *DB) CreatePost(ctx context.Context, post *Post) error {

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

	//Add New Tags
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

	return nil
}
