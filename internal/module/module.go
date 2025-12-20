package module

import (
	"os/exec"
	"time"
)

// Status represents the current state of a module.
type Status string

const (
	StatusRunning Status = "RUNNING"
	StatusStopped Status = "STOPPED"
	StatusFailed  Status = "FAILED"
)

// Module represents an external binary managed by the agent.
type Module struct {
	Name       string
	BinaryPath string
	Args       []string
	Status     Status
	LastStart  time.Time
	Restart    bool // Should this module be restarted if it fails?

	cmd *exec.Cmd
}

// NewModule creates a new module definition.
func NewModule(name, binaryPath string, args []string) *Module {
	return &Module{
		Name:       name,
		BinaryPath: binaryPath,
		Args:       args,
		Status:     StatusStopped,
		Restart:    true,
	}
}
