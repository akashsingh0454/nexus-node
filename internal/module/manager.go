package module

import (
	"context"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Manager handles the lifecycle of modules.
type Manager struct {
	modules map[string]*Module
	mu      sync.RWMutex
}

// NewManager creates a new Module Manager.
func NewManager() *Manager {
	return &Manager{
		modules: make(map[string]*Module),
	}
}

// Register adds a new module to the manager.
func (m *Manager) Register(mod *Module) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.modules[mod.Name] = mod
}

// StartAll starts all registered modules.
func (m *Manager) StartAll(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, mod := range m.modules {
		go m.runModule(ctx, mod)
	}
}

// StopAll stops all running modules.
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, mod := range m.modules {
		if mod.cmd != nil && mod.cmd.Process != nil {
			log.Printf("Stopping module: %s", mod.Name)
			// Try graceful shutdown first
			mod.cmd.Process.Signal(os.Interrupt)

			// Wait briefly then kill if needed (simplified for now)
			time.Sleep(1 * time.Second)
			mod.cmd.Process.Kill()
		}
	}
}

func (m *Manager) runModule(ctx context.Context, mod *Module) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Printf("Starting module: %s (%s)", mod.Name, mod.BinaryPath)
			mod.cmd = exec.CommandContext(ctx, mod.BinaryPath, mod.Args...)
			mod.cmd.Stdout = os.Stdout
			mod.cmd.Stderr = os.Stderr
			mod.Status = StatusRunning
			mod.LastStart = time.Now()

			err := mod.cmd.Run()

			mod.Status = StatusFailed
			log.Printf("Module %s exited: %v", mod.Name, err)

			if !mod.Restart || ctx.Err() != nil {
				log.Printf("Module %s will not restart.", mod.Name)
				mod.Status = StatusStopped
				return
			}

			// Backoff before restart
			log.Printf("Restarting module %s in 5 seconds...", mod.Name)
			time.Sleep(5 * time.Second)
		}
	}
}
