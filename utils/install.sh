#!/bin/sh

set -e

# Compile C++ "coreruncommon" lib
(
cd 3rdParty/MessageHeightTwitch/c-interop
g++ -c coreruncommon.cpp --std=c++14
ar rvs libcoreruncommon.a coreruncommon.o
)

# Compile C# library
(
cd 3rdParty/MessageHeightTwitch
dotnet publish --configuration Release -o build/
cp charmap.bin.gz build/*.dll ../../cmd/bot/
)

if [ ! -d internal/ConfusableMatcher-go-interop ] || [ ! -d internal/ConfusableMatcher-go-interop/ConfusableMatcher ]; then
    echo "you probably need to run git submodule update --init --recursive"
    exit 1
fi

# Compile ConfusableMatcher library (Release)
(
cd internal/ConfusableMatcher-go-interop/ConfusableMatcher
[ ! -d build ] && mkdir build
cd build
cmake -DCMAKE_BUILD_TYPE=Release ..
make -j"$(nproc)"
)

# Copy compiled C++ library into the go library folder
cp internal/ConfusableMatcher-go-interop/ConfusableMatcher/build/*.a internal/ConfusableMatcher-go-interop/
