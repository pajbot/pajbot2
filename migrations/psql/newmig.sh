#!/bin/sh

set -e

script_path=$(dirname "$0")
description=$1

if [ -z "$description" ]; then
    echo "usage: $0 description-of-migration"
    exit 1
fi

touch "$script_path/$(date --utc '+%s')-$description.sql"
