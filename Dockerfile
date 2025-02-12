FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod vendor && go build -o main .

FROM alpha:lastet

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config.yaml ./config.yaml

CMD ["./main"]
