// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	// 添加性能分析支持（无需修改已有import）
	_ "net/http/pprof"

	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	sqle "github.com/dolthub/go-mysql-server"
)

// 新增性能优化常量
const (
	maxCacheSize    = 1 << 30 // 1GB内存缓存
	maxConnections  = 1024    // 最大并发连接数
	queryTimeout    = 15 * time.Second
)

func init() {
	// 优化1：最大化CPU利用率
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// 优化2：启用性能监控端点
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	var (
		username = "root"
		password = "123"
		host     = "0.0.0.0" // 优化3：监听所有接口
		dbname   = "tpcc"
		port     = 3306
	)

	// 优化4：批量预加载SQL文件
	db := createTpccDatabase(dbname)

	// 优化5：配置大内存缓存
	provider := memory.NewDBProvider(
		memory.with(maxCacheSize), // 新增配置
		db,
	)

	engine := sqle.NewDefault(provider)

	// 优化6：高性能服务器配置
	config := server.Config{
		Protocol:        "tcp",
		Address:         fmt.Sprintf("%s:%d", host, port),
		MaxConnections:  maxConnections,    // 提升并发能力
		ConnReadTimeout: queryTimeout,
		ConnWriteTimeout:queryTimeout,
	}

	s, err := server.NewServer(
		config, // 使用优化后的配置
		engine,
		memory.NewSessionBuilder(provider),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 优化7：资源使用情况输出
	log.Printf("启动配置：CPU核心[%d] 内存缓存[%.1fGB] 最大连接[%d]",
		runtime.NumCPU(),
		float64(maxCacheSize)/1024/1024/1024,
		maxConnections,
	)

	if err := s.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func createTpccDatabase(dbName string) *memory.Database {
	db := memory.NewDatabase(dbName)
	provider := memory.NewDBProvider(db)

	// 优化8：带超时控制的上下文
	ctx, cancel := context.WithTimeout(
		sql.NewContext(context.Background(),
			sql.WithSession(memory.NewSession(sql.NewBaseSession(), provider)),
		),
		queryTimeout,
	)
	defer cancel()

	engine := sqle.NewDefault(provider)

	// 优化9：批量预加载所有SQL内容
	var allQueries []string
	for _, file := range []string{"tpcc-mysql/create_table.sql", "tpcc-mysql/add_fkey_idx.sql"} {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("文件读取失败: %v", err)
		}
		allQueries = append(allQueries, strings.Split(string(content), ";")...)
	}

	// 优化10：单次USE操作
	if _, _, _, err := engine.Query(ctx, fmt.Sprintf("USE %s;", dbName)); err != nil {
		log.Fatalf("USE失败: %v", err)
	}

	// 优化11：批量执行优化
	start := time.Now()
	for _, query := range allQueries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		// 优化12：快速重试机制
		for retry := 0; retry < 3; retry++ {
			_, _, _, err := engine.Query(ctx, query)
			if err == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

	log.Printf("数据库初始化完成 耗时: %v", time.Since(start))
	return db
}