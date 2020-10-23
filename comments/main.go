package main

import (
	"github.com/micro-community/micro-blog/comments/handler"
	"github.com/micro-community/micro-blog/comments/model"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("comments"),
		service.Version("latest"),
	)

	// Register Handler
	srv.Handle(&handler.Comments{
		Repository: model.NewComment(),
	})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
