package handler

import (
	"context"

	"github.com/micro-community/micro-blog/posts/model"
)

//Posts Handler of Blog
type Posts struct {
	Repository *model.Repository
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
