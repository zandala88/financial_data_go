# 第一阶段：构建 Go 应用
FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod tidy && go mod vendor && go build -o main .

ENV TimeZone=Asia/Shanghai
# 使用软连接，并且将时区配置覆盖/etc/timezone
RUN ln -snf /usr/share/zoneinfo/$TimeZone /etc/localtime && echo $TimeZone > /etc/timezone

# 第二阶段：创建最终运行环境
FROM golang:1.23

WORKDIR /app

COPY --from=builder /app/main /app/

RUN chmod +x /app/main

# 创建日志目录
RUN mkdir -p /app/runlog

# 运行应用，并将日志输出到带时间戳的文件
CMD ["sh", "-c", "./main >> /app/runlog/$(date +'%Y-%m-%d-%H:%M:%S').log 2>&1"]
