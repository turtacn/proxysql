#!/bin/bash
# 设置动态库搜索路径
export LD_LIBRARY_PATH=/usr/local/mysql/lib/mysql/

# 外部参数：
# $1: 数据库名称（必选）
# $2: WH 参数（必选）
# $3: STEP 值（可选，默认1）
# $4: 并发数，即 tpcc_load 的 level 数量（可选，默认4）
DBNAME=$1
WH=$2
STEP=${3:-1}
CONCURRENCY=${4:-4}

HOST=127.0.0.1

# 首先调用 level=1 的 tpcc_load，范围为 1 到 WH，输出到 1.out
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 1 -m 1 -n $WH >> 1.out &

x=1

# 循环处理范围分段，步长为 STEP
while [ $x -le $WH ]
do
  echo "处理范围: $x 到 $(( x + STEP - 1 ))"
  # 根据并发参数 CONCURRENCY，从 level=2 开始依次调用
  for level in $(seq 2 $CONCURRENCY); do
    ./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l $level -m $x -n $(( x + STEP - 1 )) >> ${level}_$x.out &
  done
  # 更新下一个区间
  x=$(( x + STEP ))
done
