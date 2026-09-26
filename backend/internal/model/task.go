package model

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	}
	return false
}

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	Assignee    string     `json:"assignee,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description" binding:"max=5000"`
	Status      TaskStatus `json:"status"`
	Assignee    string     `json:"assignee" binding:"max=100"`
}

type UpdateTaskRequest struct {
	Title       *string     `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string     `json:"description" binding:"omitempty,max=5000"`
	Status      *TaskStatus `json:"status"`
	Assignee    *string     `json:"assignee" binding:"omitempty,max=100"`
}

type ListTasksQuery struct {
	Status   string
	Keyword  string
	Assignee string
	Page     int
	Limit    int
	SortBy   string
	SortDir  string
}

func (q *ListTasksQuery) Normalize() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}

	allowed := map[string]bool{
		"id": true, "title": true, "status": true,
		"assignee": true, "created_at": true, "updated_at": true,
	}
	if !allowed[q.SortBy] {
		q.SortBy = "created_at"
	}

	q.SortDir = normalizeDir(q.SortDir)
}

func normalizeDir(d string) string {
	if d == "asc" || d == "ASC" {
		return "asc"
	}
	return "desc"
}

func (q *ListTasksQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}
