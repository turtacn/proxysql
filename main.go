package main

import (
	"context"
	"log"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sqle"
)

// MyAuthenticator 实现了 server.UserAuthenticator 接口，
// 只允许用户名为 "root"、密码为 "123" 的连接。
type MyAuthenticator struct{}

// Authenticate 检查用户名和密码
func (a *MyAuthenticator) Authenticate(ctx context.Context, user, host, password string) (bool, error) {
	if user == "root" && password == "123" {
		return true, nil
	}
	return false, nil
}

func main() {
	ctx := context.Background()

	// 创建一个名为 "tpcc" 的内存数据库
	tpccDB := memory.NewDatabase("tpcc")

	// 使用 NewDatabaseProvider 将数据库打包成提供者，
	// 这里 provider 内部会支持基本的 DDL/DML 操作
	provider := sqle.NewDatabaseProvider(tpccDB)
	// 添加 information_schema 数据库，便于客户端查询系统信息
	provider.AddDatabase(memory.NewInformationSchemaDatabase(provider))

	// 创建 SQL 引擎，后续所有 SQL 语句均由此引擎解析执行
	engine := sqle.NewDefault(provider)

	// 创建 MySQL 协议服务端
	// 设置监听地址、认证器以及显示的服务器版本
	srv, err := server.NewDefaultServer(engine,
		server.WithAddress("127.0.0.1:3306"),
		server.WithAuthenticator(&MyAuthenticator{}),
		server.WithServerVersion("8.0.28"),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MySQL 模拟服务器已启动，监听 127.0.0.1:3306")
	// 启动服务（此调用会阻塞）
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	// 保持 ctx 不退出（一般 Start() 已阻塞，本示例不会走到这里）
	<-ctx.Done()
}
