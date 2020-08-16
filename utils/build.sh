#!/bin/sh

set -e

basedir="$(cd "$(dirname "$0")" >/dev/null 2>&1 && pwd)"

if [ -z "$git_release" ]; then
    git_release=$(git describe --exact 2>/dev/null)
fi
if [ -z "$git_release" ]; then
    git_release="git"
fi
if [ -z "$git_hash" ]; then
    git_hash=$(git rev-parse --short HEAD)
fi
if [ -z "$git_branch" ]; then
    git_branch=$(git rev-parse --abbrev-ref HEAD)
fi

>&2 echo " * Building pajbot2 with the following flags: git_release=$git_release, git_hash=$git_hash, git_branch=$git_branch"

cd "$basedir/../web" && npm i && npm run build && mv static/build/index.html views/

go build -ldflags "\
    -X \"main.buildTime=$(date +%Y-%m-%dT%H:%M:%S%:z)\" \
    -X \"main.buildRelease=$git_release\" \
    -X \"main.buildHash=$git_hash\" \
    -X \"main.buildBranch=$git_branch\" \
    " \
    "$@" \
    -o "$basedir/../cmd/bot/" \
    "$basedir/../cmd/bot/"
