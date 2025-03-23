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

	_ "net/http/pprof" // 性能分析

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
)

const (
	maxCacheSize    = 1 << 30      // 1GB内存缓存
	maxConnections  = 1024         // 最大并发连接数
	queryTimeout    = 30 * time.Second
	initTimeout     = 5 * time.Minute // 延长初始化超时
)

func init() {
	// 设置最大CPU并行度
	if runtime.GOMAXPROCS(0) < runtime.NumCPU() {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}

func main() {
	// 启动性能监控服务
	go startMetricsServer()

	// 初始化数据库核心组件
	db, provider := initializeDatabase("tpcc")
	defer cleanupResources(provider)

	// 配置并启动服务器
	server := configureAndStartServer(provider)
	logServerDetails(server)

	if err := server.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// startMetricsServer 启动带自定义指标的监控服务
func startMetricsServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Fprintf(w, "heap_alloc_bytes %d\n", m.HeapAlloc)
		fmt.Fprintf(w, "cache_size_bytes %d\n", maxCacheSize)
	})
	log.Println("监控端点已启动 :6060")
	log.Fatal(http.ListenAndServe(":6060", mux))
}

// initializeDatabase 初始化数据库实例
func initializeDatabase(dbName string) (*memory.Database, *memory.DBProvider) {
	db := memory.NewDatabase(dbName).WithCacheSize(maxCacheSize)
	provider := memory.NewDBProvider(db)

	ctx, cancel := context.WithTimeout(context.Background(), initTimeout)
	defer cancel()

	if err := executeBootstrapSQL(ctx, provider, dbName); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 预热缓存
	if _, _, _, err := sqle.NewDefault(provider).Query(
		sql.NewContext(ctx, sql.WithSession(newSession(provider))),
		"SELECT 1 FROM DUAL"); err != nil {
		log.Printf("缓存预热警告: %v", err)
	}

	return db, provider
}

// executeBootstrapSQL 执行初始化SQL脚本
func executeBootstrapSQL(ctx context.Context, provider *memory.DBProvider, dbName string) error {
	engine := sqle.NewDefault(provider)
	session := newSession(provider)
	sqlCtx := sql.NewContext(ctx, sql.WithSession(session))

	// 切换数据库上下文
	if _, _, _, err := engine.Query(sqlCtx, fmt.Sprintf("USE %s;", dbName)); err != nil {
		return fmt.Errorf("USE操作失败: %w", err)
	}

	// 加载并执行SQL文件
	queries, err := loadSQLFiles([]string{
		"tpcc-mysql/create_table.sql",
		"tpcc-mysql/add_fkey_idx.sql",
	})
	if err != nil {
		return err
	}

	// 带重试机制的批量执行
	return executeWithRetry(sqlCtx, engine, queries, 3, 100*time.Millisecond)
}

// loadSQLFiles 加载并解析SQL文件
func loadSQLFiles(files []string) ([]string, error) {
	var queries []string
	for _, file := range files {
		content, err := os.ReadFile(file)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("文件不存在: %w", err)
		} else if err != nil {
			return nil, fmt.Errorf("文件读取失败: %w", err)
		}

		for _, q := range strings.Split(string(content), ";") {
			if cleaned := strings.TrimSpace(q); cleaned != "" {
				queries = append(queries, cleaned)
			}
		}
	}
	return queries, nil
}

// executeWithRetry 带指数退避的重试执行
func executeWithRetry(ctx context.Context, engine *sqle.Engine, queries []string, maxRetries int, baseDelay time.Duration) error {
	for i, query := range queries {
		for retry := 0; ; retry++ {
			queryCtx, cancel := context.WithTimeout(ctx, queryTimeout)
			_, _, _, err := engine.Query(sql.NewContext(queryCtx, sql.WithSession(ctx.Session)), query)
			cancel()

			if err == nil {
				break
			}

			if retry == maxRetries {
				return fmt.Errorf("执行失败[%d/%d] %q: %w", i+1, len(queries), query, err)
			}

			delay := time.Duration(retry*retry) * baseDelay
			time.Sleep(delay)
			log.Printf("重试中 (%d/%d): %v", retry+1, maxRetries, err)
		}
	}
	return nil
}

// configureAndStartServer 配置服务器实例
func configureAndStartServer(provider *memory.DBProvider) *server.Server {
	config := server.Config{
		Protocol:         "tcp",
		Address:         "0.0.0.0:3306",
		MaxConnections:  maxConnections,
		ConnReadTimeout: queryTimeout,
		ConnWriteTimeout: queryTimeout,
		Auth:            memory.NewAuth("root", "123"),
	}

	server, err := server.NewServer(
		config,
		sqle.NewDefault(provider),
		memory.NewSessionBuilder(provider),
		nil,
	)
	if err != nil {
		log.Fatalf("服务器配置失败: %v", err)
	}
	return server
}

// logServerDetails 记录启动配置信息
func logServerDetails(s *server.Server) {
	log.Printf(`启动配置:
  CPU核心      : %d
  内存缓存    : %.1fGB
  最大连接数  : %d
  监听地址    : %s`,
		runtime.NumCPU(),
		float64(maxCacheSize)/(1<<30),
		maxConnections,
		s.Address(),
	)
}

// newSession 创建统一会话实例
func newSession(provider *memory.DBProvider) sql.Session {
	return memory.NewSession(sql.NewBaseSession(), provider)
}

// cleanupResources 资源清理钩子
func cleanupResources(provider *memory.DBProvider) {
	// 未来扩展资源释放逻辑
}