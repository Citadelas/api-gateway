package app

import (
	prometheus2 "github.com/Citadelas/api-gateway/internal/app/prometheus"
	"github.com/Citadelas/api-gateway/internal/handlers/sso"
	"github.com/Citadelas/api-gateway/internal/handlers/task"
	"github.com/Citadelas/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *App) setupRoutes() {

	prometheus.MustRegister(prometheus2.RequestsTotal, prometheus2.RequestDuration)

	a.router = gin.Default()

	// Add middleware
	a.router.Use(gin.Recovery())
	a.router.Use(gin.Logger())

	a.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := a.router.Group("/api/v1")

	// Public routes
	a.setupAuthRoutes(api)

	// Protected routes
	a.setupProtectedRoutes(api)
}

// setupAuthRoutes configures authentication routes
func (a *App) setupAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", sso.LoginHandler(a.log, a.ssoClient))
		auth.POST("/register", sso.RegisterHandler(a.log, a.ssoClient))
		auth.POST("/refresh", sso.RefreshToken(a.log, a.ssoClient))
		auth.POST("/isadmin", sso.IsAdmin(a.log, a.ssoClient))
	}
}

// setupProtectedRoutes configures protected routes
func (a *App) setupProtectedRoutes(api *gin.RouterGroup) {
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(a.ssoClient))
	protected.Use(middleware.CacheMiddleware(a.log, a.redis))
	// Task routes
	tasks := protected.Group("/tasks")
	{
		tasks.POST("", task.CreateTaskHandler(a.log, a.taskClient))
		tasks.GET("/:id", task.GetTaskHandler(a.log, a.taskClient))
		tasks.PUT("/:id", task.UpdateTaskHandler(a.log, a.taskClient))
		tasks.DELETE("/:id", task.DeleteTaskHandler(a.log, a.taskClient))
		tasks.PATCH("/:id/status", task.UpdateStatusHandler(a.log, a.taskClient))
	}
}
