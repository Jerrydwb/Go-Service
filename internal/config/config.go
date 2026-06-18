package config

// AppConfig contiene los límites de procesamiento optimizados para memoria.
type AppConfig struct {
	ItemsPerBatch        int  // Máximo de insumos por batch
	MaxMovementsPerBatch int  // Máximo de movimientos por batch
	RowsPerTableChunk    int  // Filas máximas antes de dividir tabla
	SingleFileThreshold  int  // Movimientos máximos para archivo único
	MergeThreshold       int  // Hasta este límite se hace merge de PDFs
	ZipCompressionLevel  int  // Nivel de compresión (1-9)
	ProcessCleanupDelay  int  // Delay antes de limpiar proceso (ms)
	UseZip               bool // Si es true, usa estrategia ZIP para datasets grandes
	Port                 string
}

// Default devuelve la configuración por defecto.
func Default() *AppConfig {
	return &AppConfig{
		ItemsPerBatch:        50,
		MaxMovementsPerBatch: 1500,
		RowsPerTableChunk:    500,
		SingleFileThreshold:  1000,
		MergeThreshold:       5000,
		ZipCompressionLevel:  6,
		ProcessCleanupDelay:  5000,
		UseZip:               false,
		Port:                 "8080",
	}
}

// Update aplica campos no-nil de un partial config.
func (c *AppConfig) Update(update ConfigUpdate) {
	if update.ItemsPerBatch != nil {
		c.ItemsPerBatch = *update.ItemsPerBatch
	}
	if update.MaxMovementsPerBatch != nil {
		c.MaxMovementsPerBatch = *update.MaxMovementsPerBatch
	}
	if update.SingleFileThreshold != nil {
		c.SingleFileThreshold = *update.SingleFileThreshold
	}
	if update.MergeThreshold != nil {
		c.MergeThreshold = *update.MergeThreshold
	}
	if update.UseZip != nil {
		c.UseZip = *update.UseZip
	}
}

// ConfigUpdate representa campos actualizables vía API.
type ConfigUpdate struct {
	ItemsPerBatch        *int  `json:"itemsPerBatch,omitempty"`
	MaxMovementsPerBatch *int  `json:"maxMovementsPerBatch,omitempty"`
	SingleFileThreshold  *int  `json:"singleFileThreshold,omitempty"`
	MergeThreshold       *int  `json:"mergeThreshold,omitempty"`
	UseZip               *bool `json:"useZip,omitempty"`
}

// ToMap convierte la config en un map para responses JSON.
func (c *AppConfig) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"items_per_batch":         c.ItemsPerBatch,
		"max_movements_per_batch": c.MaxMovementsPerBatch,
		"single_file_threshold":   c.SingleFileThreshold,
		"merge_threshold":         c.MergeThreshold,
		"use_zip":                 c.UseZip,
	}
}
