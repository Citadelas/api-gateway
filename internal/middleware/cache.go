package middleware

import (
	"bytes"
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func CacheMiddleware(log *slog.Logger, client *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != "GET" {
			ctx.Next()
			return
		}
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(401, "Unauthorized")
			return
		}
		cacheKey := ctx.Request.URL.String() + ":" + strconv.Itoa(int(userID.(uint64)))
		cached, err := client.Get(context.Background(), cacheKey).Result()
		if err == nil {
			ctx.Header("X-Cache-Status", "HIT")
			ctx.Header("Content-Type", "application/json")
			ctx.String(200, cached)
			ctx.Abort()
			return
		}
		ctx.Header("X-Cache-Status", "MISS")
		blw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		ctx.Next()
		if blw.Status() == 200 && blw.body.Len() > 0 {
			log.Info("Saving to cache, key:", cacheKey, "body length:", blw.body.Len())
			err := client.SetEx(context.Background(), cacheKey, blw.body.String(), time.Minute).Err()
			if err != nil {
				log.Error("Failed to save to Redis:", err)
			} else {
				log.Info("Successfully saved to Redis")
			}
		}
	}
}
