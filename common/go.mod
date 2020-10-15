module github.com/micro-community/micro-blog/common

go 1.15

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/golang/protobuf v1.4.3
	github.com/micro/micro/v3 v3.0.0-beta.6.0.20201014170732-9bd296d435bc
	google.golang.org/protobuf v1.25.0
)
