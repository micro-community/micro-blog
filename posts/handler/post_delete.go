package handler

import (
	"context"
	"fmt"

	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

// Delete a post
func (p *Posts) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {

	logger.Info("Received Post.Delete request")
	if len(req.Id) == 0 {
		return errors.BadRequest("posts.Save.input-check", "ID is missing")
	}

	post, err := p.Repository.CheckByPostID(req.Id)

	if err != nil {
		return err
	}

	if post == nil {
		return fmt.Errorf("Post with ID %v not found", req.Id)
	}

	// Delete by ID
	if err = p.Repository.DeletePostByID(ctx, post.ID); err != nil {
		return err
	}

	// Delete by slug
	if err = p.Repository.DeletePostBySlug(ctx, post.Slug); err != nil {
		return err
	}

	// Delete by slug
	return p.Repository.DeletePostByTimeStamp(ctx, post.Created)
}
