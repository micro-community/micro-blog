package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro-community/micro-blog/posts/model"
	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

// Query the posts
func (p *Posts) Query(ctx context.Context, req *pb.QueryRequest, rsp *pb.QueryResponse) error {

	logger.Info("Received Post.Query request")

	var records []*store.Record
	var err error

	if len(req.Slug) > 0 { // first to search by slug
		key := fmt.Sprintf("%v:%v", slugPrefix, req.Slug)
		logger.Infof("Reading post by slug: %v", req.Slug)
		records, err = store.Read("", store.Prefix(key))
	} else if len(req.Id) > 0 { //then by id
		key := fmt.Sprintf("%v:%v", idPrefix, req.Id)
		logger.Infof("Reading post by id: %v", req.Id)
		records, err = store.Read("", store.Prefix(key))
	} else { //last by timestamp
		key := fmt.Sprintf("%v:", timeStampPrefix)
		var limit uint
		limit = 20 //default if without limition in req
		if req.Limit > 0 {
			limit = uint(req.Limit)
		}
		logger.Infof("Listing posts, offset: %v, limit: %v", req.Offset, limit)
		records, err = store.Read("", store.Prefix(key),
			store.Offset(uint(req.Offset)),
			store.Limit(limit))
	}

	if err != nil {
		return errors.BadRequest("posts.Query.store-read", "Failed to read from store: %v", err.Error())
	}
	// serialize the response list
	rsp.Posts = make([]*pb.Post, len(records))
	for i, record := range records {

		//dto proc to handle po to bo
		postRecord := &model.Post{}
		if err := json.Unmarshal(record.Value, postRecord); err != nil {
			return err
		}

		rsp.Posts[i] = &pb.Post{
			Id:      postRecord.ID,
			Title:   postRecord.Title,
			Slug:    postRecord.Slug,
			Content: postRecord.Content,
			Tags:    postRecord.Tags,
		}
	}
	return nil
}
