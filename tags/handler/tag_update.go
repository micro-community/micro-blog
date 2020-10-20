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

	for _, title := range req.Titles {

		tagSlug := slug.Make(title)
		oldTag, err := t.DB.CheckBySlug(tagSlug)
		if err != nil {
			return err
		}
		//find no old tag
		if oldTag == nil {
			rsp.Results = append(rsp.Results, false)
			//return fmt.Errorf("Tag with slug '%v' not found, nothing to update", tagSlug)
			continue
		}

		//exist the target
		newTag := &model.Tag{
			Type:            req.Type,
			Title:           title,
			Count:           oldTag.Count,
			Slug:            slug.Make(title),
			CreateTimestamp: time.Now().Unix(),
		}
		if err := t.DB.UpdateTag(ctx, oldTag, newTag); err != nil {
			rsp.Results = append(rsp.Results, false)
			//return errors.InternalServerError("tags.Update.UpdateTag", "Failed to Update new tag: %v", err.Error())
		}
		rsp.Results = append(rsp.Results, true)
	}
	return nil
}
