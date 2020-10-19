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

//Add a tag
func (t *Tags) Add(ctx context.Context, req *pb.AddRequest, rsp *pb.AddResponse) error {
	logger.Info("Received Tags.Save request")

	if len(req.ResourceID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.Add.input-check", "ID or Type is missing")
	}

	for _, title := range req.Titles {

		existTag, err := t.DB.CheckBySlug(title)

		if err != nil {
			rsp.Results = append(rsp.Results, false)
			return err
		}
		//already exist
		if existTag != nil {
			rsp.Results = append(rsp.Results, true)
			continue
		}
		//find no old tag
		newTag := &model.Tag{
			Title:           title,
			Slug:            slug.Make(title),
			CreateTimestamp: time.Now().Unix(),
		}
		if err := t.DB.CreateTag(ctx, newTag); err != nil {
			rsp.Results = append(rsp.Results, false)
			return errors.InternalServerError("tags.save.tag-save", "Failed to save new tag: %v", err.Error())

		}
		rsp.Results = append(rsp.Results, true)

	}

	return nil
}
