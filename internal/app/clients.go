package app

import (
	ssov1 "github.com/Citadelas/protos/golang/sso"
	taskv1 "github.com/Citadelas/protos/golang/task"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func (a *App) mustInitClients() error {
	// Initialize SSO client
	ssoConn := mustGenerateClient(a.cfg.Services.SSO.Endpoint, a.cfg.Services.SSO.Timeout)

	a.ssoClient = ssov1.NewAuthClient(ssoConn)

	// Initialize Task client
	taskConn := mustGenerateClient(a.cfg.Services.Task.Endpoint, a.cfg.Services.Task.Timeout)

	a.taskClient = taskv1.NewTaskServiceClient(taskConn)

	return nil
}

func mustGenerateClient(endpoint string, timeout time.Duration) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
		grpc_retry.WithMax(5),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
	)))
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithConnectParams(grpc.ConnectParams{
		MinConnectTimeout: timeout,
	}))
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		panic(err)
	}
	return conn
}
