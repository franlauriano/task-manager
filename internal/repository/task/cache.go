package task

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"taskmanager/internal/entity/task"
	"taskmanager/internal/platform/cache"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const cacheKeyPrefix = "tasks:list:"

// cachedDatasource wraps a Persistent implementation with Redis cache-aside logic.
type cachedDatasource struct {
	next   Persistent
	client *redis.Client
	ttl    time.Duration
}

// NewCachedPersist creates a cache-aside decorator around the given Persistent implementation.
func NewCachedPersist(next Persistent, client *redis.Client, ttl time.Duration) Persistent {
	return &cachedDatasource{
		next:   next,
		client: client,
		ttl:    ttl,
	}
}

// ListPaginated checks the cache first; on miss, queries the database and caches the result.
func (c *cachedDatasource) ListPaginated(ctx context.Context, statusFilter *task.TaskStatus, page, limit int) (*task.ListTasks, error) {
	key := listCacheKey(statusFilter, page, limit)

	result, err := cache.Get[task.ListTasks](ctx, c.client, key)
	if err != nil {
		slog.Warn("Cache get error, falling back to database", "key", key, "error", err)
	}
	if result != nil {
		return result, nil
	}

	result, err = c.next.ListPaginated(ctx, statusFilter, page, limit)
	if err != nil {
		return nil, err
	}

	if err := cache.Set(ctx, c.client, key, result, c.ttl); err != nil {
		slog.Warn("Cache set error", "key", key, "error", err)
	}

	return result, nil
}

// Create delegates to the next implementation and invalidates the list cache.
func (c *cachedDatasource) Create(ctx context.Context, t *task.Task) error {
	if err := c.next.Create(ctx, t); err != nil {
		return err
	}
	c.invalidateListCache(ctx)
	return nil
}

// Update delegates to the next implementation and invalidates the list cache.
func (c *cachedDatasource) Update(ctx context.Context, taskUUID uuid.UUID, t *task.Task) error {
	if err := c.next.Update(ctx, taskUUID, t); err != nil {
		return err
	}
	c.invalidateListCache(ctx)
	return nil
}

// Delete delegates to the next implementation and invalidates the list cache.
func (c *cachedDatasource) Delete(ctx context.Context, taskUUID uuid.UUID) error {
	if err := c.next.Delete(ctx, taskUUID); err != nil {
		return err
	}
	c.invalidateListCache(ctx)
	return nil
}

// UpdateStatus delegates to the next implementation and invalidates the list cache.
func (c *cachedDatasource) UpdateStatus(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
	if err := c.next.UpdateStatus(ctx, taskUUID, updates); err != nil {
		return err
	}
	c.invalidateListCache(ctx)
	return nil
}

// RetrieveByUUID delegates directly to the next implementation (no cache).
func (c *cachedDatasource) RetrieveByUUID(ctx context.Context, taskUUID uuid.UUID) (*task.Task, error) {
	return c.next.RetrieveByUUID(ctx, taskUUID)
}

// ListByTeamID delegates directly to the next implementation (no cache).
func (c *cachedDatasource) ListByTeamID(ctx context.Context, teamID uint) ([]task.Task, error) {
	return c.next.ListByTeamID(ctx, teamID)
}

// invalidateListCache removes all cached list entries.
func (c *cachedDatasource) invalidateListCache(ctx context.Context) {
	if err := cache.DeleteByPrefix(ctx, c.client, cacheKeyPrefix); err != nil {
		slog.Warn("Cache invalidation error", "prefix", cacheKeyPrefix, "error", err)
	}
}

// listCacheKey builds a deterministic cache key for a paginated list query.
func listCacheKey(statusFilter *task.TaskStatus, page, limit int) string {
	status := "all"
	if statusFilter != nil {
		status = string(*statusFilter)
	}
	return fmt.Sprintf("%sstatus=%s:page=%d:limit=%d", cacheKeyPrefix, status, page, limit)
}
