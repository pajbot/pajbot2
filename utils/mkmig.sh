#!/bin/bash

migrate -database mysql://pajbot2:password@/pajbot2_test -verbose create -ext sql -seq  -dir migrations $1
