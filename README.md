# proxysql

当前就一个功能，通过一个golang mysql server in-memory 去接管实际的mysql server 完成tpcc测试（https://github.com/Percona-Lab/tpcc-mysql.git），作为tpcc的性能上限判定

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
...
PRIMARY KEY(s_w_id, s_i_id) ) Engine=InnoDB
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS

Executed SQL from tpcc-mysql/create_table.sql
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
CREATE INDEX idx_customer ON customer (c_w_id,c_d_id,c_last,c_first)
CREATE INDEX idx_orders ON orders (o_w_id,o_d_id,o_c_id,o_id)
CREATE INDEX fkey_stock_2 ON stock (s_i_id)
CREATE INDEX fkey_order_line_2 ON order_line (ol_supply_w_id,ol_i_id)
ALTER TABLE district  ADD CONSTRAINT fkey_district_1 FOREIGN KEY(d_w_id) REFERENCES warehouse(w_id)
ALTER TABLE customer ADD CONSTRAINT fkey_customer_1 FOREIGN KEY(c_w_id,c_d_id) REFERENCES district(d_w_id,d_id)
ALTER TABLE history  ADD CONSTRAINT fkey_history_1 FOREIGN KEY(h_c_w_id,h_c_d_id,h_c_id) REFERENCES customer(c_w_id,c_d_id,c_id)
ALTER TABLE history  ADD CONSTRAINT fkey_history_2 FOREIGN KEY(h_w_id,h_d_id) REFERENCES district(d_w_id,d_id)
ALTER TABLE new_orders ADD CONSTRAINT fkey_new_orders_1 FOREIGN KEY(no_w_id,no_d_id,no_o_id) REFERENCES orders(o_w_id,o_d_id,o_id)
ALTER TABLE orders ADD CONSTRAINT fkey_orders_1 FOREIGN KEY(o_w_id,o_d_id,o_c_id) REFERENCES customer(c_w_id,c_d_id,c_id)
ALTER TABLE order_line ADD CONSTRAINT fkey_order_line_1 FOREIGN KEY(ol_w_id,ol_d_id,ol_o_id) REFERENCES orders(o_w_id,o_d_id,o_id)
ALTER TABLE order_line ADD CONSTRAINT fkey_order_line_2 FOREIGN KEY(ol_supply_w_id,ol_i_id) REFERENCES stock(s_w_id,s_i_id)
ALTER TABLE stock ADD CONSTRAINT fkey_stock_1 FOREIGN KEY(s_w_id) REFERENCES warehouse(w_id)
ALTER TABLE stock ADD CONSTRAINT fkey_stock_2 FOREIGN KEY(s_i_id) REFERENCES item(i_id)
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS

Executed SQL from tpcc-mysql/add_fkey_idx.sql
TPCC database and tables created.
MySQL server listening on localhost:3306
```


## 兼容性示例（tpcc_load部分）

```shell script
ubuntu $ ./tpcc_load -h 127.0.0.1 -P 3306 -u root -p123 -d tpcc -w 1
*************************************
*** TPCC-mysql Data Loader        ***
*************************************
option h with value '127.0.0.1'
option P with value '3306'
option u with value 'root'
option p with value '123'
option d with value 'tpcc'
option w with value '1'
<Parameters>
     [server]: 127.0.0.1
     [port]: 3306
     [DBname]: tpcc
       [user]: root
       [pass]: 123
  [warehouse]: 1
TPCC Data Load Started...
Loading Item 
..............
```


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