package main

import (
	"github.com/micro-community/micro-blog/posts/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("posts"),
		service.Version("latest"),
	)

	srv.Handle(handler.NewPost(srv.Client()))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
