package service

import (
	"context"
	"strings"

	"github.com/MadeRaditya/task-management/internal/apperror"
	"github.com/MadeRaditya/task-management/internal/model"
	"github.com/MadeRaditya/task-management/internal/repository"
)

type TaskService interface {
	Create(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error)
	GetByID(ctx context.Context, id int64) (*model.Task, error)
	List(ctx context.Context, q model.ListTasksQuery) (*ListResult, error)
	Update(ctx context.Context, id int64, req model.UpdateTaskRequest) (*model.Task, error)
	Delete(ctx context.Context, id int64) error
}

type ListResult struct {
	Items      []model.Task `json:"items"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	TotalPages int          `json:"total_pages"`
}

type TaskCacheReader interface {
	Get(ctx context.Context, q model.ListTasksQuery) (*ListResult, bool, error)
	Set(ctx context.Context, q model.ListTasksQuery, result *ListResult) error
	InvalidateTasks(ctx context.Context) error
}

type taskService struct {
	repo  repository.TaskRepository
	cache TaskCacheReader
}

func NewTaskService(repo repository.TaskRepository, cache TaskCacheReader) TaskService {
	if cache == nil {
		cache = &noopCache{}
	}
	return &taskService{repo: repo, cache: cache}
}

type noopCache struct{}

func (noopCache) Get(ctx context.Context, q model.ListTasksQuery) (*ListResult, bool, error) {
	return nil, false, nil
}
func (noopCache) Set(ctx context.Context, q model.ListTasksQuery, r *ListResult) error { return nil }
func (noopCache) InvalidateTasks(ctx context.Context) error                            { return nil }

func (s *taskService) Create(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, apperror.BadRequest("INVALID_TITLE", "Title is required")
	}

	exists, err := s.repo.ExistsActiveTitle(ctx, title, 0)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	if exists {
		return nil, apperror.Conflict("DUPLICATE_TITLE", "Title already exists")
	}

	status := req.Status
	if status == "" {
		status = model.StatusTodo
	}
	if !status.IsValid() {
		return nil, apperror.BadRequest("INVALID_STATUS", "Status must be todo, in_progress, or done")
	}

	t := &model.Task{
		Title:       title,
		Description: strings.TrimSpace(req.Description),
		Status:      status,
		Assignee:    strings.TrimSpace(req.Assignee),
	}
	id, err := s.repo.Create(ctx, t)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	t.ID = id

	created, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	if err := s.cache.InvalidateTasks(ctx); err != nil {
		_ = err
	}

	return created, nil
}

func (s *taskService) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, apperror.NotFound("TASK_NOT_FOUND", "Task not found")
		}
		return nil, apperror.Internal(err)
	}
	return t, nil
}

func (s *taskService) List(ctx context.Context, q model.ListTasksQuery) (*ListResult, error) {
	if q.Status != "" && !model.TaskStatus(q.Status).IsValid() {
		return nil, apperror.BadRequest("INVALID_STATUS", "Status filter must be todo, in_progress, or done")
	}
	q.Normalize()

	if cached, hit, err := s.cache.Get(ctx, q); err == nil && hit {
		return cached, nil
	}

	tasks, total, err := s.repo.List(ctx, q)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(q.Limit) - 1) / int64(q.Limit))
	}

	result := &ListResult{
		Items:      tasks,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}

	_ = s.cache.Set(ctx, q, result)

	return result, nil
}

func (s *taskService) Update(ctx context.Context, id int64, req model.UpdateTaskRequest) (*model.Task, error) {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, apperror.NotFound("TASK_NOT_FOUND", "Task not found")
		}
		return nil, apperror.Internal(err)
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, apperror.BadRequest("INVALID_TITLE", "Title cannot be empty")
		}
		if title != current.Title {
			exists, err := s.repo.ExistsActiveTitle(ctx, title, id)
			if err != nil {
				return nil, apperror.Internal(err)
			}
			if exists {
				return nil, apperror.Conflict("DUPLICATE_TITLE", "Title already exists")
			}
		}
		current.Title = title
	}
	if req.Description != nil {
		current.Description = strings.TrimSpace(*req.Description)
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, apperror.BadRequest("INVALID_STATUS", "Status must be todo, in_progress, or done")
		}
		current.Status = *req.Status
	}
	if req.Assignee != nil {
		current.Assignee = strings.TrimSpace(*req.Assignee)
	}

	if err := s.repo.Update(ctx, current); err != nil {
		if err == repository.ErrNotFound {
			return nil, apperror.NotFound("TASK_NOT_FOUND", "Task not found")
		}
		return nil, apperror.Internal(err)
	}

	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	if err := s.cache.InvalidateTasks(ctx); err != nil {
		_ = err
	}
	return updated, nil
}

func (s *taskService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			return apperror.NotFound("TASK_NOT_FOUND", "Task not found")
		}
		return apperror.Internal(err)
	}
	if err := s.cache.InvalidateTasks(ctx); err != nil {
		_ = err
	}
	return nil
}
