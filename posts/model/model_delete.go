package model

import (
	"context"
	"fmt"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

//DeletePostByID from db
func (p *DB) DeletePostByID(ctx context.Context, id string) error {

	return nil
}

//DeletePostBySlug from db
func (p *DB) DeletePostBySlug(ctx context.Context, slug string) error {

	key := fmt.Sprintf("%v:%v", SlugPrefix, slug)
	logger.Infof("Reading post by slug: %v", slug)
	records, err := store.Read("", store.Prefix(key))
	return nil
}

//DeletePostByTimeStamp from db
func (p *DB) DeletePostByTimeStamp(ctx context.Context, slug string) error {

	return nil
}
