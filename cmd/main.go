package main

import (
	"admin/internal/chat"
	"admin/internal/chat/adapter"
	"admin/pkg/apperror"
	"admin/pkg/database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Initialize database connection (example with PostgreSQL)
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "postgres")
	os.Setenv("DB_HOST", "172.233.25.33") // Use your database IP address
	os.Setenv("DB_PORT", "5432")

	// Initialize database connection pool and run migrations
	if err := database.InitPool(); err != nil {
		log.Fatal("Failed to initialize database pool:", err)
	}
	defer database.ClosePool()
	db := database.GetPool() // Retrieve the initialized database pool
	// Initialize repository, AI chat adapter, and chat service
	repo := adapter.NewSQLChatSessionRepository(db)
	ai := adapter.NewChatAdapter("http://localhost") // Update the URL to your AI service if needed
	service := chat.NewService(repo, ai, "EAAEccmazn6IBO0td4CEL98CA7w3NSjst4Iz85ZB7pjDbZBo7IlF2eKh9a628IJxvjMA4KmCwtZCgTTePrExe58LZB1szK4VDlpjZBoo0TY61n75MQKuAL93romYWGb9VfmQEpNPULipanwAagdds34V4BpsifWDPWo3rviSgMR8vPaQCKuonUcNHOZASrQk7qgGlq7KEpqU43CYJBiCPqX1AZDZD")

	// Setup HTTP server using Fiber
	app := fiber.New()
	app.Use(apperror.ErrorHandler) // Add error handling middleware

	// Initialize and register chat handler with routes
	handler := chat.NewHandler(*service)
	handler.RegisterRoutes(app)

	// Start the server
	log.Println("Starting server on port 8080...")
	if err := app.Listen(":8080"); err != nil {
		log.Panic("Error starting server:", err)
	}
}
