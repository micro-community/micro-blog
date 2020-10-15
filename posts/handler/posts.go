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
)

const (
	idPrefix        = "id"
	slugPrefix      = "slug"
	timestampPrefix = "timestamp"
)

//Posts Handler of Blog
type Posts struct{}

//Save a post
func (p *Posts) Save(ctx context.Context, req *pb.SaveRequest, rsp *pb.SaveResponse) error {
	logger.Info("Received Posts.Save request")

	if len(req.Post.Id) == 0 || len(req.Post.Title) == 0 || len(req.Post.Content) == 0 {
		return errors.BadRequest("posts.Save", "ID, title or content is missing")
	}

	// read by parent ID so we can check if it exists without slug changes getting in the way.
	records, err := store.Read(fmt.Sprintf("%v:%v", idPrefix, req.Post.Id))
	if err != nil && err != store.ErrNotFound {
		return err
	}
	postSlug := slug.Make(req.Post.Title)

	// If no existing record is found, create a new one
	if len(records) == 0 {
		return p.savePost(ctx, nil, &model.Post{
			ID:              req.Post.Id,
			Title:           req.Post.Title,
			Content:         req.Post.Content,
			TagNames:        req.Post.TagNames,
			Slug:            postSlug,
			CreateTimestamp: time.Now().Unix(),
		})
	}

	record := records[0]

	oldPost := &model.Post{}
	if err := json.Unmarshal(record.Value, oldPost); err != nil {
		return err
	}

	post := &model.Post{
		ID:              req.Post.Id,
		Title:           req.Post.Title,
		Content:         req.Post.Content,
		Slug:            slug.Make(req.Post.Title),
		TagNames:        req.Post.TagNames,
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
		Key:   fmt.Sprintf("%v:%v", timestampPrefix, math.MaxInt64-post.CreateTimestamp),
		Value: bytes,
	})
}
