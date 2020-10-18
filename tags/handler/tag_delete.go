package handler

import (
	"context"
	"fmt"

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
	tag, err := p.DB.CheckByTagID(req.Id)

	if err != nil {
		return err
	}

	if tag == nil {
		return fmt.Errorf("Tag with ID %v not found", req.Id)
	}

	// Delete by ID
	if err = p.DB.DeleteTagByID(ctx, tag.ID); err != nil {
		return err
	}

	// Delete by slug
	return p.DB.DeleteTagBySlug(ctx, tag.Slug)

}
