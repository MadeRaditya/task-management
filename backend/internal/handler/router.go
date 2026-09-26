package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/MadeRaditya/task-management/internal/service"
)

func RegisterRoutes(r *gin.Engine, taskSvc service.TaskService, dbPinger func() error, redisPinger func() error) {
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
