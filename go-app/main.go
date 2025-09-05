// SeedFlow Knowledge Management Tool - Main Application
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// AppConfig holds application configuration
type AppConfig struct {
	Port          string
	AIServiceHost string
	AIServicePort string
	DataDir       string
	LogDir        string
	Debug         bool
}

// loadConfig loads configuration from environment variables
func loadConfig() *AppConfig {
	// Load .env file if it exists
	godotenv.Load()

	config := &AppConfig{
		Port:          getEnv("KNOWLEDGE_APP_PORT", "8080"),
		AIServiceHost: getEnv("AI_SERVICE_HOST", "localhost"),
		AIServicePort: getEnv("AI_SERVICE_PORT", "8001"),
		DataDir:       getEnv("DATA_DIR", "./data"),
		LogDir:        getEnv("LOG_DIR", "./logs"),
		Debug:         getEnv("DEBUG", "false") == "true",
	}

	return config
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setupRouter sets up the Gin router
func setupRouter(config *AppConfig) *gin.Engine {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Routes
	api := r.Group("/api")
	{
		api.GET("/health", healthCheck(config))
		api.HEAD("/health", healthCheck(config))
		api.GET("/version", versionHandler)
		
		// Knowledge routes (placeholder)
		knowledge := api.Group("/knowledge")
		{
			knowledge.GET("/", listKnowledge)
			knowledge.POST("/", createKnowledge)
			knowledge.GET("/:id", getKnowledge)
			knowledge.PUT("/:id", updateKnowledge)
			knowledge.DELETE("/:id", deleteKnowledge)
		}

		// AI routes (proxy to AI service)
		ai := api.Group("/ai")
		{
			ai.POST("/process", proxyToAIService(config, "/ai/process"))
			ai.POST("/extract-url", proxyToAIService(config, "/ai/extract-url"))
		}
	}

	// Static files and UI routes
	r.Static("/static", "./static")
	
	// Try to load templates, fallback to JSON response if not found
	if _, err := os.Stat("templates"); err == nil {
		r.LoadHTMLGlob("templates/*")
		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "SeedFlow Knowledge Management",
			})
		})
	} else {
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"name":    "SeedFlow Knowledge Management",
				"version": "1.0.0",
				"status":  "running",
				"api":     "/api/health",
			})
		})
	}

	return r
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// healthCheck returns application health status
func healthCheck(config *AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set health status headers for both GET and HEAD requests
		c.Header("X-Health-Status", "healthy")
		c.Header("X-Version", "1.0.0")
		
		// For HEAD requests, just return 200 status
		if c.Request.Method == "HEAD" {
			c.Status(http.StatusOK)
			return
		}
		
		// Check AI service connectivity (only for GET requests)
		aiURL := fmt.Sprintf("http://%s:%s/ai/health", config.AIServiceHost, config.AIServicePort)
		client := &http.Client{Timeout: 5 * time.Second}
		
		aiStatus := "unavailable"
		if resp, err := client.Get(aiURL); err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				aiStatus = "healthy"
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     "healthy",
			"version":    "1.0.0",
			"timestamp":  time.Now().Format(time.RFC3339),
			"ai_service": aiStatus,
			"database":   "connected", // TODO: Check actual database connection
		})
	}
}

// versionHandler returns application version
func versionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    "SeedFlow",
		"version": "1.0.0",
		"build":   "docker",
	})
}

// Placeholder handlers for knowledge management
func listKnowledge(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"knowledge": []map[string]interface{}{},
		"total":     0,
	})
}

func createKnowledge(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Knowledge created (placeholder)",
		"id":      1,
	})
}

func getKnowledge(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"title": "Sample Knowledge",
		"type":  "placeholder",
	})
}

func updateKnowledge(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Knowledge updated (placeholder)",
		"id":      id,
	})
}

func deleteKnowledge(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Knowledge deleted (placeholder)",
		"id":      id,
	})
}

// proxyToAIService forwards requests to the AI service
func proxyToAIService(config *AppConfig, endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple proxy implementation
		aiURL := fmt.Sprintf("http://%s:%s%s", config.AIServiceHost, config.AIServicePort, endpoint)
		
		c.JSON(http.StatusOK, gin.H{
			"message": "AI service proxy (placeholder)",
			"url":     aiURL,
		})
	}
}

func main() {
	config := loadConfig()

	// Create necessary directories
	os.MkdirAll(config.DataDir, 0755)
	os.MkdirAll(config.LogDir, 0755)

	// Setup router
	router := setupRouter(config)

	// Setup server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("🌱 SeedFlow server started on port %s", config.Port)
	log.Printf("   Web UI: http://localhost:%s", config.Port)
	log.Printf("   API: http://localhost:%s/api", config.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server shutdown complete")
}