package handler

import (
	"context"

	"github.com/micro-community/micro-blog/tags/model"
)

//Tags Handler of Blog
type Tags struct {
	DB *model.DB
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
