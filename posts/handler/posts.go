package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/gosimple/slug"
	"github.com/micro-community/micro-blog/posts/model"
	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	tagProtos "github.com/micro-community/micro-blog/tags/proto"
)

const (
	tagType         = "post-tag"
	slugPrefix      = "slug"
	idPrefix        = "id"
	timeStampPrefix = "timestamp"
)

//Posts Handler of Blog
type Posts struct {
	Tags tagProtos.TagsService
}

//Save a post
func (p *Posts) Save(ctx context.Context, req *pb.SaveRequest, rsp *pb.SaveResponse) error {
	logger.Info("Received Posts.Save request")

	if len(req.Id) == 0 || len(req.Title) == 0 || len(req.Content) == 0 {
		return errors.BadRequest("posts.Save", "ID, title or content is missing")
	}

	// read by parent ID so we can check if it exists without slug changes getting in the way.
	records, err := store.Read(fmt.Sprintf("%v:%v", idPrefix, req.Id))
	if err != nil && err != store.ErrNotFound {
		return err
	}
	postSlug := slug.Make(req.Title)

	// If no existing record is found, create a new one
	if len(records) == 0 {
		return p.savePost(ctx, nil, &model.Post{
			ID:              req.Id,
			Title:           req.Title,
			Content:         req.Content,
			Tags:            req.Tags,
			Slug:            postSlug,
			CreateTimestamp: time.Now().Unix(),
		})
	}

	record := records[0]
	oldPost := &model.Post{}
	err = json.Unmarshal(record.Value, oldPost)
	if err != nil {
		return errors.InternalServerError("posts.save.unmarshal", "Failed to unmarshal old post: %v", err.Error())
	}
	post := &model.Post{
		ID:              req.Id,
		Title:           req.Title,
		Content:         req.Content,
		Slug:            slug.Make(req.Title),
		Tags:            req.Tags,
		CreateTimestamp: time.Now().Unix(),
		UpdateTimestamp: time.Now().Unix(),
	}

	// Check if slug exists
	recordsBySlug, err := store.Read(fmt.Sprintf("%v:%v", slugPrefix, postSlug))
	if err != nil && err != store.ErrNotFound {
		return err
	}

	otherSlugPost := &model.Post{}
	if err := json.Unmarshal(record.Value, otherSlugPost); err != nil {
		return err
	}
	if len(recordsBySlug) > 0 && oldPost.ID != otherSlugPost.ID {
		return errors.BadRequest("posts.Save", "An other post with this slug already exists")
	}

	return p.savePost(ctx, oldPost, post)

}

func (p *Posts) savePost(ctx context.Context, oldPost, post *model.Post) error {

	bytes, err := json.Marshal(post)
	if err != nil {
		return err
	}

	// Save post by content ID
	record := &store.Record{
		Key:   fmt.Sprintf("%v:%v", idPrefix, post.ID),
		Value: bytes,
	}
	if err := store.Write(record); err != nil {
		return err
	}

	// Delete old by slug index if the slug has changed
	if oldPost.Slug != post.Slug {
		if err := store.Delete(fmt.Sprintf("%v:%v", slugPrefix, post.Slug)); err != nil {
			return err
		}
	}

	// Save post by slug
	slugRecord := &store.Record{
		Key:   fmt.Sprintf("%v:%v", slugPrefix, post.Slug),
		Value: bytes,
	}
	if err := store.Write(slugRecord); err != nil {
		return err
	}

	// Save post by timeStamp
	return store.Write(&store.Record{
		// We revert the timestamp so the order is chronologically reversed
		Key:   fmt.Sprintf("%v:%v", timeStampPrefix, math.MaxInt64-post.CreateTimestamp),
		Value: bytes,
	})
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
