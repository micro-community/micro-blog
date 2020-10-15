package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gosimple/slug"
	"github.com/micro-community/micro-blog/posts/model"
	pb "github.com/micro-community/micro-blog/posts/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

//Posts Handler of Blog
type Posts struct{}

//Save a post
func (p *Posts) Save(ctx context.Context, req *pb.SaveRequest, rsp *pb.SaveResponse) error {
	logger.Info("Received Posts.Save request")

	if len(req.Post.Id) == 0 || len(req.Post.Title) == 0 || len(req.Post.Content) == 0 {
		return errors.BadRequest("posts.Save", "ID, title or content is missing")
	}

	post := &model.Post{
		ID:              req.Post.Id,
		Title:           req.Post.Title,
		Content:         req.Post.Content,
		Slug:            slug.Make(req.Post.Title),
		TagNames:        req.Post.TagNames,
		CreateTimestamp: time.Now().Unix(),
		UpdateTimestamp: time.Now().Unix(),
	}

	//serialize the post
	bytes, err := json.Marshal(post)
	if err != nil {
		return err
	}

	//store the article
	return store.Write(&store.Record{
		Key:   post.Slug,
		Value: bytes,
	})

}
