package jwt

import (
	"context"
	"fmt"
	"strings"
	"time"

	ssov1 "github.com/Citadelas/protos/golang/sso"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomClaims struct {
	UserID uint64 `json:"uid"`
	Email  string `json:"email"`
	AppID  int32  `json:"app_id"`
	jwt.RegisteredClaims
}

func ExtractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func ValidateToken(ctx context.Context, ssoClient ssov1.AuthClient, tokenString string) (uint64, error) {
	if tokenString == "" {
		return 0, fmt.Errorf("empty token")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &CustomClaims{})
	if err != nil {
		return 0, fmt.Errorf("malformed token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims structure")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return 0, fmt.Errorf("token expired")
	}

	_, err = ssoClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: int64(claims.UserID),
	})

	if err != nil {
		grpcErr, ok := status.FromError(err)
		if ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return 0, fmt.Errorf("invalid token: user not found")
			case codes.Unauthenticated:
				return 0, fmt.Errorf("invalid token: authentication failed")
			case codes.DeadlineExceeded:
				return 0, fmt.Errorf("token validation timeout")
			default:
				return 0, fmt.Errorf("token validation failed: %w", err)
			}
		}
		return 0, fmt.Errorf("SSO service unavailable: %w", err)
	}

	return claims.UserID, nil
}
