package repository

import (
	"strconv"
	"testing"

	"github.com/jahapanah123/todo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryRepository(t *testing.T) {
	repo := NewInMemoryRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.data)
}

func TestRepository_Create(t *testing.T) {
	repo := NewInMemoryRepository()

	todo := domain.Todo{
		ID:          "1",
		Title:       "Test Todo",
		Description: "Test Description",
	}

	err := repo.Create(todo)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(repo.data))
}

func TestRepository_Get(t *testing.T) {
	repo := NewInMemoryRepository()

	todo := domain.Todo{
		ID:          "1",
		Title:       "Test Todo",
		Description: "Test Description",
	}
	err := repo.Create(todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	t.Run("Get existing todo", func(t *testing.T) {
		result, err := repo.Get("1")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Todo", result.Title)
		assert.Equal(t, "Test Description", result.Description)
	})

	t.Run("Get non-existing todo", func(t *testing.T) {
		result, err := repo.Get("999")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestRepository_Update(t *testing.T) {
	repo := NewInMemoryRepository()

	todo := domain.Todo{
		ID:          "1",
		Title:       "Test Todo",
		Description: "Test Description",
	}
	err := repo.Create(todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	t.Run("Update existing todo", func(t *testing.T) {
		err := repo.Update("1", "Updated Title", "Updated Description")
		assert.NoError(t, err)

		updated, _ := repo.Get("1")
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, "Updated Description", updated.Description)
	})

	t.Run("Update non-existing todo", func(t *testing.T) {
		err := repo.Update("999", "Title", "Description")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestRepository_Delete(t *testing.T) {
	repo := NewInMemoryRepository()

	todo := domain.Todo{
		ID:          "1",
		Title:       "Test Todo",
		Description: "Test Description",
	}
	err := repo.Create(todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	t.Run("Delete existing todo", func(t *testing.T) {
		err := repo.Delete("1")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(repo.data))
	})

	t.Run("Delete non-existing todo", func(t *testing.T) {
		err := repo.Delete("999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestRepository_ConcurrentAccess(t *testing.T) {
	repo := NewInMemoryRepository()

	// Test concurrent writes
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(id int) {
			todo := domain.Todo{
				ID:    strconv.Itoa(id),
				Title: "Concurrent Todo",
			}
			err := repo.Create(todo)
			if err != nil {
				t.Errorf("Create failed: %v", err)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}

	assert.Equal(t, 100, len(repo.data))
}
