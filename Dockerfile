FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod vendor && go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/
COPY dockerNohupRun.sh /app/nohupRun.sh

RUN chmod +x /app/nohupRun.sh
RUN chmod +x /app/main

# 安装必要的依赖
RUN apk add --no-cache bash

CMD ["./nohupRun.sh"]
