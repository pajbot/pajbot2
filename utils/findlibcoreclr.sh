#!/bin/bash

LIBCORECLR_PATH=$(dotnet --list-runtimes | grep Microsoft.NETCore.App | tail -1 | awk '{gsub(/\[|\]/, "", $3); print $3 "/" $2}')
if [ -z "$LIBCORECLR_PATH" ]; then
    echo "Unable to find path to libcoreclr. Ensure dotnet is installed"
    exit 1
fi

echo "Found libcoreclr.so at $LIBCORECLR_PATH"
echo "To run the bot with csharp modules enabled, build the bot with '-tags csharp' and set the path with LIBCOREFOLDER=$LIBCORECLR_PATH"
