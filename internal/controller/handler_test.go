package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jahapanah123/todo/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) Create(todo domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockService) Get(id string) (*domain.Todo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Todo), args.Error(1)
}

func (m *MockService) Update(id string, title string, description string) error {
	args := m.Called(id, title, description)
	return args.Error(0)
}

func (m *MockService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestHandler_Create(t *testing.T) {
	t.Run("Valid create request", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.POST("/todos", handler.Create)

		todo := domain.Todo{
			ID:          "1",
			Title:       "Test Todo",
			Description: "Test Description",
		}

		mockService.On("Create", todo).Return(nil)

		reqBody := CreateTodoRequest{
			ID:          "1",
			Title:       "Test Todo",
			Description: "Test Description",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.POST("/todos", handler.Create)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing required fields", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.POST("/todos", handler.Create)

		reqBody := CreateTodoRequest{
			Description: "Test Description",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service error", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.POST("/todos", handler.Create)

		todo := domain.Todo{
			ID:          "1",
			Title:       "Test Todo",
			Description: "Test Description",
		}

		mockService.On("Create", todo).Return(fmt.Errorf("service error"))

		reqBody := CreateTodoRequest{
			ID:          "1",
			Title:       "Test Todo",
			Description: "Test Description",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHandler_Get(t *testing.T) {
	t.Run("Get existing todo", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.GET("/todos/:id", handler.Get)

		expectedTodo := &domain.Todo{
			ID:          "1",
			Title:       "Test Todo",
			Description: "Test Description",
		}

		mockService.On("Get", "1").Return(expectedTodo, nil)

		req, _ := http.NewRequest("GET", "/todos/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		assert.Equal(t, expectedTodo.ID, response.ID)
		assert.Equal(t, expectedTodo.Title, response.Title)

		mockService.AssertExpectations(t)
	})

	t.Run("Get non-existing todo", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.GET("/todos/:id", handler.Get)

		mockService.On("Get", "999").Return(nil, fmt.Errorf("not found"))

		req, _ := http.NewRequest("GET", "/todos/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHandler_Update(t *testing.T) {
	t.Run("Valid update", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.PUT("/todos/:id", handler.Update)

		mockService.On("Update", "1", "Updated Title", "Updated Description").Return(nil)

		reqBody := UpdateTodoRequest{
			Title:       "Updated Title",
			Description: "Updated Description",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("PUT", "/todos/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.PUT("/todos/:id", handler.Update)

		req, _ := http.NewRequest("PUT", "/todos/1", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service error", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.PUT("/todos/:id", handler.Update)

		mockService.On("Update", "1", "Title", "Description").Return(fmt.Errorf("update failed"))

		reqBody := UpdateTodoRequest{
			Title:       "Title",
			Description: "Description",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("PUT", "/todos/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("Valid delete", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.DELETE("/todos/:id", handler.Delete)

		mockService.On("Delete", "1").Return(nil)

		req, _ := http.NewRequest("DELETE", "/todos/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Service error", func(t *testing.T) {
		mockService := new(MockService)
		handler := NewHandler(mockService)
		router := setupRouter()

		router.DELETE("/todos/:id", handler.Delete)

		mockService.On("Delete", "999").Return(fmt.Errorf("not found"))

		req, _ := http.NewRequest("DELETE", "/todos/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHealthCheckHandler(t *testing.T) {
	router := setupRouter()
	router.GET("/health", HealthCheckHandler)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "ok", response["status"])
}
