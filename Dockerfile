FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod vendor && go build -o main .

FROM golang:1.23

WORKDIR /app

# 确保 main 被正确复制到 /app 目录
COPY --from=builder /app/main /app/

# 确保 main 文件可以被执行
RUN chmod +x /app/main

CMD ["./main"]