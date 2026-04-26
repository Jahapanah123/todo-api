package service

import (
	"fmt"

	"github.com/jahapanah123/todo/internal/domain"
)

type Repository interface {
	Create(todo domain.Todo) error
	Get(id string) (*domain.Todo, error)
	Update(id string, title string, description string) error
	Delete(id string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(todo domain.Todo) error {
	if todo.ID == "" || todo.Title == "" {
		return fmt.Errorf("ID and Title are required fields")
	}
	if err := s.repo.Create(todo); err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}

func (s *Service) Get(id string) (*domain.Todo, error) {

	if id == "" {
		return nil, fmt.Errorf("ID is required")
	}
	todo, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	return todo, nil
}

func (s *Service) Update(id string, title string, description string) error {
	if id == "" {
		return fmt.Errorf("ID is required")
	}
	if title == "" && description == "" {
		return fmt.Errorf("at least one of Title or Description must be provided")
	}
	if err := s.repo.Update(id, title, description); err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}

func (s *Service) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("ID is required")
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}
