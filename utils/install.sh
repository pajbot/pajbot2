#!/bin/sh

set -e

# Compile C# library

if [ -n "$1" ]; then
    RID="$1"
else
    RID=linux-x64
fi


cd 3rdParty/MessageHeightTwitch/
dotnet publish MessageHeightTwitch.csproj -c Release -r $RID -o ../../cmd/bot/MessageHeightTwitch
