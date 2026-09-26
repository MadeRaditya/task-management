package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MadeRaditya/task-management/internal/model"
)

func newMockRepo(t *testing.T) (TaskRepository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return NewTaskRepository(db), mock
}

func TestList_SearchByKeyword_AndStatus_AndPagination(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM tasks")).
		WithArgs(
			"todo",
			"%meeting%",
			"%meeting%",
			"alice",
		).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery("SELECT id, title, description, status, assignee, created_at, updated_at").
		WithArgs(
			"todo", "%meeting%", "%meeting%", "alice", 10, 20,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "description", "status", "assignee", "created_at", "updated_at",
		}).
			AddRow(1, "Meeting A", "desc", "todo", "alice", now(), now()).
			AddRow(2, "Meeting B", "desc", "todo", "alice", now(), now()))

	tasks, total, err := repo.List(context.Background(), model.ListTasksQuery{
		Status:   "todo",
		Keyword:  "meeting",
		Assignee: "alice",
		Page:     3,
		Limit:    10,
		SortBy:   "created_at",
		SortDir:  "desc",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, tasks, 2)
	assert.Equal(t, "Meeting A", tasks[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_OnlySoftDeleteFilter_WhenNoFiltersProvided(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM tasks")).
		WithArgs(). // tidak ada filter
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery("SELECT id, title, description, status, assignee, created_at, updated_at").
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "description", "status", "assignee", "created_at", "updated_at",
		}))

	tasks, total, err := repo.List(context.Background(), model.ListTasksQuery{
		Page: 1, Limit: 10, SortBy: "created_at", SortDir: "desc",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, tasks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSoftDelete_SetsDeletedAt(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL")).
		WithArgs(int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.SoftDelete(context.Background(), 7)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSoftDelete_NotFound(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectExec("UPDATE tasks SET deleted_at").
		WithArgs(int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.SoftDelete(context.Background(), 7)
	assert.ErrorIs(t, err, ErrNotFound)
}
