#!/bin/sh

basedir="$(cd "$(dirname "$0")" >/dev/null 2>&1 && pwd)"

[ -z "$git_release" ] && git_release=$(git describe --exact 2>/dev/null)
[ -z "$git_release" ] && git_release="git"
[ -z "$git_hash" ] && git_hash=$(git rev-parse --short HEAD)
[ -z "$git_branch" ] && git_branch=$(git rev-parse --abbrev-ref HEAD)

>&2 echo " * Building pajbot2 with the following flags: git_release=$git_release, git_hash=$git_hash, git_branch=$git_branch"

go build -ldflags "\
    -X \"main.buildTime=$(date +%Y-%m-%dT%H:%M:%S%:z)\" \
    -X \"main.buildRelease=$git_release\" \
    -X \"main.buildHash=$git_hash\" \
    -X \"main.buildBranch=$git_branch\" \
    " \
    "$@" \
    -o "$basedir/../cmd/bot/" \
    "$basedir/../cmd/bot/"
