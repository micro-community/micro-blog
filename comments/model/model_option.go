package model

import "context"

// Some Const data for comments
const (
	IDPrefix        = "id"
	TimeStampPrefix = "timestamp"
)

// QueryType type
type QueryType int

// QueryType type
const (
	QueryByID QueryType = iota
	QueryBySlug
	QueryByTimestamp
)

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
