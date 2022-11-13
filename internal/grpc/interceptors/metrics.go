package interceptors

import (
	"context"
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/metrics"
	"google.golang.org/grpc"
)

func ServerMetricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	defer metrics.GrpcServerResponseTime(time.Since(start).Seconds(), info.FullMethod)

	m, err := handler(ctx, req)

	return m, err
}

func ClientMetricsInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	defer metrics.GrpcClientResponseTime(time.Since(start).Seconds(), method)

	return invoker(ctx, method, req, reply, cc, opts...)
}
