package model

import (
	"context"
	"fmt"
	"math"

	"github.com/micro/micro/v3/service/store"
)

//DeleteTagByID from db
func (p *DB) DeleteTagByID(ctx context.Context, postID string) error {
	// Delete by ID
	if err := store.Delete(fmt.Sprintf("%v:%v", IDPrefix, postID)); err != nil {
		return err
	}
	return nil
}

//DeleteTagBySlug from db
func (p *DB) DeleteTagBySlug(ctx context.Context, slug string) error {

	if err := store.Delete(fmt.Sprintf("%v:%v", SlugPrefix, slug)); err != nil {
		return err
	}
	return nil
}

//DeleteTagByTimeStamp from db
func (p *DB) DeleteTagByTimeStamp(ctx context.Context, createdtimestamp int64) error {

	if err := store.Delete(fmt.Sprintf("%v:%v", TimeStampPrefix, math.MaxInt64-createdtimestamp)); err != nil {
		return err
	}
	return nil
}
