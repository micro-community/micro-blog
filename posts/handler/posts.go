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

		for _, tagName := range post.Tags {

			if _, err := p.Tags.Add(ctx, &tags.AddRequest{
				ResourceID: post.ID,
				Type:       tagType,
				Title:      tagName,
			}); err != nil {
				return err
			}

		}
		// this is all
		return nil
	}

	//update tags
	return p.diffTags(ctx, post.ID, oldPost.Tags, post.Tags)

}

//diffTags to update tags
func (p *Posts) diffTags(ctx context.Context, parentID string, oldTagNames, newTagNames []string) error {

	oldTags := map[string]struct{}{}
	for _, v := range oldTagNames {
		oldTags[v] = struct{}{}
	}

	newTags := map[string]struct{}{}
	for _, v := range newTagNames {
		newTags[v] = struct{}{}
	}

	return nil

}

// Query the posts
func (p *Posts) Query(ctx context.Context, req *pb.QueryRequest, rsp *pb.QueryResponse) error {
	var records []*store.Record
	var err error
	if len(req.Slug) > 0 {
		key := fmt.Sprintf("%v:%v", slugPrefix, req.Slug)
		logger.Infof("Reading post by slug: %v", req.Slug)
		records, err = store.Read("", store.Prefix(key))
	} else {
		key := fmt.Sprintf("%v:", timeStampPrefix)
		var limit uint
		limit = 20
		if req.Limit > 0 {
			limit = uint(req.Limit)
		}
		logger.Infof("Listing posts, offset: %v, limit: %v", req.Offset, limit)
		records, err = store.Read("", store.Prefix(key),
			store.Offset(uint(req.Offset)),
			store.Limit(limit))
	}

	if err != nil {
		return errors.BadRequest("posts.query.store-read", "Failed to read from store: %v", err.Error())
	}
	// serialize the response
	rsp.Posts = make([]*pb.Post, len(records))
	for i, record := range records {
		postRecord := &model.Post{}
		if err := json.Unmarshal(record.Value, postRecord); err != nil {
			return err
		}

		rsp.Posts[i] = &pb.Post{
			Id:      postRecord.ID,
			Title:   postRecord.Title,
			Slug:    postRecord.Slug,
			Content: postRecord.Content,
			Tags:    postRecord.Tags,
		}
	}
	return nil
}

// Delete a post
func (p *Posts) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	records, err := store.Read(fmt.Sprintf("%v:%v", idPrefix, req.Id))
	if err == store.ErrNotFound {
		return errors.NotFound("posts.Delete", "Post not found")
	} else if err != nil {
		return err
	}

	post := &model.Post{}
	if err := json.Unmarshal(records[0].Value, post); err != nil {
		return err
	}

	// Delete by ID
	if err = store.Delete(fmt.Sprintf("%v:%v", idPrefix, post.ID)); err != nil {
		return err
	}

	// Delete by slug
	if err := store.Delete(fmt.Sprintf("%v:%v", slugPrefix, post.Slug)); err != nil {
		return err
	}

	// Delete by timeStamp
	return store.Delete(fmt.Sprintf("%v:%v", timeStampPrefix, post.CreateTimestamp))
}
