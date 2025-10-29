package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	config "github.com/tudormiron/medical-reports/internal/configs"
	"github.com/tudormiron/medical-reports/internal/services"
)

type Server struct {
	config   *config.Config
	handlers *Handlers
	router   *gin.Engine
}

func NewServer(cfg *config.Config, reportService *services.ReportService, referenceService *services.ReferenceService, authService *services.AuthService) *Server {
	handlers := NewHandlers(reportService, referenceService, authService)
	router := setupRouter(handlers)

	return &Server{
		config:   cfg,
		handlers: handlers,
		router:   router,
	}
}

func setupRouter(handlers *Handlers) *gin.Engine {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Apply middleware
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware())
	router.Use(CORSMiddleware())
	router.Use(ErrorHandlingMiddleware())

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication (public routes)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Reports
		reports := v1.Group("/reports")
		{
			reports.POST("", handlers.CreateReport)
			reports.GET("", handlers.ListReports)
			reports.GET("/:id", handlers.GetReport)
			reports.PUT("/:id/content", handlers.UpdateReportContent)
			reports.PUT("/:id/status", handlers.UpdateReportStatus)
			reports.DELETE("/:id", handlers.DeleteReport)
			reports.GET("/:id/versions", handlers.GetReportVersions)
		}

		// Reference data
		reference := v1.Group("/reference")
		{
			reference.GET("/icd10", handlers.SearchICD10)
			reference.GET("/medications", handlers.SearchMedications)
		}
	}

	return router
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	log.Printf("Starting server on %s", addr)
	return s.router.Run(addr)
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
