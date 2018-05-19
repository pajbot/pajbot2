#!/bin/bash

# User configuration
PAJBOT1_PATH="/c/dev/git/pajbot"

# Calculated paths
PROTO_PATH="${PAJBOT1_PATH}/proto/pajbot"

if [[ ! -e "${PROTO_PATH}/grpc/pajbot.proto" ]]; then
    echo "pajbot.proto not found under given PAJBOT1_PATH. Configure build-proto.sh yourself"
    exit 1
fi

protoc -I ${PROTO_PATH} grpc/pajbot.proto --go_out=plugins=grpc:.
