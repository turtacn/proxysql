package main

import (
	"context"
	"log"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
)

// MyAuthenticator 实现 server.UserAuthenticator 接口，
// 仅允许用户名为 "root"，密码为 "123" 的连接。
type MyAuthenticator struct{}

// Authenticate 检查用户名、主机和密码
func (a *MyAuthenticator) Authenticate(ctx context.Context, user, host, password string) (bool, error) {
	if user == "root" && password == "123" {
		return true, nil
	}
	return false, nil
}

func main() {
	ctx := context.Background()

	// 创建一个内存数据库 "tpcc"
	// 注意：tpcc-mysql 脚本会执行 "CREATE DATABASE IF NOT EXISTS tpcc"，
	// 本示例预先创建该数据库，若数据库已存在则 DDL 不会报错。
	tpccDB := memory.NewDatabase("tpcc")

	// 使用预先创建的数据库构建引擎
	// NewDefault 接受一个或多个 sql.Database，内部会生成一个 Catalog，
	// 后续支持 CREATE DATABASE、DDL、DML 等操作。
	engine, err := sql.NewDefault(tpccDB)
	if err != nil {
		log.Fatal(err)
	}

	// 如果需要支持 information_schema 查询，可以将其加入 Catalog，
	// 例如：
	// engine.Catalog.RegisterDatabase(memory.NewInformationSchemaDatabase(engine.Catalog))

	// 创建 MySQL 协议服务端，设置监听地址、认证器和显示的服务器版本
	srv, err := server.NewDefaultServer(engine,
		server.WithAddress("127.0.0.1:3306"),
		server.WithAuthenticator(&MyAuthenticator{}),
		server.WithServerVersion("8.0.28"),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MySQL 模拟服务器已启动，监听 127.0.0.1:3306")
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
