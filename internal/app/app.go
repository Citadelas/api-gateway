package app

import (
	"context"
	"github.com/Citadelas/api-gateway/internal/config"
	ssov1 "github.com/Citadelas/protos/golang/sso"
	taskv1 "github.com/Citadelas/protos/golang/task"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	cfg        *config.Config
	log        *slog.Logger
	ssoClient  ssov1.AuthClient
	taskClient taskv1.TaskServiceClient
	//ssoConn    *grpc.ClientConn
	//taskConn   *grpc.ClientConn
	router *gin.Engine
}

func NewApp() (*App, error) {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("Starting app",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
	)

	app := &App{
		cfg: cfg,
		log: log,
	}

	if err := app.mustInitClients(); err != nil {
		return nil, err
	}

	app.setupRoutes()
	return app, nil
}

func (a *App) Run() {
	srv := &http.Server{
		Addr:    a.cfg.Addr,
		Handler: a.router,
	}

	go func() {
		a.log.Info("Starting HTTP server", slog.String("addr", a.cfg.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("Failed to start server", slog.String("error", err.Error()))
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.log.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	a.log.Info("Server exited")
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
