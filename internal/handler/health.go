package handler

import (
	"runtime"

	"github.com/gin-gonic/gin"

	"kardex-pdf-service/internal/config"
)

// HealthHandler responde con estado del servicio y métricas del runtime.
func HealthHandler(cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		c.JSON(200, gin.H{
			"status":         "ok",
			"version":        "2.0.0",
			"alloc_mb":       m.Alloc / 1024 / 1024,
			"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
			"sys_mb":         m.Sys / 1024 / 1024,
			"num_gc":         m.NumGC,
			"num_goroutine":  runtime.NumGoroutine(),
			"config":         cfg.ToMap(),
		})
	}
}
