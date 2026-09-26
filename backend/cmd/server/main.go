package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/MadeRaditya/task-management/internal/config"
	"github.com/MadeRaditya/task-management/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("gagal load config: %v", err)
	}

	db, err := database.NewMySQL(cfg.DSN())
	if err != nil {
		log.Fatalf("gagal konek MySQL: %v", err)
	}
	defer db.Close()
	log.Println("MySQL connected")

	rdb, err := database.NewRedis(cfg.RedisAddr(), cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("gagal konek Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("Redis connected")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		dbStatus := "ok"
		if err := db.PingContext(ctx); err != nil {
			dbStatus = "error: " + err.Error()
		}

		redisStatus := "ok"
		if err := rdb.Ping(ctx).Err(); err != nil {
			redisStatus = "error: " + err.Error()
		}

		status := http.StatusOK
		if dbStatus != "ok" || redisStatus != "ok" {
			status = http.StatusServiceUnavailable
		}

		c.JSON(status, gin.H{
			"mysql": dbStatus,
			"redis": redisStatus,
		})
	})

	log.Printf("Server running on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
