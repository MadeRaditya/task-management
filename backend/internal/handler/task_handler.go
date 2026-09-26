package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MadeRaditya/task-management/internal/model"
	"github.com/MadeRaditya/task-management/internal/service"
)

type TaskHandler struct {
	svc service.TaskService
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) List(c *gin.Context) {
	q := model.ListTasksQuery{
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
		Assignee: c.Query("assignee"),
		SortBy:   c.Query("sort_by"),
		SortDir:  c.Query("sort_dir"),
	}
	if p, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		q.Page = p
	}
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil {
		q.Limit = l
	}
	q.Normalize()

	result, err := h.svc.List(c.Request.Context(), q)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		respondError(c, err)
		return
	}
	t, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Data: t})
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req model.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, badRequestFromBinding(err))
		return
	}
	t, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{Data: t})
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		respondError(c, err)
		return
	}
	var req model.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, badRequestFromBinding(err))
		return
	}
	t, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Data: t})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		respondError(c, err)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return 0, apperrorBadRequest()
	}
	return id, nil
}
