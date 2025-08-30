package main

import (
	"github.com/Citadelas/api-gateway/internal/config"
	"github.com/Citadelas/api-gateway/internal/handlers/sso"
	ssov1 "github.com/Citadelas/protos/golang/sso"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Starting app",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
	)

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
	api.POST("/auth/register", sso.RegisterHandler(ssoClient))
	api.POST("/auth/refresh", sso.RefreshToken(ssoClient))
	api.POST("/auth/isadmin", sso.IsAdmin(ssoClient))

	err := r.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
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
