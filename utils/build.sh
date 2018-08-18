#!/bin/bash

go build -ldflags "-X main.buildTime=`date -u +.%Y%m%d.%H%M%S`"
