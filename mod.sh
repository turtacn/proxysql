go mod tidy
go mod vendor
go build -mod=mod -o  proxysql  *.go