package sso

import (
	ssov1 "github.com/Citadelas/protos/golang/sso"
	"github.com/gin-gonic/gin"
)

type loginReq struct {
	AppId    int32  `json:"app_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(client ssov1.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		grpcReq := ssov1.LoginRequest{
			AppId:    req.AppId,
			Email:    req.Email,
			Password: req.Password,
		}
		resp, err := client.Login(c, &grpcReq)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		c.JSON(200, resp)
	}
}
