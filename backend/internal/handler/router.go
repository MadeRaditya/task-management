package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/MadeRaditya/task-management/internal/service"
)

func RegisterRoutes(r *gin.Engine, taskSvc service.TaskService, dbPinger func() error, redisPinger func() error) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "http://127.0.0.1:8081", "http://localhost:19006"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/health", healthHandler(dbPinger, redisPinger))

	api := r.Group("/api")
	{
		tasks := api.Group("/tasks")
		th := NewTaskHandler(taskSvc)
		tasks.GET("", th.List)
		tasks.POST("", th.Create)
		tasks.GET("/:id", th.GetByID)
		tasks.PUT("/:id", th.Update)
		tasks.DELETE("/:id", th.Delete)
	}
}
