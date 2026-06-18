package server

import (
	"log"

	"github.com/gin-gonic/gin"

	"kardex-pdf-service/internal/config"
	"kardex-pdf-service/internal/handler"
	"kardex-pdf-service/internal/server/middleware"
	"kardex-pdf-service/internal/tracker"
	kardexHandler "kardex-pdf-service/services/kardex"
)

// Server envuelve el router Gin y sus dependencias.
type Server struct {
	Router  *gin.Engine
	Config  *config.AppConfig
	Tracker *tracker.ProcessTracker
}

// New crea y configura el servidor con todas las rutas.
func New(cfg *config.AppConfig) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORS())

	trk := tracker.New(cfg.ProcessCleanupDelay)

	s := &Server{
		Router:  r,
		Config:  cfg,
		Tracker: trk,
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	// Health check
	s.Router.GET("/health", handler.HealthHandler(s.Config))

	// API routes
	api := s.Router.Group("/api")
	{
		// Kardex PDF service
		kardexGroup := api.Group("/pdf")
		kardexHandler.RegisterRoutes(kardexGroup, s.Config, s.Tracker)

		// Config
		api.PUT("/config", handler.ConfigHandler(s.Config))
	}
}

// Run inicia el servidor HTTP.
func (s *Server) Run() error {
	log.Printf("🚀 Kardex PDF Service v2.0 iniciado en puerto %s", s.Config.Port)
	log.Printf("📊 Config: ItemsPerBatch=%d, MaxMovementsPerBatch=%d, SingleFileThreshold=%d, MergeThreshold=%d",
		s.Config.ItemsPerBatch, s.Config.MaxMovementsPerBatch, s.Config.SingleFileThreshold, s.Config.MergeThreshold)

	if err := s.Router.Run(":" + s.Config.Port); err != nil {
		return err
	}
	return nil
}
