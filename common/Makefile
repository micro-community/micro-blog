#vars
GOPATH:=$(shell go env GOPATH)
GOROOT:=$(shell go env GOROOT)
empty :=
space := $(empty) $(empty)

#go imports
VALIDATE_IMPORT := Mvalidate/validate.proto=github.com/envoyproxy/protoc-gen-validate/validate
GO_IMPORT_SPACES := ${VALIDATE_IMPORT},\
	Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any,\
	Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,\
	Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,\
	Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp,\
	Mgoogle/protobuf/wrappers.proto=github.com/golang/protobuf/ptypes/wrappers,\
	Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor


GO_IMPORT:=$(subst $(space),,$(GO_IMPORT_SPACES))

all: build

vendor:
	go mod vendor

.PHONY: init
init:
	go get -d github.com/envoyproxy/protoc-gen-validate
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro


.PHONY: proto
proto:
	protoc \
	-I "protos/${p}"  \
	--go_out="${GO_IMPORT}:protos/${p}"  \
	--micro_out="${GO_IMPORT}:protos/${p}"   \
	--validate_out="lang=go:protos/${p}"   \
  	"${p}.proto"
