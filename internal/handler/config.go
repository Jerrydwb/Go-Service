package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kardex-pdf-service/internal/config"
)

// ConfigHandler permite actualizar la configuración en runtime.
func ConfigHandler(cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var update config.ConfigUpdate
		if err := c.BindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cfg.Update(update)

		c.JSON(http.StatusOK, gin.H{
			"message": "Configuración actualizada",
			"config":  cfg.ToMap(),
		})
	}
}
