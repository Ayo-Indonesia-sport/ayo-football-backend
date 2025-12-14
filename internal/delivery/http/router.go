package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zenkriztao/ayo-football-backend/internal/delivery/http/handler"
	"github.com/zenkriztao/ayo-football-backend/internal/delivery/http/middleware"
	"github.com/zenkriztao/ayo-football-backend/internal/infrastructure/security"
)

// Router holds all HTTP handlers
type Router struct {
	authHandler   *handler.AuthHandler
	teamHandler   *handler.TeamHandler
	playerHandler *handler.PlayerHandler
	matchHandler  *handler.MatchHandler
	reportHandler *handler.ReportHandler
	jwtService    security.JWTService
}

// NewRouter creates a new Router instance
func NewRouter(
	authHandler *handler.AuthHandler,
	teamHandler *handler.TeamHandler,
	playerHandler *handler.PlayerHandler,
	matchHandler *handler.MatchHandler,
	reportHandler *handler.ReportHandler,
	jwtService security.JWTService,
) *Router {
	return &Router{
		authHandler:   authHandler,
		teamHandler:   teamHandler,
		playerHandler: playerHandler,
		matchHandler:  matchHandler,
		reportHandler: reportHandler,
		jwtService:    jwtService,
	}
}

// Setup configures all routes
func (r *Router) Setup(engine *gin.Engine) {
	// Global middlewares
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.RecoveryMiddleware())

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "ayo-football-api",
		})
	})

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/register", r.authHandler.Register)
		}

		// Protected auth routes
		authProtected := v1.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(r.jwtService))
		{
			authProtected.GET("/profile", r.authHandler.GetProfile)
		}

		// Team routes
		teams := v1.Group("/teams")
		{
			// Public routes
			teams.GET("", r.teamHandler.GetAll)
			teams.GET("/:id", r.teamHandler.GetByID)

			// Protected routes (Admin only)
			teamsAdmin := teams.Group("")
			teamsAdmin.Use(middleware.AuthMiddleware(r.jwtService))
			teamsAdmin.Use(middleware.AdminMiddleware())
			{
				teamsAdmin.POST("", r.teamHandler.Create)
				teamsAdmin.PUT("/:id", r.teamHandler.Update)
				teamsAdmin.DELETE("/:id", r.teamHandler.Delete)
			}
		}

		// Player routes
		players := v1.Group("/players")
		{
			// Public routes
			players.GET("", r.playerHandler.GetAll)
			players.GET("/:id", r.playerHandler.GetByID)

			// Protected routes (Admin only)
			playersAdmin := players.Group("")
			playersAdmin.Use(middleware.AuthMiddleware(r.jwtService))
			playersAdmin.Use(middleware.AdminMiddleware())
			{
				playersAdmin.POST("", r.playerHandler.Create)
				playersAdmin.PUT("/:id", r.playerHandler.Update)
				playersAdmin.DELETE("/:id", r.playerHandler.Delete)
			}
		}

		// Match routes
		matches := v1.Group("/matches")
		{
			// Public routes
			matches.GET("", r.matchHandler.GetAll)
			matches.GET("/:id", r.matchHandler.GetByID)

			// Protected routes (Admin only)
			matchesAdmin := matches.Group("")
			matchesAdmin.Use(middleware.AuthMiddleware(r.jwtService))
			matchesAdmin.Use(middleware.AdminMiddleware())
			{
				matchesAdmin.POST("", r.matchHandler.Create)
				matchesAdmin.PUT("/:id", r.matchHandler.Update)
				matchesAdmin.DELETE("/:id", r.matchHandler.Delete)
				matchesAdmin.POST("/:id/result", r.matchHandler.RecordResult)
			}
		}

		// Report routes (public)
		reports := v1.Group("/reports")
		{
			reports.GET("/matches", r.reportHandler.GetAllMatchReports)
			reports.GET("/matches/:id", r.reportHandler.GetMatchReport)
			reports.GET("/top-scorers", r.reportHandler.GetTopScorers)
		}
	}
}
