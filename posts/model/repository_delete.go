package model

import (
	"context"
	"fmt"
	"math"

	"github.com/micro/micro/v3/service/store"
)

//DeletePostByID from db
func (r *Repository) DeletePostByID(ctx context.Context, postID string) error {
	// Delete by ID
	if err := store.Delete(fmt.Sprintf("%v:%v", IDPrefix, postID)); err != nil {
		return err
	}
	return nil
}

//DeletePostBySlug from db
func (r *Repository) DeletePostBySlug(ctx context.Context, slug string) error {

	if err := store.Delete(fmt.Sprintf("%v:%v", SlugPrefix, slug)); err != nil {
		return err
	}
	return nil
}

//DeletePostByTimeStamp from db
func (r *Repository) DeletePostByTimeStamp(ctx context.Context, createdtimestamp int64) error {

	if err := store.Delete(fmt.Sprintf("%v:%v", TimeStampPrefix, math.MaxInt64-createdtimestamp)); err != nil {
		return err
	}
	return nil
}
