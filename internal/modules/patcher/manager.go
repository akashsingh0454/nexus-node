package patcher

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"nexus/internal/safety"
)

type Manager struct {
	DryRun      bool
	SafetyChain *safety.Chain
}

func NewManager(dryRun bool) *Manager {
	// Default safety chain: Patching allowed 24/7 (0-24) and CPU < 90%
	chain := safety.NewChain(
		&safety.TimeWindowCheck{AllowedStartHour: 0, AllowedEndHour: 24},
		&safety.CPUCheck{MaxUsage: 90.0},
	)

	return &Manager{
		DryRun:      dryRun,
		SafetyChain: chain,
	}
}

// ListUpdates checks for available updates.
func (m *Manager) ListUpdates() ([]string, error) {
	// winget upgrade
	log.Println("Checking for updates via winget...")
	cmd := exec.Command("winget", "upgrade", "--include-unknown")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Winget often returns non-zero exit codes for "no updates found" etc.
		// log.Printf("winget warning: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var updates []string
	for _, line := range lines {
		// Simplistic parsing: just capture the output for now
		if strings.TrimSpace(line) != "" {
			updates = append(updates, line)
		}
	}
	return updates, nil
}

// InstallUpdate attempts to upgrade a specific package ID.
func (m *Manager) InstallUpdate(packageID string) error {
	// 1. Run Safety Checks
	if err := m.SafetyChain.Execute(); err != nil {
		log.Printf("[SAFETY BLOCK] Skipping patch for %s: %v", packageID, err)
		return err
	}

	// 2. Execute (or Dry Run)
	args := []string{"upgrade", "--id", packageID, "--silent", "--accept-source-agreements", "--accept-package-agreements"}

	if m.DryRun {
		log.Printf("[DRY RUN] Would execute: winget %s", strings.Join(args, " "))
		return nil
	}

	log.Printf("Installing patch for %s...", packageID)
	cmd := exec.Command("winget", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("winget failed: %v, output: %s", err, string(output))
	}

	log.Printf("Successfully patched %s", packageID)
	return nil
}
