FROM ubuntu:20.04

# 安装运行所需的依赖库
RUN apt-get update && apt-get install -y vim sysbench libmysqlclient21 libssl1.1 && rm -rf /var/lib/apt/lists/*

# 设置工作目录为 /app，此目录中包含 proxysql 可执行文件和 tpcc-mysql 目录
WORKDIR /app

# 将 proxysql 可执行文件复制到 /app 下
COPY proxysql .

# 将整个 tpcc-mysql 目录复制到 /app 下，保持原目录结构
COPY tpcc-mysql ./tpcc-mysql

# 设置默认入口为当前目录下的 proxysql
ENTRYPOINT ["./proxysql"]
