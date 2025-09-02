package app

import (
	ssov1 "github.com/Citadelas/protos/golang/sso"
	taskv1 "github.com/Citadelas/protos/golang/task"
	"google.golang.org/grpc"
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
	conn, err := grpc.NewClient(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: timeout,
		}),
	)
	if err != nil {
		panic(err)
	}
	return conn
}
