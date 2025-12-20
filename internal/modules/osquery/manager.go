package osquery

import (
	"log"
	"os/exec"
	"path/filepath"
)

// Manager handles the lifecycle of the embedded osquery process.
type Manager struct {
	BinaryPath string
}

// NewManager creates a new Osquery manager.
func NewManager(binaryPath string) *Manager {
	return &Manager{
		BinaryPath: binaryPath,
	}
}

// Start launches osqueryd in extension mode (daemon).
func (m *Manager) Start() error {
	absPath, err := filepath.Abs(m.BinaryPath)
	if err != nil {
		return err
	}

	log.Printf("Starting Osquery from: %s", absPath)

	// Example: osqueryd --pidfile ... --config_path ...
	cmd := exec.Command(absPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Osquery start check failed: %v, Output: %s", err, string(output))
		// For now, we don't return error to keep the agent running during dev without the binary
		return nil
	}

	log.Printf("Osquery detected: %s", string(output))
	return nil
}

// TODO: function to run ad-hoc queries
// func (m *Manager) Query(q string) ([]map[string]string, error) { ... }
