package main

import (
	"log"
	"os"

	"docger/internal/handler"
	"docger/internal/repository"
	"docger/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "docger_secret_key_2024"
	}

	// Initialize services
	userService := service.NewUserService(userRepo, jwtSecret)
	ticketService := service.NewTicketService(ticketRepo, userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	ticketHandler := handler.NewTicketHandler(ticketService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public routes
	r.POST("/api/auth/register", userHandler.Register)
	r.POST("/api/auth/login", userHandler.Login)

	// Protected routes
	authMiddleware := handler.AuthMiddleware(userService)
	protected := r.Group("/api")
	protected.Use(authMiddleware)
	{
		// User routes
		protected.GET("/profile", userHandler.GetProfile)

		// Ticket routes
		protected.POST("/tickets", ticketHandler.CreateTicket)
		protected.GET("/tickets/my", ticketHandler.GetMyTickets)
		protected.GET("/tickets/:id", ticketHandler.GetTicket)
		protected.PUT("/tickets/:id", ticketHandler.UpdateTicket)
		protected.DELETE("/tickets/:id", ticketHandler.DeleteTicket)

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(handler.AdminMiddleware())
		{
			admin.GET("/users", userHandler.GetAllUsers)
			admin.GET("/tickets", ticketHandler.GetAllTickets)
			admin.POST("/tickets/:id/assign", ticketHandler.AssignTicket)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
