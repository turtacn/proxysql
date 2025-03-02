# proxysql

当前就一个功能，通过一个golang mysql server in-memory 去接管实际的mysql server 完成tpcc测试，作为tpcc的性能上限判定

## 构建

指定goalng 版本 1.20.2，见 gobin.sh
指 github.com/dolthub/go-mysql-server 版本 v0.11.0, go.mod

```shell script
go mod tidy 
go build -o proxysql main.go
```

## 运行

```shell script
./proxysql
```


## 等价的过程
```shell script
# 创建 tpcc 数据库
echo "创建 tpcc 数据库..."
mysql -h 127.0.0.1 -P 3306 -u root -p123 -e "CREATE DATABASE IF NOT EXISTS tpcc;"
echo "初始化 tpcc 表结构..."
mysql -h 127.0.0.1 -P 3306 -u root -p123 -D tpcc < create_table.sql
mysql -h 127.0.0.1 -P 3306 -u root -p123 -D tpcc < add_fkey_idx.sql
 
# 快速验证 -w 100 ==> -w 1
# 加载 TPC-C 数据 (100 仓库)
echo "加载 TPC-C 数据 (100 仓库)..."
./tpcc_load -h 127.0.0.1 -P 3306 -u root -p123 -d tpcc -w 100
 
# 运行 TPC-C 测试 (100 终端并发)
echo "运行 TPC-C 测试 (100 终端并发)..."
./tpcc_start -h 127.0.0.1 -P 3306 -u root -p123 -d tpcc -w 100 -c 100 -r 10 > tpcc_100_terminals.log
 
# 运行 TPC-C 测试 (300 终端并发)
echo "运行 TPC-C 测试 (300 终端并发)..."
./tpcc_start -h 127.0.0.1 -P 3306 -u root -p123 -d tpcc -w 100 -c 300 -r 10 > tpcc_300_terminals.log
 
# 运行 TPC-C 测试 (500 终端并发)
echo "运行 TPC-C 测试 (500 终端并发)..."
./tpcc_start -h 127.0.0.1 -P 3306 -u root -p123 -d tpcc -w 100 -c 500 -r 10 > tpcc_500_terminals.log

```