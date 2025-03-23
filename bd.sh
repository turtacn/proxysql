#!/bin/sh
export GOPROXY=http://mirrors.sangfor.org/nexus/repository/go-proxy
export GOPROXY=https://goproxy.cn,direct
export GOPRIVATE=sangfor.local
export GO111MODULE=on
export GOSUMDB=off
#export GOINSECURE=sangfor.local
export GIT_SSL_NO_VERIFY=1



go mod tidy
go mod vendor
go build -mod=mod -o  proxysql  *.go
exit



## 下面走容器编译

#docker run --net host --rm -t -i --add-host sangfor.local:192.168.198.1 -v $(pwd):$(pwd) -w $(pwd) -e GIT_SSL_NO_VERIFY=1 -e GO111MODULE=on -e GOSUMDB=off -e GOPRIVATE=sangfor.local -e GOPROXY=https://goproxy.cn,direct -e GOINSECURE=sangfor.local golang:sxf sh $(pwd)/mod.sh

docker run --net host --rm -t -i --add-host sangfor.local:192.168.198.1 -v $(pwd):$(pwd) -w $(pwd) -e GIT_SSL_NO_VERIFY=1 -e GO111MODULE=on -e GOSUMDB=off -e GOPRIVATE=sangfor.local -e GOPROXY=https://goproxy.cn,direct -e GOINSECURE=sangfor.local golang:1.23.6-alpine sh $(pwd)/mod.sh

exit

