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

//Update a tag
func (t *Tags) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	logger.Info("Received Tags.Update request")

	if len(req.Titles) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.Add.input-check", "Titles or Type is missing")
	}

	oldTag, err := p.DB.CheckByTagID(req.Titles)

	if err != nil {
		return err
	}

	//find no old tag
	if oldTag == nil {
		newTag := &model.Tag{
			Title:           req.Title,
			Slug:            slug.Make(req.Title),
			CreateTimestamp: time.Now().Unix(),
		}
		if err := t.DB.CreateTag(ctx, newTag); err != nil {
			return errors.InternalServerError("tags.save.tag-save", "Failed to save new tag: %v", err.Error())
		}
		return nil

	}

	//new tag content from old
	tag := &model.Tag{
		ID:              req.ResourceID,
		Title:           oldTag.Title,
		Slug:            oldTag.Slug,
		CreateTimestamp: oldTag.CreateTimestamp,
		UpdateTimestamp: time.Now().Unix(),
	}

	//update article content

	if len(req.Slug) > 0 {
		tag.Slug = req.Slug
	}

	// Check if slug exists
	tagSlug := slug.Make(req.Titles)

	if err := t.DB.CheckBySlug(tagSlug, oldTag.ID); err != nil {
		return err
	}

	return t.DB.UpdateTag(ctx, oldTag, tag)

}
