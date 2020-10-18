package handler

import (
	"context"
	"fmt"

	pb "github.com/micro-community/micro-blog/tags/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

// Delete a tag
func (p *Tags) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {

	logger.Info("Received Tag.Delete request")
	if len(req.Id) == 0 {
		return errors.BadRequest("posts.Save.input-check", "ID is missing")
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
	if err = p.DB.DeleteTagBySlug(ctx, tag.Slug); err != nil {
		return err
	}

	// Delete by slug
	return p.DB.DeleteTagByTimeStamp(ctx, tag.CreateTimestamp)
}
