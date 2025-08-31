package sso

import (
	"github.com/Citadelas/api-gateway/internal/helpers/grpc"
	ssov1 "github.com/Citadelas/protos/golang/sso"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Req struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(log *slog.Logger, client ssov1.AuthClient) gin.HandlerFunc {
	const op = "handlers.sso.Login"
	return func(c *gin.Context) {
		var req Req
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
		}
		grpcReq := ssov1.LoginRequest{
			AppId:    1,
			Email:    req.Email,
			Password: req.Password,
		}
		resp, err := client.Login(c, &grpcReq)
		if err != nil {
			grpc.HandleGRPCError(c, err)
			return
		}
		c.JSON(200, resp)
	}
}

func RegisterHandler(log *slog.Logger, client ssov1.AuthClient) gin.HandlerFunc {
	const op = "handlers.sso.Register"
	return func(c *gin.Context) {
		var req Req
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
		}
		grpcReq := ssov1.RegisterRequest{
			Email:    req.Email,
			Password: req.Password,
		}
		resp, err := client.Register(c, &grpcReq)
		if err != nil {
			grpc.HandleGRPCError(c, err)
			return
		}
		c.JSON(200, resp)
	}
}

type refReq struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshToken(log *slog.Logger, client ssov1.AuthClient) gin.HandlerFunc {
	const op = "handlers.sso.Refresh"
	return func(c *gin.Context) {
		var req refReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
		}
		grpcReq := ssov1.RefreshTokenRequest{
			AppId:        1,
			RefreshToken: req.RefreshToken,
		}
		resp, err := client.RefreshToken(c, &grpcReq)
		if err != nil {
			grpc.HandleGRPCError(c, err)
			return
		}
		c.JSON(200, resp)
	}
}

type adminReq struct {
	UserId int64 `json:"user_id"`
}

func IsAdmin(log *slog.Logger, client ssov1.AuthClient) gin.HandlerFunc {
	const op = "handlers.sso.Login"
	return func(c *gin.Context) {
		var req adminReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
		}
		grpcReq := ssov1.IsAdminRequest{
			UserId: req.UserId,
		}
		resp, err := client.IsAdmin(c, &grpcReq)
		if err != nil {
			grpc.HandleGRPCError(c, err)
			return
		}
		c.JSON(200, resp.GetIsAdmin())
	}
}
