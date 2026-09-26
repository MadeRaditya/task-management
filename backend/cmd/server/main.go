package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/MadeRaditya/task-management/internal/cache"
	"github.com/MadeRaditya/task-management/internal/config"
	"github.com/MadeRaditya/task-management/internal/database"
	"github.com/MadeRaditya/task-management/internal/handler"
	"github.com/MadeRaditya/task-management/internal/repository"
	"github.com/MadeRaditya/task-management/internal/service"
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

	taskRepo := repository.NewTaskRepository(db)
	taskCache := cache.NewTaskCache(rdb)
	taskSvc := service.NewTaskService(taskRepo, taskCache)

	r := gin.Default()
	handler.RegisterRoutes(
		r,
		taskSvc,
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return db.PingContext(ctx)
		},
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return rdb.Ping(ctx).Err()
		},
	)

	log.Printf("Server running on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
