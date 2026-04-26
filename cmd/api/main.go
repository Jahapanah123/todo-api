package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jahapanah123/todo/internal/config"
	"github.com/jahapanah123/todo/internal/controller"
	"github.com/jahapanah123/todo/internal/repository"
	"github.com/jahapanah123/todo/internal/service"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Initialize repository and service and handler
	repo := repository.NewInMemoryRepository()
	service := service.NewService(repo)
	handler := controller.NewHandler(service)

	// Set up Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", controller.HealthCheckHandler)

	// Todo endpoints
	todoGroup := router.Group("/todos")
	{
		todoGroup.POST("/", handler.Create)
		todoGroup.GET("/:id", handler.Get)
		todoGroup.PUT("/:id", handler.Update)
		todoGroup.DELETE("/:id", handler.Delete)
	}
	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	log.Printf("Starting server on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
