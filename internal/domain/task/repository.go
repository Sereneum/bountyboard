package task

import (
	"context"
	"time"
)

type Repository interface {
	Create(task *Task) error
	ListByUser(userID string) ([]*Task, error)
	Complete(id string) error
}

// Cache интерфейс кеша для задач
type Cache interface {
	Get(userID string) ([]*Task, bool)
	Set(userID string, tasks []*Task)
	Invalidate(userID string)
}

// FileCache дополненный интерфейс кеша для задач с реализацией сохранения кэша в файл
type FileCache interface {
	Cache // Базовый интерфейс
	SaveToFile(filename string) error
	LoadFromFile(filename string) error
	RunAutoSave(ctx context.Context, filepath string, interval time.Duration)
	Janitor(ctx context.Context, interval time.Duration)
}
