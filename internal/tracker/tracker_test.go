package tracker

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	trk := New(0) // no cleanup
	defer trk.Stop()

	if trk == nil {
		t.Fatal("New returned nil")
	}
}

func TestUpdateAndGet(t *testing.T) {
	trk := New(0)
	defer trk.Stop()

	trk.Update("test-1", 50, "in-progress", "Procesando...")

	p := trk.Get("test-1")
	if p == nil {
		t.Fatal("Get returned nil")
	}
	if p.Percent != 50 {
		t.Errorf("Percent = %d, want 50", p.Percent)
	}
	if p.Status != "in-progress" {
		t.Errorf("Status = %q, want %q", p.Status, "in-progress")
	}
	if p.Message != "Procesando..." {
		t.Errorf("Message = %q, want %q", p.Message, "Procesando...")
	}
}

func TestGetNotFound(t *testing.T) {
	trk := New(0)
	defer trk.Stop()

	p := trk.Get("nonexistent")
	if p != nil {
		t.Errorf("Get nonexistent = %v, want nil", p)
	}
}

func TestDelete(t *testing.T) {
	trk := New(0)
	defer trk.Stop()

	trk.Update("test-1", 100, "completed", "Done")
	trk.Delete("test-1")

	p := trk.Get("test-1")
	if p != nil {
		t.Errorf("Get after Delete = %v, want nil", p)
	}
}

func TestUpdateOverwrite(t *testing.T) {
	trk := New(0)
	defer trk.Stop()

	trk.Update("test-1", 50, "in-progress", "Halfway")
	trk.Update("test-1", 100, "completed", "Done")

	p := trk.Get("test-1")
	if p.Percent != 100 {
		t.Errorf("Percent after overwrite = %d, want 100", p.Percent)
	}
}

func TestAutoCleanup(t *testing.T) {
	// Cleanup each 50ms for testing
	trk := New(50)
	defer trk.Stop()

	trk.Update("completed-1", 100, "completed", "Done")
	trk.Update("inprogress-1", 50, "in-progress", "Running")

	// Wait for cleanup cycle
	time.Sleep(120 * time.Millisecond)

	// Completed should be cleaned
	if p := trk.Get("completed-1"); p != nil {
		t.Errorf("completed process should be cleaned, got %v", p)
	}

	// In-progress should remain
	if p := trk.Get("inprogress-1"); p == nil {
		t.Error("in-progress process should remain")
	}
}
