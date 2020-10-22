package handler

import (
	"context"

	"github.com/gosimple/slug"
	pb "github.com/micro-community/micro-blog/tags/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

// Remove a tag
func (p *Tags) Remove(ctx context.Context, req *pb.RemoveRequest, rsp *pb.RemoveResponse) error {

	logger.Info("Received Tag.Delete request")
	if len(req.ResourceID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.Delete.input-check", "ID or Type is missing")
	}

	for _, title := range req.GetTitles() {

		tag, err := p.DB.CheckBySlug(slug.Make(title))
		if err != nil {
			rsp.Results = append(rsp.Results, false)
			return err
		}
		if tag == nil {
			rsp.Results = append(rsp.Results, true)
			//return fmt.Errorf("Tag with ID %v not found", req.Id)
		}
		// Delete by ID
		if err = p.DB.DeleteTagByResourceID(ctx, req.GetResourceID(), tag); err != nil {
			rsp.Results = append(rsp.Results, false)
			return err
		}
		rsp.Results = append(rsp.Results, true)
	}

	return nil

}
