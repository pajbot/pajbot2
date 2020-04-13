#!/bin/sh

[ -z "$git_release" ] && git_release=$(git describe --exact 2>/dev/null) || git_release="dev"
[ -z "$git_hash" ] && git_hash=$(git rev-parse --short HEAD)
[ -z "$git_branch" ] && git_branch=$(git rev-parse --abbrev-ref HEAD)

go build -ldflags "\
    -X \"main.buildTime=$(date +%Y-%m-%dT%H:%M:%S%:z)\" \
    -X \"main.buildRelease=$git_release\" \
    -X \"main.buildHash=$git_hash\" \
    -X \"main.buildBranch=$git_branch\" \
    " \
 ./cmd/bot
