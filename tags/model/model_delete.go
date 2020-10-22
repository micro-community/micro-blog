package model

import (
	"context"
	"fmt"

	"github.com/micro/micro/v3/service/store"
)

//DeleteTagByResourceID from db
func (p *DB) DeleteTagByResourceID(ctx context.Context, resourceID string, tag *Tag) error {
	// Delete by ID
	if err := store.Delete(fmt.Sprintf("%v:%v:%v", slugPrefix, tag.Slug, resourceID)); err != nil {
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

//DecreaseTagCount increase  a tag count basing on all post with same tag
func (p *DB) DecreaseTagCount(resourceID string, tag *Tag) error {

	//tagCountPrefix:tagslug:resourceID ,add resource in some slug
	if err := store.Delete(fmt.Sprintf("%v:%v:%v", tagCountPrefix, tag.Slug, resourceID)); err != nil {
		return err
	}

	// get tag count
	recs, err := store.List(store.Prefix(fmt.Sprintf("%v:%v", tagCountPrefix, tag.Slug)), store.Limit(1000))
	if err != nil {
		return err
	}

	tag.Count = int64(len(recs))

	//save the tag
	return p.saveTag(tag)

}
