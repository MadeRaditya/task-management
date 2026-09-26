package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/MadeRaditya/task-management/internal/model"
	"github.com/MadeRaditya/task-management/internal/service"
)

const (
	taskCachePrefix = "tasks:list:"
	taskCacheTTL    = 60 * time.Second
	tasksVersionKey = "tasks:version"
)

type TaskCache struct {
	rdb *redis.Client
}

func NewTaskCache(rdb *redis.Client) *TaskCache {
	return &TaskCache{rdb: rdb}
}

func (c *TaskCache) cacheKey(ctx context.Context, q model.ListTasksQuery) (string, error) {
	version, err := c.rdb.Get(ctx, tasksVersionKey).Int64()
	if err != nil && err != redis.Nil {
		return "", err
	}

	raw := fmt.Sprintf(
		"v=%d|status=%s|keyword=%s|assignee=%s|page=%d|limit=%d|sort_by=%s|sort_dir=%s",
		version, q.Status, q.Keyword, q.Assignee, q.Page, q.Limit, q.SortBy, q.SortDir,
	)
	return taskCachePrefix + raw, nil
}

func (c *TaskCache) Get(ctx context.Context, q model.ListTasksQuery) (*service.ListResult, bool, error) {
	key, err := c.cacheKey(ctx, q)
	if err != nil {
		return nil, false, err
	}

	raw, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var result service.ListResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, false, err
	}
	return &result, true, nil
}

func (c *TaskCache) Set(ctx context.Context, q model.ListTasksQuery, result *service.ListResult) error {
	key, err := c.cacheKey(ctx, q)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, payload, taskCacheTTL).Err()
}

func (c *TaskCache) InvalidateTasks(ctx context.Context) error {
	return c.rdb.Incr(ctx, tasksVersionKey).Err()
}
