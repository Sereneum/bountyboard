package memory

import (
	"bountyboard/internal/domain/task"
	"context"
	"encoding/gob"
	"log/slog"
	"os"
	"sync"
	"time"
)

type cacheEntry struct {
	Tasks      []*task.Task
	Expiration time.Time
}

type TaskCache struct {
	data map[string]cacheEntry
	mu   sync.RWMutex
	ttl  time.Duration
}

func NewTaskCache(ttl time.Duration) *TaskCache {
	return &TaskCache{
		ttl:  ttl,
		data: make(map[string]cacheEntry),
	}
}

// Get возвращает кэшированные задачи для userID, если они существуют и не устарели
func (c *TaskCache) Get(userID string) ([]*task.Task, bool) {
	c.mu.RLock()
	entry, found := c.data[userID]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}
	if time.Now().After(entry.Expiration) {
		// Кэш просрочен, но мы не удаляем его здесь
		return nil, false
	}
	return entry.Tasks, true
}

// Set сохраняет задачи в кэш с TTL
func (c *TaskCache) Set(userID string, tasks []*task.Task) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[userID] = cacheEntry{
		Tasks:      tasks,
		Expiration: time.Now().Add(c.ttl),
	}
}

// Invalidate удаляет кэш для конкретного пользователя
func (c *TaskCache) Invalidate(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, userID)
}

// Janitor чистит устаревший кэш
func (c *TaskCache) Janitor(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpired()
		case <-ctx.Done():
			return
		}
	}
}

func (c *TaskCache) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for userID, entry := range c.data {
		if now.After(entry.Expiration) {
			delete(c.data, userID)
		}
	}
}

func (c *TaskCache) SaveToFile(path string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	return enc.Encode(c.data)
}

func (c *TaskCache) LoadFromFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	return dec.Decode(&c.data)
}

func (c *TaskCache) RunAutoSave(ctx context.Context, filepath string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			if err := c.SaveToFile(filepath); err != nil {
				slog.Error("auto-save failed", slog.String("err", err.Error()))
			} else {
				slog.Debug("cache auto-saved", slog.String("path", filepath))
			}
		}
	}
}
