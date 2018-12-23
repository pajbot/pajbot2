#!/bin/sh

set -e

# Compile C++ "coreruncommon" lib
cd 3rdParty/MessageHeightTwitch/c-interop
g++ -c coreruncommon.cpp --std=c++14
ar rvs libcoreruncommon.a coreruncommon.o

# Compile C# library
cd ../
dotnet publish --configuration Release -o build/

cp charmap.bin.gz build/*.dll ../../cmd/bot/
