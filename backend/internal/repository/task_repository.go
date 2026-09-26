package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/MadeRaditya/task-management/internal/model"
)

var ErrNotFound = errors.New("task not found")

type TaskRepository interface {
	Create(ctx context.Context, t *model.Task) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Task, error)
	List(ctx context.Context, q model.ListTasksQuery) ([]model.Task, int64, error)
	Update(ctx context.Context, t *model.Task) error
	SoftDelete(ctx context.Context, id int64) error
	ExistsActiveTitle(ctx context.Context, title string, excludeID int64) (bool, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, t *model.Task) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (title, description, status, assignee) VALUES (?, ?, ?, ?)`,
		t.Title, t.Description, t.Status, t.Assignee,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *taskRepository) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, status, assignee, created_at, updated_at
		 FROM tasks
		 WHERE id = ? AND deleted_at IS NULL`, id)

	var t model.Task
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Assignee, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *taskRepository) List(ctx context.Context, q model.ListTasksQuery) ([]model.Task, int64, error) {
	q.Normalize()

	where := []string{"deleted_at IS NULL"}
	args := []any{}

	if q.Status != "" {
		where = append(where, "status = ?")
		args = append(args, q.Status)
	}
	if q.Keyword != "" {
		where = append(where, "(title LIKE ? OR description LIKE ?)")
		kw := "%" + q.Keyword + "%"
		args = append(args, kw, kw)
	}
	if q.Assignee != "" {
		where = append(where, "assignee = ?")
		args = append(args, q.Assignee)
	}

	whereSQL := "WHERE " + strings.Join(where, " AND ")

	var total int64
	countSQL := "SELECT COUNT(*) FROM tasks " + whereSQL
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf(
		`SELECT id, title, description, status, assignee, created_at, updated_at
		 FROM tasks
		 %s
		 ORDER BY %s %s
		 LIMIT ? OFFSET ?`,
		whereSQL, q.SortBy, strings.ToUpper(q.SortDir),
	)

	listArgs := append(args, q.Limit, q.Offset())
	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0, q.Limit)
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Assignee, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) Update(ctx context.Context, t *model.Task) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE tasks
		 SET title = ?, description = ?, status = ?, assignee = ?
		 WHERE id = ? AND deleted_at IS NULL`,
		t.Title, t.Description, t.Status, t.Assignee, t.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		if _, err := r.GetByID(ctx, t.ID); err != nil {
			return err
		}
	}
	return nil
}

func (r *taskRepository) SoftDelete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *taskRepository) ExistsActiveTitle(ctx context.Context, title string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tasks
		 WHERE title = ? AND deleted_at IS NULL AND id <> ?`,
		title, excludeID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
