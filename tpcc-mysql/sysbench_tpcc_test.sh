#!/bin/bash

# ---------------------- 配置 ----------------------
MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_USER="root"
MYSQL_PASSWORD="123"
MYSQL_DATABASE="tpcc"

# 测试参数 (您可以根据需要修改这些参数)
TABLE_SIZE=100000  # 每个表的数据行数
TABLES=1          # 测试表的数量
THREADS=10        # 并发线程数
TIME=60           # 测试运行时间（秒）
OLTP_TEST="oltp_read_write" # 测试类型，例如：oltp_read_only, oltp_write_only, oltp_read_write

# ---------------------- 函数 ----------------------

prepare_data() {
  echo "准备测试数据..."
  sysbench $OLTP_TEST --mysql-host=$MYSQL_HOST --mysql-port=$MYSQL_PORT --mysql-user=$MYSQL_USER --mysql-password=$MYSQL_PASSWORD --mysql-db=$MYSQL_DATABASE --table-size=$TABLE_SIZE --tables=$TABLES prepare
  echo "测试数据准备完成。"
}

run_test() {
  echo "开始运行 $THREADS 个线程的 $OLTP_TEST 测试，持续 $TIME 秒..."
  sysbench $OLTP_TEST --mysql-host=$MYSQL_HOST --mysql-port=$MYSQL_PORT --mysql-user=$MYSQL_USER --mysql-password=$MYSQL_PASSWORD --mysql-db=$MYSQL_DATABASE --table-size=$TABLE_SIZE --tables=$TABLES --threads=$THREADS --time=$TIME run
  echo "测试运行完成。"
}

cleanup_data() {
  echo "清理测试数据..."
  sysbench $OLTP_TEST --mysql-host=$MYSQL_HOST --mysql-port=$MYSQL_PORT --mysql-user=$MYSQL_USER --mysql-password=$MYSQL_PASSWORD --mysql-db=$MYSQL_DATABASE --table-size=$TABLE_SIZE --tables=$TABLES cleanup
  echo "测试数据清理完成。"
}

# ---------------------- 执行 ----------------------

prepare_data
run_test
cleanup_data

echo "测试脚本执行完毕。"