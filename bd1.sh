#!/bin/sh
#export GOPROXY=http://mirrors.sangfor.org/nexus/repository/go-proxy
#export GOPROXY=https://goproxy.cn,direct
#export GOPRIVATE=sangfor.local
export GO111MODULE=on
export GOSUMDB=off
#export GOINSECURE=sangfor.local
export GIT_SSL_NO_VERIFY=1

go version
set -e
go mod vendor
go mod tidy
go mod verify
go build -mod=mod -o proxysql *.go
