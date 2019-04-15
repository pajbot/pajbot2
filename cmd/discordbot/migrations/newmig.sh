#!/bin/bash

RELATIVE_PATH=$(dirname $0)

if [ -z $1 ]; then
    echo "usage: $0 description-of-migration"
    exit 1
fi

touch "$RELATIVE_PATH/$(date --utc '+%s')-$1.sql"
