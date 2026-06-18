package tracker

import (
	"sync"
	"time"
)

// ProcessProgress tracking de progreso de un proceso.
type ProcessProgress struct {
	Percent   int       `json:"percent"`
	Status    string    `json:"status"` // "in-progress" | "completed" | "error"
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"-"`
}

// ProcessTracker para seguimiento de procesos (thread-safe).
type ProcessTracker struct {
	mu        sync.RWMutex
	processes map[string]*ProcessProgress
	cleanupMs int // milisegundos antes de limpiar procesos completados
	stopCh    chan struct{}
}

// New crea una nueva instancia del tracker con auto-cleanup.
func New(cleanupMs int) *ProcessTracker {
	t := &ProcessTracker{
		processes: make(map[string]*ProcessProgress),
		cleanupMs: cleanupMs,
		stopCh:    make(chan struct{}),
	}

	if cleanupMs > 0 {
		go t.autoCleanup()
	}

	return t
}

// Update actualiza el progreso de un proceso.
func (pt *ProcessTracker) Update(key string, percent int, status, message string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.processes[key] = &ProcessProgress{
		Percent:   percent,
		Status:    status,
		Message:   message,
		UpdatedAt: time.Now(),
	}
}

// Get obtiene el progreso de un proceso.
func (pt *ProcessTracker) Get(key string) *ProcessProgress {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	return pt.processes[key]
}

// Delete elimina un proceso del tracker.
func (pt *ProcessTracker) Delete(key string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	delete(pt.processes, key)
}

// Stop detiene el goroutine de cleanup.
func (pt *ProcessTracker) Stop() {
	close(pt.stopCh)
}

func (pt *ProcessTracker) autoCleanup() {
	interval := time.Duration(pt.cleanupMs) * time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pt.mu.Lock()
			now := time.Now()
			for key, p := range pt.processes {
				if (p.Status == "completed" || p.Status == "error") && now.Sub(p.UpdatedAt) > interval {
					delete(pt.processes, key)
				}
			}
			pt.mu.Unlock()
		case <-pt.stopCh:
			return
		}
	}
}
