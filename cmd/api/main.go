package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zenkriztao/ayo-football-backend/internal/config"
	httpDelivery "github.com/zenkriztao/ayo-football-backend/internal/delivery/http"
	"github.com/zenkriztao/ayo-football-backend/internal/delivery/http/handler"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/usecase"
	"github.com/zenkriztao/ayo-football-backend/internal/infrastructure/database"
	"github.com/zenkriztao/ayo-football-backend/internal/infrastructure/security"
)

func main() {
	log.Println("=== AYO Football API Starting ===")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Server Port: %s", cfg.Server.Port)
	log.Printf("Database Driver: %s", cfg.Database.Driver)
	log.Printf("GIN Mode: %s", cfg.Server.Mode)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	log.Println("Connecting to database...")
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully!")

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	teamRepo := database.NewTeamRepository(db)
	playerRepo := database.NewPlayerRepository(db)
	matchRepo := database.NewMatchRepository(db)
	goalRepo := database.NewGoalRepository(db)

	// Initialize services
	jwtService := security.NewJWTService(cfg)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtService)
	teamUseCase := usecase.NewTeamUseCase(teamRepo)
	playerUseCase := usecase.NewPlayerUseCase(playerRepo, teamRepo)
	matchUseCase := usecase.NewMatchUseCase(matchRepo, teamRepo, playerRepo, goalRepo)
	reportUseCase := usecase.NewReportUseCase(matchRepo, goalRepo, teamRepo)

	// Create default admin user
	ctx := context.Background()
	if err := authUseCase.CreateDefaultAdmin(ctx, cfg.Admin.Email, cfg.Admin.Password); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	} else {
		log.Printf("Default admin user ensured: %s", cfg.Admin.Email)
	}

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	teamHandler := handler.NewTeamHandler(teamUseCase)
	playerHandler := handler.NewPlayerHandler(playerUseCase)
	matchHandler := handler.NewMatchHandler(matchUseCase)
	reportHandler := handler.NewReportHandler(reportUseCase)

	// Initialize router
	router := httpDelivery.NewRouter(
		authHandler,
		teamHandler,
		playerHandler,
		matchHandler,
		reportHandler,
		jwtService,
	)

	// Setup Gin engine
	engine := gin.Default()
	router.Setup(engine)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
