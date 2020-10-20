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

//Add a tag for a post
func (t *Tags) Add(ctx context.Context, req *pb.AddRequest, rsp *pb.AddResponse) error {
	logger.Info("Received Tags.Add request")

	if len(req.ResourceID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.Add.input-check", "ID or Type is missing")
	}
	for _, title := range req.Titles {

		tagSlug := slug.Make(title)
		existTag, err := t.DB.CheckBySlug(tagSlug)

		if err != nil {
			rsp.Results = append(rsp.Results, false)
			return err
		}
		tag := existTag
		//no exist
		if tag == nil {
			//find no old tag
			tag = &model.Tag{
				Title:           title,
				Type:            req.Type,
				Slug:            slug.Make(title),
				CreateTimestamp: time.Now().Unix(),
			}
		}

		//Create tag for a post
		if err := t.DB.CreateTag(ctx, tag); err != nil {
			rsp.Results = append(rsp.Results, false)
			return errors.InternalServerError("tags.Add.tag-create", "Failed to create new tag: %v", err.Error())

		}
		rsp.Results = append(rsp.Results, true)
		t.DB.IncreseTagCount(req.ResourceID, tag)
	}

	return nil
}
