package main

import (
	"github.com/micro-community/micro-blog/common/protos/tags"
	"github.com/micro-community/micro-blog/tags/handler"
	"github.com/micro-community/micro-blog/tags/model"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("posts"),
		service.Version("latest"),
	)

	srv.Handle(&handler.Tags{
		DB: model.NewService(tags.NewTagsService("tags", srv.Client())),
	})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
