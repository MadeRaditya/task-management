package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func healthHandler(dbPing func() error, redisPing func() error) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		dbStatus := "ok"
		if err := dbPing(); err != nil {
			dbStatus = "error: " + err.Error()
		}
		redisStatus := "ok"
		if err := redisPing(); err != nil {
			redisStatus = "error: " + err.Error()
		}

		status := http.StatusOK
		if dbStatus != "ok" || redisStatus != "ok" {
			status = http.StatusServiceUnavailable
		}

		_ = ctx
		c.JSON(status, gin.H{"mysql": dbStatus, "redis": redisStatus})
	}
}
