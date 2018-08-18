#!/bin/bash

COMMAND=${1:-up}

echo $COMMAND

$GOPATH/bin/migrate -verbose -database mysql://root:penis123@/pajbot2_test -path ./migrations  ${COMMAND}
