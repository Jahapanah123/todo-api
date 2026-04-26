package service

import (
	"fmt"
	"testing"

	"github.com/jahapanah123/todo/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(todo domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockRepository) Get(id string) (*domain.Todo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Todo), args.Error(1)
}

func (m *MockRepository) Update(id string, title string, description string) error {
	args := m.Called(id, title, description)
	return args.Error(0)
}

func (m *MockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestNewService(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	assert.NotNil(t, svc)
}

func TestService_Create(t *testing.T) {
	t.Run("Valid todo creation", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		todo := domain.Todo{
			ID:          "1",
			Title:       "Test Todo",
			Description: "Test Description",
		}

		mockRepo.On("Create", todo).Return(nil)

		err := svc.Create(todo)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Missing ID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		todo := domain.Todo{
			Title:       "Test Todo",
			Description: "Test Description",
		}

		err := svc.Create(todo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID and Title are required")
	})

	t.Run("Missing Title", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		todo := domain.Todo{
			ID:          "1",
			Description: "Test Description",
		}

		err := svc.Create(todo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID and Title are required")
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		todo := domain.Todo{
			ID:    "1",
			Title: "Test Todo",
		}

		mockRepo.On("Create", todo).Return(fmt.Errorf("database error"))

		err := svc.Create(todo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create todo")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	t.Run("Get existing todo", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		expectedTodo := &domain.Todo{
			ID:          "1",
			Title:       "Test Todo",
			Description: "Test Description",
		}

		mockRepo.On("Get", "1").Return(expectedTodo, nil)

		todo, err := svc.Get("1")
		assert.NoError(t, err)
		assert.Equal(t, expectedTodo, todo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		todo, err := svc.Get("")
		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		mockRepo.On("Get", "999").Return(nil, fmt.Errorf("not found"))

		todo, err := svc.Get("999")
		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Contains(t, err.Error(), "failed to get todo")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("Valid update", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		mockRepo.On("Update", "1", "New Title", "New Description").Return(nil)

		err := svc.Update("1", "New Title", "New Description")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		err := svc.Update("", "Title", "Description")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("Empty title and description", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		err := svc.Update("1", "", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one of Title or Description must be provided")
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		mockRepo.On("Update", "1", "Title", "Description").Return(fmt.Errorf("update failed"))

		err := svc.Update("1", "Title", "Description")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update todo")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("Valid delete", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		mockRepo.On("Delete", "1").Return(nil)

		err := svc.Delete("1")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		err := svc.Delete("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		svc := NewService(mockRepo)

		mockRepo.On("Delete", "999").Return(fmt.Errorf("not found"))

		err := svc.Delete("999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete todo")
		mockRepo.AssertExpectations(t)
	})
}
