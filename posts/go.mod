module github.com/micro-community/micro-blog/posts

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/gosimple/slug v1.9.0
	github.com/micro-community/micro-blog/common v0.0.1
	github.com/micro/micro/v3 v3.0.0-beta.6.0.20201014170732-9bd296d435bc
	google.golang.org/protobuf v1.25.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

//use local mod for debug
replace github.com/micro-community/micro-blog/common => ../common
