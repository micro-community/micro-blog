package handler

import (
	"context"

	"github.com/micro-community/micro-blog/posts/model"
	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

// Query the posts
func (p *Posts) Query(ctx context.Context, req *pb.QueryRequest, rsp *pb.QueryResponse) error {

	logger.Info("Received Post.Query request")

	if len(req.Id) == 0 {
		return errors.BadRequest("posts.Query.input-check", "ID is missing")
	}

	var err error
	var records []*model.Post

	if len(req.Slug) > 0 { // first to search by slug
		records, err = p.DB.QueryPostBySlug(ctx, req.Slug)
	} else if len(req.Id) > 0 { //then by id
		records, err = p.DB.QueryPostByID(ctx, req.Id)
	} else { //last by timestamp
		records, err = p.DB.QueryPostByTimeStamp(ctx, req.Limit, req.Offset)
	}

	if err != nil {
		return errors.BadRequest("posts.Query.store-read", "Failed to read from db: %v", err.Error())
	}
	// serialize the response list
	rsp.Posts = make([]*pb.Post, len(records))
	for i, record := range records {

		rsp.Posts[i] = &pb.Post{
			Id:      record.ID,
			Title:   record.Title,
			Slug:    record.Slug,
			Content: record.Content,
			Tags:    record.Tags,
		}
	}
	return nil
}
