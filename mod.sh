# 通过 golang 容器化编译
go mod tidy
go mod vendor
go build -mod=mod -o  proxysql  *.go