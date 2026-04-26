package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jahapanah123/todo/internal/domain"
)

type Service interface {
	Create(todo domain.Todo) error
	Get(id string) (*domain.Todo, error)
	Update(id string, title string, description string) error
	Delete(id string) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

type CreateTodoRequest struct {
	ID          string `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo := domain.Todo{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
	}
	if err := h.service.Create(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	todo, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Update(id, req.Title, req.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
