#!/bin/bash

# Locate `libcoreclr.so` (requires locate)
LIBCORECLR_PATH=$(locate /libcoreclr.so | grep -v cache | head -n 1 | sed -ne 's/\/libcoreclr\.so//p')
if [ -z "$LIBCORECLR_PATH" ]; then
    echo "Unable to find path to libcoreclr. Ensure you've run 'sudo updatedb' and that dotnet and mlocate is installed"
    exit 1
fi

echo "Found libcoreclr.so at $LIBCORECLR_PATH"
echo "To run the bot with csharp modules enabled, build the bot with '-tags csharp' and set the path with LIBCOREFOLDER=$LIBCORECLR_PATH"
