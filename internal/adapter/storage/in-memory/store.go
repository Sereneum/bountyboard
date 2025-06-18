package inmemory

import (
	"bountyboard/internal/domain/task"
	"sync"
)

type Repo struct {
	mu    sync.RWMutex
	tasks map[string]*task.Task
}

func NewRepo() *Repo {
	return &Repo{tasks: make(map[string]*task.Task)}
}

func (r *Repo) Create(t *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[t.ID] = t
	return nil
}

func (r *Repo) ListByUser(userID string) ([]*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*task.Task
	for _, t := range r.tasks {
		if t.UserID == userID {
			res = append(res, t)
		}
	}
	return res, nil
}

func (r *Repo) Complete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, ok := r.tasks[id]; ok {
		t.Done = true
	}
	return nil
}
