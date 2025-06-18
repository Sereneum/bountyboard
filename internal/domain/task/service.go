package task

import (
	"github.com/google/uuid"
)

type Service interface {
	CreateTask(userID, title, description string, amount int) error
	ListTasks(userID string) ([]*Task, error)
	CompleteTask(id string) error

	Cache() Cache
}

type service struct {
	repo  Repository
	cache Cache
}

func New(r Repository, c Cache) Service {
	return &service{repo: r, cache: c}
}

func (s *service) CreateTask(userID, title, description string, amount int) error {
	task := &Task{
		ID:           uuid.NewString(),
		Title:        title,
		Description:  description,
		Done:         false,
		BountyAmount: amount,
		UserID:       userID,
	}

	if err := s.repo.Create(task); err != nil {
		return err
	}

	if s.cache != nil {
		s.cache.Invalidate(userID)
	}

	return nil
}

func (s *service) ListTasks(userID string) ([]*Task, error) {
	if s.cache != nil {
		if cached, ok := s.cache.Get(userID); ok {
			return cached, nil
		}
	}

	tasks, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		s.cache.Set(userID, tasks)
	}

	return tasks, nil
}

func (s *service) CompleteTask(id string) error {
	return s.repo.Complete(id)
}

func (s *service) Cache() Cache {
	return s.cache
}
