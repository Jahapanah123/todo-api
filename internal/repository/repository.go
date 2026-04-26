package repository

import (
	"fmt"
	"sync"

	"github.com/jahapanah123/todo/internal/domain"
)

type Repository struct {
	data map[string]domain.Todo
	mu   sync.RWMutex
}

func NewInMemoryRepository() *Repository {
	return &Repository{
		data: make(map[string]domain.Todo),
	}
}

func (r *Repository) Create(todo domain.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[todo.ID] = todo
	return nil

}

func (r *Repository) Get(id string) (*domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.data[id]

	if !exists {
		return nil, fmt.Errorf("Todo with ID %s not found", id)
	}
	return &todo, nil
}

func (r *Repository) Update(id string, title string, description string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return fmt.Errorf("Todo with ID %s not found", id)
	}
	todo := r.data[id]
	todo.Title = title
	todo.Description = description
	r.data[id] = todo
	return nil
}

func (r *Repository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return fmt.Errorf("Todo with ID %s not found", id)
	}
	delete(r.data, id)
	return nil
}
