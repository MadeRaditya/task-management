package cache

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MadeRaditya/task-management/internal/model"
	"github.com/MadeRaditya/task-management/internal/service"
)

func setupCache(t *testing.T) (*TaskCache, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	return NewTaskCache(rdb), mr
}

func TestTaskCache_SetGet(t *testing.T) {
	c, _ := setupCache(t)
	ctx := context.Background()

	q := model.ListTasksQuery{Page: 1, Limit: 10, SortBy: "created_at", SortDir: "desc"}
	result := &service.ListResult{
		Items:      []model.Task{{ID: 1, Title: "A", Status: model.StatusTodo}},
		Total:      1,
		Page:       1,
		Limit:      10,
		TotalPages: 1,
	}

	_, hit, err := c.Get(ctx, q)
	require.NoError(t, err)
	assert.False(t, hit, "cache harus miss sebelum di-Set")

	require.NoError(t, c.Set(ctx, q, result))

	got, hit, err := c.Get(ctx, q)
	require.NoError(t, err)
	require.True(t, hit)
	assert.Equal(t, int64(1), got.Total)
	assert.Len(t, got.Items, 1)
	assert.Equal(t, "A", got.Items[0].Title)
}

func TestTaskCache_DifferentQuery_DifferentKey(t *testing.T) {
	c, _ := setupCache(t)
	ctx := context.Background()

	q1 := model.ListTasksQuery{Page: 1, Limit: 10, SortBy: "created_at", SortDir: "desc"}
	q2 := model.ListTasksQuery{Page: 2, Limit: 10, SortBy: "created_at", SortDir: "desc"}

	require.NoError(t, c.Set(ctx, q1, &service.ListResult{Total: 1, Page: 1}))

	_, hit, err := c.Get(ctx, q2)
	require.NoError(t, err)
	assert.False(t, hit, "query berbeda harus menghasilkan cache key berbeda")
}

func TestTaskCache_Invalidate(t *testing.T) {
	c, _ := setupCache(t)
	ctx := context.Background()

	q := model.ListTasksQuery{Page: 1, Limit: 10, SortBy: "created_at", SortDir: "desc"}
	require.NoError(t, c.Set(ctx, q, &service.ListResult{Total: 5}))

	require.NoError(t, c.InvalidateTasks(ctx))

	_, hit, err := c.Get(ctx, q)
	require.NoError(t, err)
	assert.False(t, hit, "cache harus miss setelah invalidate")
}

func TestTaskCache_TTL(t *testing.T) {
	c, mr := setupCache(t)
	ctx := context.Background()

	q := model.ListTasksQuery{Page: 1, Limit: 10, SortBy: "created_at", SortDir: "desc"}
	require.NoError(t, c.Set(ctx, q, &service.ListResult{Total: 1}))

	mr.FastForward(61 * 1e9)

	_, hit, err := c.Get(ctx, q)
	require.NoError(t, err)
	assert.False(t, hit, "cache harus expired setelah 60 detik")
}
