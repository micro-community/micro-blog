package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro/micro/v3/service/store"
)

//UpdateTag to db
func (p *DB) UpdateTag(ctx context.Context, oldTag, tag *Tag) error {

	return p.saveTag(tag)

}

func (p *DB) saveTag(tag *Tag) error {

	key := fmt.Sprintf("%v:%v", slugPrefix, tag.Slug)
	typeKey := fmt.Sprintf("%v:%v:%v", typePrefix, tag.Type, tag.Slug)

	bytes, err := json.Marshal(tag)
	if err != nil {
		return err
	}

	// write resourceId:slug to enable prefix listing based on type
	err = store.Write(&store.Record{
		Key:   key,
		Value: bytes,
	})
	if err != nil {
		return err
	}
	return store.Write(&store.Record{
		Key:   typeKey,
		Value: bytes,
	})
}
