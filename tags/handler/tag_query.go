package handler

import (
	"context"

	"github.com/micro-community/micro-blog/tags/model"
	pb "github.com/micro-community/micro-blog/tags/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

// Query the posts
func (p *Tags) Query(ctx context.Context, req *pb.QueryRequest, rsp *pb.QueryResponse) error {

	logger.Info("Received Tag.Query request")

	if len(req.Id) == 0 {
		return errors.BadRequest("posts.Query.input-check", "ID is missing")
	}

	var err error
	var records []*model.Tag

	if len(req.Slug) > 0 { // first to search by slug
		records, err = p.DB.QueryTagBySlug(ctx, req.Slug)
	} else if len(req.Id) > 0 { //then by id
		records, err = p.DB.QueryTagByID(ctx, req.Id)
	} else { //last by timestamp
		records, err = p.DB.QueryTagByTimeStamp(ctx, req.Limit, req.Offset)
	}

	if err != nil {
		return errors.BadRequest("posts.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	rsp.Tags = make([]*pb.Tag, len(records))
	for i, record := range records {

		rsp.Tags[i] = &pb.Tag{
			Id:      record.ID,
			Title:   record.Title,
			Slug:    record.Slug,
			Content: record.Content,
			Tags:    record.Tags,
		}
	}
	return nil
}
