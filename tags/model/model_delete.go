package model

import (
	"context"
	"fmt"

	"github.com/micro/micro/v3/service/store"
)

//DeleteTagByID from db
func (p *DB) DeleteTagByID(ctx context.Context, tagID string) error {
	// Delete by ID
	if err := store.Delete(fmt.Sprintf("%v:%v", idPrefix, tagID)); err != nil {
		return err
	}
	return nil
}

//DeleteTagBySlug from db
func (p *DB) DeleteTagBySlug(ctx context.Context, slug string) error {

	if err := store.Delete(fmt.Sprintf("%v:%v", slugPrefix, slug)); err != nil {
		return err
	}
	return nil
}
