package handler

import (
	"context"

	"github.com/micro-community/micro-blog/tags/model"
	pb "github.com/micro-community/micro-blog/tags/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

// List the tags
func (p *Tags) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {

	logger.Info("Received Tag.List request")

	var err error
	var records []*model.Tag

	if len(req.ResourceID) > 0 { // first to search by slug
		records, err = p.DB.QueryTagByID(ctx, req.ResourceID)
	} else if len(req.Type) > 0 { //last by timestamp
		records, err = p.DB.QueryTagsByType(ctx, req.Type)
	}

	if err != nil {
		return errors.BadRequest("tags.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	rsp.Tags = make([]*pb.Tag, len(records))
	for i, record := range records {

		rsp.Tags[i] = &pb.Tag{
			Title: record.Title,
			Slug:  record.Slug,
			Type:  record.Type,
			Count: record.Count,
		}
	}
	return nil
}
