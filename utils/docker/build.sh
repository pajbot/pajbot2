#!/usr/bin/env bash

if [ ! -d utils/docker ]; then
    echo "This script needs to be called from the root folder, i.e. ./utils/docker/build.sh"
    exit 1
fi

IMAGE_NAME=pajbot2:latest

echo docker build --pull -t "$IMAGE_NAME" .
