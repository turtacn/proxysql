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

// 定义配置结构体，提高可读性和可维护性
type Config struct {
	Username        string
	Password        string
	Host            string
	Port            int
	DBName          string
	MaxCacheSize    int64
	MaxConnections  int64
	QueryTimeout    time.Duration
	ProfilerAddress string
}

// 定义常量，集中管理
const (
	defaultUsername        = "root"
	defaultPassword        = "123"
	defaultHost            = "0.0.0.0"
	defaultDBName          = "tpcc"
	defaultPort            = 3306
	defaultProfilerAddress = ":6060"

	cacheSizeGB            = 1    // 内存缓存大小，单位GB
	maxConnectionsDefault  = 1024 // 最大并发连接数
	queryTimeoutSeconds    = 15   // 查询超时时间，单位秒
	sqlFileCreateTable     = "tpcc-mysql/create_table.sql"
	sqlFileAddForeignKey   = "tpcc-mysql/add_fkey_idx.sql"
	maxRetries             = 3
	retryDelay             = 10 * time.Millisecond
)

// 初始化函数，设置全局配置
func init() {
	// 优化1：最大化CPU利用率
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// 加载配置
	cfg := loadConfig()

	// 优化2：启用性能监控端点
	startProfiler(cfg.ProfilerAddress)

	// 初始化数据库
	dbProvider := initializeDatabase(cfg.DBName, cfg.QueryTimeout)

	// 创建并启动服务器
	startServer(cfg, dbProvider)
}

// loadConfig 加载应用程序配置
func loadConfig() Config {
	return Config{
		Username:        defaultUsername,
		Password:        defaultPassword,
		Host:            defaultHost,
		Port:            defaultPort,
		DBName:          defaultDBName,
		MaxCacheSize:    cacheSizeGB << 30, // 1GB内存缓存
		MaxConnections:  maxConnectionsDefault,
		QueryTimeout:    queryTimeoutSeconds * time.Second,
		ProfilerAddress: defaultProfilerAddress,
	}
}

// startProfiler 启动性能分析服务
func startProfiler(address string) {
	go func() {
		log.Printf("性能分析服务监听在 %s", address)
		if err := http.ListenAndServe(address, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动性能分析服务失败: %v", err)
		}
	}()
}

// initializeDatabase 初始化数据库
func initializeDatabase(dbName string, queryTimeout time.Duration) *memory.DbProvider {
	db := memory.NewDatabase(dbName)
	provider := memory.NewDBProvider(db)
	engine := sqle.NewDefault(provider)

	// 创建带超时控制的上下文
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	sqlCtx := sql.NewContext(ctx,sql.WithSession(memory.NewSession(sql.NewBaseSession(), provider)))

	// 批量预加载SQL内容
	sqlFiles := []string{sqlFileCreateTable, sqlFileAddForeignKey}
	allQueries := loadSQLFiles(sqlFiles)

	// 单次USE操作
	useDatabase(sqlCtx, engine, dbName)

	// 批量执行SQL查询
	executeSQLQueries(sqlCtx, engine, allQueries)

	log.Printf("数据库初始化完成")
	return provider
}

// loadSQLFiles 从文件中加载SQL语句
func loadSQLFiles(files []string)[]string {
	var allQueries []string
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("读取SQL文件失败: %s - %v", file, err)
		}
		queries := strings.Split(string(content), ";")
		for _, q := range queries {
			trimmedQuery := strings.TrimSpace(q)
			if trimmedQuery != "" {
				allQueries = append(allQueries, trimmedQuery)
			}
		}
	}
	return allQueries
}

// useDatabase 执行USE数据库操作
func useDatabase(ctx *sql.Context, engine *sqle.Engine, dbName string) {
	_, _, _, err := engine.Query(ctx, fmt.Sprintf("USE %s;", dbName))
	if err != nil {
		log.Fatalf("执行USE %s 失败: %v", dbName, err)
	}
}

// executeSQLQueries 批量执行SQL查询，带重试机制
func executeSQLQueries(ctx *sql.Context, engine *sqle.Engine, queries []string) {
	start := time.Now()
	for _, query := range queries {
		for retry := 0; retry < maxRetries; retry++ {
			_, _, _, err := engine.Query(ctx, query)
			if err == nil {
				break // 执行成功，跳出重试
			}
			log.Printf("执行SQL失败 (重试 %d/%d): %v\nQuery: %s", retry+1, maxRetries, err, query)
			time.Sleep(retryDelay)
		}
	}
	log.Printf("SQL执行完成 耗时: %v", time.Since(start))
}

// startServer 启动MySQL服务器
func startServer(cfg Config, dbProvider *memory.DbProvider) {
	engine := sqle.NewDefault(dbProvider)

	// 优化6：高性能服务器配置
	config := server.Config{
		Protocol:        "tcp",
		Address:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		MaxConnections:  uint64(cfg.MaxConnections), // 提升并发能力
		ConnReadTimeout: cfg.QueryTimeout,
		ConnWriteTimeout: cfg.QueryTimeout,
	}

	s, err := server.NewServer(
		config, // 使用优化后的配置
		engine,
		memory.NewSessionBuilder(dbProvider),
		nil,
	)
	if err != nil {
		log.Fatalf("创建服务器实例失败: %v", err)
	}

	// 优化7：资源使用情况输出
	log.Printf("启动配置：CPU核心[%d] 内存缓存[%.1fGB] 最大连接[%d]",
		runtime.NumCPU(),
		float64(cfg.MaxCacheSize)/1024/1024/1024,
		cfg.MaxConnections,
	)

	dsn := fmt.Sprintf(
		"MySQL server listening on %s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=Local&parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	log.Println(dsn)
	if err := s.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}