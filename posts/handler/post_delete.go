package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro-community/micro-blog/posts/model"
	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

// Delete a post
func (p *Posts) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {

	logger.Info("Received Post.Delete request")

	records, err := store.Read(fmt.Sprintf("%v:%v", model.IDPrefix, req.Id))
	if err != nil && err != store.ErrNotFound {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("Post with ID %v not found", req.Id)
	}
	post := &model.Post{}
	if err := json.Unmarshal(records[0].Value, post); err != nil {
		return err
	}
	// Delete by ID
	if err = store.Delete(fmt.Sprintf("%v:%v", model.IDPrefix, post.ID)); err != nil {
		return err
	}

	// Delete by slug
	if err := store.Delete(fmt.Sprintf("%v:%v", model.SlugPrefix, post.Slug)); err != nil {
		return err
	}

	// Delete by timeStamp
	return store.Delete(fmt.Sprintf("%v:%v", model.TimeStampPrefix, post.CreateTimestamp))
}
