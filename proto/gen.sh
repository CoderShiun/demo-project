#!/usr/bin/env bash

GRPC_GW_PATH=`go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway`
GRPC_GW_PATH="${GRPC_GW_PATH}/../third_party/googleapis"

#LS_PATH_NS=`go list -f '{{ .Dir }}' github.com/mxc-foundation/lpwan-server/proto/ns`
#LS_PATH_NS="${LS_PATH_NS}/../.."

PROTOBUF_PATH=`go list -f '{{ .Dir }}' github.com/golang/protobuf/ptypes`

# generate the gRPC code
protoc -I. -I${GRPC_GW_PATH} -I${PROTOBUF_PATH} --go_out=plugins=grpc:. \
    grpc.proto \
    restful.proto

# generate the JSON interface code
protoc -I. -I${GRPC_GW_PATH} -I${PROTOBUF_PATH} --grpc-gateway_out=logtostderr=true:. \
    grpc.proto \
    restful.proto

# generate the swagger definitions
#protoc -I. -I${GRPC_GW_PATH} -I${PROTOBUF_PATH} --swagger_out=json_names_for_fields=true:.swagger \
#    grpc.proto \
#    restful.proto

# merge the swagger code into one file
#go run swagger/main.go swagger > ../static/swagger/proto.swagger.json
