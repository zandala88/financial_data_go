package python

import (
	"context"
	"financia/config"
	pb "financia/server/python/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var (
	rpcConn *grpc.ClientConn
	rpcCli  pb.PredictorClient
)

// NewGRPCClient
// 初始化 gRPC 负载均衡连接
func NewGRPCClient() {
	// 负载均衡：round_robin
	servers := config.Configs.Python.Url
	conn, err := grpc.NewClient(
		servers,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Error("[NewGRPCClient] [err] = ", err.Error())
		return
	}

	client := pb.NewPredictorClient(conn)
	rpcConn = conn
	rpcCli = client
	return
}

func SendPredictRequest(req *pb.PredictRequest) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if len(req.Data) != 31 {
		return 0, nil
	}

	resp, err := rpcCli.Predict(ctx, req)
	if err != nil {
		zap.S().Error("[SendPredictRequest] [err] = ", err.Error())
		return 0, err
	}
	return resp.Val, nil
}
