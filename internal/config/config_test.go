package config

import (
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.ItemsPerBatch != 50 {
		t.Errorf("Default ItemsPerBatch = %d, want 50", cfg.ItemsPerBatch)
	}
	if cfg.MaxMovementsPerBatch != 1500 {
		t.Errorf("Default MaxMovementsPerBatch = %d, want 1500", cfg.MaxMovementsPerBatch)
	}
	if cfg.SingleFileThreshold != 1000 {
		t.Errorf("Default SingleFileThreshold = %d, want 1000", cfg.SingleFileThreshold)
	}
	if cfg.UseZip != false {
		t.Errorf("Default UseZip = %v, want false", cfg.UseZip)
	}
	if cfg.Port != "8080" {
		t.Errorf("Default Port = %q, want %q", cfg.Port, "8080")
	}
}

func TestUpdate(t *testing.T) {
	cfg := Default()

	itemsPerBatch := 100
	useZip := true

	cfg.Update(ConfigUpdate{
		ItemsPerBatch: &itemsPerBatch,
		UseZip:        &useZip,
	})

	if cfg.ItemsPerBatch != 100 {
		t.Errorf("Update ItemsPerBatch = %d, want 100", cfg.ItemsPerBatch)
	}
	if cfg.UseZip != true {
		t.Errorf("Update UseZip = %v, want true", cfg.UseZip)
	}
	// Unchanged fields should remain
	if cfg.MaxMovementsPerBatch != 1500 {
		t.Errorf("Unchanged MaxMovementsPerBatch = %d, want 1500", cfg.MaxMovementsPerBatch)
	}
}

func TestUpdateNilFields(t *testing.T) {
	cfg := Default()
	original := cfg.ItemsPerBatch

	cfg.Update(ConfigUpdate{}) // all nil

	if cfg.ItemsPerBatch != original {
		t.Errorf("Update with nil changed ItemsPerBatch from %d to %d", original, cfg.ItemsPerBatch)
	}
}

func TestToMap(t *testing.T) {
	cfg := Default()
	m := cfg.ToMap()

	if m["items_per_batch"] != 50 {
		t.Errorf("ToMap items_per_batch = %v, want 50", m["items_per_batch"])
	}
	if m["use_zip"] != false {
		t.Errorf("ToMap use_zip = %v, want false", m["use_zip"])
	}
}
