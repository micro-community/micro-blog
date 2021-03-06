package handler

import (
	"context"
	"time"

	"github.com/gosimple/slug"
	"github.com/micro-community/micro-blog/posts/model"
	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

//Save a post
func (p *Posts) Save(ctx context.Context, req *pb.SaveRequest, rsp *pb.SaveResponse) error {
	logger.Info("Received Posts.Save request")

	if len(req.Id) == 0 {
		//this is a new post
		return errors.BadRequest("posts.Save.input-check", "ID is missing")
	}

	oldPost, err := p.Repository.CheckByPostID(req.Id)

	if err != nil {
		return err
	}

	//find no old post
	if oldPost == nil {
		newPost := &model.Post{
			ID:      req.Id,
			Title:   req.Title,
			Content: req.Content,
			Tags:    req.Tags,
			Slug:    slug.Make(req.Title),
			Created: time.Now().Unix(),
		}
		if err := p.Repository.CreatePost(ctx, newPost); err != nil {
			return errors.InternalServerError("posts.save.post-save", "Failed to save new post: %v", err.Error())
		}
		return nil

	}

	//new post content from old
	post := &model.Post{
		ID:      req.Id,
		Title:   oldPost.Title,
		Content: oldPost.Content,
		Slug:    oldPost.Slug,
		Tags:    oldPost.Tags,
		Created: oldPost.Created,
		Updated: time.Now().Unix(),
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
	postSlug := slug.Make(req.Title)

	if err := p.Repository.CheckBySlug(postSlug, oldPost.ID); err != nil {
		return err
	}

	return p.Repository.UpdatePost(ctx, oldPost, post)

}
