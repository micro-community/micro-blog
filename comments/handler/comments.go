package handler

import (
	"context"

	"github.com/micro-community/micro-blog/comments/model"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/blog/comments/proto"
)

//Comments of blog post
type Comments struct {
	Repository *model.Comment
}

//Options for handler
type Options struct {
	Namespace string
	Ctx       context.Context
}

//Option to set option
type Option func(o *Options)

//WithContext to set model for db
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Ctx = ctx
	}
}

// Save is a single request handler called via client.Call or the generated client code
func (c *Comments) Save(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logger.Info("Not yet implemented")
	return nil
}
