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
	// 1. Initialize database (Postgres/SQL)
	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("PostgreSQL connected")

	// 2. Configuration (Port & Secret)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "docger_secret_key_2024"
	}

	// 3. Initialize Layers (Dependency Injection)
	userRepo := repository.NewUserRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	userService := service.NewUserService(userRepo, jwtSecret)
	ticketService := service.NewTicketService(ticketRepo, userRepo)

	userHandler := handler.NewUserHandler(userService)
	ticketHandler := handler.NewTicketHandler(ticketService)

	// 4. Setup Gin router
	r := gin.Default()

	// Global Middleware - CORS
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

	// Public Routes
	r.POST("/api/auth/register", userHandler.Register)
	r.POST("/api/auth/login", userHandler.Login)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Protected Routes Group
	authMiddleware := handler.AuthMiddleware(userService)
	protected := r.Group("/api")
	protected.Use(authMiddleware)
	{
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

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
