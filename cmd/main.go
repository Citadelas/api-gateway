package main

import (
	"github.com/Citadelas/api-gateway/internal/config"
	"github.com/Citadelas/api-gateway/internal/handlers/sso"
	ssov1 "github.com/Citadelas/protos/golang/sso"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func mustGenerateClient(endpoint string, timeout time.Duration) *grpc.ClientConn {
	conn, err := grpc.NewClient(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: timeout,
		}),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func main() {
	cfg := config.MustLoad()
	ssoConn := mustGenerateClient(cfg.Services.SSO.Endpoint, cfg.Services.SSO.Timeout)
	defer ssoConn.Close()
	ssoClient := ssov1.NewAuthClient(ssoConn)

	taskConn := mustGenerateClient(cfg.Services.Task.Endpoint, cfg.Services.Task.Timeout)
	defer ssoConn.Close()
	taskClient := ssov1.NewAuthClient(taskConn)
	_ = taskClient

	r := gin.Default()
	api := r.Group("/api/v1")
	api.POST("/auth/login", sso.LoginHandler(ssoClient))
	err := r.Run(cfg.Addr)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
