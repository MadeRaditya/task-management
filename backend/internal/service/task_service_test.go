package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MadeRaditya/task-management/internal/apperror"
	"github.com/MadeRaditya/task-management/internal/model"
	"github.com/MadeRaditya/task-management/internal/repository"
)

type fakeRepo struct {
	tasks  map[int64]*model.Task
	nextID int64
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{tasks: map[int64]*model.Task{}, nextID: 1}
}

func (f *fakeRepo) Create(ctx context.Context, t *model.Task) (int64, error) {
	id := f.nextID
	f.nextID++
	cp := *t
	cp.ID = id
	f.tasks[id] = &cp
	return id, nil
}

func (f *fakeRepo) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	t, ok := f.tasks[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (f *fakeRepo) List(ctx context.Context, q model.ListTasksQuery) ([]model.Task, int64, error) {
	out := make([]model.Task, 0, len(f.tasks))
	for _, t := range f.tasks {
		out = append(out, *t)
	}
	return out, int64(len(out)), nil
}

func (f *fakeRepo) Update(ctx context.Context, t *model.Task) error {
	if _, ok := f.tasks[t.ID]; !ok {
		return repository.ErrNotFound
	}
	cp := *t
	f.tasks[t.ID] = &cp
	return nil
}

func (f *fakeRepo) SoftDelete(ctx context.Context, id int64) error {
	if _, ok := f.tasks[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.tasks, id)
	return nil
}

func (f *fakeRepo) ExistsActiveTitle(ctx context.Context, title string, excludeID int64) (bool, error) {
	for id, t := range f.tasks {
		if id != excludeID && t.Title == title {
			return true, nil
		}
	}
	return false, nil
}

type fakeCache struct {
	invalidateCalls int
}

func (c *fakeCache) Get(ctx context.Context, q model.ListTasksQuery) (*ListResult, bool, error) {
	return nil, false, nil
}
func (c *fakeCache) Set(ctx context.Context, q model.ListTasksQuery, r *ListResult) error { return nil }
func (c *fakeCache) InvalidateTasks(ctx context.Context) error {
	c.invalidateCalls++
	return nil
}

func asAppError(t *testing.T, err error) *apperror.AppError {
	t.Helper()
	var appErr *apperror.AppError
	require.True(t, errors.As(err, &appErr), "expected *AppError, got %T", err)
	return appErr
}

func TestUpdate_Success(t *testing.T) {
	repo := newFakeRepo()
	cache := &fakeCache{}
	svc := NewTaskService(repo, cache)

	_, _ = repo.Create(context.Background(), &model.Task{Title: "Task A", Status: model.StatusTodo})

	newStatus := model.StatusDone
	updated, err := svc.Update(context.Background(), 1, model.UpdateTaskRequest{Status: &newStatus})

	require.NoError(t, err)
	assert.Equal(t, model.StatusDone, updated.Status)
	assert.Equal(t, 1, cache.invalidateCalls, "cache harus di-invalidate setelah update sukses")
}

func TestUpdate_DuplicateTitle_Returns409(t *testing.T) {
	repo := newFakeRepo()
	cache := &fakeCache{}
	svc := NewTaskService(repo, cache)

	_, _ = repo.Create(context.Background(), &model.Task{Title: "Task A"})
	_, _ = repo.Create(context.Background(), &model.Task{Title: "Task B"})

	newTitle := "Task A"
	_, err := svc.Update(context.Background(), 2, model.UpdateTaskRequest{Title: &newTitle})

	require.Error(t, err)
	appErr := asAppError(t, err)
	assert.Equal(t, 409, appErr.Status)
	assert.Equal(t, "DUPLICATE_TITLE", appErr.Code)
	assert.Equal(t, 0, cache.invalidateCalls, "cache TIDAK boleh di-invalidate kalau update gagal")
}

func TestUpdate_NotFound(t *testing.T) {
	repo := newFakeRepo()
	svc := NewTaskService(repo, &fakeCache{})

	newStatus := model.StatusDone
	_, err := svc.Update(context.Background(), 999, model.UpdateTaskRequest{Status: &newStatus})

	require.Error(t, err)
	appErr := asAppError(t, err)
	assert.Equal(t, 404, appErr.Status)
	assert.Equal(t, "TASK_NOT_FOUND", appErr.Code)
}

func TestDelete_InvalidatesCache(t *testing.T) {
	repo := newFakeRepo()
	cache := &fakeCache{}
	svc := NewTaskService(repo, cache)

	_, _ = repo.Create(context.Background(), &model.Task{Title: "Task A"})

	err := svc.Delete(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, cache.invalidateCalls)
}

func TestDelete_NotFound(t *testing.T) {
	repo := newFakeRepo()
	svc := NewTaskService(repo, &fakeCache{})

	err := svc.Delete(context.Background(), 42)
	require.Error(t, err)
	appErr := asAppError(t, err)
	assert.Equal(t, 404, appErr.Status)
}

func TestCreate_InvalidStatus(t *testing.T) {
	repo := newFakeRepo()
	svc := NewTaskService(repo, &fakeCache{})

	_, err := svc.Create(context.Background(), model.CreateTaskRequest{
		Title:  "Task X",
		Status: "invalid_status",
	})
	require.Error(t, err)
	appErr := asAppError(t, err)
	assert.Equal(t, 400, appErr.Status)
	assert.Equal(t, "INVALID_STATUS", appErr.Code)
}
