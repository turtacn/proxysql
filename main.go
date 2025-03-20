// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	sqle "github.com/dolthub/go-mysql-server"
)

func main() {
	var (
		username = "root"
		password = "123"
		host     = "localhost"
		dbname   = "tpcc"
		port     = 3306
	)

	// 创建数据库和表
	db := createTpccDatabase(dbname)

	// 创建数据库提供者
	provider := memory.NewDBProvider(db)

	// 创建引擎时需要指定正确的数据库提供者
	engine := sqle.NewDefault(provider)

	// 创建服务器配置
	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", host, port),
	}

	// 使用memory的会话构建器
	s, err := server.NewServer(
		config,
		engine,
		memory.NewSessionBuilder(provider),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MySQL server listening on %s:%d\n", host, port)
	if err := s.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func createTpccDatabase(dbName string) *memory.Database {
	db := memory.NewDatabase(dbName)
	provider := memory.NewDBProvider(db)

	// 创建正确的会话上下文


	session := sql.WithSession(memory.NewSession(sql.NewBaseSession(),provider))
	ctx := sql.NewContext(context.Background(), session)


	engine := sqle.NewDefault(provider)

	sqlFiles := []string{"tpcc-mysql/create_table.sql", "tpcc-mysql/add_fkey_idx.sql"}

	for _, file := range sqlFiles {
		sqlContent, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Error reading SQL file %s: %v", file, err)
		}

		// 先切换到目标数据库
		_, _, _, err = engine.Query(ctx, fmt.Sprintf("USE %s;", dbName))
		if err != nil {
			log.Fatalf("Error using database: %v", err)
		}

		queries := strings.Split(string(sqlContent), ";")
		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}

			// 执行前打印调试信息
			log.Printf("Executing query: %s\n", query)
			_, _, _, err = engine.Query(ctx, query)
			if err != nil {
				log.Fatalf("Error executing query '%s': %v", query, err)
			}
		}
		log.Printf("Successfully executed %s", file)
	}

	log.Println("TPCC database and tables created successfully.")
	return db
}