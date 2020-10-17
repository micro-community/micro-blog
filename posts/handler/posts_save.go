package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/micro-community/micro-blog/common/protos/tags"
	"github.com/micro-community/micro-blog/posts/model"
	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

//Posts Handler of Blog
type Posts struct {
	Tags tags.TagsService
	DB   *model.DB
}

//Save a post
func (p *Posts) Save(ctx context.Context, req *pb.SaveRequest, rsp *pb.SaveResponse) error {
	logger.Info("Received Posts.Save request")

	if len(req.Id) == 0 {
		return errors.BadRequest("posts.Save.input-check", "ID is missing")
	}

	// read by post id.
	records, err := store.Read(fmt.Sprintf("%v:%v", model.IDPrefix, req.Id))
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
	recordsBySlug, err := store.Read(fmt.Sprintf("%v:%v", model.SlugPrefix, postSlug))
	if err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("posts.Save.store-read", "Failed to read post by slug: %v", err.Error())
	}

	if len(recordsBySlug) > 0 {
		otherSlugPost := &model.Post{}
		err := json.Unmarshal(recordsBySlug[0].Value, otherSlugPost)
		if oldPost.ID != otherSlugPost.ID {
			if err != nil {
				return errors.InternalServerError("posts.Save.slug-unmarshal", "Error un-marshalling other post with same slug: %v", err.Error())
			}
		}
		return errors.BadRequest("posts.Save.slug-check", "An other post with this slug already exists")
	}

	return p.savePost(ctx, oldPost, post)

}
