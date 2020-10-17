package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/gosimple/slug"
	"github.com/micro-community/micro-blog/common/protos/tags"
	"github.com/micro-community/micro-blog/posts/model"
	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

const (
	tagType         = "post-tag"
	slugPrefix      = "slug"
	idPrefix        = "id"
	timeStampPrefix = "timestamp"
)

//Posts Handler of Blog
type Posts struct {
	Tags tags.TagsService
}

//Save a post
func (p *Posts) Save(ctx context.Context, req *pb.SaveRequest, rsp *pb.SaveResponse) error {
	logger.Info("Received Posts.Save request")

	if len(req.Id) == 0 {
		return errors.BadRequest("posts.Save.input-check", "ID is missing")
	}

	// read by post id.
	records, err := store.Read(fmt.Sprintf("%v:%v", idPrefix, req.Id))
	if err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("posts.Save.store-id-read", "Failed to read post by id: %v", err.Error())
	}

	postSlug := slug.Make(req.Title)

	// If no existing record is found, create a new one
	if len(records) == 0 {
		post := &model.Post{
			ID:              req.Id,
			Title:           req.Title,
			Content:         req.Content,
			Tags:            req.Tags,
			Slug:            postSlug,
			CreateTimestamp: time.Now().Unix(),
		}

		err := p.savePost(ctx, nil, post)
		if err != nil {
			return errors.InternalServerError("posts.save.post-save", "Failed to save new post: %v", err.Error())
		}
		return nil
	}

	// there is some posts with this id, so we update current post
	record := records[0]
	oldPost := &model.Post{}
	err = json.Unmarshal(record.Value, oldPost)
	if err != nil {
		return errors.InternalServerError("posts.save.unmarshal", "Failed to unmarshal old post: %v", err.Error())
	}
	//new post from old
	post := &model.Post{
		ID:              req.Id,
		Title:           oldPost.Title,
		Content:         oldPost.Content,
		Slug:            oldPost.Slug,
		Tags:            oldPost.Tags,
		CreateTimestamp: oldPost.CreateTimestamp,
		UpdateTimestamp: time.Now().Unix(),
	}

	//update article content
	if len(req.Title) > 0 {
		post.Title = req.Title
		post.Slug = slug.Make(post.Title)
	}
	if len(req.Slug) > 0 {
		post.Slug = req.Slug
	}
	if len(req.Content) > 0 {
		post.Content = req.Content
	}
	if len(req.Tags) > 0 {
		//update :only remove the tags
		if len(req.Tags) == 1 && req.Tags[0] == "" {
			post.Tags = []string{}
		} else {
			post.Tags = req.Tags
		}
	}

	// Check if slug exists
	recordsBySlug, err := store.Read(fmt.Sprintf("%v:%v", slugPrefix, postSlug))
	if err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("posts.Save.store-read", "Failed to read post by slug: %v", err.Error())
	}

	if len(recordsBySlug) > 0 {
		otherSlugPost := &model.Post{}
		err := json.Unmarshal(recordsBySlug[0].Value, otherSlugPost)
		if oldPost.ID != otherSlugPost.ID {
			if err != nil {
				return errors.InternalServerError("posts.Save.slug-unmarshal", "Error unmarshaling other post with same slug: %v", err.Error())
			}
		}
		return errors.BadRequest("posts.Save.slug-check", "An other post with this slug already exists")
	}

	return p.savePost(ctx, oldPost, post)

}

func (p *Posts) savePost(ctx context.Context, oldPost, post *model.Post) error {

	bytes, err := json.Marshal(post)
	if err != nil {
		return err
	}

	// Save post by content ID
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", idPrefix, post.ID),
		Value: bytes,
	}); err != nil {
		return err
	}

	// Delete old by slug index if the slug has changed
	if oldPost != nil && oldPost.Slug != post.Slug {
		if err := store.Delete(fmt.Sprintf("%v:%v", slugPrefix, oldPost.Slug)); err != nil {
			return err
		}
	}

	// Save post by slug
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", slugPrefix, post.Slug),
		Value: bytes,
	}); err != nil {
		return err
	}

	// Save post by timeStamp
	if err := store.Write(&store.Record{
		// We revert the timestamp so the order is chronologically reversed
		Key:   fmt.Sprintf("%v:%v", timeStampPrefix, math.MaxInt64-post.CreateTimestamp),
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
			Type:       tagType,
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
func (p *Posts) diffTags(ctx context.Context, resourceID string, oldTagNames, newTagNames []string) error {

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
			Type:       tagType,
			Titles:     tags2remove,
		})
		if err != nil {
			logger.Errorf("Error decreasing count for tag '%v' with type '%v' for Post '%v'", tags2remove, tagType, resourceID)
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
			Type:       tagType,
			Titles:     tags2add,
		})

		if err != nil {
			logger.Errorf("Error increasing count for tag '%v' with type '%v' for parent '%v': %v", tags2add, tagType, resourceID, err)
		}

	}

	return nil

}
