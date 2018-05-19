#!/bin/bash

$GOPATH/bin/migrate -database mysql://pajbot2:password@/pajbot2_test -path ./migrations  create -dir ./migrations $1
