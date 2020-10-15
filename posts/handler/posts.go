package handler

import (
	"context"

	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/logger"
)

//Posts of Blog
type Posts struct{}

//Save a post
func (p *Posts) Save(ctx context.Context, req *pb.SaveRequest, rsp *pb.SaveResponse) error {
	logger.Info("Received Posts.Save request")
	return nil
}
