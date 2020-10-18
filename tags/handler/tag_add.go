package handler

import (
	"context"
	"time"

	"github.com/gosimple/slug"
	"github.com/micro-community/micro-blog/tags/model"
	pb "github.com/micro-community/micro-blog/tags/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

//Save a tag
func (p *Tags) Save(ctx context.Context, req *pb.AddRequest, rsp *pb.AddResponse) error {
	logger.Info("Received Tags.Save request")

	if len(req.Id) == 0 {
		return errors.BadRequest("posts.Save.input-check", "ID is missing")
	}

	oldTag, err := p.DB.CheckByTagID(req.Id)

	if err != nil {
		return err
	}

	//find no old tag
	if oldTag == nil {
		newTag := &model.Tag{
			ID:              req.Id,
			Title:           req.Title,
			Content:         req.Content,
			Tags:            req.Tags,
			Slug:            slug.Make(req.Title),
			CreateTimestamp: time.Now().Unix(),
		}
		if err := p.DB.CreateTag(ctx, newTag); err != nil {
			return errors.InternalServerError("posts.save.tag-save", "Failed to save new tag: %v", err.Error())
		}
		return nil

	}

	//new tag content from old
	tag := &model.Tag{
		ID:              req.Id,
		Title:           oldTag.Title,
		Content:         oldTag.Content,
		Slug:            oldTag.Slug,
		Tags:            oldTag.Tags,
		CreateTimestamp: oldTag.CreateTimestamp,
		UpdateTimestamp: time.Now().Unix(),
	}

	//update article content
	if len(req.Title) > 0 {
		tag.Title = req.Title
		tag.Slug = slug.Make(tag.Title)
	}
	if len(req.Slug) > 0 {
		tag.Slug = req.Slug
	}
	if len(req.Content) > 0 {
		tag.Content = req.Content
	}
	if len(req.Tags) > 0 {
		//update :only remove the tags
		if len(req.Tags) == 1 && req.Tags[0] == "" {
			tag.Tags = []string{}
		} else {
			tag.Tags = req.Tags
		}
	}

	// Check if slug exists
	postSlug := slug.Make(req.Title)

	if err := p.DB.CheckBySlug(postSlug, oldTag.ID); err != nil {
		return err
	}

	return p.DB.UpdateTag(ctx, oldTag, tag)

}
