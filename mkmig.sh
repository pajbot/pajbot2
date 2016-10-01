#!/bin/bash

$GOPATH/bin/migrate -url mysql://pajbot2:password@/pajbot2_test -path ./migrations create $1
